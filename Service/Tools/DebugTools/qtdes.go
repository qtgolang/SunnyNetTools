package DebugTools

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"errors"
	"math/rand"
	"time"
)

type qtDes struct {
	desIv     []byte
	desKey    []byte
	pacs      int
	encMethod int
	err       error
	fill      bool
}

// 获取运行时的错误
func (qtdes *qtDes) GetErr() string {
	if qtdes.err != nil {
		return qtdes.err.Error()
	}
	return ""
}

// 设置填充模式，如果为true,则按填充位填充，若为false,则按0填充
//
// 默认按位填充
//
// 则按填充位填充意思是，比如数据长度10，每次加密16 则填充6个6字节
func (qtdes *qtDes) SetFill(fill bool) {
	qtdes.fill = fill
}

// 设置IV,自动识别是字符串还是字节数组
//
// 若不设置IV则ECB，若设置IV则CBC
func (qtdes *qtDes) SetIv(iv interface{}) {
	switch tmp := iv.(type) {
	case []byte:
		qtdes.desIv = tmp
	case string:
		qtdes.desIv = []byte(tmp)
	}
}

// 设置加密类型，请使用 [Type.AES_DES_选择]默认ECB
func (qtdes *qtDes) SetEncMethod(Method int) {
	qtdes.encMethod = Method
}

// 设置key,自动识别是字符串还是字节数组
//
// # Key 若为16位则加密为128
//
// # Key 若为24位则加密为192
//
// Key 若为32位则加密为256
func (qtdes *qtDes) SetKey(key interface{}) {
	switch tmp := key.(type) {
	case []byte:
		qtdes.desKey = tmp
	case string:
		qtdes.desKey = []byte(tmp)
	}
}

