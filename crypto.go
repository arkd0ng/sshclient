package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	// AES-256 requires 32 bytes key
	keySize   = 32
	saltSize  = 32
	nonceSize = 12
	// PBKDF2 iterations
	iterations = 100000
	// Hard-coded passphrase for automatic encryption/decryption
	// This is not the most secure, but provides basic protection from
	// casual viewing of the config file
	internalPassphrase = "sshclient-internal-key-v1-do-not-share-2024"
)

// DeriveKey derives a cryptographic key from a password using PBKDF2
func DeriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, iterations, keySize, sha256.New)
}

// Encrypt encrypts plaintext using AES-256-GCM with a password
// Returns base64 encoded: salt + nonce + ciphertext
func Encrypt(plaintext, password string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	// Generate random salt
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive key from password
	key := DeriveKey(password, salt)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	// Combine: salt + nonce + ciphertext
	result := make([]byte, saltSize+len(nonce)+len(ciphertext))
	copy(result[0:], salt)
	copy(result[saltSize:], nonce)
	copy(result[saltSize+len(nonce):], ciphertext)

	// Encode to base64
	return base64.StdEncoding.EncodeToString(result), nil
}

// Decrypt decrypts base64 encoded ciphertext using AES-256-GCM with a password
func Decrypt(encoded, password string) (string, error) {
	if encoded == "" {
		return "", nil
	}

	// Decode from base64
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// Check minimum size
	if len(data) < saltSize+nonceSize {
		return "", fmt.Errorf("invalid encrypted data: too short")
	}

	// Extract salt, nonce, and ciphertext
	salt := data[0:saltSize]
	nonce := data[saltSize : saltSize+nonceSize]
	ciphertext := data[saltSize+nonceSize:]

	// Derive key from password
	key := DeriveKey(password, salt)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: incorrect password or corrupted data")
	}

	return string(plaintext), nil
}

// EncryptAuto encrypts plaintext using the internal passphrase
// This provides basic protection from casual viewing of the config file
func EncryptAuto(plaintext string) (string, error) {
	return Encrypt(plaintext, internalPassphrase)
}

// DecryptAuto decrypts ciphertext using the internal passphrase
func DecryptAuto(encoded string) (string, error) {
	return Decrypt(encoded, internalPassphrase)
}
