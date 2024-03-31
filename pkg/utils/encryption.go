package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
)

type Encryption interface {
	Encrypt(plainText string) (string, error)
	Decrypt(encryptedText string) (string, error)
	HashSHA256(plainText string, context string) string
}

type Encryptor struct {
	Passphrase string
}

func NewEncryptor(passphrase string) *Encryptor {
	return &Encryptor{
		Passphrase: passphrase,
	}
}

// pad applies PKCS#7 padding to the plaintext.
func pad(plaintext []byte, blockSize int) []byte {
	padding := blockSize - len(plaintext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plaintext, padtext...)
}

// unpad removes PKCS#7 padding from the plaintext.
func unpad(plaintext []byte) ([]byte, error) {
	length := len(plaintext)
	if length == 0 {
		return nil, errors.New("plaintext too short")
	}
	padding := int(plaintext[length-1])
	if padding > length || padding == 0 {
		return nil, errors.New("invalid padding")
	}
	for i := length - padding; i < length; i++ {
		if plaintext[i] != byte(padding) {
			return nil, errors.New("invalid padding")
		}
	}
	return plaintext[:length-padding], nil
}

// encrypt encrypts plain text string using AES encryption algorithm.
func (e *Encryptor) Encrypt(plainText string) (string, error) {
	key := []byte(e.Passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	plainTextBytes := pad([]byte(plainText), block.BlockSize())
	cipherText := make([]byte, aes.BlockSize+len(plainTextBytes))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], plainTextBytes)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// decrypt decrypts the encrypted string to original string using the AES encryption algorithm.
func (e *Encryptor) Decrypt(encryptedText string) (string, error) {
	key := []byte(e.Passphrase)
	cipherText, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	if len(cipherText) < aes.BlockSize {
		return "", errors.New("cipherText too short")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)

	plainTextBytes, err := unpad(cipherText)
	if err != nil {
		return "", err
	}
	return string(plainTextBytes), nil
}

func (e *Encryptor) HashSHA256(text string, context string) string {
	hasher := sha256.New()
	hasher.Write([]byte(text))
	hashBytes := hasher.Sum([]byte(context))
	return hex.EncodeToString(hashBytes)
}
