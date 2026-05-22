package WithSalt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"math/rand"
)

var PrePadPatterns [aes.BlockSize + 1][]byte

var SaltHeader = []byte("Salted__")

func init() {
	for i := 0; i < len(PrePadPatterns); i++ {
		PrePadPatterns[i] = bytes.Repeat([]byte{byte(i)}, i)
	}
}

func extract256Aes(password, salt []byte) (key, iv []byte) {
	var sC [48]byte
	m := sC[:]
	buf := make([]byte, 0, 16+len(password)+len(salt))
	var prevSum [16]byte
	for i := 0; i < 3; i++ {
		n := 0
		if i > 0 {
			n = 16
		}
		buf = buf[:n+len(password)+len(salt)]
		copy(buf, prevSum[:])
		copy(buf[n:], password)
		copy(buf[n+len(password):], salt)
		prevSum = md5.Sum(buf)
		copy(m[i*16:], prevSum[:])
	}
	return sC[:32], sC[32:]
}
func extract192Aes(password, salt []byte) (key, iv []byte) {
	var sC [40]byte
	m := sC[:]
	buf := make([]byte, 0, 16+len(password)+len(salt))
	var prevSum [16]byte
	for i := 0; i < 3; i++ {
		n := 0
		if i > 0 {
			n = 16
		}
		buf = buf[:n+len(password)+len(salt)]
		copy(buf, prevSum[:])
		copy(buf[n:], password)
		copy(buf[n+len(password):], salt)
		prevSum = md5.Sum(buf)
		copy(m[i*16:], prevSum[:])
	}
	return sC[:24], sC[24:]
}
func extract128Aes(password, salt []byte) (key, iv []byte) {
	var sC [32]byte
	m := sC[:]
	buf := make([]byte, 0, 16+len(password)+len(salt))
	var prevSum [16]byte
	for i := 0; i < 3; i++ {
		n := 0
		if i > 0 {
			n = 16
		}
		buf = buf[:n+len(password)+len(salt)]
		copy(buf, prevSum[:])
		copy(buf[n:], password)
		copy(buf[n+len(password):], salt)
		prevSum = md5.Sum(buf)
		copy(m[i*16:], prevSum[:])
	}
	return sC[:16], sC[16:]
}
func AesCbc256WithSaltDecrypt(encrypt, SecretPass []byte) []byte {
	if len(encrypt) < aes.BlockSize {
		return nil
	}
	saltHeader := encrypt[:aes.BlockSize]
	if !bytes.Equal(saltHeader[:8], SaltHeader) {
		return nil
	}
	key, iv := extract256Aes(SecretPass, saltHeader[8:])

	if len(encrypt) == 0 || len(encrypt)%aes.BlockSize != 0 {
		return nil
	}
	cc, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	cbc := cipher.NewCBCDecrypter(cc, iv)
	decrypted := make([]byte, len(encrypt)-16)
	cbc.CryptBlocks(decrypted, encrypt[16:])
	return _pkcs7Unpading(decrypted)
}
func AesCbc256WithSaltEncrypt(encrypt, SecretPass []byte) []byte {
	var salt [8]byte
	for i := 0; i < 8; i++ {
		salt[i] = byte(rand.Intn(250) + 10)
	}
	key, iv := extract256Aes(SecretPass, salt[:])
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	cbc := cipher.NewCBCEncrypter(block, iv)
	padded := _pkcs7Padding(encrypt, block.BlockSize())
	encrypted := make([]byte, len(padded))
	cbc.CryptBlocks(encrypted, padded)
	result := append(append(SaltHeader, salt[:]...), encrypted...)
	return result
}

func AesCbc192WithSaltDecrypt(encrypt, SecretPass []byte) []byte {
	if len(encrypt) < aes.BlockSize {
		return nil
	}
	saltHeader := encrypt[:aes.BlockSize]
	if !bytes.Equal(saltHeader[:8], SaltHeader) {
		return nil
	}
	key, iv := extract192Aes(SecretPass, saltHeader[8:])

	if len(encrypt) == 0 || len(encrypt)%aes.BlockSize != 0 {
		return nil
	}
	cc, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	cbc := cipher.NewCBCDecrypter(cc, iv)
	decrypted := make([]byte, len(encrypt)-16)
	cbc.CryptBlocks(decrypted, encrypt[16:])
	return _pkcs7Unpading(decrypted)
}
func AesCbc128WithSaltDecrypt(encrypt, SecretPass []byte) []byte {
	if len(encrypt) < aes.BlockSize {
		return nil
	}
	saltHeader := encrypt[:aes.BlockSize]
	if !bytes.Equal(saltHeader[:8], SaltHeader) {
		return nil
	}
	key, iv := extract128Aes(SecretPass, saltHeader[8:])

	if len(encrypt) == 0 || len(encrypt)%aes.BlockSize != 0 {
		return nil
	}
	cc, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	cbc := cipher.NewCBCDecrypter(cc, iv)
	decrypted := make([]byte, len(encrypt)-16)
	cbc.CryptBlocks(decrypted, encrypt[16:])
	return _pkcs7Unpading(decrypted)
}
func AesCbc192WithSaltEncrypt(encrypt, SecretPass []byte) []byte {
	var salt [8]byte
	for i := 0; i < 8; i++ {
		salt[i] = byte(rand.Intn(250) + 10)
	}
	key, iv := extract192Aes(SecretPass, salt[:])
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	cbc := cipher.NewCBCEncrypter(block, iv)
	padded := _pkcs7Padding(encrypt, block.BlockSize())
	encrypted := make([]byte, len(padded))
	cbc.CryptBlocks(encrypted, padded)
	result := append(append(SaltHeader, salt[:]...), encrypted...)
	return result
}
func AesCbc128WithSaltEncrypt(encrypt, SecretPass []byte) []byte {
	var salt [8]byte
	for i := 0; i < 8; i++ {
		salt[i] = byte(rand.Intn(250) + 10)
	}
	key, iv := extract128Aes(SecretPass, salt[:])
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	cbc := cipher.NewCBCEncrypter(block, iv)
	padded := _pkcs7Padding(encrypt, block.BlockSize())
	encrypted := make([]byte, len(padded))
	cbc.CryptBlocks(encrypted, padded)
	result := append(append(SaltHeader, salt[:]...), encrypted...)
	return result
}
func _pkcs7Unpading(data []byte) []byte {
	if len(data)%aes.BlockSize != 0 || len(data) == 0 {
		return nil
	}
	padlen := int(data[len(data)-1])
	if padlen > aes.BlockSize || padlen == 0 {
		return nil
	}
	if !bytes.Equal(PrePadPatterns[padlen], data[len(data)-padlen:]) {
		return nil
	}
	return data[:len(data)-padlen]
}
