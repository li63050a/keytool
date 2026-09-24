package menu

import (
	"fmt"

	"github.com/li63050a/keytool/internal/config"
	"github.com/li63050a/keytool/internal/core"
	"github.com/li63050a/keytool/internal/i18n"
)

func MenuSSH() {
	last := 0
	for {
		fmt.Println()
		fmt.Println("---- " + i18n.T("ssh.title") + " ----")
		opts := []string{
			i18n.T("ssh.gen"),
			i18n.T("ssh.list"),
			i18n.T("ssh.info"),
			i18n.T("ssh.delete"),
			i18n.T("ssh.copy"),
			i18n.T("ssh.config"),
			i18n.T("menu.back"),
		}
		choice := SelectWithDefault("请选择 / Select: ", opts, last)
		if choice < 0 {
			return
		}
		last = choice
		switch choice {
		case 0:
			sshGen()
		case 1:
			sshList()
		case 2:
			sshInfo()
		case 3:
			sshDelete()
		case 4:
			sshCopy()
		case 5:
			sshConfig()
		case 6:
			return
		}
	}
}

func sshGen() {
	Section(i18n.T("ssh.gen"))
	name := ReadLine("名称 / Name: ")
	fmt.Println("类型 / Type:")
	t := Select("请选择 / Select: ", []string{"ed25519 (推荐)", "rsa 4096"})
	keyType := "ed25519"
	bits := 0
	if t == 1 {
		keyType = "rsa"
		bits = 4096
	} else if t < 0 {
		return
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
	groups, err := core.ListSSHGroups(config.SSHDir())
	if err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	if len(groups) == 0 {
		fmt.Println("  （还没有 SSH 密钥）")
		Pause()
		return
	}
	for i, g := range groups {
		fmt.Printf("  %d. %s\n", i+1, g.Display())
	}
	Pause()
}

func sshInfo() {
	Section(i18n.T("ssh.info"))
	g := PickSSHGroup(i18n.T("ssh.info"))
	if g == nil {
		return
	}
	if !g.HasPrivate() {
		fmt.Println("这组密钥没有私钥，无法读取详情。")
		Pause()
		return
	}
	info, err := core.InspectSSHKey(g.PrivatePath)
	if err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("名称 / Name       :", g.Name)
	fmt.Println("类型 / Type       :", info.Type)
	fmt.Println("指纹 / Fingerprint:", info.Fingerprint)
	fmt.Println("位数 / BitSize    :", info.BitSize)
	fmt.Println("私钥 / Private    :", g.PrivatePath)
	fmt.Println("公钥 / Public     :", g.PublicPath)
	Pause()
}

func sshDelete() {
	Section(i18n.T("ssh.delete"))
	g := PickSSHGroup(i18n.T("ssh.delete"))
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
		if err := core.DeleteSSHKey(g.PrivatePath); err != nil {
			fmt.Println("删除私钥失败:", err)
		}
	}
	fmt.Println("已删除 / Deleted:", g.Name)
	Pause()
}

func sshCopy() {
	Section(i18n.T("ssh.copy"))
	g := PickSSHGroup(i18n.T("ssh.copy"))
	if g == nil {
		return
	}
	if !g.HasPublic() {
		fmt.Println("这组密钥没有公钥文件，无法复制。")
		Pause()
		return
	}
	fmt.Println("公钥:", g.PublicPath)
	user := ReadLine("用户名 / User: ")
	host := ReadLine("主机 / Host: ")
	port := ReadLine("端口 / Port (默认22): ")
	if port == "" {
		port = "22"
	}
	if err := core.CopyPubKeyToServer(g.PublicPath, user, host, port); err != nil {
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
	g := PickSSHGroup("选择 SSH 密钥")
	if g == nil {
		return
	}
	if !g.HasPrivate() {
		fmt.Println("这组密钥没有私钥，无法用于 config。")
		Pause()
		return
	}
	cfg, err := core.GenerateSSHConfig(name, host, user, port, g.PrivatePath)
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
