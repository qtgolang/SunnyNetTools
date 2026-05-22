package main

import (
	"changeme/Service"
	"embed"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

//go:embed all:frontend/dist
var assets embed.FS

func init() {
	if runtime.GOOS == "darwin" {
		if os.Geteuid() != 0 {
			if strings.Contains(os.Args[0], "/Contents/MacOS/") {
				cmd := exec.Command("osascript", "-e", `
						tell application "Terminal"
    						activate
    						do script "echo -e '\\033]0; SunnyNet Pro \\007' && \
        						clear && \
        						echo '' && \
        						echo '------------------------------------' && \
        						echo -e '\\033[31m        需要超级用户权限以继续\\033[0m' && \
        						echo '------------------------------------' && \
        						echo '' && \
        						echo '由于系统代理配置或 utun 驱动加载，需要授权才能继续操作。' && \
        						echo '' && \
        						echo '请输入您的超级用户密码以授予权限...' && \
        						echo '' && \
        						sudo `+os.Args[0]+` && \
        						exit"
						end tell
						`)
				_ = cmd.Start()
				os.Exit(0)
			} else {
				//终端程序
				fmt.Println("请使用 sudo 启动程序")
				_, _ = fmt.Scanln()
				os.Exit(0)
			}
		}
	}
}
func main() {
	err := Service.CreateMainWindow(assets).Run()
	if err != nil {
		log.Fatal(err)
	}
}
