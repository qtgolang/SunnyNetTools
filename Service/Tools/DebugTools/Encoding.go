package DebugTools

import (
	"bytes"
	"changeme/Service/Session"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"github.com/qtgolang/SunnyNet/src/Compress"
	"github.com/qtgolang/SunnyNet/src/encoding/hex"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"hash/crc32"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func GetMD5Hash(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()

	hash := md5.New()
	io.Copy(hash, file)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func GetCRC32Hash(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()

	hash := crc32.NewIEEE()
	io.Copy(hash, file)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func GetSHA1Hash(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()

	hash := sha1.New()
	io.Copy(hash, file)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func GetSHA256Hash(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()

	hash := sha256.New()
	io.Copy(hash, file)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func GetSHA512Hash(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()

	hash := sha512.New()
	io.Copy(hash, file)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

var fileHashCodeFunc = map[string]func(filePath string) string{
	"CRC32":   GetCRC32Hash,
	"MD5":     GetMD5Hash,
	"SHA-1":   GetSHA1Hash,
	"SHA-256": GetSHA256Hash,
	"SHA-512": GetSHA512Hash,
}

type Encoding struct {
}
type EncodingInfo struct {
	Children []string `json:"children"`
	Data     string   `json:"Data"`
	output   bytes.Buffer
}

func (e *EncodingInfo) DataBytes() []byte {
	if len(e.Children) < 1 {
		return nil
	}
	switch e.Children[0] {
	case "字符串":
		return []byte(e.Data)
	case "HEX":
		data, _ := hex.DecodeString(strings.ReplaceAll(strings.TrimSpace(e.Data), " ", ""))
		return data
	case "Base64":
		data, _ := base64.StdEncoding.DecodeString(strings.ReplaceAll(strings.TrimSpace(e.Data), " ", ""))
		return data
	default:
		return nil
	}
}
func (e *EncodingInfo) Base(Data []byte) {
	if len(e.Children) < 2 {
		return
	}
	switch e.Children[1] {
	case "URL编码":
		e.output.Write([]byte(url.QueryEscape(string(Data))))
		return
	case "URL解码":
		U, _ := url.QueryUnescape(string(Data))
		e.output.Write([]byte(U))
		return
	case "解码到字符串":
		e.output.Write(Data)
		return
	case "转到Base64", "解码到Base64":
		e.output.WriteString(base64.StdEncoding.EncodeToString(Data))
		return
	case "转到HEX", "解码到Hex":
		e.output.WriteString(hex.EncodeToString(Data))
		return
	default:
		return
	}
}
func (e *EncodingInfo) ExtendFile() {
	if len(e.Children) < 2 {
		return
	}
	if e.Children[1] == "到字符串" || e.Children[1] == "到图片预览" {
		bs, _ := os.ReadFile(strings.ReplaceAll(strings.TrimSpace(e.Data), "\"", ""))
		e.output.Write(bs)
		return
	}
	if e.Children[1] == "到Hex" {
		bs, _ := os.ReadFile(strings.ReplaceAll(strings.TrimSpace(e.Data), "\"", ""))
		e.output.WriteString(strings.ToUpper(hex.EncodeToString(bs)))
		return
	}
	if e.Children[1] == "效验" {
		if len(e.Children) < 3 {
			return
		}
		f := fileHashCodeFunc[e.Children[2]]
		if f == nil {
			return
		}
		e.output.WriteString(strings.ToUpper(f(strings.ReplaceAll(strings.TrimSpace(e.Data), "\"", ""))))
		return
	}
	if e.Children[1] == "到Base64" {
		bs, _ := os.ReadFile(strings.ReplaceAll(strings.TrimSpace(e.Data), "\"", ""))
		e.output.WriteString(base64.StdEncoding.EncodeToString(bs))
		return
	}
}
func (e *EncodingInfo) ExtendString() {
	if len(e.Children) < 2 {
		return
	}
	switch e.Children[1] {
	case "到大写":
		e.output.WriteString(strings.ToUpper(e.Data))
		return
	case "到小写":
		e.output.WriteString(strings.ToLower(e.Data))
		return
	case "删除":
		if len(e.Children) < 2 {
			return
		}
		switch e.Children[1] {
		case "全部空格":
			e.output.WriteString(strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(e.Data), "\t", ""), " ", ""))
			return
		case "全部换行及空格":
			e.output.WriteString(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(e.Data), "\r", ""), "\n", ""), "\t", ""), " ", ""))
			return
		case "全部换行":
			e.output.WriteString(strings.ReplaceAll(strings.ReplaceAll(e.Data, "\r", ""), "\n", ""))
			return
		}
		return
	case "取长度":
		reader := transform.NewReader(strings.NewReader(e.Data), simplifiedchinese.GBK.NewEncoder())
		gbkStr, _ := io.ReadAll(reader)

		data := []byte(e.Data)
		_hexLen, _ := hex.DecodeString(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(e.Data), "\r", ""), "\n", ""), "\t", ""), " ", ""))
		_base64Len, _ := base64.StdEncoding.DecodeString(e.Data)

		e.output.WriteString("字节长度(GBK):" + strconv.Itoa(len(gbkStr)) + "\n")
		e.output.WriteString("字节长度(UTF8):" + strconv.Itoa(len(data)) + "\n")
		e.output.WriteString("-----------------------------------------------------\n")
		e.output.WriteString("字符长度(GBK):" + strconv.Itoa(len(string(gbkStr))) + "\n")
		e.output.WriteString("字符长度(UTF8):" + strconv.Itoa(len(string(data))) + "\n")
		e.output.WriteString("-----------------------------------------------------\n")
		e.output.WriteString("Hex解码后字节长度:" + strconv.Itoa(len(_hexLen)) + "\n")
		e.output.WriteString("Base64解码后字节长度:" + strconv.Itoa(len(_base64Len)) + "\n")

		return
	default:
		return
	}
}
func (e *EncodingInfo) eBytesToBytes(_aa string) []byte {
	_s := _aa
	_s = strings.ReplaceAll(_s, "到文本", "")
	_s = strings.ReplaceAll(_s, "字符", "")
	_s = strings.ReplaceAll(_s, "{", "")
	_s = strings.ReplaceAll(_s, "(", "")
	_s = strings.ReplaceAll(_s, "}", "")
	_s = strings.ReplaceAll(_s, ")", "")
	_s = strings.ReplaceAll(_s, " ", "")
	_s = strings.ReplaceAll(_s, "，", ",")
	_s = strings.ReplaceAll(_s, "+", ",")
	_s = strings.ReplaceAll(_s, "＋", ",")
	_s = strings.ReplaceAll(_s, "\t", "")
	arr := strings.Split(_s, ",")
	_s = ""
	bs := make([]byte, 0)
	for _, k := range arr {
		_b, e1 := strconv.Atoi(k)
		if e1 != nil {
			break
		}
		bs = append(bs, byte(_b))
	}
	return bs
}
func (e *EncodingInfo) ExtendStringConv() {
	if len(e.Children) < 2 {
		return
	}
	switch e.Children[1] {
	case "字节集", "字符":
		if len(e.Children) < 3 {
			return
		}
		data := []byte(e.Data)
		switch e.Children[2] {
		case "文本到字节集":
			_s := ""
			for _, k := range data {
				if _s == "" {
					_s = strconv.Itoa(int(k))
				} else {
					_s += "," + strconv.Itoa(int(k))
				}
			}
			e.output.WriteString("到文本({" + _s + "})")
			return
		case "文本到字符":
			_s := ""
			for _, k := range data {
				if _s == "" {
					_s = "字符(" + strconv.Itoa(int(k)) + ")"
				} else {
					_s += " + " + "字符(" + strconv.Itoa(int(k)) + ")"
				}
			}
			e.output.WriteString(_s)
			return
		case "字节集到文本", "字符到文本":
			e.output.Write(e.eBytesToBytes(e.Data))
			return
		case "字节集到HEX", "字符到HEX":
			e.output.WriteString(strings.ToUpper(hex.EncodeToString(e.eBytesToBytes(e.Data))))
			return
		case "字节集到Base64", "字符到Base64":
			e.output.WriteString(base64.StdEncoding.EncodeToString(e.eBytesToBytes(e.Data)))
			return
		}
		return
	default:
		return
	}
}
func (e *EncodingInfo) ExtendCompress() {
	if len(e.Children) < 5 {
		return
	}
	var data []byte
	if e.Children[3] == "输入为GBK" {
		reader := transform.NewReader(strings.NewReader(e.Data), simplifiedchinese.GBK.NewEncoder())
		gbkStr, _ := io.ReadAll(reader)
		data = gbkStr
	} else if e.Children[3] == "输入为HEX" {
		data, _ = hex.DecodeString(e.Data)
	} else if e.Children[3] == "输入为Base64" {
		data, _ = base64.StdEncoding.DecodeString(e.Data)
	} else {
		data = []byte(e.Data)
	}
	WriteOutput := func(Data []byte) {
		if e.Children[4] == "到HEX" {
			e.output.WriteString(strings.ToUpper(hex.EncodeToString(Data)))
		} else if e.Children[4] == "到文本" {
			e.output.Write(Data)
		} else {
			e.output.WriteString(base64.StdEncoding.EncodeToString(Data))
		}
	}
	operation := e.Children[2] == "压缩"
	call := func(Data []byte) {
		switch e.Children[1] {
		case "GZIP":
			if operation {
				WriteOutput(Compress.GzipCompress(Data))
			} else {
				WriteOutput(Compress.GzipUnCompress(Data))
			}
			return
		case "ZLIB":
			if operation {
				WriteOutput(Compress.ZlibCompress(Data))
			} else {
				WriteOutput(Compress.ZlibUnCompress(Data))
			}
			return
		case "Brotli":
			if operation {
				WriteOutput(Compress.BrCompress(Data))
			} else {
				WriteOutput(Compress.BrUnCompress(Data))
			}
			return
		case "Deflate":
			if operation {
				WriteOutput(Compress.DeflateCompress(Data))
			} else {
				WriteOutput(Compress.DeflateUnCompress(Data))
			}
			return
		case "ZSTD":
			if operation {
				WriteOutput(Compress.ZSTDCompress(Data))
			} else {
				WriteOutput(Compress.ZSTDDecompress(Data))
			}
			return
		}
	}
	call(data)
}
func (e *EncodingInfo) ExtendOther() {
	if len(e.Children) < 2 {
		return
	}
	switch e.Children[1] {
	case "获取当前时间戳":
		now := time.Now()
		timestamp13 := now.UnixNano() / 1000000
		e.output.WriteString(fmt.Sprintf("%d", timestamp13))
		return
	case "参数排序":
		ds := SortQueryParams(e.Data, e.Children[2])
		hash := md5.Sum([]byte(ds))
		md5Str := hex.EncodeToString(hash[:])
		hash1 := sha1.Sum([]byte(ds))
		sha1Str := hex.EncodeToString(hash1[:])
		d := ds + "\n\n排序后的数据哈希值"
		d += "\n\tMD5:\t" + md5Str
		d += "\n\tSHA-1:\t" + sha1Str
		e.output.WriteString(d)
		return
	case "本地文件":
		if len(e.Children) < 3 {
			return
		}
		if e.Children[2] == "取Base64" {
			bs, _ := os.ReadFile(strings.ReplaceAll(strings.TrimSpace(e.Data), "\"", ""))
			e.output.WriteString(base64.StdEncoding.EncodeToString(bs))
			return
		}
		if e.Children[2] == "取HEX" {
			bs, _ := os.ReadFile(e.Data)
			e.output.WriteString(strings.ToUpper(hex.EncodeToString(bs)))
			return
		}
		if e.Children[2] == "DEX重命名" {
			unmap := make(map[string]string)
			// 构建匹配模式
			pattern := filepath.Join(e.Data, "*.dex")

			timestamp13 := time.Now().UnixNano() / int64(time.Millisecond)
			// 将13位时间戳转换为字符串
			timestamp13Str := fmt.Sprintf("%d", timestamp13)
			// 匹配.dex文件
			dexFiles, err := filepath.Glob(pattern)
			if err != nil {
				e.output.WriteString("无法匹配.dex文件:" + err.Error())
				return
			}
			for nm, dexFile := range dexFiles {
				fileName := filepath.Base(dexFile)
				NewFileName := "dex_" + timestamp13Str + "_" + strconv.Itoa(nm) + ".dex"
				unmap[fileName] = strings.ReplaceAll(dexFile, fileName, NewFileName)
				_ = os.Rename(dexFile, unmap[fileName])
			}
			res := ""
			i := 0
			ok := 0
			f := 0
			for _, dexFile := range unmap {
				fileName := filepath.Base(dexFile)
				NewFileName := "classes"
				i++
				if i == 1 {
					NewFileName += ".dex"
				} else {
					NewFileName += strconv.Itoa(i) + ".dex"
				}
				err1 := os.Rename(dexFile, strings.ReplaceAll(dexFile, fileName, NewFileName))
				if err1 != nil {
					f++
				} else {
					ok++
				}
			}
			res = "DEX重命名完成:\n\n成功: " + strconv.Itoa(ok) + " 个\n\n失败: " + strconv.Itoa(f) + " 个"
			e.output.WriteString(res)
			return
		}
		return
	default:
		return
	}
}
func (e *EncodingInfo) Extend() []byte {
	if len(e.Children) < 1 {
		return e.output.Bytes()
	}
	switch e.Children[0] {
	case "文本操作":
		e.ExtendString()
		return e.output.Bytes()
	case "文本转换":
		e.ExtendStringConv()
		return e.output.Bytes()
	case "压缩转换":
		e.ExtendCompress()
		return e.output.Bytes()
	case "其他操作":
		e.ExtendOther()
		return e.output.Bytes()
	default:
		return e.output.Bytes()
	}
}
func (e *EncodingInfo) apply() []byte {
	if len(e.Children) < 1 {
		return e.output.Bytes()
	}
	if len(e.Children) < 2 {
		switch e.Children[0] {
		case "Ansi转USC2":
			ansi := []rune(e.Data)
			_s := ""
			for _, r := range ansi {
				if r >= 0xFF {
					x := strconv.FormatInt(int64(r), 16)
					if len(x) == 1 {
						x = "000" + x
					} else if len(x) == 2 {
						x = "00" + x
					} else if len(x) == 3 {
						x = "0" + x
					}
					_s += "%u" + x
				} else {
					_s += string(r)
				}
			}
			return []byte(_s)
		case "USC2转Ansi":
			arr := strings.Split(e.Data, "\n")
			for n, k := range arr {
				dd := strings.ReplaceAll(k, "%u", "\\u")
				decodedStr, _ := strconv.Unquote(`"` + dd + `"`)
				arr[n] = decodedStr
			}
			return []byte(strings.Join(arr, "\n"))
		case "引号转换":
			return []byte(Session.ConvertELangFormat(e.Data))
		default:
			return nil
		}
		return nil
	}
	switch e.Children[0] {
	case "文件":
		e.ExtendFile()
		return e.output.Bytes()
	case "字符串", "Base64", "HEX":
		e.Base(e.DataBytes())
		return e.output.Bytes()
	default:
		return e.Extend()
	}
}
func (e *Encoding) AppEncodingCall(obj EncodingInfo) []byte {
	return obj.apply()
}

