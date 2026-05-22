package DebugTools

import (
	"bytes"
	"changeme/Service/Tools/DebugTools/WithSalt"
	"encoding/base64"
	"github.com/deatil/go-cryptobin/cryptobin/sm2"
	"github.com/qtgolang/SunnyNet/src/encoding/hex"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"os"
	"strings"
)

type Encryption struct {
}
type EncryptionInfo struct {
	Children  []string `json:"children"`
	Key       string   `json:"key"`
	Iv        string   `json:"iv"`
	IvType    string   `json:"ivType"`
	KeyType   string   `json:"keyType"`
	IsEncrypt bool     `json:"isEncrypt"`
	SM2Mode   string   `json:"SM2Mode"`
	Data      string   `json:"Data"`
	output    bytes.Buffer
}

func (e *EncryptionInfo) key() []byte {
	switch e.KeyType {
	case "String":
		return []byte(e.Key)
	case "Base64":
		bs, _ := base64.StdEncoding.DecodeString(strings.ReplaceAll(strings.TrimSpace(e.Key), " ", ""))
		return bs
	case "HEX":
		bs, _ := hex.DecodeString(strings.ReplaceAll(strings.TrimSpace(e.Key), " ", ""))
		return bs
	default:
		return nil
	}
}
func (e *EncryptionInfo) Public() []byte {
	s := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(e.Key, "\\n", "\n"), " ----", "\n----"), "---- ", "----\n")
	arr := strings.Split(s, "\n")
	st := ""
	for _, k := range arr {
		m := strings.ReplaceAll(strings.TrimSpace(k), " ", "")
		if strings.HasPrefix(m, "----") {
			continue
		}
		st += m
	}
	return []byte("-----BEGIN PUBLIC KEY-----\n" + st + "\n-----END PUBLIC KEY-----")
}
func (e *EncryptionInfo) Private() []byte {
	s := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(e.Key, "\\n", "\n"), " ----", "\n----"), "---- ", "----\n")
	arr := strings.Split(s, "\n")
	st := ""
	for _, k := range arr {
		m := strings.ReplaceAll(strings.TrimSpace(k), " ", "")
		if strings.HasPrefix(m, "----") {
			continue
		}
		st += m
	}
	return []byte("-----BEGIN PRIVATE KEY-----\n" + st + "\n-----END PRIVATE KEY-----")
}
func (e *EncryptionInfo) iv() []byte {
	switch e.IvType {
	case "String":
		return []byte(e.Iv)
	case "Base64":
		bs, _ := base64.StdEncoding.DecodeString(strings.ReplaceAll(strings.TrimSpace(e.Iv), " ", ""))
		return bs
	case "HEX":
		bs, _ := hex.DecodeString(strings.ReplaceAll(strings.TrimSpace(e.Iv), " ", ""))
		return bs
	default:
		return nil
	}
}
func (e *EncryptionInfo) WriteOutput(Data []byte) {
	if len(e.Children) < 2 {
		return
	}
	switch e.Children[1] {
	case "到字符串":
		e.output.Write(Data)
		return
	case "到Hex":
		e.output.WriteString(hex.EncodeToString(Data))
		return
	case "到Base64":
		e.output.WriteString(base64.StdEncoding.EncodeToString(Data))
		return
	default:
		return
	}
}
func (e *EncryptionInfo) DataBytes() []byte {
	if len(e.Children) < 1 {
		return nil
	}
	switch e.Children[0] {
	case "输入为GBK字符串":
		reader := transform.NewReader(strings.NewReader(e.Data), simplifiedchinese.GBK.NewEncoder())
		data, _ := io.ReadAll(reader)
		return data
	case "输入为UTF8字符串":
		return []byte(e.Data)
	case "从HEX":
		data, _ := hex.DecodeString(strings.ReplaceAll(strings.TrimSpace(e.Data), " ", ""))
		return data
	case "从Base64":
		data, _ := base64.StdEncoding.DecodeString(strings.ReplaceAll(strings.TrimSpace(e.Data), " ", ""))
		return data
	case "从文件":
		data, _ := os.ReadFile(e.Data)
		return data
	default:
		return nil
	}
}
func (e *EncryptionInfo) apply() []byte {
	if len(e.Children) < 3 {
		return e.output.Bytes()
	}
	switch e.Children[2] {
	case "使用Hash":
		e.hash()
		return e.output.Bytes()
	case "使用Hmac-Hash":
		e.hmac()
		return e.output.Bytes()
	case "使用AES":
		e.aes()
		return e.output.Bytes()
	case "使用DES":
		e.des()
		return e.output.Bytes()
	case "使用3DES":
		e.desEde()
		return e.output.Bytes()
	case "使用SM2":
		e.sm2()
		return e.output.Bytes()
	case "使用RC4":
		e.rc4()
		return e.output.Bytes()
	case "使用RC4-WithSalt":
		e.rc4WithSalt()
		return e.output.Bytes()
	default:
		return e.output.Bytes()
	}
}
func (e *EncryptionInfo) rc4WithSalt() {
	if e.IsEncrypt {
		e.WriteOutput(WithSalt.RC4WithSaltEncrypt(e.DataBytes(), e.key()))
	} else {
		e.WriteOutput(WithSalt.RC4WithSaltEncrypt(e.DataBytes(), e.key()))
	}
}
func (e *EncryptionInfo) rc4() {
	if e.IsEncrypt {
		e.WriteOutput(Rc4Encrypt(e.DataBytes(), e.key()))
	} else {
		e.WriteOutput(Rc4Decrypt(e.DataBytes(), e.key()))
	}
}
func (e *EncryptionInfo) sm2() {
	obj := sm2.New()
	if e.IsEncrypt {
		_, err := obj.ParsePublicKeyFromPEM(e.Public())
		if err != nil {
			e.output.WriteString("公钥格式错误")
			return
		}
		var deData = obj.FromBytes(e.DataBytes()).FromPublicKey(e.Public()).WithEncoding(1).SetMode(e.SM2Mode).Encrypt()
		e.WriteOutput(deData.ToBytes())
		return
	}
	_, err := obj.ParsePKCS8PrivateKeyFromPEM(e.Private())
	if err != nil {
		e.output.WriteString("私钥格式错误")
		return
	}
	var deData = obj.FromBytes(e.DataBytes()).FromPrivateKey(e.Private()).WithEncoding(1).SetMode(e.SM2Mode).Decrypt()
	e.WriteOutput(deData.ToBytes())
	return
}
func (e *EncryptionInfo) desEde() {

	if len(e.Children) < 4 {
		return
	}
	switch e.Children[3] {
	case "WithSalt":
		if e.IsEncrypt {
			e.WriteOutput(WithSalt.TripleDesCbcWithSaltEncrypt(e.DataBytes(), e.key()))
		} else {
			e.WriteOutput(WithSalt.TripleDesCbcWithSaltDecrypt(e.DataBytes(), e.key()))
		}
		return
	case "ECB":
		a := NewDes()
		a.SetEncMethod(Type_Const_DES_3DES_ECB)
		a.SetKey(e.key())
		if e.IsEncrypt {
			e.WriteOutput(a.Encrypt(e.DataBytes()))
		} else {
			e.WriteOutput(a.Decrypt(e.DataBytes()))
		}
		return
	case "CBC":
		a := NewDes()
		a.SetEncMethod(Type_Const_DES_3DES_CBC)
		a.SetKey(e.key())
		a.SetIv(e.iv())
		if e.IsEncrypt {
			e.WriteOutput(a.Encrypt(e.DataBytes()))
		} else {
			e.WriteOutput(a.Decrypt(e.DataBytes()))
		}
		return
	default:
		return
	}
}
func (e *EncryptionInfo) des() {
	if len(e.Children) < 4 {
		return
	}
	switch e.Children[3] {
	case "WithSalt":
		if e.IsEncrypt {
			e.WriteOutput(WithSalt.DesCbcWithSaltEncrypt(e.DataBytes(), e.key()))
		} else {
			e.WriteOutput(WithSalt.DesCbcWithSaltDecrypt(e.DataBytes(), e.key()))
		}
		return
	case "ECB":
		a := NewDes()
		a.SetEncMethod(Type_Const_AES_DES_ECB)
		a.SetKey(e.key())
		if e.IsEncrypt {
			e.WriteOutput(a.Encrypt(e.DataBytes()))
		} else {
			e.WriteOutput(a.Decrypt(e.DataBytes()))
		}
		return
	case "CBC":
		a := NewDes()
		a.SetEncMethod(Type_Const_AES_DES_CBC)
		a.SetKey(e.key())
		a.SetIv(e.iv())
		if e.IsEncrypt {
			e.WriteOutput(a.Encrypt(e.DataBytes()))
		} else {
			e.WriteOutput(a.Decrypt(e.DataBytes()))
		}
		return
	default:
		return
	}
}
func (e *EncryptionInfo) aes() {
	if len(e.Children) < 4 {
		return
	}
	switch e.Children[3] {
	case "WithSalt-256":
		if e.IsEncrypt {
			e.WriteOutput(WithSalt.AesCbc256WithSaltEncrypt(e.DataBytes(), e.key()))
		} else {
			e.WriteOutput(WithSalt.AesCbc256WithSaltDecrypt(e.DataBytes(), e.key()))
		}
		return
	case "WithSalt-192":
		if e.IsEncrypt {
			e.WriteOutput(WithSalt.AesCbc192WithSaltEncrypt(e.DataBytes(), e.key()))
		} else {
			e.WriteOutput(WithSalt.AesCbc192WithSaltDecrypt(e.DataBytes(), e.key()))
		}
		return
	case "WithSalt-128":
		if e.IsEncrypt {
			e.WriteOutput(WithSalt.AesCbc128WithSaltEncrypt(e.DataBytes(), e.key()))
		} else {
			e.WriteOutput(WithSalt.AesCbc128WithSaltDecrypt(e.DataBytes(), e.key()))
		}
		return
	case "ECB":
		a := NewAes()
		a.SetEncMethod(Type_Const_AES_DES_ECB)
		a.SetKey(e.key())
		if e.IsEncrypt {
			e.WriteOutput(a.Encrypt(e.DataBytes()))
		} else {
			e.WriteOutput(a.Decrypt(e.DataBytes()))
		}
		return
	case "CBC":
		a := NewAes()
		a.SetEncMethod(Type_Const_AES_DES_CBC)
		a.SetKey(e.key())
		a.SetIv(e.iv())
		if e.IsEncrypt {
			e.WriteOutput(a.Encrypt(e.DataBytes()))
		} else {
			e.WriteOutput(a.Decrypt(e.DataBytes()))
		}
		return
	default:
		return
	}
}
func (e *EncryptionInfo) hmac() {
	if len(e.Children) < 4 {
		return
	}
	switch e.Children[3] {
	case "MD2":
		e.WriteOutput(CalculateHMACMD2(e.DataBytes(), e.key()))
		return
	case "MD4":
		e.WriteOutput(CalculateHMACMD4(e.DataBytes(), e.key()))
		return
	case "MD5":
		e.WriteOutput(CalculateHMACMD5(e.DataBytes(), e.key()))
		return
	case "SHA-1":
		e.WriteOutput(CalculateHMACSHA1(e.DataBytes(), e.key()))
		return
	case "SHA-256":
		e.WriteOutput(CalculateHMACSHA256(e.DataBytes(), e.key()))
		return
	case "SHA-512":
		e.WriteOutput(CalculateHMACSHA512(e.DataBytes(), e.key()))
		return
	default:
		return
	}
}
func (e *EncryptionInfo) hash() {
	if len(e.Children) < 4 {
		return
	}
	switch e.Children[3] {
	case "MD2":
		e.WriteOutput(CalculateMD2(e.DataBytes()))
		return
	case "MD4":
		e.WriteOutput(CalculateMD4(e.DataBytes()))
		return
	case "MD5":
		e.WriteOutput(CalculateMD5(e.DataBytes()))
		return
	case "SHA-1":
		e.WriteOutput(CalculateSHA1(e.DataBytes()))
		return
	case "SHA-256":
		e.WriteOutput(CalculateSHA256(e.DataBytes()))
		return
	case "SHA-512":
		e.WriteOutput(CalculateSHA512(e.DataBytes()))
		return
	default:
		return
	}
}
func (g *Encryption) AppEncryptCall(obj EncryptionInfo) []byte {
	return obj.apply()
}
