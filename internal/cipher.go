package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"wombat/pkg"
)

func AesGcmEncrypt(key []byte, plaintext string) (ciphertext, nonce []byte) {
	block, err := aes.NewCipher(key)
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

func AesGcmDecrypt(key, ciphertext, nonce []byte) (plaintext string) {
	block, err := aes.NewCipher(key)
	pkg.Throw(err)
	aesGcm, err := cipher.NewGCM(block)
	pkg.Throw(err)
	plaintextBytes, err := aesGcm.Open(nil, nonce, ciphertext, nil)
	pkg.Throw(err)
	plaintext = string(plaintextBytes)
	return
}
