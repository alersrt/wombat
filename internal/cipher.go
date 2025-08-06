package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"wombat/pkg"
)

type AesGcmCipher struct {
	key []byte
}

func NewAesGcmCipher(key []byte) *AesGcmCipher {
	return &AesGcmCipher{key}
}

func (a *AesGcmCipher) AesGcmEncrypt(plaintext string) (ciphertext []byte, nonce []byte) {
	block, err := aes.NewCipher(a.key)
	pkg.Throw(err)
	nonce = make([]byte, 12)
	_, err = io.ReadFull(rand.Reader, nonce)
	pkg.Throw(err)
	aesGcm, err := cipher.NewGCM(block)
	pkg.Throw(err)
	plaintextBytes := []byte(plaintext)
	ciphertext = aesGcm.Seal(nil, nonce, plaintextBytes, nil)
	return
}

func (a *AesGcmCipher) AesGcmDecrypt(ciphertext, nonce []byte) (plaintext string) {
	block, err := aes.NewCipher(a.key)
	pkg.Throw(err)
	aesGcm, err := cipher.NewGCM(block)
	pkg.Throw(err)
	plaintextBytes, err := aesGcm.Open(nil, nonce, ciphertext, nil)
	pkg.Throw(err)
	plaintext = string(plaintextBytes)
	return
}
