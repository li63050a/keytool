package menu

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/li63050a/keytool/internal/config"
	"github.com/li63050a/keytool/internal/core"
	"github.com/li63050a/keytool/internal/i18n"
)

func MenuTools() {
	for {
		fmt.Println()
		fmt.Println("---- " + i18n.T("tools.title") + " ----")
		opts := []string{
			i18n.T("tools.hash"),
			i18n.T("tools.b64e"),
			i18n.T("tools.b64d"),
			i18n.T("tools.backup"),
			i18n.T("tools.restore"),
			i18n.T("menu.back"),
		}
		switch Select("请选择 / Select: ", opts) {
		case 0:
			toolHash()
		case 1:
			toolB64e()
		case 2:
			toolB64d()
		case 3:
			toolBackup()
		case 4:
			toolRestore()
		case 5, -1:
			return
		}
	}
}

func toolHash() {
	Section(i18n.T("tools.hash"))
	in := ReadLine("文件 / File: ")
	fmt.Println("算法 / Algo:")
	algoOpts := []string{"md5", "sha1", "sha256", "sha512"}
	algo := "sha256"
	switch Select("请选择 / Select: ", algoOpts) {
	case 0:
		algo = "md5"
	case 1:
		algo = "sha1"
	case 2:
		algo = "sha256"
	case 3:
		algo = "sha512"
	default:
		return
	}
	h, err := core.HashFile(in, algo)
	if err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Printf("%s: %s\n", algo, h)
	Pause()
}

func toolB64e() {
	Section(i18n.T("tools.b64e"))
	in := ReadLine("输入文件 / Input: ")
	out := ReadLine("输出文件 / Output: ")
	if err := core.Base64EncodeFile(in, out); err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("完成 / Done:", out)
	Pause()
}

func toolB64d() {
	Section(i18n.T("tools.b64d"))
	in := ReadLine("输入文件 / Input: ")
	out := ReadLine("输出文件 / Output: ")
	if err := core.Base64DecodeFile(in, out); err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("完成 / Done:", out)
	Pause()
}

func toolBackup() {
	Section(i18n.T("tools.backup"))
	src := config.HomeDir()
	dst := filepath.Join(config.BackupDir(), "backup-"+time.Now().Format("20060102-150405")+".tar.gz")
	if err := core.CreateBackup(src, dst); err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("完成 / Done:", dst)
	Pause()
}

func toolRestore() {
	Section(i18n.T("tools.restore"))
	src := ReadLine("备份文件 / Archive: ")
	dst := config.HomeDir()
	if err := core.RestoreBackup(src, dst); err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("完成 / Done:", dst)
	Pause()
}