// 请使用 [Type_Const_Padding_]来选择
func (qtdes *qtDes) SetPadding(pkcs int) {
	qtdes.pacs = pkcs
}
func (qtdes *qtDes) zeropadding(ciphertext []byte, blockSize int) []byte {
	Bytelen := len(ciphertext)
	if Bytelen < 1 {
		return []byte{}
	}
	padding := blockSize - Bytelen%blockSize
	tmp := ciphertext
	if padding%blockSize != 0 {
		if qtdes.fill {
			padtext := bytes.Repeat([]byte{byte(padding)}, padding)
			tmp = append(ciphertext, padtext...)
		} else {
			padtext := bytes.Repeat([]byte{byte(0)}, padding)
			tmp = append(ciphertext, padtext...)
		}

	}
	return tmp
}
func iso10126padding(ciphertext []byte, blockSize int) []byte {
	Bytelen := len(ciphertext)
	if Bytelen < 1 {
		return []byte{}
	}
	padding := blockSize - Bytelen%blockSize
	tmp := ciphertext
	for i := 0; i < padding-1; i++ {
		x := rand.Intn(254) + 1
		padtext := bytes.Repeat([]byte{byte(x)}, 1)
		tmp = append(tmp, padtext...)
	}
	padtext := bytes.Repeat([]byte{byte(padding)}, 1)
	tmp = append(tmp, padtext...)
	return tmp
}
func (qtdes *qtDes) ansiX923padding(ciphertext []byte, blockSize int) []byte {
	Bytelen := len(ciphertext)
	if Bytelen < 1 {
		return []byte{}
	}
	padding := blockSize - Bytelen%blockSize
	tmp := ciphertext
	if qtdes.fill {
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
func (qtdes *qtDes) keyFactory(key []byte) []byte {
	_keyLen := 0
	STMP := []byte{}
	if qtdes.encMethod == Type_Const_DES_3DES_CBC || qtdes.encMethod == Type_Const_DES_3DES_ECB || qtdes.encMethod == Type_Const_DES_3DES_FCB {
		_keyLen = len(key)
		if _keyLen < 8 {
			tmp := append(STMP, bytes.Repeat([]byte{byte(0)}, 8-_keyLen)...)
			tmp = append(key, tmp...)
			tmp = append(tmp, tmp[:8]...)
			tmp = append(tmp, tmp[:8]...)
			return tmp
		} else if _keyLen < 16 {
			tmp := append(STMP, key[:8]...)
			tmp = append(tmp, tmp[:8]...)
			tmp = append(tmp, tmp[:8]...)
			return tmp
		} else if _keyLen < 24 {
			tmp := append(STMP, key[:16]...)
			tmp1 := append(STMP, key[:8]...)
			tmp = append(tmp, tmp1...)
			return tmp
		}
		return key[:24]
	}
	_keyLen = len(key)
	if _keyLen < 8 {
		tmp := append(STMP, bytes.Repeat([]byte{byte(0)}, 8-_keyLen)...)
		tmp = append(key, tmp...)
		return tmp
	}
	return key[:8]
}
func (qtdes *qtDes) ivFactory(iv []byte) []byte {
	_ivLen := len(iv)
	STMP := []byte{}
	if _ivLen < 8 {
		tmp := append(STMP, bytes.Repeat([]byte{byte(0)}, 8-_ivLen)...)
		tmp = append(iv, tmp...)
		return tmp
	}
	return iv[:8]
}

// 默认PKSC#7
//
// # Key 若为16位则加密为128
//
// # Key 若为24位则加密为192
//
// Key 若为32位则加密为256
func (qtdes *qtDes) Encrypt(Data []byte) []byte {
	var block cipher.Block
	var err error
	if qtdes.encMethod == Type_Const_DES_3DES_CBC || qtdes.encMethod == Type_Const_DES_3DES_ECB || qtdes.encMethod == Type_Const_DES_3DES_FCB {
		block, err = des.NewTripleDESCipher(qtdes.keyFactory(qtdes.desKey))
	} else {
		block, err = des.NewCipher(qtdes.keyFactory(qtdes.desKey))
	}

	if err != nil {
		qtdes.err = err
		return nil
	}
	in := []byte{}
	blockSize := block.BlockSize()
	switch qtdes.pacs {
	case Type_Const_Padding_NoPadding:
		in = Data
	case Type_Const_Padding_Zero:
		in = qtdes.zeropadding(Data, blockSize)
	case Type_Const_Padding_Iso10126:
		in = iso10126padding(Data, blockSize)
	case Type_Const_Padding_AnsiX923:
		in = qtdes.ansiX923padding(Data, blockSize)
	case Type_Const_Padding_Iso97971:
		in = qtdes.zeropadding(append(Data, bytes.Repeat([]byte{byte(128)}, 1)...), blockSize)
	case Type_Const_Padding_Pkcs5, Type_Const_Padding_Pkcs7:
		in = qtdes.padding(Data, blockSize)
	default:
		in = qtdes.padding(Data, blockSize)
	}
	if qtdes.encMethod == Type_Const_AES_DES_ECB || qtdes.encMethod == Type_Const_DES_3DES_ECB {
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
	if qtdes.encMethod == Type_Const_AES_DES_FCB || qtdes.encMethod == Type_Const_DES_3DES_FCB {
		blockMode := cipher.NewCFBEncrypter(block, qtdes.ivFactory(qtdes.desIv))
		crypted := make([]byte, len(in))
		blockMode.XORKeyStream(crypted, in)
		return crypted
	}
	if qtdes.encMethod == Type_Const_AES_DES_CBC || qtdes.encMethod == Type_Const_DES_3DES_CBC {
		blockMode := cipher.NewCBCEncrypter(block, qtdes.ivFactory(qtdes.desIv))
		crypted := make([]byte, len(in))
		blockMode.CryptBlocks(crypted, in)
		return crypted
	}
	qtdes.err = errors.New("encMethod err ")
	return nil
}

// 默认PKSC#7
//
// # Key 若为16位则加密为128
//
// # Key 若为24位则加密为192
//
// Key 若为32位则加密为256
func (qtdes *qtDes) Decrypt(Data []byte) []byte {
	var block cipher.Block
	var err error
	if qtdes.encMethod == Type_Const_DES_3DES_CBC || qtdes.encMethod == Type_Const_DES_3DES_ECB || qtdes.encMethod == Type_Const_DES_3DES_FCB {
		block, err = des.NewTripleDESCipher(qtdes.keyFactory(qtdes.desKey))
	} else {
		block, err = des.NewCipher(qtdes.keyFactory(qtdes.desKey))
	}
	if err != nil {
		qtdes.err = err
		return nil
	}
	blockSize := block.BlockSize()
	if len(Data)%blockSize != 0 {
		qtdes.err = NotFullBlocks
		return nil
	}
	if qtdes.encMethod == Type_Const_AES_DES_ECB || qtdes.encMethod == Type_Const_DES_3DES_ECB {
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
		if qtdes.pacs != Type_Const_Padding_NoPadding {
			decryptData = qtdes.unpadding(decryptData)
		}
		return decryptData
	}
	if qtdes.encMethod == Type_Const_AES_DES_CBC || qtdes.encMethod == Type_Const_DES_3DES_CBC {
		blockMode := cipher.NewCBCDecrypter(block, qtdes.ivFactory(qtdes.desIv))
		origData := make([]byte, len(Data))
		blockMode.CryptBlocks(origData, Data)
		if qtdes.pacs != Type_Const_Padding_NoPadding {
			origData = qtdes.unpadding(origData)
		}
		return origData
	}
	if qtdes.encMethod == Type_Const_AES_DES_FCB || qtdes.encMethod == Type_Const_DES_3DES_FCB {
		blockMode := cipher.NewCFBDecrypter(block, qtdes.ivFactory(qtdes.desIv))
		origData := make([]byte, len(Data))
		blockMode.XORKeyStream(origData, Data)
		if qtdes.pacs != Type_Const_Padding_NoPadding {
			origData = qtdes.unpadding(origData)
		}
		return origData
	}
	qtdes.err = errors.New("encMethod err ")
	return nil
}
func (qtdes *qtDes) padding(ciphertext []byte, blockSize int) []byte {
	if qtdes.fill {
		padding := blockSize - len(ciphertext)%blockSize
		padtext := bytes.Repeat([]byte{byte(padding)}, padding)
		return append(ciphertext, padtext...)
	}
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(0)}, padding)
	return append(ciphertext, padtext...)
}
func (qtdes *qtDes) removeZero(data []byte) []byte {
	max := len(data)
	if max < 8 {
		return data
	}
	EndIndex := bytes.Index(data[8:], []byte{0})
	if EndIndex < 0 {
		EndIndex = max
	} else {
		EndIndex += 8
	}
	return data[0:EndIndex]
}
func (qtdes *qtDes) unpadding(origData []byte) []byte {
	length := len(origData)
	if length < 1 {
		return []byte{}
	}
	unpadding := int(origData[length-1])
	if qtdes.pacs == Type_Const_Padding_Iso97971 {
		l := (length - unpadding) - 1
		if len(origData) >= l {
			return qtdes.removeZero(origData[:l])
		}
		return qtdes.removeZero(origData)
	}
	if !qtdes.fill {
		var tmp []byte
		i := len(origData)
		max := i
		for {
			if max >= i {
				if origData[i-1] != 0 {
					tmp = origData[:i]
					break
				}
			}
			i--
			if i < 0 {
				break
			}
		}
		return tmp
	}
	l := length - unpadding
	if len(origData) >= l && l >= 0 {
		return qtdes.removeZero(origData[:l])
	}
	return qtdes.removeZero(origData)

}

// 获取一个des实例对象
func NewDes() *qtDes {
	rand.Seed(time.Now().UnixNano())
	tmp := new(qtDes)
	tmp.fill = true
	return tmp
}

const (
	// Hash_Stat Hash开始处的标志
	Type_Const_Hash_Stat       = 1
	Type_Const_Hash_Sha1       = 1
	Type_Const_Hash_Md4        = 2
	Type_Const_Hash_Md5        = 3
	Type_Const_Hash_Sha224     = 4
	Type_Const_Hash_Sha256     = 5
	Type_Const_Hash_Sha384     = 6
	Type_Const_Hash_Sha512     = 7
	Type_Const_Hash_Sha512_224 = 8
	Type_Const_Hash_Sha512_256 = 9
	// Hash_End Hash结束处的标志
	Type_Const_Hash_End = 9

	Type_Const_Padding_Pkcs5     = 0
	Type_Const_Padding_Pkcs7     = 1
	Type_Const_Padding_NoPadding = 2
	Type_Const_Padding_Zero      = 3
	Type_Const_Padding_Iso10126  = 4
	Type_Const_Padding_Iso97971  = 5
	Type_Const_Padding_AnsiX923  = 6
	Type_Const_AES_DES_ECB       = 0
	Type_Const_AES_DES_CBC       = 1
	Type_Const_AES_DES_FCB       = 2
	Type_Const_AES_keySize_128   = 128
	Type_Const_AES_keySize_192   = 192
	Type_Const_AES_keySize_256   = 256
	Type_Const_AES_keySize_Auto  = 0
	Type_Const_DES_3DES_ECB      = 3
	Type_Const_DES_3DES_CBC      = 4
	Type_Const_DES_3DES_FCB      = 5
)
