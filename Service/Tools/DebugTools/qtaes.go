package DebugTools

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"math/rand"
	"time"
)

type qtAes struct {
	aesIv     []byte
	aesKey    []byte
	pacs      int
	encMethod int
	err       error
	fill      bool
	keySize   int
}

// SetKeySize 设置 KeySize
//
// 请使用 [ Type_Const_AES_keySize_ ]选择
//
// Type_Const_AES_keySize_128
//
// Type_Const_AES_keySize_192
//
// Type_Const_AES_keySize_256
//
// Type_Const_AES_keySize_Auto
func (qtaes *qtAes) SetKeySize(KeySize int) {

	if KeySize != Type_Const_AES_keySize_128 && KeySize != Type_Const_AES_keySize_192 && KeySize != Type_Const_AES_keySize_256 {
		qtaes.keySize = Type_Const_AES_keySize_Auto
	} else {
		qtaes.keySize = KeySize
	}
}

// 获取运行时的错误
func (qtaes *qtAes) GetErr() string {
	if qtaes.err != nil {
		return qtaes.err.Error()
	}
	return ""
}

// 设置填充模式，如果为true,则按填充位填充，若为false,则按0填充
//
// 默认按位填充
//
// 则按填充位填充意思是，比如数据长度10，每次加密16 则填充6个6字节
func (qtaes *qtAes) SetFill(fill bool) {
	qtaes.fill = fill
}

// 设置IV,自动识别是字符串还是字节数组
//
// 若不设置IV则ECB，若设置IV则CBC
func (qtaes *qtAes) SetIv(iv interface{}) {
	switch tmp := iv.(type) {
	case []byte:
		qtaes.aesIv = tmp
	case string:
		qtaes.aesIv = []byte(tmp)
	}
}

// 设置加密类型，请使用 [Type_Const_AES_DES_]选择 默认ECB
func (qtaes *qtAes) SetEncMethod(Method int) {
	qtaes.encMethod = Method
}

// 设置key,自动识别是字符串还是字节数组
//
// # Key 若为16位则加密为128
//
// # Key 若为24位则加密为192
//
// Key 若为32位则加密为256
func (qtaes *qtAes) SetKey(key interface{}) {
	switch tmp := key.(type) {
	case []byte:
		qtaes.aesKey = tmp
	case string:
		qtaes.aesKey = []byte(tmp)
	}
}

// 请使用 [Type_Const_Padding_]来选择
func (qtaes *qtAes) SetPadding(pkcs int) {
	qtaes.pacs = pkcs
}
func (qtaes *qtAes) zeropadding(ciphertext []byte, blockSize int) []byte {
	Bytelen := len(ciphertext)
	if Bytelen < 1 {
		return []byte{}
	}
	padding := blockSize - Bytelen%blockSize
	tmp := ciphertext
	if padding%blockSize != 0 {
		if qtaes.fill {
			padtext := bytes.Repeat([]byte{byte(padding)}, padding)
			tmp = append(ciphertext, padtext...)
		} else {
			padtext := bytes.Repeat([]byte{byte(0)}, padding)
			tmp = append(ciphertext, padtext...)
		}

	}
	return tmp
}

func (qtaes *qtAes) ansiX923padding(ciphertext []byte, blockSize int) []byte {
	Bytelen := len(ciphertext)
	if Bytelen < 1 {
		return []byte{}
	}
	padding := blockSize - Bytelen%blockSize
	tmp := ciphertext
	if qtaes.fill {
		padtext := bytes.Repeat([]byte{byte(padding)}, padding-1)
		tmp = append(tmp, padtext...)
		padtext = bytes.Repeat([]byte{byte(padding)}, 1)
		tmp = append(tmp, padtext...)
	} else {
		padtext := bytes.Repeat([]byte{byte(0)}, padding-1)
		tmp = append(tmp, padtext...)
		padtext = bytes.Repeat([]byte{byte(padding)}, 1)
		tmp = append(tmp, padtext...)
	}
	return tmp
}
func (qtaes *qtAes) keyFactory(key []byte) []byte {
	_keyLen := 0
	_keySize := 0
	STMP := []byte{}
	_keyLen = len(key)
	if qtaes.keySize == Type_Const_AES_keySize_Auto {
		if _keyLen >= 32 {
			_keySize = 8
		} else if _keyLen >= 24 {
			_keySize = 6
		} else {
			_keySize = 4
		}

	} else if qtaes.keySize == Type_Const_AES_keySize_256 {
		_keySize = 8
	} else if qtaes.keySize == Type_Const_AES_keySize_192 {
		_keySize = 6
	} else {
		_keySize = 4
	}
	if _keyLen < _keySize*4 {
		tmp := append(STMP, bytes.Repeat([]byte{byte(0)}, _keySize*4-_keyLen)...)
		key = append(key, tmp...)
	}
	return key[:_keySize*4]
}

