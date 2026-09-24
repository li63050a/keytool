package menu

import (
	"fmt"

	"github.com/li63050a/keytool/internal/config"
	"github.com/li63050a/keytool/internal/i18n"
)

func MenuSettings() {
	for {
		Section(i18n.T("settings.title"))
		fmt.Printf("  %s\n", fmt.Sprintf(i18n.T("settings.lang"), i18n.CurrentName()))
		fmt.Printf("  [1] %s\n", fmt.Sprintf(i18n.T("settings.switch"), i18n.OtherLangName()))
		fmt.Printf("  [0] %s\n", i18n.T("menu.back"))

		switch ReadLine(i18n.T("menu.select")) {
		case "1":
			newLang := "en"
			if i18n.CurrentCode() == "en" {
				newLang = "zh"
			}
			i18n.Set(newLang)
			config.SaveLang(newLang)
			fmt.Println(fmt.Sprintf(i18n.T("settings.saved"), i18n.CurrentName()))
			Pause()
		case "0":
			return
		default:
			fmt.Println(i18n.T("menu.invalid"))
			Pause()
		}
	}
}
