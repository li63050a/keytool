// KeyTool 是一个命令行密钥管理工具。
//
// 支持生成 GPG 和 SSH 密钥、AES 文件加解密、GPG 签名、SSL 证书查看。
//
// 安装:
//
//go install github.com/li63050a/keytool@latest
//
// 使用:
//
//./keytool
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
