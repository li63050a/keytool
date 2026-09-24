// KeyTool 是一个命令行密钥管理工具。
//
// 支持 GPG、SSH、AES 加解密、SSL 证书、哈希、Base64、备份等。
//
// 安装:
//
//	go install github.com/li63050a/keytool@latest
//
// 数据目录:
//
//	~/.li/keytool/
package main

import (
	"fmt"
	"os"

	"github.com/li63050a/keytool/internal/config"
	"github.com/li63050a/keytool/internal/menu"
)

func main() {
	if err := config.EnsureDirs(); err != nil {
		fmt.Fprintln(os.Stderr, "初始化目录失败:", err)
		os.Exit(1)
	}
	menu.Run()
}
