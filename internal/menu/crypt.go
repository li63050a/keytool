package menu

import (
	"fmt"

	"github.com/li63050a/keytool/internal/core"
	"github.com/li63050a/keytool/internal/i18n"
)

func MenuCrypt() {
	for {
		Section(i18n.T("crypt.title"))
		fmt.Printf("  [1] %s\n", i18n.T("crypt.enc"))
		fmt.Printf("  [2] %s\n", i18n.T("crypt.dec"))
		fmt.Printf("  [0] %s\n", i18n.T("menu.back"))

		switch ReadLine(i18n.T("menu.select")) {
		case "1":
			aesEnc()
		case "2":
			aesDec()
		case "0":
			return
		default:
			fmt.Println(i18n.T("menu.invalid"))
			Pause()
		}
	}
}

func aesEnc() {
	Section(i18n.T("crypt.enc"))
	in := ReadLine("输入 / Input: ")
	out := ReadLine("输出 / Output: ")
	pw := ReadLine("密码 / Password: ")
	if err := core.EncryptFileAES(pw, in, out); err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("完成 / Done:", out)
	Pause()
}

func aesDec() {
	Section(i18n.T("crypt.dec"))
	in := ReadLine("加密文件 / Encrypted: ")
	out := ReadLine("输出 / Output: ")
	pw := ReadLine("密码 / Password: ")
	if err := core.DecryptFileAES(pw, in, out); err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("完成 / Done:", out)
	Pause()
}
