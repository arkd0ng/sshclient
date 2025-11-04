# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Pure Go SSH client for macOS that implements SSH protocol from scratch without using OS ssh commands.

**Current Version**: v1.2.0
**Location**: `~/sshclient`

## Architecture

### Three Input Parsing Modes (main.go)

The CLI supports three distinct parsing modes with priority order:

1. **Profile mode** (`@profile`): Line 38-94 in main.go
   - Checks for `@` prefix first
   - Loads from `~/.sshclient/config.yaml` or `~/.ssh/config`
   - Priority: custom profiles → SSH config

2. **Traditional SSH** (`user@host`): Line 95-137 in main.go
   - Parses `user@host` format using `parseUserHost()`
   - Remaining args become remote command or flags
   - Auto-enables interactive mode if no command

3. **Flag mode**: Line 139 onwards
   - Explicit `-host`, `-user`, `-port`, etc.
   - Requires `-i` flag for interactive mode

All modes converge to create `SSHClient` via password or key authentication.

### Profile System (config.go + profile.go)

**Two-tier profile lookup**:
- `FindProfile()` (config.go:264): Searches custom profiles first, falls back to SSH config
- Custom profiles: YAML at `~/.sshclient/config.yaml` (0600 permissions)
- SSH config: Parsed from `~/.ssh/config` using `ParseSSHConfig()` (config.go:158)

**Profile struct** (config.go:14-22):
```go
type Profile struct {
    Host, User, Port string
    Key               string  // SSH private key path
    Password          string  // Deprecated: plain text
    EncryptedPassword string  // AES-256-GCM encrypted (recommended)
}
```

### Password Encryption (crypto.go)

**Automatic AES-256-GCM encryption** (added in v1.2.0):
- `EncryptAuto()`: Encrypts passwords using internal passphrase
- `DecryptAuto()`: Decrypts passwords automatically
- Uses **PBKDF2** (100,000 iterations) for key derivation
- **No user interaction required** - master password removed in v1.2.0
- Internal passphrase hardcoded in `internalPassphrase` constant

**Implementation details**:
- Salt: 32 bytes (random per encryption)
- Nonce: 12 bytes (GCM standard)
- Key size: 32 bytes (AES-256)
- Encoding: Base64 for YAML storage

### SSH Client Implementation (client.go)

**Direct SSH protocol implementation** using `golang.org/x/crypto/ssh`:
- `NewSSHClient()`: Password auth
- `NewSSHClientWithKey()`: Public key auth (RSA/Ed25519/ECDSA)
- `Connect()`: TCP dial + SSH handshake
- `StartInteractiveShell()`: Sets terminal raw mode, handles PTY
- `RunCommand()`: Single command execution
- `CopyFile()` / `DownloadFile()`: SCP implementation

**Key point**: Uses `ssh.InsecureIgnoreHostKey()` (client.go:31,61) for host key verification - testing only, not production-ready.

## Build & Test Commands

```bash
# Build
go build -o sshclient

# Quick test
./sshclient -version

# Test all three modes
./sshclient @myprofile              # Profile mode
./sshclient user@host uptime        # Traditional SSH
./sshclient -host h -user u -i      # Flag mode

# Profile management
./sshclient profile add NAME        # Interactive creation
./sshclient profile list            # Shows both custom + SSH config
./sshclient profile remove NAME     # Deletes from custom only
```

## Dependencies

```
golang.org/x/crypto/ssh       # SSH protocol
golang.org/x/term            # Terminal control (raw mode, password input)
gopkg.in/yaml.v3             # Config file parsing
```

## File Responsibilities

- **main.go**: CLI parsing, mode detection, client initialization
- **client.go**: SSH protocol implementation (connect, shell, commands, SCP)
- **config.go**: Profile storage/retrieval, SSH config parser
- **profile.go**: Profile management commands (add/list/remove/show)
- **crypto.go**: AES-256-GCM password encryption/decryption (v1.2.0)

## Critical Implementation Details

1. **Argument processing** (main.go:74-86, 109-124):
   - Separates flags from command arguments
   - Reconstructs `os.Args` for `flag.Parse()`
   - Command args joined with spaces

2. **SSH config parsing** (config.go:158-240):
   - Case-insensitive keys
   - Expands `~` in IdentityFile paths
   - Default port 22 if not specified
   - Uses HostName if set, otherwise uses Host name

3. **Terminal handling** (client.go:131-146):
   - Saves/restores terminal state with `terminal.MakeRaw()`
   - Auto-detects terminal size, defaults to 80x24
   - PTY type: `xterm-256color`

4. **File permissions**:
   - Config directory: 0700 (config.go:84)
   - Config file: 0600 (config.go:93)

5. **Default SSH key auto-detection** (client.go:169-190, main.go:230-256):
   - `GetDefaultKeyPath()`: Searches for keys in priority order
     1. `~/.ssh/id_rsa`
     2. `~/.ssh/id_ed25519`
     3. `~/.ssh/id_ecdsa`
   - Auto-applied when:
     - No `-key` flag specified
     - No key in profile
     - No password provided
   - Falls back to password prompt if key auth fails
   - **Profile creation**: Suggests default key path, auto-fills on empty input

6. **Help system** (main.go:139-164, profile.go:204-230):
   - Custom `flag.Usage` function with usage examples
   - Context-aware error messages with helpful hints
   - `PrintProfileHelp()`: Comprehensive profile command guide
   - Each profile subcommand has dedicated help text

## Security Warnings

⚠️ **Not production-ready**:
- No host key verification (uses `InsecureIgnoreHostKey`)
- Passwords stored in plain text in YAML
- No connection timeouts beyond initial dial (10s)

## Common Modifications

**Adding a new profile source**:
1. Create parser function in config.go
2. Update `FindProfile()` to check new source
3. Update `ProfileList()` to display new source

**Adding new authentication method**:
1. Add new `NewSSHClientWith*()` constructor in client.go
2. Add CLI flag in main.go
3. Add profile field in config.go Profile struct

**Supporting new terminal features**:
1. Modify `ssh.TerminalModes` in `StartInteractiveShell()` (client.go:124)
2. Adjust PTY parameters in `RequestPty()` call (client.go:144)
