// Package cryptox provides cryptographic utilities for encrypting and decrypting
// text and binary data using AES-GCM (Advanced Encryption Standard with Galois/Counter Mode).
// It also includes functionality for generating secure secret phrases and deriving keys.
package cryptox

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// CryptoModule is a struct that provides methods for encryption, decryption,
// and secure secret phrase generation.
type CryptoModule struct{}

// NewCryproModule initializes and returns a new instance of CryptoModule.
// Note: The function name contains a typo ("Crypro" instead of "Crypto").
func NewCryproModule() *CryptoModule {
	return &CryptoModule{}
}

// Encrypt encrypts a plaintext string using AES-GCM and returns the result as a base64-encoded string.
// The key is used to derive the encryption key.
// Returns an error if the encryption process fails.
func (c *CryptoModule) Encrypt(plaintext, key string) (string, error) {
	keyBytes := c.deriveKey(key)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nil, nonce, []byte(plaintext), nil)
	result := append(nonce, ciphertext...)

	return base64.StdEncoding.EncodeToString(result), nil
}

// Decrypt decrypts a base64-encoded encrypted string using AES-GCM and returns the plaintext.
// The key is used to derive the decryption key.
// Returns an error if the decryption process fails or if the input data is invalid.
func (c *CryptoModule) Decrypt(encryptedText, key string) (string, error) {
	keyBytes := c.deriveKey(key)

	data, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("invalid data")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// GenerateSecretPhrase generates a secure secret phrase by hashing the input text
// using SHA-256 and encoding the result in base64. The returned phrase is truncated
// to 14 characters for convenience.
func (c *CryptoModule) GenerateSecretPhrase(txt string) string {
	hash := sha256.Sum256([]byte(txt))

	encoded := base64.StdEncoding.EncodeToString(hash[:])

	return encoded[:14]
}

// deriveKey derives a 32-byte key from the provided password by truncating or padding
// the password to 32 bytes. This is a simple key derivation method and may not be
// suitable for all use cases.
func (c *CryptoModule) deriveKey(password string) []byte {
	hash := make([]byte, 32)
	copy(hash, []byte(password))
	return hash
}

// EncryptBinaryData encrypts binary data using AES-GCM and returns the encrypted data
// with the nonce prepended. The key is used to derive the encryption key.
// Returns an error if the encryption process fails.
func (c *CryptoModule) EncryptBinaryData(plaintext []byte, key string) ([]byte, error) {
	keyBytes := c.deriveKey(key)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)
	return append(nonce, ciphertext...), nil
}

// DecryptBinaryData decrypts binary data encrypted with AES-GCM. The nonce is expected
// to be prepended to the encrypted data. The key is used to derive the decryption key.
// Returns an error if the decryption process fails or if the input data is invalid.
func (c *CryptoModule) DecryptBinaryData(encryptedData []byte, key string) ([]byte, error) {
	keyBytes := c.deriveKey(key)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, errors.New("invalid data")
	}

	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("cipher: message authentication failed: %w", err)
	}

	return plaintext, nil
}
