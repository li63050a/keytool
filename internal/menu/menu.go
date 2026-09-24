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
		choice := readLine("请选择 [1-4]: ")
		switch choice {
		case "1":
			handleGPG()
		case "2":
			handleSSH()
		case "3":
			handleStats()
		case "4":
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
	fmt.Println("  4. 退出")
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

func readLine(prompt string) string {
	fmt.Print(prompt)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func pause() {
	fmt.Print("按回车返回主菜单...")
	reader.ReadString('\n')
}
