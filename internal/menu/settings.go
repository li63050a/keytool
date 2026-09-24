package menu

import (
	"fmt"

	"github.com/li63050a/keytool/internal/config"
	"github.com/li63050a/keytool/internal/i18n"
)

func MenuSettings() {
	for {
		fmt.Println()
		fmt.Println("---- " + i18n.T("settings.title") + " ----")
		fmt.Printf("  %s\n", fmt.Sprintf(i18n.T("settings.lang"), i18n.CurrentName()))

		opts := []string{
			fmt.Sprintf(i18n.T("settings.switch"), i18n.OtherLangName()),
			i18n.T("menu.back"),
		}
		switch Select("请选择 / Select: ", opts) {
		case 0:
			newLang := "en"
			if i18n.CurrentCode() == "en" {
				newLang = "zh"
			}
			i18n.Set(newLang)
			config.SaveLang(newLang)
			fmt.Println(fmt.Sprintf(i18n.T("settings.saved"), i18n.CurrentName()))
			Pause()
		case 1, -1:
			return
		}
	}
}
