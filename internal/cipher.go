package internal

import (
	"crypto/aes"
	"crypto/cipher"
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
	aesGcm, err := cipher.NewGCM(block)
	pkg.Throw(err)
	plaintextBytes := []byte(plaintext)
	nonce := zeros()
	ciphertext = aesGcm.Seal(nil, nonce, plaintextBytes, nil)
	return
}

func (a *AesGcmCipher) AesGcmDecrypt(ciphertext, nonce []byte) (plaintext string) {
	block, err := aes.NewCipher(a.key)
	pkg.Throw(err)
	aesGcm, err := cipher.NewGCM(block)
	pkg.Throw(err)
	nonce = zeros()
	plaintextBytes, err := aesGcm.Open(nil, nonce, ciphertext, nil)
	pkg.Throw(err)
	plaintext = string(plaintextBytes)
	return
}

func zeros() (nonce []byte) {
	nonce = make([]byte, 12)
	for i := range nonce {
		nonce[i] = 0
	}
	return
}
