package cryptox

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCryptoModule_EncryptDecrypt(t *testing.T) {
	cryptoModule := NewCryproModule()

	plaintext := "Hello, World!"
	key := "supersecretkey"

	encrypted, err := cryptoModule.Encrypt(plaintext, key)
	require.NoError(t, err, "Encryption should not return an error")

	decrypted, err := cryptoModule.Decrypt(encrypted, key)
	require.NoError(t, err, "Decryption should not return an error")
	assert.Equal(t, plaintext, decrypted, "Decrypted text should be the same as original plaintext")
}

func TestCryptoModule_EncryptDecrypt_Binary(t *testing.T) {
	cryptoModule := NewCryproModule()

	plaintext := []byte("Binary data test")
	key := "supersecretkey"

	encrypted, err := cryptoModule.EncryptBinaryData(plaintext, key)
	require.NoError(t, err, "Encryption should not return an error")

	decrypted, err := cryptoModule.DecryptBinaryData(encrypted, key)
	require.NoError(t, err, "Decryption should not return an error")
	assert.Equal(t, plaintext, decrypted, "Decrypted data should be the same as original plaintext")
}

func TestCryptoModule_GenerateSecretPhrase(t *testing.T) {
	cryptoModule := NewCryproModule()

	txt := "supersecretpassword"

	secret := cryptoModule.GenerateSecretPhrase(txt)

	assert.Len(t, secret, 14, "Secret phrase should be 14 characters long")

	expectedSecret := "WsFStvi9uLsSlZ"
	assert.Equal(t, expectedSecret, secret, "Generated secret phrase should match the expected value")
}

func TestCryptoModule_EncryptDecrypt_InvalidData(t *testing.T) {
	cryptoModule := NewCryproModule()

	plaintext := "Hello, World!"
	key := "supersecretkey"

	encrypted, err := cryptoModule.Encrypt(plaintext, key)
	require.NoError(t, err, "Encryption should not return an error")

	_, err = cryptoModule.Decrypt(encrypted, "wrongkey")
	assert.Error(t, err, "Decryption with incorrect key should return an error")
}

func TestCryptoModule_Decrypt_InvalidData(t *testing.T) {
	cryptoModule := NewCryproModule()

	invalidData := []byte("invalid-encrypted-data")
	key := "supersecretkey"

	_, err := cryptoModule.DecryptBinaryData(invalidData, key)
	assert.Error(t, err, "Decryption with invalid data should return an error")
}
