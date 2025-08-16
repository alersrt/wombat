package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/pkg/errors"
	"io"
)

type AesGcmCipher struct {
	nonceSize int
	gcm       cipher.AEAD
}

// NewAesGcmCipher Creates new GCM with standard nonce size (12)
func NewAesGcmCipher(key []byte) (*AesGcmCipher, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return &AesGcmCipher{gcm: gcm}, nil
}

func (a *AesGcmCipher) Encrypt(plaintext string) ([]byte, error) {
	nonce := make([]byte, a.gcm.NonceSize())
	_, err := io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	ciphertext := a.gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return ciphertext, nil
}

func (a *AesGcmCipher) Decrypt(ciphertext []byte) (string, error) {
	nonceSize := a.gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext and nonce size mismatch")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintextBytes, err := a.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.New(err.Error())
	}
	plaintext := string(plaintextBytes)
	return plaintext, nil
}
