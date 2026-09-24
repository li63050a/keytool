package menu

import (
	"fmt"

	"github.com/li63050a/keytool/internal/config"
	"github.com/li63050a/keytool/internal/core"
)

func PickCA(title string) *core.CAInfo {
	cas, err := core.ListCAs(config.HomeDir())
	if err != nil {
		fmt.Println("读取 CA 目录失败:", err)
		Pause()
		return nil
	}
	if len(cas) == 0 {
		fmt.Println()
		fmt.Println("还没有任何 CA，请先生成。")
		Pause()
		return nil
	}

	labels := make([]string, len(cas))
	for i, c := range cas {
		kind := "中间 CA"
		if c.IsRoot {
			kind = "根 CA"
		}
		labels[i] = fmt.Sprintf("[%s] %s → %s", kind, c.Subject, c.NotAfter.Format("2006-01-02"))
	}

	fmt.Println()
	fmt.Println("---- " + title + " ----")
	fmt.Println("目录:", config.CADir())
	idx := Select("请选择 CA / Select CA: ", labels)
	if idx < 0 {
		return nil
	}
	return &cas[idx]
}
