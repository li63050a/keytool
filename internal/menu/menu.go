package menu

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/li63050a/keytool/internal/config"
	"github.com/li63050a/keytool/internal/core"
)

var reader = bufio.NewReader(os.Stdin)

func Run() {
	for {
		showMain()
		choice := readLine("请选择 [1-8]: ")
		switch choice {
		case "1":
			handleGPG()
		case "2":
			handleSSH()
		case "3":
			handleStats()
		case "4":
			handleAESEncrypt()
		case "5":
			handleAESDecrypt()
		case "6":
			handleGPGEncrypt()
		case "7":
			handleGPGSign()
		case "8":
			handleCertView()
		case "0":
			fmt.Println("再见。")
			return
		default:
			fmt.Println("无效选项，请重新选择。")
			pause()
		}
	}
}

func showMain() {
	fmt.Println()
	fmt.Println("============================")
	fmt.Println("      KeyTool 密钥管理")
	fmt.Println("============================")
	fmt.Println("  1. 生成 GPG 密钥")
	fmt.Println("  2. 生成 SSH 密钥")
	fmt.Println("  3. 查看密钥统计")
	fmt.Println("  4. 文件加密 (AES-256)")
	fmt.Println("  5. 文件解密 (AES-256)")
	fmt.Println("  6. GPG 加密文件")
	fmt.Println("  7. GPG 签名文件")
	fmt.Println("  8. 查看 SSL 证书")
	fmt.Println("  0. 退出")
	fmt.Println("============================")
}

func handleGPG() {
	fmt.Println()
	fmt.Println("---- 生成 GPG 密钥 ----")
	name := readLine("姓名: ")
	if name == "" {
		fmt.Println("姓名不能为空。")
		pause()
		return
	}
	email := readLine("邮箱: ")
	if email == "" {
		fmt.Println("邮箱不能为空。")
		pause()
		return
	}

	fmt.Println("生成中...")
	res, err := core.GenerateGPGKey(name, email, config.GPGDir())
	if err != nil {
		fmt.Println("失败:", err)
		pause()
		return
	}
	fmt.Println()
	fmt.Println("完成！")
	fmt.Println("  公钥:", res.PublicPath)
	fmt.Println("  私钥:", res.PrivatePath)
	fmt.Println("  指纹:", res.Fingerprint)
	pause()
}

func handleSSH() {
	fmt.Println()
	fmt.Println("---- 生成 SSH 密钥 ----")
	name := readLine("名称（用作文件名）: ")
	if name == "" {
		fmt.Println("名称不能为空。")
		pause()
		return
	}

	fmt.Println("类型:")
	fmt.Println("  1. ed25519（推荐）")
	fmt.Println("  2. rsa")
	t := readLine("请选择 [1-2]: ")

	keyType := "ed25519"
	bits := 0
	switch t {
	case "1":
		keyType = "ed25519"
	case "2":
		keyType = "rsa"
		bits = 4096
	default:
		fmt.Println("无效选择，默认使用 ed25519。")
	}

	fmt.Println("生成中...")
	res, err := core.GenerateSSHKey(name, keyType, bits, config.SSHDir())
	if err != nil {
		fmt.Println("失败:", err)
		pause()
		return
	}
	fmt.Println()
	fmt.Println("完成！")
	fmt.Println("  私钥:", res.PrivatePath)
	fmt.Println("  公钥:", res.PublicPath)
	fmt.Println("  指纹:", res.Fingerprint)
	pause()
}

func handleStats() {
	fmt.Println()
	fmt.Println("---- 密钥统计 ----")

	g, err := core.CountFiles(config.GPGDir())
	if err != nil {
		fmt.Println("读取 GPG 目录失败:", err)
		pause()
		return
	}
	s, err := core.CountFiles(config.SSHDir())
	if err != nil {
		fmt.Println("读取 SSH 目录失败:", err)
		pause()
		return
	}

	fmt.Printf("  GPG 密钥: %d 个\n", g)
	fmt.Printf("  SSH 密钥: %d 个\n", s)
	fmt.Println("  存放目录:", config.KeysDir())
	pause()
}

func handleAESEncrypt() {
	fmt.Println()
	fmt.Println("---- 文件加密 (AES-256) ----")
	in := readLine("输入文件路径: ")
	out := readLine("输出文件路径: ")
	pw := readLine("密码: ")

	if err := core.EncryptFileAES(pw, in, out); err != nil {
		fmt.Println("失败:", err)
		pause()
		return
	}
	fmt.Println("完成！加密文件:", out)
	pause()
}

func handleAESDecrypt() {
	fmt.Println()
	fmt.Println("---- 文件解密 (AES-256) ----")
	in := readLine("输入加密文件路径: ")
	out := readLine("输出文件路径: ")
	pw := readLine("密码: ")

	if err := core.DecryptFileAES(pw, in, out); err != nil {
		fmt.Println("失败:", err)
		pause()
		return
	}
	fmt.Println("完成！解密文件:", out)
	pause()
}

func handleGPGEncrypt() {
	fmt.Println()
	fmt.Println("---- GPG 加密文件 ----")
	pub := readLine("公钥路径 (.asc): ")
	in := readLine("输入文件路径: ")
	out := readLine("输出文件路径: ")

	if err := core.GPGEncryptFile(pub, in, out); err != nil {
		fmt.Println("失败:", err)
		pause()
		return
	}
	fmt.Println("完成！加密文件:", out)
	pause()
}

func handleGPGSign() {
	fmt.Println()
	fmt.Println("---- GPG 签名文件 ----")
	priv := readLine("私钥路径 (.asc): ")
	in := readLine("输入文件路径: ")
	out := readLine("输出签名路径 (.asc): ")

	if err := core.GPGSignFile(priv, in, out); err != nil {
		fmt.Println("失败:", err)
		pause()
		return
	}
	fmt.Println("完成！签名文件:", out)
	pause()
}

func handleCertView() {
	fmt.Println()
	fmt.Println("---- 查看 SSL 证书 ----")
	path := readLine("证书文件路径 (.pem/.crt): ")

	info, err := core.ParseCertFile(path)
	if err != nil {
		fmt.Println("失败:", err)
		pause()
		return
	}

	fmt.Println()
	fmt.Println("证书信息:")
	fmt.Println("  主体:", info.Subject)
	fmt.Println("  颁发者:", info.Issuer)
	fmt.Println("  生效:", info.NotBefore.Format("2006-01-02 15:04:05"))
	fmt.Println("  过期:", info.NotAfter.Format("2006-01-02 15:04:05"))
	fmt.Println("  是否 CA:", info.IsCA)
	if len(info.DNSNames) > 0 {
		fmt.Println("  域名:", strings.Join(info.DNSNames, ", "))
	}
	pause()
}

func readLine(prompt string) string {
	fmt.Print(prompt)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func pause() {
	fmt.Print("按回车返回主菜单...")
	reader.ReadString('\n')
}
