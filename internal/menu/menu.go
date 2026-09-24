package menu

import (
	"os"

	"golang.org/x/term"
	"fmt"

	"github.com/li63050a/keytool/internal/config"
	"github.com/li63050a/keytool/internal/i18n"
)

func Run() {
	i18n.Set(config.LoadLang())
	for {
		fmt.Println()
		fmt.Println("============================")
		fmt.Printf("      %s\n", i18n.T("app.title"))
		fmt.Println("============================")

		opts := []string{
			i18n.T("menu.gpg"),
			i18n.T("menu.ssh"),
			i18n.T("menu.crypt"),
			i18n.T("menu.cert"),
			i18n.T("menu.tools"),
			i18n.T("menu.list"),
			i18n.T("menu.settings"),
			i18n.T("menu.exit"),
		}
		choice := Select("请选择 / Select: ", opts)
		switch choice {
		case 0:
			MenuGPG()
		case 1:
			MenuSSH()
		case 2:
			MenuCrypt()
		case 3:
			MenuCert()
		case 4:
			MenuTools()
		case 5:
			MenuList()
		case 6:
			MenuSettings()
		case 7, -1:
			fmt.Println(i18n.T("menu.bye"))
			return
		}
	}
}

func Pause() {
	fmt.Print(i18n.T("menu.pause"))
	ReadLineRaw()
}

func Section(title string) {
	fmt.Println()
	fmt.Println("---- " + title + " ----")
}

// ReadLineRaw 只等待回车，不接受输入
func ReadLineRaw() {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return
	}
	defer term.Restore(fd, oldState)
	buf := make([]byte, 1)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			return
		}
		if buf[0] == '\r' || buf[0] == '\n' || buf[0] == 'q' || buf[0] == 'Q' {
			fmt.Print("\r\n")
			return
		}
	}
}
