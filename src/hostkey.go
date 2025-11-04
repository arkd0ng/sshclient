package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

// GetKnownHostsPath returns the path to sshclient's known_hosts file
func GetKnownHostsPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "known_hosts"), nil
}

// InitKnownHosts initializes the known_hosts file
// If ~/.ssh/known_hosts exists, it copies it to ~/.sshclient/known_hosts
func InitKnownHosts() error {
	sshClientKnownHosts, err := GetKnownHostsPath()
	if err != nil {
		return fmt.Errorf("failed to get known_hosts path: %w", err)
	}

	// If sshclient's known_hosts already exists, nothing to do
	if _, err := os.Stat(sshClientKnownHosts); err == nil {
		return nil
	}

	// Check if ~/.ssh/known_hosts exists
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	sshKnownHosts := filepath.Join(home, ".ssh", "known_hosts")

	if _, err := os.Stat(sshKnownHosts); err == nil {
		// Copy ~/.ssh/known_hosts to ~/.sshclient/known_hosts
		fmt.Printf("📋 Copying existing known_hosts from ~/.ssh/known_hosts\n")

		input, err := os.ReadFile(sshKnownHosts)
		if err != nil {
			return fmt.Errorf("failed to read ~/.ssh/known_hosts: %w", err)
		}

		// Ensure config directory exists
		configDir := filepath.Dir(sshClientKnownHosts)
		if err := os.MkdirAll(configDir, 0700); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}

		if err := os.WriteFile(sshClientKnownHosts, input, 0600); err != nil {
			return fmt.Errorf("failed to write known_hosts: %w", err)
		}

		fmt.Printf("✅ Copied %d bytes to ~/.sshclient/known_hosts\n", len(input))
	} else {
		// Create empty known_hosts file
		configDir := filepath.Dir(sshClientKnownHosts)
		if err := os.MkdirAll(configDir, 0700); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}

		if err := os.WriteFile(sshClientKnownHosts, []byte{}, 0600); err != nil {
			return fmt.Errorf("failed to create known_hosts: %w", err)
		}
	}

	return nil
}

// GetHostKeyCallback returns a callback function for host key verification
func GetHostKeyCallback() (ssh.HostKeyCallback, error) {
	// Initialize known_hosts if needed
	if err := InitKnownHosts(); err != nil {
		return nil, err
	}

	knownHostsPath, err := GetKnownHostsPath()
	if err != nil {
		return nil, err
	}

	// Create a callback that handles unknown hosts
	callback, err := knownhosts.New(knownHostsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create known_hosts callback: %w", err)
	}

	// Wrap the callback to handle unknown hosts
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		err := callback(hostname, remote, key)
		if err == nil {
			// Host key is already known and matches
			return nil
		}

		// Check if this is an unknown host
		if keyErr, ok := err.(*knownhosts.KeyError); ok && len(keyErr.Want) == 0 {
			// Unknown host - ask user
			return handleUnknownHost(hostname, remote, key, knownHostsPath)
		}

		// Host key has changed or other error
		if keyErr, ok := err.(*knownhosts.KeyError); ok && len(keyErr.Want) > 0 {
			// Key mismatch - potential security issue
			return handleKeyMismatch(hostname, remote, key, keyErr, knownHostsPath)
		}

		return err
	}, nil
}

// handleUnknownHost prompts the user to accept a new host key
func handleUnknownHost(hostname string, remote net.Addr, key ssh.PublicKey, knownHostsPath string) error {
	fingerprint := ssh.FingerprintSHA256(key)

	fmt.Fprintf(os.Stderr, "\n⚠️  WARNING: Unknown host!\n")
	fmt.Fprintf(os.Stderr, "The authenticity of host '%s (%s)' can't be established.\n", hostname, remote.String())
	fmt.Fprintf(os.Stderr, "%s key fingerprint is %s\n", key.Type(), fingerprint)
	fmt.Fprintf(os.Stderr, "\n")

	// Ask user
	fmt.Fprintf(os.Stderr, "Are you sure you want to continue connecting (yes/no)? ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read user input: %w", err)
	}

	response = strings.ToLower(strings.TrimSpace(response))
	if response != "yes" {
		return fmt.Errorf("host key verification failed: user rejected")
	}

	// Add to known_hosts
	if err := addHostKey(hostname, remote, key, knownHostsPath); err != nil {
		return fmt.Errorf("failed to add host key: %w", err)
	}

	fmt.Fprintf(os.Stderr, "✅ Host '%s' added to known_hosts\n\n", hostname)
	return nil
}

