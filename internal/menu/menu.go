package menu

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/li63050a/keytool/internal/config"
	"github.com/li63050a/keytool/internal/i18n"
)

var reader = bufio.NewReader(os.Stdin)

func Run() {
	i18n.Set(config.LoadLang())
	for {
		showMain()
		switch ReadLine(i18n.T("menu.select")) {
		case "1":
			MenuGPG()
		case "2":
			MenuSSH()
		case "3":
			MenuCrypt()
		case "4":
			MenuCert()
		case "5":
			MenuTools()
		case "6":
			MenuList()
		case "7":
			MenuSettings()
		case "0":
			fmt.Println(i18n.T("menu.bye"))
			return
		default:
			fmt.Println(i18n.T("menu.invalid"))
			Pause()
		}
	}
}

func showMain() {
	fmt.Println()
	fmt.Println("============================")
	fmt.Printf("      %s\n", i18n.T("app.title"))
	fmt.Println("============================")
	fmt.Printf("  [1] %s\n", i18n.T("menu.gpg"))
	fmt.Printf("  [2] %s\n", i18n.T("menu.ssh"))
	fmt.Printf("  [3] %s\n", i18n.T("menu.crypt"))
	fmt.Printf("  [4] %s\n", i18n.T("menu.cert"))
	fmt.Printf("  [5] %s\n", i18n.T("menu.tools"))
	fmt.Printf("  [6] %s\n", i18n.T("menu.list"))
	fmt.Printf("  [7] %s\n", i18n.T("menu.settings"))
	fmt.Printf("  [0] %s\n", i18n.T("menu.exit"))
	fmt.Println("============================")
}

func ReadLine(prompt string) string {
	fmt.Print(prompt)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func Pause() {
	fmt.Print(i18n.T("menu.pause"))
	reader.ReadString('\n')
}

func Section(title string) {
	fmt.Println()
	fmt.Println("---- " + title + " ----")
}
