package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"golang.org/x/term"
)


// ProfileAdd adds a new profile interactively
func ProfileAdd(name string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Creating new profile: %s\n\n", name)

	// Get host
	fmt.Print("Host (hostname or IP): ")
	host, _ := reader.ReadString('\n')
	host = strings.TrimSpace(host)
	if host == "" {
		return fmt.Errorf("host is required")
	}

	// Get user
	fmt.Print("User: ")
	user, _ := reader.ReadString('\n')
	user = strings.TrimSpace(user)
	if user == "" {
		return fmt.Errorf("user is required")
	}

	// Get port
	fmt.Print("Port (default: 22): ")
	port, _ := reader.ReadString('\n')
	port = strings.TrimSpace(port)
	if port == "" {
		port = "22"
	}

	// Ask for authentication method
	fmt.Print("\nAuthentication method (1: SSH key, 2: Password): ")
	method, _ := reader.ReadString('\n')
	method = strings.TrimSpace(method)

	profile := Profile{
		Name: name,
		Host: host,
		User: user,
		Port: port,
	}

	if method == "1" {
		// SSH key
		defaultKey := GetDefaultKeyPath()
		if defaultKey != "" {
			fmt.Printf("Path to SSH key (default: %s): ", defaultKey)
		} else {
			fmt.Print("Path to SSH key (e.g., ~/.ssh/id_rsa): ")
		}
		key, _ := reader.ReadString('\n')
		key = strings.TrimSpace(key)
		if key != "" {
			profile.Key = key
		} else if defaultKey != "" {
			// Use default key if user pressed enter
			profile.Key = defaultKey
			fmt.Printf("Using default key: %s\n", defaultKey)
		}
	} else if method == "2" {
		// Password
		fmt.Print("Password (leave empty to prompt on connect): ")
		passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
		fmt.Println()
		password := string(passwordBytes)
		if password != "" {
			fmt.Println("\n🔐 Password will be encrypted using AES-256-GCM")

			// Encrypt the password automatically
			encrypted, err := EncryptAuto(password)
			if err != nil {
				return fmt.Errorf("failed to encrypt password: %w", err)
			}

			profile.EncryptedPassword = encrypted
			fmt.Println("✅ Password encrypted and will be stored securely")
		}
	}

	// Save profile
	if err := AddProfile(profile); err != nil {
		return fmt.Errorf("failed to save profile: %w", err)
	}

	fmt.Printf("\n✓ Profile '%s' created successfully!\n", name)
	return nil
}

// ProfileList lists all available profiles
func ProfileList() error {
	// Get custom profiles
	customProfiles, err := LoadProfiles()
	if err != nil {
		return fmt.Errorf("failed to load custom profiles: %w", err)
	}

	// Get SSH config profiles
	sshProfiles, err := ParseSSHConfig()
	if err != nil {
		// Don't fail if SSH config doesn't exist
		sshProfiles = make(map[string]Profile)
	}

	// Print custom profiles
	if len(customProfiles.Profiles) > 0 {
		fmt.Println("📋 Custom Profiles (~/.sshclient/config.yaml):")
		fmt.Println(strings.Repeat("─", 70))

		names := make([]string, 0, len(customProfiles.Profiles))
		for name := range customProfiles.Profiles {
			names = append(names, name)
		}
		sort.Strings(names)

		for _, name := range names {
			profile := customProfiles.Profiles[name]
			auth := "password"
			if profile.Key != "" {
				auth = fmt.Sprintf("key: %s", profile.Key)
			}
			fmt.Printf("  @%-15s %s@%s:%s (%s)\n",
				name, profile.User, profile.Host, profile.Port, auth)
		}
		fmt.Println()
	}

	// Print SSH config profiles
	if len(sshProfiles) > 0 {
		fmt.Println("🔧 SSH Config Profiles (~/.ssh/config):")
		fmt.Println(strings.Repeat("─", 70))

		names := make([]string, 0, len(sshProfiles))
		for name := range sshProfiles {
			names = append(names, name)
		}
		sort.Strings(names)

		for _, name := range names {
			profile := sshProfiles[name]
			auth := "password"
			if profile.Key != "" {
				auth = fmt.Sprintf("key: %s", profile.Key)
			}
			port := profile.Port
			if port == "" {
				port = "22"
			}
			user := profile.User
			if user == "" {
				user = "(default)"
			}
			fmt.Printf("  @%-15s %s@%s:%s (%s)\n",
				name, user, profile.Host, port, auth)
		}
		fmt.Println()
	}

	if len(customProfiles.Profiles) == 0 && len(sshProfiles) == 0 {
		fmt.Println("No profiles found.")
		fmt.Println("\nCreate a profile with: sshclient profile add <name>")
		fmt.Println("Or add hosts to ~/.ssh/config")
	}

	return nil
}