// handleKeyMismatch handles the case where a host key has changed
func handleKeyMismatch(hostname string, remote net.Addr, key ssh.PublicKey, keyErr *knownhosts.KeyError, knownHostsPath string) error {
	fingerprint := ssh.FingerprintSHA256(key)

	fmt.Fprintf(os.Stderr, "\n❌ WARNING: HOST KEY HAS CHANGED!\n")
	fmt.Fprintf(os.Stderr, "IT IS POSSIBLE THAT SOMEONE IS DOING SOMETHING NASTY!\n")
	fmt.Fprintf(os.Stderr, "Someone could be eavesdropping on you right now (man-in-the-middle attack)!\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Host: %s (%s)\n", hostname, remote.String())
	fmt.Fprintf(os.Stderr, "Key type: %s\n", key.Type())
	fmt.Fprintf(os.Stderr, "New fingerprint: %s\n", fingerprint)
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Expected one of:\n")
	for _, wantKey := range keyErr.Want {
		fmt.Fprintf(os.Stderr, "  %s\n", ssh.FingerprintSHA256(wantKey.Key))
	}
	fmt.Fprintf(os.Stderr, "\n")

	// Ask user if they want to update (risky!)
	fmt.Fprintf(os.Stderr, "Do you want to update the host key? This is DANGEROUS! (yes/no)? ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read user input: %w", err)
	}

	response = strings.ToLower(strings.TrimSpace(response))
	if response != "yes" {
		return fmt.Errorf("host key verification failed: key mismatch")
	}

	// Remove old keys and add new one
	if err := removeHostKey(hostname, knownHostsPath); err != nil {
		return fmt.Errorf("failed to remove old host key: %w", err)
	}

	if err := addHostKey(hostname, remote, key, knownHostsPath); err != nil {
		return fmt.Errorf("failed to add new host key: %w", err)
	}

	fmt.Fprintf(os.Stderr, "⚠️  Host key updated for '%s'\n\n", hostname)
	return nil
}

// addHostKey adds a host key to the known_hosts file
func addHostKey(hostname string, remote net.Addr, key ssh.PublicKey, knownHostsPath string) error {
	// Format: hostname ssh-rsa AAAAB3...
	host := hostname
	// Also add IP address if different
	if tcpAddr, ok := remote.(*net.TCPAddr); ok {
		ip := tcpAddr.IP.String()
		if ip != hostname {
			host = hostname + "," + ip
		}
	}

	line := knownhosts.Line([]string{host}, key)

	// Append to file
	f, err := os.OpenFile(knownHostsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(line + "\n"); err != nil {
		return err
	}

	return nil
}

// removeHostKey removes all entries for a hostname from known_hosts
func removeHostKey(hostname string, knownHostsPath string) error {
	input, err := os.ReadFile(knownHostsPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(input), "\n")
	var newLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			newLines = append(newLines, line)
			continue
		}

		// Check if this line contains the hostname
		parts := strings.Fields(line)
		if len(parts) > 0 {
			hosts := strings.Split(parts[0], ",")
			match := false
			for _, h := range hosts {
				if h == hostname {
					match = true
					break
				}
			}
			if !match {
				newLines = append(newLines, line)
			}
		}
	}

	output := strings.Join(newLines, "\n")
	return os.WriteFile(knownHostsPath, []byte(output), 0600)
}
