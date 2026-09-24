package menu

import (
	"fmt"

	"github.com/li63050a/keytool/internal/config"
	"github.com/li63050a/keytool/internal/core"
	"github.com/li63050a/keytool/internal/i18n"
)

func MenuSSH() {
	for {
		Section(i18n.T("ssh.title"))
		fmt.Printf("  [1] %s\n", i18n.T("ssh.gen"))
		fmt.Printf("  [2] %s\n", i18n.T("ssh.list"))
		fmt.Printf("  [3] %s\n", i18n.T("ssh.info"))
		fmt.Printf("  [4] %s\n", i18n.T("ssh.delete"))
		fmt.Printf("  [5] %s\n", i18n.T("ssh.copy"))
		fmt.Printf("  [6] %s\n", i18n.T("ssh.config"))
		fmt.Printf("  [0] %s\n", i18n.T("menu.back"))

		switch ReadLine(i18n.T("menu.select")) {
		case "1":
			sshGen()
		case "2":
			sshList()
		case "3":
			sshInfo()
		case "4":
			sshDelete()
		case "5":
			sshCopy()
		case "6":
			sshConfig()
		case "0":
			return
		default:
			fmt.Println(i18n.T("menu.invalid"))
			Pause()
		}
	}
}

func sshGen() {
	Section(i18n.T("ssh.gen"))
	name := ReadLine("名称 / Name: ")
	fmt.Println("类型 / Type:")
	fmt.Println("  [1] ed25519 (推荐)")
	fmt.Println("  [2] rsa 4096")
	t := ReadLine("请选择 / Select: ")
	keyType := "ed25519"
	bits := 0
	if t == "2" {
		keyType = "rsa"
		bits = 4096
	}
	fmt.Println("生成中... / Generating...")
	res, err := core.GenerateSSHKey(name, keyType, bits, config.SSHDir())
	if err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("完成 / Done")
	fmt.Println("  私钥 / Private:", res.PrivatePath)
	fmt.Println("  公钥 / Public :", res.PublicPath)
	fmt.Println("  指纹 / Finger :", res.Fingerprint)
	Pause()
}

func sshList() {
	Section(i18n.T("ssh.list"))
	keys, err := core.ListKeys(config.SSHDir())
	if err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Print(core.FormatKeyList(keys))
	Pause()
}

func sshInfo() {
	Section(i18n.T("ssh.info"))
	path := ReadLine("私钥路径 / Private key: ")
	info, err := core.InspectSSHKey(path)
	if err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("类型 / Type       :", info.Type)
	fmt.Println("指纹 / Fingerprint:", info.Fingerprint)
	fmt.Println("位数 / BitSize    :", info.BitSize)
	fmt.Println("公钥 / Public     :", info.PublicPath)
	Pause()
}

func sshDelete() {
	Section(i18n.T("ssh.delete"))
	path := ReadLine("私钥路径 / Private key: ")
	confirm := ReadLine("确认删除？输入 yes: ")
	if confirm != "yes" {
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

func sshCopy() {
	Section(i18n.T("ssh.copy"))
	pub := ReadLine("公钥路径 / Public key: ")
	user := ReadLine("用户名 / User: ")
	host := ReadLine("主机 / Host: ")
	port := ReadLine("端口 / Port (默认22): ")
	if port == "" {
		port = "22"
	}
	if err := core.CopyPubKeyToServer(pub, user, host, port); err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("完成 / Done")
	Pause()
}

func sshConfig() {
	Section(i18n.T("ssh.config"))
	name := ReadLine("别名 / Alias: ")
	host := ReadLine("主机 / Host: ")
	user := ReadLine("用户 / User: ")
	port := ReadLine("端口 / Port (22): ")
	if port == "" {
		port = "22"
	}
	key := ReadLine("私钥路径 / Key: ")
	cfg, err := core.GenerateSSHConfig(name, host, user, port, key)
	if err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println()
	fmt.Println("---- 复制到 ~/.ssh/config ----")
	fmt.Println(cfg)
	fmt.Println("---- 结束 ----")
	Pause()
}
