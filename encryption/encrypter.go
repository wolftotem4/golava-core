package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

const gcmTagSize = 12

type Encrypter struct {
	Key []byte
}

func NewEncrypter(key []byte) *Encrypter {
	return &Encrypter{Key: key}
}

func (e *Encrypter) Encrypt(value []byte) ([]byte, error) {
	nonce, err := generateNonce(gcmTagSize)
	if err != nil {
		return nil, err
	}

	ciphertext, err := e.gcmEncrypt(value, nonce)
	if err != nil {
		return nil, err
	}

	// 合併 nonce 和 ciphertext
	return append(nonce, ciphertext...), nil
}

func (e *Encrypter) Decrypt(value []byte) ([]byte, error) {
	if len(value) < gcmTagSize {
		return nil, io.ErrUnexpectedEOF
	}

	nonce := value[:gcmTagSize]
	ciphertext := value[gcmTagSize:]

	return e.gcmDecrypt(ciphertext, nonce)
}

func (e *Encrypter) gcmEncrypt(value []byte, nonce []byte) ([]byte, error) {
	c, err := aes.NewCipher(e.Key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nil, nonce, value, nil), nil
}

func (e *Encrypter) cbcEncrypt(value []byte, iv []byte) ([]byte, error) {
	c, err := aes.NewCipher(e.Key)
	if err != nil {
		return nil, err
	}

	data := PKCS7Padding(value, c.BlockSize())
	ciphertext := make([]byte, len(data))
	copy(ciphertext[:aes.BlockSize], iv)

	cbc := cipher.NewCBCEncrypter(c, iv)
	cbc.CryptBlocks(ciphertext, data)

	return ciphertext, nil
}

func (e *Encrypter) gcmDecrypt(ciphertext []byte, nonce []byte) ([]byte, error) {
	c, err := aes.NewCipher(e.Key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	return gcm.Open(nil, nonce, ciphertext, nil)
}

func generateNonce(size int) ([]byte, error) {
	nonce := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, nonce)
	return nonce, err
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func GenerateKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	return key, err
}
