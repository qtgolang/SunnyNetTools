package Service

import (
	"changeme/Service/Session"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func bridgeRecordsImport(app *AppMain, m map[string]any) (any, error) {
	filePath := strings.TrimSpace(argString(m, "filePath"))
	if filePath == "" {
		return nil, errors.New("filePath 必填（.sy4 记录文件绝对路径）")
	}
	errMsg, rows := app.AppImport(filePath)
	if errMsg != "" {
		return nil, errors.New(errMsg)
	}
	emitMCPRecordsImport(rows)
	return map[string]any{
		"ok":    true,
		"count": len(rows),
		"path":  filePath,
	}, nil
}

func bridgeRecordsExport(app *AppMain, m map[string]any) (any, error) {
	filePath := strings.TrimSpace(argString(m, "filePath"))
	if filePath == "" {
		return nil, errors.New("filePath 必填（.sy4 保存路径）")
	}
	if !strings.HasSuffix(strings.ToLower(filePath), ".sy4") {
		filePath += ".sy4"
	}
	ids, err := argTheologyList(m)
	if err != nil {
		if strings.Contains(err.Error(), "至少填一项") {
			ids = collectAllStoredTheologies()
		} else {
			return nil, err
		}
	}
	if len(ids) == 0 {
		return nil, errors.New("没有可导出的会话")
	}
	if msg := app.AppExport(ids, filePath); msg != "" {
		return nil, errors.New(msg)
	}
	return map[string]any{
		"ok":    true,
		"count": len(ids),
		"path":  filePath,
	}, nil
}

func bridgeSessionPackExport(app *AppMain, m map[string]any) (any, error) {
	_ = app
	filePath := strings.TrimSpace(argString(m, "filePath"))
	if filePath == "" {
		return nil, errors.New("filePath 必填（.bin 会话包保存路径）")
	}
	th, err := argTheologyOne(m)
	if err != nil {
		return nil, err
	}
	obj := Session.GetAppSession(th)
	if obj == nil {
		return nil, errors.New("会话不存在")
	}
	if !obj.IsWebsocket() && !obj.IsTCP() {
		return nil, errors.New("会话包导出仅支持 WebSocket 或 TCP 会话")
	}
	tmpPath, err := Session.ExportMessage(obj)
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpPath)
	if err := copyFile(tmpPath, filePath); err != nil {
		return nil, err
	}
	return map[string]any{
		"ok":       true,
		"path":     filePath,
		"theology": th,
		"rowId": strconv.Itoa(th),
	}, nil
}

func collectAllStoredTheologies() []int {
	var all []int
	Session.Session.Range(func(key, value any) bool {
		th, ok := key.(int)
		if !ok {
			return true
		}
		if _, ok := value.(Session.AppSession); ok {
			all = append(all, th)
		}
		return true
	})
	sort.Ints(all)
	return all
}

func emitMCPRecordsImport(rows []Session.Insert) {
	if len(rows) == 0 {
		return
	}
	emitMCPMainJSON("recordsimport", map[string]any{
		"rows":  rows,
		"count": len(rows),
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if dir := filepath.Dir(dst); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
