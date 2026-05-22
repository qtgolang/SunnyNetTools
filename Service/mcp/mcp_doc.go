package mcp

import (
	_ "embed"
	"net/http"
)

//go:embed doc.html
var doc []byte

// handleDoc 返回一个简单的 HTML 文档页，方便快速查看接口和工具说明
func handleDoc(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "GET only", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")

	_, _ = w.Write(doc)
}
