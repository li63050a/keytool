package menu

import (
	"fmt"

	"github.com/li63050a/keytool/internal/config"
	"github.com/li63050a/keytool/internal/core"
	"github.com/li63050a/keytool/internal/i18n"
)

func MenuList() {
	Section(i18n.T("menu.list"))
	fmt.Println()
	fmt.Println("【GPG】", config.GPGDir())
	gpgKeys, _ := core.ListKeys(config.GPGDir())
	fmt.Print(core.FormatKeyList(gpgKeys))
	fmt.Println()
	fmt.Println("【SSH】", config.SSHDir())
	sshKeys, _ := core.ListKeys(config.SSHDir())
	fmt.Print(core.FormatKeyList(sshKeys))
	fmt.Println()
	fmt.Println("【Cert】", config.CertDir())
	certKeys, _ := core.ListKeys(config.CertDir())
	fmt.Print(core.FormatKeyList(certKeys))
	Pause()
}
