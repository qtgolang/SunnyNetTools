package MapHash

import (
	"bytes"
	"github.com/andybalholm/brotli"
	"io"
	"io/ioutil"
)

// BrUnCompress br解压缩
func BrUnCompress(bin []byte) []byte {
	if len(bin) < 1 {
		return make([]byte, 0)
	}
	b, _ := io.ReadAll(brotli.NewReader(ioutil.NopCloser(bytes.NewBuffer(bin))))
	return b
}

// BrCompress br压缩
func BrCompress(bin []byte) []byte {
	// 创建一个Brotli Writer
	var compressed bytes.Buffer
	writer := brotli.NewWriter(&compressed)
	// 将数据写入压缩器
	_, err := writer.Write(bin)
	if err != nil {
		return make([]byte, 0)
	}
	// 关闭压缩器以确保所有数据都被写入
	err = writer.Close()
	if err != nil {
		return make([]byte, 0)
	}
	return compressed.Bytes()
}
