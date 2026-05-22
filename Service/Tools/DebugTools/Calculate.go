package DebugTools

import (
	"changeme/Service/Tools/DebugTools/md2"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rc4"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"golang.org/x/crypto/md4"
)

func Rc4Encrypt(data, key []byte) []byte {
	cipher, _ := rc4.NewCipher(key)
	encrypted := make([]byte, len(data))
	cipher.XORKeyStream(encrypted, data)
	return encrypted
}

func Rc4Decrypt(encrypted, key []byte) []byte {
	cipher, _ := rc4.NewCipher(key)
	decrypted := make([]byte, len(encrypted))
	cipher.XORKeyStream(decrypted, encrypted)
	return decrypted
}

// 计算MD2哈希值
func CalculateMD2(data []byte) []byte {
	hash := md2.New()
	hash.Write(data)
	return hash.Sum(nil)
}

// 计算MD4哈希值
func CalculateMD4(data []byte) []byte {
	hash := md4.New()
	hash.Write(data)
	return hash.Sum(nil)
}

// 计算MD5哈希值
func CalculateMD5(data []byte) []byte {
	hash := md5.New()
	hash.Write(data)
	return hash.Sum(nil)
}

// 计算SHA-1哈希值
func CalculateSHA1(data []byte) []byte {
	hash := sha1.New()
	hash.Write(data)
	return hash.Sum(nil)
}

// 计算SHA-256哈希值
func CalculateSHA256(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

// 计算SHA-512哈希值
func CalculateSHA512(data []byte) []byte {
	hash := sha512.New()
	hash.Write(data)
	return hash.Sum(nil)
}

// 计算HMAC-MD2哈希值
func CalculateHMACMD2(data, key []byte) []byte {
	hash := hmac.New(md2.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}

// 计算HMAC-MD4哈希值
func CalculateHMACMD4(data, key []byte) []byte {
	hash := hmac.New(md4.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}

// 计算HMAC-MD5哈希值
func CalculateHMACMD5(data, key []byte) []byte {
	hash := hmac.New(md5.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}

// 计算HMAC-SHA-1哈希值
func CalculateHMACSHA1(data, key []byte) []byte {
	hash := hmac.New(sha1.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}

// 计算HMAC-SHA-256哈希值
func CalculateHMACSHA256(data, key []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}

// 计算HMAC-SHA-512哈希值
func CalculateHMACSHA512(data, key []byte) []byte {
	hash := hmac.New(sha512.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}