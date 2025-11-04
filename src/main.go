package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

const (
	version = "1.2.1"
)

// parseUserHost parses "user@host" format and returns user, host
func parseUserHost(arg string) (user, host string, ok bool) {
	parts := strings.Split(arg, "@")
	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return parts[0], parts[1], true
	}
	return "", "", false
}

func main() {
	// Command line flags
	host := flag.String("host", "", "SSH server hostname or IP address")
	port := flag.String("port", "22", "SSH server port (default: 22)")
	user := flag.String("user", "", "SSH username")
	password := flag.String("password", "", "SSH password (not recommended, use -key instead)")
	keyPath := flag.String("key", "", "Path to SSH private key file")
	cmd := flag.String("cmd", "", "Command to execute on remote server")
	interactive := flag.Bool("i", false, "Start interactive shell session")
	showVersion := flag.Bool("version", false, "Show version information")

	// Check for profile command
	if len(os.Args) > 1 && os.Args[1] == "profile" {
		if err := HandleProfileCommand(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Check for @profile format
	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "@") {
		profileName := strings.TrimPrefix(os.Args[1], "@")
		profile, err := FindProfile(profileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Set values from profile
		*host = profile.Host
		*user = profile.User
		if profile.Port != "" {
			*port = profile.Port
		}
		if profile.Key != "" {
			*keyPath = profile.Key
		}

		// Handle password (encrypted or plain text)
		if profile.EncryptedPassword != "" {
			// Decrypt the password automatically
			decrypted, err := DecryptAuto(profile.EncryptedPassword)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error decrypting password: %v\n", err)
				os.Exit(1)
			}
			*password = decrypted
		} else if profile.Password != "" {
			// Legacy plain text password
			*password = profile.Password
		}

		// Process remaining arguments
		remainingArgs := os.Args[2:]
		newArgs := []string{os.Args[0]}
		cmdArgs := []string{}

		for i := 0; i < len(remainingArgs); i++ {
			arg := remainingArgs[i]

			if strings.HasPrefix(arg, "-") {
				newArgs = append(newArgs, arg)
				if i+1 < len(remainingArgs) && !strings.HasPrefix(remainingArgs[i+1], "-") {
					i++
					newArgs = append(newArgs, remainingArgs[i])
				}
			} else {
				cmdArgs = append(cmdArgs, arg)
			}
		}

		if len(cmdArgs) > 0 {
			*cmd = strings.Join(cmdArgs, " ")
		} else {
			*interactive = true
		}

		os.Args = newArgs
	} else if len(os.Args) > 1 && !strings.HasPrefix(os.Args[1], "-") {
		// Parse traditional SSH format: user@host [command...]
		if u, h, ok := parseUserHost(os.Args[1]); ok {
			// Found user@host format
			*user = u
			*host = h

			// Process remaining arguments
			remainingArgs := os.Args[2:]

			// Reconstruct os.Args for flag parsing
			newArgs := []string{os.Args[0]}
			cmdArgs := []string{}

			for i := 0; i < len(remainingArgs); i++ {
				arg := remainingArgs[i]

				// Check if it's a flag
				if strings.HasPrefix(arg, "-") {
					newArgs = append(newArgs, arg)
					// If flag has a value, add it too
					if i+1 < len(remainingArgs) && !strings.HasPrefix(remainingArgs[i+1], "-") {
						i++
						newArgs = append(newArgs, remainingArgs[i])
					}
				} else {
					// It's a command argument
					cmdArgs = append(cmdArgs, arg)
				}
			}

			// If there are command arguments, join them
			if len(cmdArgs) > 0 {
				*cmd = strings.Join(cmdArgs, " ")
			} else {
				// No command specified, default to interactive mode
				*interactive = true
			}

			// Update os.Args for flag parsing
			os.Args = newArgs
		}
	}

	// Custom usage function
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "SSH Client v%s - Pure Go SSH client\n\n", version)
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  sshclient @profile [command...]          # Use saved profile\n")
		fmt.Fprintf(os.Stderr, "  sshclient user@host [command...]         # Traditional SSH style\n")
		fmt.Fprintf(os.Stderr, "  sshclient [flags]                        # Flag-based style\n")
		fmt.Fprintf(os.Stderr, "  sshclient profile <command> [args]       # Manage profiles\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  # Profile style\n")
		fmt.Fprintf(os.Stderr, "  sshclient @myserver\n")
		fmt.Fprintf(os.Stderr, "  sshclient @myserver uptime\n\n")
		fmt.Fprintf(os.Stderr, "  # Traditional SSH style\n")
		fmt.Fprintf(os.Stderr, "  sshclient user@example.com\n")
		fmt.Fprintf(os.Stderr, "  sshclient user@example.com -key ~/.ssh/id_rsa\n")
		fmt.Fprintf(os.Stderr, "  sshclient user@example.com ls -la\n\n")
		fmt.Fprintf(os.Stderr, "  # Flag style\n")
		fmt.Fprintf(os.Stderr, "  sshclient -host example.com -user myuser -i\n")
		fmt.Fprintf(os.Stderr, "  sshclient -host example.com -user myuser -cmd \"uptime\"\n\n")
		fmt.Fprintf(os.Stderr, "  # Profile management\n")
		fmt.Fprintf(os.Stderr, "  sshclient profile add myserver\n")
		fmt.Fprintf(os.Stderr, "  sshclient profile list\n\n")
		fmt.Fprintf(os.Stderr, "For more info: sshclient profile help\n")
	}

	flag.Parse()

	// Show version
	if *showVersion {
		fmt.Printf("SSH Client v%s\n", version)
		fmt.Println("A simple, standalone SSH client written in pure Go")
		fmt.Println("No root/sudo permissions required")
		os.Exit(0)
	}

	// Validate required flags
	if *host == "" {
		fmt.Fprintln(os.Stderr, "Error: Connection target not specified")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "You can connect using:")
		fmt.Fprintln(os.Stderr, "  sshclient @profile           # Use a saved profile")
		fmt.Fprintln(os.Stderr, "  sshclient user@host          # Traditional SSH style")
		fmt.Fprintln(os.Stderr, "  sshclient -host H -user U    # Explicit flags")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Examples:")
		fmt.Fprintln(os.Stderr, "  sshclient @myserver")
		fmt.Fprintln(os.Stderr, "  sshclient root@example.com")
		fmt.Fprintln(os.Stderr, "  sshclient -host example.com -user root -i")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "For full help: sshclient -h")
		os.Exit(1)
	}

	if *user == "" {
		fmt.Fprintln(os.Stderr, "Error: Username not specified")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Specify username using:")
		fmt.Fprintln(os.Stderr, "  sshclient user@host          # In user@host format")
		fmt.Fprintln(os.Stderr, "  sshclient -host H -user U    # Using -user flag")
		fmt.Fprintln(os.Stderr, "  sshclient @profile           # From saved profile")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Examples:")
		fmt.Fprintln(os.Stderr, "  sshclient root@example.com")
		fmt.Fprintln(os.Stderr, "  sshclient -host example.com -user root -i")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "For full help: sshclient -h")
		os.Exit(1)
	}

	// Determine authentication method
	var client *SSHClient
	var err error

	if *keyPath != "" {
		// Key-based authentication (explicit key path)
		fmt.Printf("Connecting to %s@%s:%s using key authentication...\n", *user, *host, *port)
		client, err = NewSSHClientWithKey(*host, *port, *user, *keyPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create SSH client: %v\n", err)
			os.Exit(1)
		}
	} else if *password != "" {
		// Password authentication (from command line)
		fmt.Printf("Connecting to %s@%s:%s using password authentication...\n", *user, *host, *port)
		client, err = NewSSHClient(*host, *port, *user, *password)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create SSH client: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Try default SSH key first
		defaultKey := GetDefaultKeyPath()
		if defaultKey != "" {
			fmt.Printf("Trying default SSH key: %s\n", defaultKey)
			client, err = NewSSHClientWithKey(*host, *port, *user, defaultKey)
			if err == nil {
				fmt.Printf("Using key authentication with %s\n", defaultKey)
			}
		}

		// If no default key or key auth failed, prompt for password
		if client == nil {
			fmt.Printf("Password for %s@%s: ", *user, *host)
			passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println() // New line after password input
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to read password: %v\n", err)
				os.Exit(1)
			}
			client, err = NewSSHClient(*host, *port, *user, string(passwordBytes))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to create SSH client: %v\n", err)
				os.Exit(1)
			}
		}
	}

	// Connect to server
	if err := client.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	fmt.Println("Connected successfully!")

	// Execute command or start interactive shell
	if *interactive {
		// Interactive shell
		fmt.Println("Starting interactive shell... (Press Ctrl+D or type 'exit' to quit)")
		if err := client.StartInteractiveShell(); err != nil {
			fmt.Fprintf(os.Stderr, "Shell session failed: %v\n", err)
			os.Exit(1)
		}
	} else if *cmd != "" {
		// Execute single command
		output, err := client.RunCommand(*cmd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Command failed: %v\n", err)
			fmt.Print(output) // Print output even if command failed
			os.Exit(1)
		}
		fmt.Print(output)
	} else {
		// No command specified, show help
		fmt.Println("\nNo command or interactive mode specified.")
		fmt.Println("Use -i for interactive shell or provide a command.")
		fmt.Println("\nExamples:")
		fmt.Println("  Profile style:")
		fmt.Printf("    sshclient @myserver\n")
		fmt.Printf("    sshclient @myserver ls -la\n\n")
		fmt.Println("  Traditional SSH style:")
		fmt.Printf("    sshclient user@example.com\n")
		fmt.Printf("    sshclient user@example.com ls -la\n")
		fmt.Printf("    sshclient user@example.com -key ~/.ssh/id_rsa\n\n")
		fmt.Println("  Flag style:")
		fmt.Printf("    sshclient -host example.com -user myuser -i\n")
		fmt.Printf("    sshclient -host example.com -user myuser -cmd \"ls -la\"\n\n")
		fmt.Println("  Profile management:")
		fmt.Printf("    sshclient profile add myserver\n")
		fmt.Printf("    sshclient profile list\n")
		fmt.Printf("    sshclient profile remove myserver\n\n")
	}
}

// Helper functions

func promptYesNo(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s (y/n): ", prompt)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

func printBanner() {
	fmt.Println("╔═══════════════════════════════════════════╗")
	fmt.Printf("║   SSH Client v%-27s ║\n", version)
	fmt.Println("║   Pure Go Implementation                  ║")
	fmt.Println("║   No root/sudo required                   ║")
	fmt.Println("╚═══════════════════════════════════════════╝")
	fmt.Println()
}
