package menu

import (
	"fmt"

	"github.com/li63050a/keytool/internal/config"
	"github.com/li63050a/keytool/internal/core"
	"github.com/li63050a/keytool/internal/i18n"
)

func MenuGPG() {
	last := 0
	for {
		fmt.Println()
		fmt.Println("---- " + i18n.T("gpg.title") + " ----")
		opts := []string{
			i18n.T("gpg.gen"),
			i18n.T("gpg.enc"),
			i18n.T("gpg.dec"),
			i18n.T("gpg.sign"),
			i18n.T("gpg.verify"),
			i18n.T("gpg.import"),
			i18n.T("gpg.info"),
			i18n.T("gpg.delete"),
			i18n.T("menu.back"),
		}
		choice := SelectWithDefault("请选择 / Select: ", opts, last)
		if choice < 0 {
			return
		}
		last = choice
		switch choice {
		case 0:
			gpgGen()
		case 1:
			gpgEnc()
		case 2:
			gpgDec()
		case 3:
			gpgSign()
		case 4:
			gpgVerify()
		case 5:
			gpgImport()
		case 6:
			gpgInfo()
		case 7:
			gpgDelete()
		case 8:
			return
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
	g := PickGPGGroup("选择接收方公钥")
	if g == nil {
		return
	}
	if !g.HasPublic() {
		fmt.Println("这组密钥没有公钥，无法加密。")
		Pause()
		return
	}
	in := ReadLine("输入文件 / Input: ")
	out := ReadLine("输出文件 / Output: ")
	if err := core.GPGEncryptFile(g.PublicPath, in, out); err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("完成 / Done:", out)
	Pause()
}

func gpgDec() {
	Section(i18n.T("gpg.dec"))
	g := PickGPGGroup("选择解密私钥")
	if g == nil {
		return
	}
	if !g.HasPrivate() {
		fmt.Println("这组密钥没有私钥，无法解密。")
		Pause()
		return
	}
	in := ReadLine("加密文件 / Encrypted: ")
	out := ReadLine("输出文件 / Output: ")
	if err := core.GPGDecryptFile(g.PrivatePath, in, out); err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("完成 / Done:", out)
	Pause()
}

func gpgSign() {
	Section(i18n.T("gpg.sign"))
	g := PickGPGGroup("选择签名私钥")
	if g == nil {
		return
	}
	if !g.HasPrivate() {
		fmt.Println("这组密钥没有私钥，无法签名。")
		Pause()
		return
	}
	in := ReadLine("输入文件 / Input: ")
	out := ReadLine("签名输出 / Signature (.asc): ")
	if err := core.GPGSignFile(g.PrivatePath, in, out); err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("完成 / Done:", out)
	Pause()
}

func gpgVerify() {
	Section(i18n.T("gpg.verify"))
	g := PickGPGGroup("选择签名方公钥")
	if g == nil {
		return
	}
	if !g.HasPublic() {
		fmt.Println("这组密钥没有公钥，无法验证。")
		Pause()
		return
	}
	in := ReadLine("原文件 / Original: ")
	sig := ReadLine("签名文件 / Signature: ")
	if err := core.GPGVerifyFile(g.PublicPath, in, sig); err != nil {
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
	g := PickGPGGroup(i18n.T("gpg.info"))
	if g == nil {
		return
	}
	var path string
	if g.HasPrivate() {
		path = g.PrivatePath
	} else {
		path = g.PublicPath
	}
	info, err := core.InspectGPGKey(path)
	if err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("名称 / Name       :", g.Name)
	fmt.Println("指纹 / Fingerprint:", info.Fingerprint)
	fmt.Println("创建 / Created    :", info.Created)
	fmt.Println("含私钥 / Private  :", info.HasPrivate)
	fmt.Println("身份 / Identities :")
	for _, id := range info.Identities {
		fmt.Println("  -", id)
	}
	if g.HasPrivate() {
		fmt.Println("私钥 / Private    :", g.PrivatePath)
	}
	if g.HasPublic() {
		fmt.Println("公钥 / Public     :", g.PublicPath)
	}
	Pause()
}

func gpgDelete() {
	Section(i18n.T("gpg.delete"))
	g := PickGPGGroup(i18n.T("gpg.delete"))
	if g == nil {
		return
	}
	confirm := Select("确认删除整组密钥（"+g.Name+"）？", []string{"否 / No", "是 / Yes"})
	if confirm != 1 {
		fmt.Println("已取消 / Cancelled")
		Pause()
		return
	}
	if g.PrivatePath != "" {
		core.DeleteSSHKey(g.PrivatePath)
	}
	if g.PublicPath != "" {
		core.DeleteSSHKey(g.PublicPath)
	}
	fmt.Println("已删除 / Deleted:", g.Name)
	Pause()
}
