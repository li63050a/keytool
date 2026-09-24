package menu

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/li63050a/keytool/internal/config"
	"github.com/li63050a/keytool/internal/core"
	"github.com/li63050a/keytool/internal/i18n"
)

var reader = bufio.NewReader(os.Stdin)

func Run() {
	i18n.Set(config.LoadLang())
	for {
		showMain()
		choice := readLine(i18n.T("menu.select"))
		switch choice {
		case "1":
			menuGPG()
		case "2":
			menuSSH()
		case "3":
			menuCrypt()
		case "4":
			menuCert()
		case "5":
			handleListKeys()
		case "6":
			menuSettings()
		case "0":
			fmt.Println(i18n.T("menu.bye"))
			return
		default:
			fmt.Println(i18n.T("menu.invalid"))
			pause()
		}
	}
}

func showMain() {
	fmt.Println()
	fmt.Println("============================")
	fmt.Printf("      %s\n", i18n.T("app.title"))
	fmt.Println("============================")
	fmt.Printf("  [1] %s\n", i18n.T("main.gpg"))
	fmt.Printf("  [2] %s\n", i18n.T("main.ssh"))
	fmt.Printf("  [3] %s\n", i18n.T("main.crypt"))
	fmt.Printf("  [4] %s\n", i18n.T("main.cert"))
	fmt.Printf("  [5] %s\n", i18n.T("main.list"))
	fmt.Printf("  [6] %s\n", i18n.T("main.settings"))
	fmt.Printf("  [0] %s\n", i18n.T("main.exit"))
	fmt.Println("============================")
}

func menuGPG() {
	for {
		fmt.Println()
		fmt.Println("---- " + i18n.T("gpg.title") + " ----")
		fmt.Printf("  [1] %s\n", i18n.T("gpg.gen"))
		fmt.Printf("  [2] %s\n", i18n.T("gpg.enc"))
		fmt.Printf("  [3] %s\n", i18n.T("gpg.dec"))
		fmt.Printf("  [4] %s\n", i18n.T("gpg.sign"))
		fmt.Printf("  [5] %s\n", i18n.T("gpg.verify"))
		fmt.Printf("  [0] %s\n", i18n.T("menu.back"))

		switch readLine(i18n.T("menu.select")) {
		case "1":
			handleGPGGen()
		case "2":
			handleGPGEncrypt()
		case "3":
			handleGPGDecrypt()
		case "4":
			handleGPGSign()
		case "5":
			handleGPGVerify()
		case "0":
			return
		default:
			fmt.Println(i18n.T("menu.invalid"))
			pause()
		}
	}
}

func menuSSH() {
	for {
		fmt.Println()
		fmt.Println("---- " + i18n.T("ssh.title") + " ----")
		fmt.Printf("  [1] %s\n", i18n.T("ssh.gen"))
		fmt.Printf("  [2] %s\n", i18n.T("ssh.list"))
		fmt.Printf("  [0] %s\n", i18n.T("menu.back"))

		switch readLine(i18n.T("menu.select")) {
		case "1":
			handleSSHGen()
		case "2":
			handleListSSH()
		case "0":
			return
		default:
			fmt.Println(i18n.T("menu.invalid"))
			pause()
		}
	}
}

func menuCrypt() {
	for {
		fmt.Println()
		fmt.Println("---- " + i18n.T("crypt.title") + " ----")
		fmt.Printf("  [1] %s\n", i18n.T("crypt.enc"))
		fmt.Printf("  [2] %s\n", i18n.T("crypt.dec"))
		fmt.Printf("  [0] %s\n", i18n.T("menu.back"))

		switch readLine(i18n.T("menu.select")) {
		case "1":
			handleAESEncrypt()
		case "2":
			handleAESDecrypt()
		case "0":
			return
		default:
			fmt.Println(i18n.T("menu.invalid"))
			pause()
		}
	}
}

