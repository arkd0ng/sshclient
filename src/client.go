package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

// SSHClient represents an SSH client connection
type SSHClient struct {
	config *ssh.ClientConfig
	client *ssh.Client
	host   string
	port   string
}

// NewSSHClient creates a new SSH client with password authentication
func NewSSHClient(host, port, user, password string) (*SSHClient, error) {
	// Get host key callback for verification
	hostKeyCallback, err := GetHostKeyCallback()
	if err != nil {
		return nil, fmt.Errorf("failed to setup host key verification: %w", err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: hostKeyCallback,
		Timeout:         10 * time.Second,
	}

	return &SSHClient{
		config: config,
		host:   host,
		port:   port,
	}, nil
}

// NewSSHClientWithKey creates a new SSH client with key-based authentication
func NewSSHClientWithKey(host, port, user, keyPath string) (*SSHClient, error) {
	// Read private key file
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	// Parse private key
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Get host key callback for verification
	hostKeyCallback, err := GetHostKeyCallback()
	if err != nil {
		return nil, fmt.Errorf("failed to setup host key verification: %w", err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hostKeyCallback,
		Timeout:         10 * time.Second,
	}

	return &SSHClient{
		config: config,
		host:   host,
		port:   port,
	}, nil
}

// Connect establishes the SSH connection
func (c *SSHClient) Connect() error {
	addr := net.JoinHostPort(c.host, c.port)
	client, err := ssh.Dial("tcp", addr, c.config)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}
	c.client = client
	return nil
}

// Close closes the SSH connection
func (c *SSHClient) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// RunCommand executes a single command on the remote server
func (c *SSHClient) RunCommand(cmd string) (string, error) {
	if c.client == nil {
		return "", fmt.Errorf("not connected")
	}

	session, err := c.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return string(output), fmt.Errorf("command failed: %w", err)
	}

	return string(output), nil
}

// StartInteractiveShell starts an interactive shell session
func (c *SSHClient) StartInteractiveShell() error {
	if c.client == nil {
		return fmt.Errorf("not connected")
	}

	session, err := c.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Get terminal size
	fd := int(os.Stdin.Fd())
	state, err := term.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("failed to make terminal raw: %w", err)
	}
	defer term.Restore(fd, state)

	w, h, err := term.GetSize(fd)
	if err != nil {
		w, h = 80, 24 // default size
	}

	// Request pseudo terminal
	if err := session.RequestPty("xterm-256color", h, w, modes); err != nil {
		return fmt.Errorf("request for pseudo terminal failed: %w", err)
	}

	// Set up I/O
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	// Start shell
	if err := session.Shell(); err != nil {
		return fmt.Errorf("failed to start shell: %w", err)
	}

	// Wait for session to finish
	if err := session.Wait(); err != nil {
		if _, ok := err.(*ssh.ExitError); ok {
			return nil // Normal exit
		}
		return fmt.Errorf("session error: %w", err)
	}

	return nil
}

// GetDefaultKeyPath returns the default SSH key path
func GetDefaultKeyPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	// Check common key locations
	keys := []string{
		filepath.Join(home, ".ssh", "id_rsa"),
		filepath.Join(home, ".ssh", "id_ed25519"),
		filepath.Join(home, ".ssh", "id_ecdsa"),
	}

	for _, key := range keys {
		if _, err := os.Stat(key); err == nil {
			return key
		}
	}

	return ""
}

// CopyFile copies a file to the remote server using SCP
func (c *SSHClient) CopyFile(localPath, remotePath string) error {
	if c.client == nil {
		return fmt.Errorf("not connected")
	}

	// Read local file
	data, err := ioutil.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("failed to read local file: %w", err)
	}

	// Get file info
	info, err := os.Stat(localPath)
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	// Create session
	session, err := c.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Set up stdin pipe
	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	// Start scp command
	filename := filepath.Base(localPath)
	cmd := fmt.Sprintf("scp -t %s", remotePath)
	if err := session.Start(cmd); err != nil {
		return fmt.Errorf("failed to start scp: %w", err)
	}

	// Send file header
	fmt.Fprintf(stdin, "C%#o %d %s\n", info.Mode().Perm(), len(data), filename)

	// Send file content
	if _, err := stdin.Write(data); err != nil {
		return fmt.Errorf("failed to write file data: %w", err)
	}

	// Send end marker
	fmt.Fprint(stdin, "\x00")
	stdin.Close()

	// Wait for completion
	if err := session.Wait(); err != nil {
		return fmt.Errorf("scp failed: %w", err)
	}

	return nil
}

// DownloadFile downloads a file from the remote server using SCP
func (c *SSHClient) DownloadFile(remotePath, localPath string) error {
	if c.client == nil {
		return fmt.Errorf("not connected")
	}

	// Create session
	session, err := c.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Set up stdout pipe
	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	// Start scp command
	cmd := fmt.Sprintf("scp -f %s", remotePath)
	if err := session.Start(cmd); err != nil {
		return fmt.Errorf("failed to start scp: %w", err)
	}

	// Read file content
	data, err := io.ReadAll(stdout)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Write to local file
	if err := ioutil.WriteFile(localPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write local file: %w", err)
	}

	// Wait for completion
	if err := session.Wait(); err != nil {
		return fmt.Errorf("scp failed: %w", err)
	}

	return nil
}
