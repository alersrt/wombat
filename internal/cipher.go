package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"wombat/pkg"
)

type AesGcmCipher struct {
	nonceSize int
	gcm       cipher.AEAD
}

// NewAesGcmCipher Creates new GCM with standart nonce size (12)
func NewAesGcmCipher(key []byte) *AesGcmCipher {
	block, err := aes.NewCipher(key)
	pkg.Throw(err)
	gcm, err := cipher.NewGCM(block)
	pkg.Throw(err)
	return &AesGcmCipher{
		gcm: gcm,
	}
}

func (a *AesGcmCipher) AesGcmEncrypt(plaintext string) (ciphertext []byte) {
	nonce := make([]byte, a.gcm.NonceSize())
	_, err := io.ReadFull(rand.Reader, nonce)
	pkg.Throw(err)
	ciphertext = a.gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return
}

func (a *AesGcmCipher) AesGcmDecrypt(ciphertext []byte) (plaintext string) {
	nonceSize := a.gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		pkg.Throw(errors.New("ciphertext and nonce size mismatch"))
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintextBytes, err := a.gcm.Open(nil, nonce, ciphertext, nil)
	pkg.Throw(err)
	plaintext = string(plaintextBytes)
	return
}
