package utils

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	passphrase = "yourpassphrasemustbe32bytes!1234"
)

func TestEncryptionDecryption(t *testing.T) {
	t.Run("successful encryption and decryption", func(t *testing.T) {
		enctryptor := Encryptor{
			Passphrase: passphrase,
		}
		originalText := "Hello, World!"
		encryptedText, err := enctryptor.Encrypt(originalText)
		assert.Nil(t, err, "encryption should not error")

		// Ensure that we get a base64-encoded string back
		_, err = base64.StdEncoding.DecodeString(encryptedText)
		assert.Nil(t, err, "encrypted text should be base64 encoded")

		decryptedText, err := enctryptor.Decrypt(encryptedText)
		assert.Nil(t, err, "decryption should not error")
		assert.Equal(t, originalText, decryptedText, "decrypted text should match original")
	})

	t.Run("error on decryption with wrong passphrase", func(t *testing.T) {
		enctryptor := Encryptor{
			Passphrase: passphrase,
		}
		wrongPassphrase := "thisisnottherightpassphrase32b"
		originalText := "Hello, World!"

		encryptedText, err := enctryptor.Encrypt(originalText)
		assert.Nil(t, err, "encryption should not error")

		enctryptor.Passphrase = wrongPassphrase
		_, err = enctryptor.Decrypt(encryptedText)
		assert.NotNil(t, err, "should error on wrong passphrase")
	})

	t.Run("empty string encryption and decryption", func(t *testing.T) {
		enctryptor := Encryptor{
			Passphrase: passphrase,
		}
		originalText := ""
		encryptedText, err := enctryptor.Encrypt(originalText)
		assert.Nil(t, err, "encryption should not error for empty string")

		decryptedText, err := enctryptor.Decrypt(encryptedText)
		assert.Nil(t, err, "decryption should not error for empty string")
		assert.Equal(t, originalText, decryptedText, "decrypted text should match original for empty string")
	})
}
