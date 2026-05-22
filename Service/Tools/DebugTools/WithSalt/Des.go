package WithSalt

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"math/rand"
)

const DesBlockSize = 8

func extractDes(password, salt []byte) (key, iv []byte) {
	var sC [16]byte
	m := sC[:]
	buf := make([]byte, 0, 16+len(password)+len(salt))
	var prevSum [16]byte
	buf = buf[:len(password)+len(salt)]
	copy(buf, password)
	copy(buf[len(password):], salt)
	prevSum = md5.Sum(buf)
	copy(m, prevSum[:])
	return m[:8], m[8:16]
}

func DesCbcWithSaltDecrypt(encrypt, SecretPass []byte) []byte {
	if len(encrypt) < DesBlockSize {
		return nil
	}
	saltHeader := encrypt[:DesBlockSize]
	if !bytes.Equal(saltHeader[:8], SaltHeader) {
		return nil
	}
	key, iv := extractDes(SecretPass, encrypt[8:16])

	if len(encrypt) == 0 || len(encrypt)%DesBlockSize != 0 {
		return nil
	}
	cc, err := des.NewCipher(key)
	if err != nil {
		return nil
	}
	cbc := cipher.NewCBCDecrypter(cc, iv)
	decrypted := make([]byte, len(encrypt)-16)
	cbc.CryptBlocks(decrypted, encrypt[16:])

	padded := _pkcs7Unpading(decrypted)
	return padded
}

func DesCbcWithSaltEncrypt(encrypt, SecretPass []byte) []byte {
	var salt [8]byte
	for i := 0; i < 8; i++ {
		salt[i] = byte(rand.Intn(250) + 10)
	}
	key, iv := extractDes(SecretPass, salt[:])
	block, err := des.NewCipher(key)
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
