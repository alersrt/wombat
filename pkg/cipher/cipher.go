package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"wombat/pkg"
)

var ErrCipher = errors.New("cipher")

type AesGcmCipher struct {
	nonceSize int
	gcm       cipher.AEAD
}

// NewAesGcmCipher Creates new GCM with standard nonce size (12)
func NewAesGcmCipher(key []byte) (*AesGcmCipher, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, pkg.Wrap(ErrCipher, err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, pkg.Wrap(ErrCipher, err)
	}

	return &AesGcmCipher{gcm: gcm}, nil
}

func (a *AesGcmCipher) Encrypt(plaintext string) ([]byte, error) {
	nonce := make([]byte, a.gcm.NonceSize())
	_, err := io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, pkg.Wrap(ErrCipher, err)
	}
	ciphertext := a.gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return ciphertext, nil
}

func (a *AesGcmCipher) Decrypt(ciphertext []byte) (string, error) {
	nonceSize := a.gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", pkg.Wrap(ErrCipher, errors.New("size mismatch"))
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintextBytes, err := a.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", pkg.Wrap(ErrCipher, err)
	}
	plaintext := string(plaintextBytes)
	return plaintext, nil
}
