package WithSalt

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"math/rand"
)

const TripleDesKeySize = 24
const TripleDesBlockSize = 8

func extract3DesKey(password, salt []byte) (key, iv []byte) {
	var derivedKey [TripleDesKeySize + TripleDesBlockSize]byte
	m := derivedKey[:]
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

	return m[:TripleDesKeySize], m[TripleDesKeySize : TripleDesKeySize+TripleDesBlockSize]
}

func TripleDesCbcWithSaltEncrypt(encrypt, SecretPass []byte) []byte {
	var salt [8]byte
	for i := 0; i < 8; i++ {
		salt[i] = byte(rand.Intn(250) + 10)
	}
	key, iv := extract3DesKey(SecretPass, salt[:])
	block, err := des.NewTripleDESCipher(key)
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
func _pkcs7Padding(data []byte, blockSize int) []byte {
	paddle := blockSize - len(data)%blockSize
	return append(data, PrePadPatterns[paddle]...)
}
func TripleDesCbcWithSaltDecrypt(encrypt, SecretPass []byte) []byte {
	if len(encrypt) < TripleDesBlockSize {
		return nil
	}
	saltHeader := encrypt[:TripleDesBlockSize]
	if !bytes.Equal(saltHeader[:8], SaltHeader) {
		return nil
	}
	key, iv := extract3DesKey(SecretPass, encrypt[8:16])

	if len(encrypt) == 0 || len(encrypt)%TripleDesBlockSize != 0 {
		return nil
	}

	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil
	}

	cbc := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(encrypt)-16)
	cbc.CryptBlocks(decrypted, encrypt[16:])

	padded := _pkcs7Unpading(decrypted)
	return padded
}