// 默认PKSC#7
//
// # Key 若为16位则加密为128
//
// # Key 若为24位则加密为192
//
// Key 若为32位则加密为256
func (qtaes *qtAes) Encrypt(Data []byte) []byte {

	block, err := aes.NewCipher(qtaes.keyFactory(qtaes.aesKey))
	if err != nil {
		qtaes.err = err
		return nil
	}
	in := []byte{}
	blockSize := block.BlockSize()
	switch qtaes.pacs {
	case Type_Const_Padding_NoPadding:
		in = Data
	case Type_Const_Padding_Zero:
		in = qtaes.zeropadding(Data, blockSize)
	case Type_Const_Padding_Iso10126:
		in = iso10126padding(Data, blockSize)
	case Type_Const_Padding_AnsiX923:
		in = qtaes.ansiX923padding(Data, blockSize)
	case Type_Const_Padding_Iso97971:
		in = qtaes.zeropadding(append(Data, bytes.Repeat([]byte{byte(128)}, 1)...), blockSize)
	case Type_Const_Padding_Pkcs5, Type_Const_Padding_Pkcs7:
		in = qtaes.padding(Data, blockSize)
	default:
		in = qtaes.padding(Data, blockSize)
	}
	if qtaes.encMethod == Type_Const_AES_DES_ECB {
		n := 0
		encryptData := make([]byte, len(in))
		tmpData := make([]byte, blockSize)
		for index := 0; index < len(in); index += blockSize {
			block.Encrypt(tmpData, in[index:index+blockSize])
			for i := 0; i < blockSize; i++ {
				encryptData[blockSize*n+i] = tmpData[i]
			}
			n++
		}
		return encryptData
	}
	if qtaes.encMethod == Type_Const_AES_DES_FCB {
		blockMode := cipher.NewCFBEncrypter(block, qtaes.ivFactory(qtaes.aesIv))
		crypted := make([]byte, len(in))
		blockMode.XORKeyStream(crypted, in)
		return crypted
	}
	if qtaes.encMethod == Type_Const_AES_DES_CBC {

		blockMode := cipher.NewCBCEncrypter(block, qtaes.ivFactory(qtaes.aesIv))
		crypted := make([]byte, len(in))
		blockMode.CryptBlocks(crypted, in)
		return crypted
	}
	qtaes.err = errors.New("encMethod err ")
	return nil
}

func (qtaes *qtAes) ivFactory(iv []byte) []byte {
	_ivLen := len(iv)
	STMP := []byte{}
	if _ivLen < 8 {
		tmp := append(STMP, bytes.Repeat([]byte{byte(0)}, 16-_ivLen)...)
		tmp = append(iv, tmp...)
		return tmp
	}
	return iv[:16]
}

var NotFullBlocks = errors.New("crypto/cipher: input not full blocks")

// 默认PKSC#7
//
// # Key 若为16位则加密为128
//
// # Key 若为24位则加密为192
//
// Key 若为32位则加密为256
func (qtaes *qtAes) Decrypt(Data []byte) []byte {

	block, err := aes.NewCipher(qtaes.keyFactory(qtaes.aesKey))
	if err != nil {
		qtaes.err = err
		return nil
	}
	blockSize := block.BlockSize()
	if len(Data)%blockSize != 0 {
		qtaes.err = NotFullBlocks
		return nil
	}
	if qtaes.encMethod == Type_Const_AES_DES_ECB {
		n := 0
		decryptData := make([]byte, len(Data))
		tmpData := make([]byte, blockSize)
		for index := 0; index < len(Data); index += blockSize {
			block.Decrypt(tmpData, Data[index:index+blockSize])
			for i := 0; i < blockSize; i++ {
				decryptData[blockSize*n+i] = tmpData[i]
			}
			n++
		}
		if qtaes.pacs != Type_Const_Padding_NoPadding {
			decryptData = qtaes.unpadding(decryptData)
		}
		return decryptData
	}

	if qtaes.encMethod == Type_Const_AES_DES_CBC {
		blockMode := cipher.NewCBCDecrypter(block, qtaes.ivFactory(qtaes.aesIv))
		origData := make([]byte, len(Data))
		blockMode.CryptBlocks(origData, Data)
		if qtaes.pacs != Type_Const_Padding_NoPadding {
			origData = qtaes.unpadding(origData)
		}
		return origData
	}
	if qtaes.encMethod == Type_Const_AES_DES_FCB {
		blockMode := cipher.NewCFBDecrypter(block, qtaes.ivFactory(qtaes.aesIv))
		origData := make([]byte, len(Data))
		blockMode.XORKeyStream(origData, Data)
		if qtaes.pacs != Type_Const_Padding_NoPadding {
			origData = qtaes.unpadding(origData)
		}
		return origData
	}
	qtaes.err = errors.New("encMethod err ")
	return nil
}
func (qtaes *qtAes) padding(ciphertext []byte, blockSize int) []byte {
	if qtaes.fill {
		padding := blockSize - len(ciphertext)%blockSize
		padtext := bytes.Repeat([]byte{byte(padding)}, padding)
		return append(ciphertext, padtext...)
	}
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(0)}, padding)
	return append(ciphertext, padtext...)
}
func (qtaes *qtAes) removeZero(data []byte) []byte {
	max := len(data)
	if max < 16 {
		return data
	}
	EndIndex := bytes.Index(data[16:], []byte{0})
	if EndIndex < 0 {
		EndIndex = max
	} else {
		EndIndex += 16
	}
	return data[0:EndIndex]
}
func (qtaes *qtAes) unpadding(origData []byte) []byte {

	if qtaes.pacs == Type_Const_Padding_Iso97971 {
		length := len(origData)
		// 去掉最后一个字节 unpadding 次
		unpadding := int(origData[length-1])

		l := (length - unpadding) - 1
		if len(origData) >= l {
			return qtaes.removeZero(origData[:l])
		}
		return qtaes.removeZero(origData)
	}
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	l := length - unpadding
	if len(origData) >= l && l >= 0 {
		return qtaes.removeZero(origData[:l])
	}
	return qtaes.removeZero(origData)
}

// 获取一个AES实例对象
func NewAes() *qtAes {
	rand.Seed(time.Now().UnixNano())
	tmp := new(qtAes)
	tmp.fill = true
	return tmp
}
