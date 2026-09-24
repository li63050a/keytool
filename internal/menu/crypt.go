package menu

import (
	"fmt"

	"github.com/li63050a/keytool/internal/core"
	"github.com/li63050a/keytool/internal/i18n"
)

func MenuCrypt() {
	for {
		fmt.Println()
		fmt.Println("---- " + i18n.T("crypt.title") + " ----")
		opts := []string{
			i18n.T("crypt.enc"),
			i18n.T("crypt.dec"),
			i18n.T("menu.back"),
		}
		switch Select("请选择 / Select: ", opts) {
		case 0:
			aesEnc()
		case 1:
			aesDec()
		case 2, -1:
			return
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
