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
		Section(i18n.T("tools.title"))
		fmt.Printf("  [1] %s\n", i18n.T("tools.hash"))
		fmt.Printf("  [2] %s\n", i18n.T("tools.b64e"))
		fmt.Printf("  [3] %s\n", i18n.T("tools.b64d"))
		fmt.Printf("  [4] %s\n", i18n.T("tools.backup"))
		fmt.Printf("  [5] %s\n", i18n.T("tools.restore"))
		fmt.Printf("  [6] %s\n", i18n.T("tools.qr"))
		fmt.Printf("  [0] %s\n", i18n.T("menu.back"))

		switch ReadLine(i18n.T("menu.select")) {
		case "1":
			toolHash()
		case "2":
			toolB64e()
		case "3":
			toolB64d()
		case "4":
			toolBackup()
		case "5":
			toolRestore()
		case "6":
			toolQR()
		case "0":
			return
		default:
			fmt.Println(i18n.T("menu.invalid"))
			Pause()
		}
	}
}

func toolHash() {
	Section(i18n.T("tools.hash"))
	in := ReadLine("文件 / File: ")
	fmt.Println("算法 / Algo:")
	fmt.Println("  [1] md5")
	fmt.Println("  [2] sha1")
	fmt.Println("  [3] sha256")
	fmt.Println("  [4] sha512")
	algo := "sha256"
	switch ReadLine("选择 / Select (3): ") {
	case "1":
		algo = "md5"
	case "2":
		algo = "sha1"
	case "3":
		algo = "sha256"
	case "4":
		algo = "sha512"
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

func toolQR() {
	Section(i18n.T("tools.qr"))
	pub := ReadLine("公钥文件 / Public key: ")
	out := ReadLine("输出 PNG (留空显示文本): ")
	if out == "" {
		content, err := core.ExportPublicKey(pub)
		if err != nil {
			fmt.Println("失败 / Failed:", err)
			Pause()
			return
		}
		s, err := core.GenerateQRCodeString(content)
		if err != nil {
			fmt.Println("失败 / Failed:", err)
			Pause()
			return
		}
		fmt.Println(s)
	} else {
		if err := core.QRCodeFromFile(pub, out); err != nil {
			fmt.Println("失败 / Failed:", err)
			Pause()
			return
		}
		fmt.Println("完成 / Done:", out)
	}
	Pause()
}
