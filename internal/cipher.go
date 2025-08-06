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
	key []byte
}

func NewAesGcmCipher(key []byte) *AesGcmCipher {
	return &AesGcmCipher{key}
}

func (a *AesGcmCipher) AesGcmEncrypt(plaintext string) (ciphertext []byte) {
	block, err := aes.NewCipher(a.key)
	pkg.Throw(err)
	gcm, err := cipher.NewGCM(block)
	pkg.Throw(err)
	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	pkg.Throw(err)
	ciphertext = gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return
}

func (a *AesGcmCipher) AesGcmDecrypt(ciphertext []byte) (plaintext string) {
	block, err := aes.NewCipher(a.key)
	pkg.Throw(err)
	gcm, err := cipher.NewGCM(block)
	pkg.Throw(err)
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		pkg.Throw(errors.New("ciphertext and nonce size mismatch"))
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintextBytes, err := gcm.Open(nil, nonce, ciphertext, nil)
	pkg.Throw(err)
	plaintext = string(plaintextBytes)
	return
}