// ProfileRemove removes a profile
func ProfileRemove(name string) error {
	if err := RemoveProfile(name); err != nil {
		return err
	}

	fmt.Printf("✓ Profile '%s' removed successfully!\n", name)
	return nil
}

// ProfileShow displays details of a specific profile
func ProfileShow(name string) error {
	profile, err := FindProfile(name)
	if err != nil {
		return err
	}

	fmt.Printf("Profile: %s\n", name)
	fmt.Println(strings.Repeat("─", 40))
	fmt.Printf("  Host:     %s\n", profile.Host)
	fmt.Printf("  User:     %s\n", profile.User)
	fmt.Printf("  Port:     %s\n", profile.Port)

	if profile.Key != "" {
		fmt.Printf("  Key:      %s\n", profile.Key)
	}
	if profile.Password != "" {
		fmt.Printf("  Password: (stored)\n")
	}

	return nil
}

// PrintProfileHelp prints profile command usage help
func PrintProfileHelp() {
	fmt.Println("Profile Management - Manage SSH connection profiles")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  sshclient profile <command> [arguments]")
	fmt.Println()
	fmt.Println("Available Commands:")
	fmt.Println("  add <name>       Add a new profile (interactive)")
	fmt.Println("  list, ls         List all profiles (custom + SSH config)")
	fmt.Println("  show <name>      Show profile details")
	fmt.Println("  remove <name>    Remove a profile (custom only)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  sshclient profile add myserver")
	fmt.Println("  sshclient profile list")
	fmt.Println("  sshclient profile show myserver")
	fmt.Println("  sshclient profile remove myserver")
	fmt.Println()
	fmt.Println("Using Profiles:")
	fmt.Println("  sshclient @myserver              # Interactive shell")
	fmt.Println("  sshclient @myserver ls -la       # Run command")
	fmt.Println()
	fmt.Println("Profile Storage:")
	fmt.Println("  Custom profiles: ~/.sshclient/config.yaml")
	fmt.Println("  SSH config:      ~/.ssh/config (read-only)")
}

// HandleProfileCommand handles profile management commands
func HandleProfileCommand(args []string) error {
	if len(args) < 1 {
		PrintProfileHelp()
		return nil
	}

	subcommand := args[0]

	switch subcommand {
	case "add":
		if len(args) < 2 {
			fmt.Println("Error: profile add requires a name")
			fmt.Println()
			fmt.Println("Usage:")
			fmt.Println("  sshclient profile add <name>")
			fmt.Println()
			fmt.Println("Example:")
			fmt.Println("  sshclient profile add myserver")
			fmt.Println()
			fmt.Println("This will start an interactive session to configure:")
			fmt.Println("  - Host (hostname or IP address)")
			fmt.Println("  - User (SSH username)")
			fmt.Println("  - Port (default: 22)")
			fmt.Println("  - Authentication (SSH key or password)")
			return fmt.Errorf("missing profile name")
		}
		return ProfileAdd(args[1])

	case "list", "ls":
		return ProfileList()

	case "remove", "rm":
		if len(args) < 2 {
			fmt.Println("Error: profile remove requires a name")
			fmt.Println()
			fmt.Println("Usage:")
			fmt.Println("  sshclient profile remove <name>")
			fmt.Println()
			fmt.Println("Example:")
			fmt.Println("  sshclient profile remove myserver")
			fmt.Println()
			fmt.Println("Note: This only removes custom profiles from ~/.sshclient/config.yaml")
			fmt.Println("      SSH config profiles (~/.ssh/config) are read-only")
			return fmt.Errorf("missing profile name")
		}
		return ProfileRemove(args[1])

	case "show":
		if len(args) < 2 {
			fmt.Println("Error: profile show requires a name")
			fmt.Println()
			fmt.Println("Usage:")
			fmt.Println("  sshclient profile show <name>")
			fmt.Println()
			fmt.Println("Example:")
			fmt.Println("  sshclient profile show myserver")
			fmt.Println()
			fmt.Println("Tip: Use 'sshclient profile list' to see all available profiles")
			return fmt.Errorf("missing profile name")
		}
		return ProfileShow(args[1])

	case "help", "-h", "--help":
		PrintProfileHelp()
		return nil

	default:
		fmt.Printf("Error: unknown profile subcommand: %s\n", subcommand)
		fmt.Println()
		PrintProfileHelp()
		return fmt.Errorf("unknown subcommand: %s", subcommand)
	}
}
