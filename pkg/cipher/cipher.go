package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

var (
	ErrCipher        = errors.New("cipher")
	ErrCipherNew     = errors.New("cipher: new")
	ErrCipherEncrypt = errors.New("cipher: encrypt")
	ErrCipherDecrypt = errors.New("cipher: decrypt")
)

type AesGcmCipher struct {
	nonceSize int
	gcm       cipher.AEAD
}

// NewAesGcmCipher Creates new GCM with standard nonce size (12)
func NewAesGcmCipher(key []byte) (*AesGcmCipher, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCipherNew, err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCipherNew, err)
	}

	return &AesGcmCipher{gcm: gcm}, nil
}

func (a *AesGcmCipher) Encrypt(plaintext string) ([]byte, error) {
	nonce := make([]byte, a.gcm.NonceSize())
	_, err := io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCipherEncrypt, err)
	}
	ciphertext := a.gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return ciphertext, nil
}

func (a *AesGcmCipher) Decrypt(ciphertext []byte) (string, error) {
	nonceSize := a.gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("%w: size mismatch", ErrCipherDecrypt)
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintextBytes, err := a.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrCipherDecrypt, err)
	}
	plaintext := string(plaintextBytes)
	return plaintext, nil
}
