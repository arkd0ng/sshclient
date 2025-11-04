package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Profile represents an SSH connection profile
type Profile struct {
	Name              string `yaml:"-"`
	Host              string `yaml:"host"`
	User              string `yaml:"user"`
	Port              string `yaml:"port,omitempty"`
	Key               string `yaml:"key,omitempty"`
	Password          string `yaml:"password,omitempty"`           // Deprecated: plain text password
	EncryptedPassword string `yaml:"encrypted_password,omitempty"` // Encrypted password (AES-256-GCM)
}

// ProfileConfig represents the configuration file structure
type ProfileConfig struct {
	Profiles map[string]Profile `yaml:"profiles"`
}

// GetConfigDir returns the sshclient config directory path
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".sshclient"), nil
}

// GetConfigPath returns the path to the profiles config file
func GetConfigPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.yaml"), nil
}

// LoadProfiles loads profiles from the YAML config file
func LoadProfiles() (*ProfileConfig, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// If config doesn't exist, return empty config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &ProfileConfig{Profiles: make(map[string]Profile)}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config ProfileConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if config.Profiles == nil {
		config.Profiles = make(map[string]Profile)
	}

	return &config, nil
}

// SaveProfiles saves profiles to the YAML config file
func SaveProfiles(config *ProfileConfig) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// GetProfile retrieves a profile by name
func GetProfile(name string) (*Profile, error) {
	config, err := LoadProfiles()
	if err != nil {
		return nil, err
	}

	profile, ok := config.Profiles[name]
	if !ok {
		return nil, fmt.Errorf("profile '%s' not found", name)
	}

	profile.Name = name
	return &profile, nil
}

// AddProfile adds a new profile
func AddProfile(profile Profile) error {
	config, err := LoadProfiles()
	if err != nil {
		return err
	}

	config.Profiles[profile.Name] = profile
	return SaveProfiles(config)
}

// RemoveProfile removes a profile by name
func RemoveProfile(name string) error {
	config, err := LoadProfiles()
	if err != nil {
		return err
	}

	if _, ok := config.Profiles[name]; !ok {
		return fmt.Errorf("profile '%s' not found", name)
	}

	delete(config.Profiles, name)
	return SaveProfiles(config)
}

// ListProfiles returns all profile names
func ListProfiles() ([]string, error) {
	config, err := LoadProfiles()
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(config.Profiles))
	for name := range config.Profiles {
		names = append(names, name)
	}

	return names, nil
}

// ParseSSHConfig parses ~/.ssh/config file and returns profiles
func ParseSSHConfig() (map[string]Profile, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(home, ".ssh", "config")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return make(map[string]Profile), nil
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open SSH config: %w", err)
	}
	defer file.Close()

	profiles := make(map[string]Profile)
	var currentHost string
	var currentProfile Profile

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key-value pairs
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		key := strings.ToLower(parts[0])
		value := strings.Join(parts[1:], " ")

		switch key {
		case "host":
			// Save previous host if exists
			if currentHost != "" {
				currentProfile.Name = currentHost
				profiles[currentHost] = currentProfile
			}

			// Start new host
			currentHost = value
			currentProfile = Profile{
				Port: "22", // Default port
			}

		case "hostname":
			currentProfile.Host = value

		case "user":
			currentProfile.User = value

		case "port":
			currentProfile.Port = value

		case "identityfile":
			// Expand ~ to home directory
			if strings.HasPrefix(value, "~") {
				value = filepath.Join(home, value[1:])
			}
			currentProfile.Key = value
		}
	}

	// Save last host
	if currentHost != "" {
		currentProfile.Name = currentHost
		profiles[currentHost] = currentProfile
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading SSH config: %w", err)
	}

	return profiles, nil
}

// GetProfileFromSSHConfig retrieves a profile from SSH config
func GetProfileFromSSHConfig(name string) (*Profile, error) {
	profiles, err := ParseSSHConfig()
	if err != nil {
		return nil, err
	}

	profile, ok := profiles[name]
	if !ok {
		return nil, fmt.Errorf("host '%s' not found in SSH config", name)
	}

	// If hostname is not set, use the host name
	if profile.Host == "" {
		profile.Host = name
	}

	return &profile, nil
}

// FindProfile searches for a profile in both sources
// Priority: 1. Custom profiles, 2. SSH config
func FindProfile(name string) (*Profile, error) {
	// Try custom profiles first
	profile, err := GetProfile(name)
	if err == nil {
		return profile, nil
	}

	// Try SSH config
	profile, err = GetProfileFromSSHConfig(name)
	if err == nil {
		return profile, nil
	}

	return nil, fmt.Errorf("profile '%s' not found in custom profiles or SSH config", name)
}