func menuCert() {
	for {
		fmt.Println()
		fmt.Println("---- " + i18n.T("cert.title") + " ----")
		fmt.Printf("  [1] %s\n", i18n.T("cert.view"))
		fmt.Printf("  [0] %s\n", i18n.T("menu.back"))

		switch readLine(i18n.T("menu.select")) {
		case "1":
			handleCertView()
		case "0":
			return
		default:
			fmt.Println(i18n.T("menu.invalid"))
			pause()
		}
	}
}

func menuSettings() {
	for {
		fmt.Println()
		fmt.Println("---- " + i18n.T("settings.title") + " ----")
		fmt.Printf("  %s\n", fmt.Sprintf(i18n.T("settings.lang"), i18n.CurrentName()))
		fmt.Printf("  [1] %s\n", fmt.Sprintf(i18n.T("settings.switch"), i18n.OtherLangName()))
		fmt.Printf("  [0] %s\n", i18n.T("menu.back"))

		switch readLine(i18n.T("menu.select")) {
		case "1":
			newLang := "en"
			if i18n.CurrentCode() == "en" {
				newLang = "zh"
			}
			i18n.Set(newLang)
			config.SaveLang(newLang)
			fmt.Printf("%s\n", fmt.Sprintf(i18n.T("settings.saved"), i18n.CurrentName()))
			pause()
		case "0":
			return
		default:
			fmt.Println(i18n.T("menu.invalid"))
			pause()
		}
	}
}

func handleGPGGen() {
	fmt.Println()
	name := readLine("姓名 / Name: ")
	if name == "" {
		fmt.Println("姓名不能为空 / Name required")
		pause()
		return
	}
	email := readLine("邮箱 / Email: ")
	if email == "" {
		fmt.Println("邮箱不能为空 / Email required")
		pause()
		return
	}

	fmt.Println("生成中... / Generating...")
	res, err := core.GenerateGPGKey(name, email, config.GPGDir())
	if err != nil {
		fmt.Println("失败 / Failed:", err)
		pause()
		return
	}
	fmt.Println()
	fmt.Println("完成 / Done")
	fmt.Println("  公钥 / Public :", res.PublicPath)
	fmt.Println("  私钥 / Private:", res.PrivatePath)
	fmt.Println("  指纹 / Finger :", res.Fingerprint)
	pause()
}

func handleSSHGen() {
	fmt.Println()
	name := readLine("名称 / Name: ")
	if name == "" {
		fmt.Println("名称不能为空 / Name required")
		pause()
		return
	}

	fmt.Println("类型 / Type:")
	fmt.Println("  [1] ed25519 (推荐 / recommended)")
	fmt.Println("  [2] rsa 4096")
	t := readLine("请选择 / Select: ")

	keyType := "ed25519"
	bits := 0
	switch t {
	case "1":
		keyType = "ed25519"
	case "2":
		keyType = "rsa"
		bits = 4096
	default:
		fmt.Println("默认使用 ed25519 / Default ed25519")
	}

	fmt.Println("生成中... / Generating...")
	res, err := core.GenerateSSHKey(name, keyType, bits, config.SSHDir())
	if err != nil {
		fmt.Println("失败 / Failed:", err)
		pause()
		return
	}
	fmt.Println()
	fmt.Println("完成 / Done")
	fmt.Println("  私钥 / Private:", res.PrivatePath)
	fmt.Println("  公钥 / Public :", res.PublicPath)
	fmt.Println("  指纹 / Finger :", res.Fingerprint)
	pause()
}

func handleAESEncrypt() {
	fmt.Println()
	in := readLine("输入文件 / Input: ")
	out := readLine("输出文件 / Output: ")
	pw := readLine("密码 / Password: ")

	if err := core.EncryptFileAES(pw, in, out); err != nil {
		fmt.Println("失败 / Failed:", err)
		pause()
		return
	}
	fmt.Println("完成 / Done:", out)
	pause()
}

func handleAESDecrypt() {
	fmt.Println()
	in := readLine("加密文件 / Encrypted: ")
	out := readLine("输出文件 / Output: ")
	pw := readLine("密码 / Password: ")

	if err := core.DecryptFileAES(pw, in, out); err != nil {
		fmt.Println("失败 / Failed:", err)
		pause()
		return
	}
	fmt.Println("完成 / Done:", out)
	pause()
}

