package menu

import (
	"fmt"
	"strings"

	"github.com/li63050a/keytool/internal/config"
	"github.com/li63050a/keytool/internal/core"
	"github.com/li63050a/keytool/internal/i18n"
)

func MenuGPG() {
	for {
		Section(i18n.T("gpg.title"))
		fmt.Printf("  [1] %s\n", i18n.T("gpg.gen"))
		fmt.Printf("  [2] %s\n", i18n.T("gpg.enc"))
		fmt.Printf("  [3] %s\n", i18n.T("gpg.dec"))
		fmt.Printf("  [4] %s\n", i18n.T("gpg.sign"))
		fmt.Printf("  [5] %s\n", i18n.T("gpg.verify"))
		fmt.Printf("  [6] %s\n", i18n.T("gpg.import"))
		fmt.Printf("  [7] %s\n", i18n.T("gpg.info"))
		fmt.Printf("  [8] %s\n", i18n.T("gpg.delete"))
		fmt.Printf("  [0] %s\n", i18n.T("menu.back"))

		switch ReadLine(i18n.T("menu.select")) {
		case "1":
			gpgGen()
		case "2":
			gpgEnc()
		case "3":
			gpgDec()
		case "4":
			gpgSign()
		case "5":
			gpgVerify()
		case "6":
			gpgImport()
		case "7":
			gpgInfo()
		case "8":
			gpgDelete()
		case "0":
			return
		default:
			fmt.Println(i18n.T("menu.invalid"))
			Pause()
		}
	}
}

func gpgGen() {
	Section(i18n.T("gpg.gen"))
	name := ReadLine("姓名 / Name: ")
	email := ReadLine("邮箱 / Email: ")
	if name == "" || email == "" {
		fmt.Println("不能为空 / Required")
		Pause()
		return
	}
	fmt.Println("生成中... / Generating...")
	res, err := core.GenerateGPGKey(name, email, config.GPGDir())
	if err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("完成 / Done")
	fmt.Println("  公钥 / Public :", res.PublicPath)
	fmt.Println("  私钥 / Private:", res.PrivatePath)
	fmt.Println("  指纹 / Finger :", res.Fingerprint)
	Pause()
}

func gpgEnc() {
	Section(i18n.T("gpg.enc"))
	pub := ReadLine("公钥 / Public (.asc): ")
	in := ReadLine("输入 / Input: ")
	out := ReadLine("输出 / Output: ")
	if err := core.GPGEncryptFile(pub, in, out); err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("完成 / Done:", out)
	Pause()
}

func gpgDec() {
	Section(i18n.T("gpg.dec"))
	priv := ReadLine("私钥 / Private (.asc): ")
	in := ReadLine("加密文件 / Encrypted: ")
	out := ReadLine("输出 / Output: ")
	if err := core.GPGDecryptFile(priv, in, out); err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("完成 / Done:", out)
	Pause()
}

func gpgSign() {
	Section(i18n.T("gpg.sign"))
	priv := ReadLine("私钥 / Private (.asc): ")
	in := ReadLine("输入 / Input: ")
	out := ReadLine("签名 / Signature (.asc): ")
	if err := core.GPGSignFile(priv, in, out); err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("完成 / Done:", out)
	Pause()
}

func gpgVerify() {
	Section(i18n.T("gpg.verify"))
	pub := ReadLine("公钥 / Public (.asc): ")
	in := ReadLine("原文件 / Original: ")
	sig := ReadLine("签名 / Signature: ")
	if err := core.GPGVerifyFile(pub, in, sig); err != nil {
		fmt.Println("验证失败 / Verify failed:", err)
		Pause()
		return
	}
	fmt.Println("验证成功 / Signature OK")
	Pause()
}

func gpgImport() {
	Section(i18n.T("gpg.import"))
	src := ReadLine("源文件 / Source (.asc): ")
	if src == "" {
		Pause()
		return
	}
	dest, err := core.ImportGPGKey(src, config.GPGDir())
	if err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("完成 / Done:", dest)
	Pause()
}

func gpgInfo() {
	Section(i18n.T("gpg.info"))
	path := ReadLine("密钥路径 / Key path: ")
	info, err := core.InspectGPGKey(path)
	if err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("指纹 / Fingerprint:", info.Fingerprint)
	fmt.Println("创建 / Created    :", info.Created)
	fmt.Println("含私钥 / Private  :", info.HasPrivate)
	fmt.Println("身份 / Identities :")
	for _, id := range info.Identities {
		fmt.Println("  -", id)
	}
	Pause()
}

func gpgDelete() {
	Section(i18n.T("gpg.delete"))
	path := ReadLine("密钥路径 / Key path: ")
	confirm := ReadLine("确认删除？输入 yes: ")
	if strings.ToLower(confirm) != "yes" {
		fmt.Println("已取消 / Cancelled")
		Pause()
		return
	}
	if err := core.DeleteSSHKey(path); err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("已删除 / Deleted")
	Pause()
}
