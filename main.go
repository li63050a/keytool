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