func handleGPGEncrypt() {
	fmt.Println()
	pub := readLine("公钥 / Public Key (.asc): ")
	in := readLine("输入文件 / Input: ")
	out := readLine("输出文件 / Output: ")

	if err := core.GPGEncryptFile(pub, in, out); err != nil {
		fmt.Println("失败 / Failed:", err)
		pause()
		return
	}
	fmt.Println("完成 / Done:", out)
	pause()
}

func handleGPGDecrypt() {
	fmt.Println()
	priv := readLine("私钥 / Private Key (.asc): ")
	in := readLine("加密文件 / Encrypted: ")
	out := readLine("输出文件 / Output: ")

	if err := core.GPGDecryptFile(priv, in, out); err != nil {
		fmt.Println("失败 / Failed:", err)
		pause()
		return
	}
	fmt.Println("完成 / Done:", out)
	pause()
}

func handleGPGSign() {
	fmt.Println()
	priv := readLine("私钥 / Private Key (.asc): ")
	in := readLine("输入文件 / Input: ")
	out := readLine("签名输出 / Signature (.asc): ")

	if err := core.GPGSignFile(priv, in, out); err != nil {
		fmt.Println("失败 / Failed:", err)
		pause()
		return
	}
	fmt.Println("完成 / Done:", out)
	pause()
}

func handleGPGVerify() {
	fmt.Println()
	pub := readLine("公钥 / Public Key (.asc): ")
	in := readLine("原文件 / Original: ")
	sig := readLine("签名文件 / Signature: ")

	if err := core.GPGVerifyFile(pub, in, sig); err != nil {
		fmt.Println("验证失败 / Verify failed:", err)
		pause()
		return
	}
	fmt.Println("验证成功 / Signature OK")
	pause()
}

func handleListKeys() {
	fmt.Println()
	fmt.Println("---- " + i18n.T("list.title") + " ----")

	fmt.Println()
	fmt.Println("【GPG】", config.GPGDir())
	keys, err := core.ListKeys(config.GPGDir())
	if err != nil {
		fmt.Println("  读取失败 / Read failed:", err)
	} else {
		fmt.Print(core.FormatKeyList(keys))
	}

	fmt.Println()
	fmt.Println("【SSH】", config.SSHDir())
	keys, err = core.ListKeys(config.SSHDir())
	if err != nil {
		fmt.Println("  读取失败 / Read failed:", err)
	} else {
		fmt.Print(core.FormatKeyList(keys))
	}

	g, _ := core.CountFiles(config.GPGDir())
	s, _ := core.CountFiles(config.SSHDir())
	fmt.Println()
	fmt.Printf("  统计 / Stats: GPG=%d SSH=%d\n", g, s)
	pause()
}

func handleListSSH() {
	fmt.Println()
	fmt.Println("---- SSH ----")
	keys, err := core.ListKeys(config.SSHDir())
	if err != nil {
		fmt.Println("读取失败 / Read failed:", err)
		pause()
		return
	}
	fmt.Print(core.FormatKeyList(keys))
	pause()
}

func handleCertView() {
	fmt.Println()
	path := readLine("证书文件 / Cert (.pem/.crt): ")

	info, err := core.ParseCertFile(path)
	if err != nil {
		fmt.Println("失败 / Failed:", err)
		pause()
		return
	}

	fmt.Println()
	fmt.Println("主体 / Subject :", info.Subject)
	fmt.Println("颁发 / Issuer  :", info.Issuer)
	fmt.Println("生效 / NotBefore:", info.NotBefore.Format("2006-01-02 15:04:05"))
	fmt.Println("过期 / NotAfter :", info.NotAfter.Format("2006-01-02 15:04:05"))
	fmt.Println("是否CA / IsCA  :", info.IsCA)
	if len(info.DNSNames) > 0 {
		fmt.Println("域名 / DNS     :", strings.Join(info.DNSNames, ", "))
	}
	pause()
}

func readLine(prompt string) string {
	fmt.Print(prompt)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func pause() {
	fmt.Print(i18n.T("menu.pause"))
	reader.ReadString('\n')
}
