package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

// aesCrypto -
type aesCrypto struct {
	// iv - IV值
	iv []byte
	// key - aes加密key
	key []byte
}

// NewAESCrypto - 建立加密物件
func NewAESCrypto(key, iv []byte) *aesCrypto {
	return &aesCrypto{
		iv:  iv,
		key: key,
	}
}

// Encrypt - 加密
func (a *aesCrypto) AesEncrypt(src string) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}
	ecb := cipher.NewCBCEncrypter(block, a.iv)
	content := []byte(src)
	content = PKCS5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)

	return crypted, nil
}
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