// 提取所有合法范围（如 a-z、Z-A、0-9）
func extractRanges(rule string) []string {
	var ranges []string
	for i := 0; i < len(rule)-2; i++ {
		if rule[i+1] == '-' {
			start := rule[i]
			end := rule[i+2]
			// 确保 start 和 end 都是同类字符
			if (unicode.IsLower(rune(start)) && unicode.IsLower(rune(end))) ||
				(unicode.IsUpper(rune(start)) && unicode.IsUpper(rune(end))) ||
				(unicode.IsDigit(rune(start)) && unicode.IsDigit(rune(end))) {
				ranges = append(ranges, rule[i:i+3])
				i += 2
			}
		}
	}
	return ranges
}

// 构造自定义排序函数
func getSortFunc(rule string) func([]string) {
	priority := make(map[rune]int)
	index := 0

	// 按照规则生成字符优先级表
	ranges := extractRanges(rule)
	for _, r := range ranges {
		start := rune(r[0])
		end := rune(r[2])
		step := 1
		if start > end {
			step = -1
		}
		for ch := start; ch != end+rune(step); ch += rune(step) {
			if _, exists := priority[ch]; !exists {
				priority[ch] = index
				index++
			}
		}
	}

	// 用于排序 key 的排序函数
	return func(keys []string) {
		sort.Slice(keys, func(i, j int) bool {
			a, b := keys[i], keys[j]
			minLen := len(a)
			if len(b) < minLen {
				minLen = len(b)
			}
			for k := 0; k < minLen; k++ {
				ra := rune(a[k])
				rb := rune(b[k])
				wa, oka := priority[ra]
				wb, okb := priority[rb]
				if oka && okb {
					if wa != wb {
						return wa < wb
					}
				} else if oka {
					return true
				} else if okb {
					return false
				} else {
					if ra != rb {
						return ra < rb
					}
				}
			}
			return len(a) < len(b)
		})
	}
}

// 主函数：对 URL 参数按规则排序
func SortQueryParams(raw, rule string) string {
	values, err := url.ParseQuery(raw)
	if err != nil {
		return raw // fallback
	}

	var keys []string
	for k := range values {
		keys = append(keys, k)
	}

	sorter := getSortFunc(rule)
	sorter(keys)

	var sortedParts []string
	for _, k := range keys {
		for _, v := range values[k] {
			sortedParts = append(sortedParts, fmt.Sprintf("%s=%s", k, v))
		}
	}
	return strings.Join(sortedParts, "&")
}
