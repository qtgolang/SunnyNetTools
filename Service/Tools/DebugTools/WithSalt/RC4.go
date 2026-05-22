package WithSalt

import (
	"bytes"
	"crypto/md5"
	"crypto/rc4"
	"math/rand"
)

func extractRC4Key(password, salt []byte) []byte {
	var sC [32]byte
	m := sC[:]
	buf := make([]byte, 0, 16+len(password)+len(salt))
	var prevSum [16]byte
	for i := 0; i < 2; i++ {
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
	return sC[:32]
}

func RC4WithSaltDecrypt(encrypt, SecretPass []byte) []byte {
	if len(encrypt) < 16 {
		return nil
	}
	saltHeader := encrypt[:16]
	if !bytes.Equal(saltHeader[:8], SaltHeader) {
		return nil
	}

	salt := saltHeader[8:16]
	key := extractRC4Key(SecretPass, salt)
	cipher, err := rc4.NewCipher(key)
	if err != nil {
		return nil
	}
	decrypted := make([]byte, len(encrypt)-16)
	cipher.XORKeyStream(decrypted, encrypt[16:])
	return decrypted
}
func RC4WithSaltEncrypt(encrypt, SecretPass []byte) []byte {
	var salt [8]byte
	for i := 0; i < 8; i++ {
		salt[i] = byte(rand.Intn(250) + 10)
	}
	key := extractRC4Key(SecretPass, salt[:])
	cipher, err := rc4.NewCipher(key)
	if err != nil {
		return nil
	}
	encrypted := make([]byte, len(encrypt))
	cipher.XORKeyStream(encrypted, encrypt)
	result := append(append(SaltHeader, salt[:]...), encrypted...)
	return result
}
