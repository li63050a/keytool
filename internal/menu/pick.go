package menu

import (
	"fmt"

	"github.com/li63050a/keytool/internal/config"
	"github.com/li63050a/keytool/internal/core"
)

// PickSSHGroup 选择一组 SSH 密钥（私钥 + 公钥）
func PickSSHGroup(title string) *core.KeyGroup {
	groups, err := core.ListSSHGroups(config.SSHDir())
	if err != nil {
		fmt.Println("读取目录失败:", err)
		Pause()
		return nil
	}
	if len(groups) == 0 {
		fmt.Println()
		fmt.Println("还没有 SSH 密钥，请先生成。")
		Pause()
		return nil
	}

	labels := make([]string, len(groups))
	for i, g := range groups {
		labels[i] = g.Display()
	}

	fmt.Println()
	fmt.Println("---- " + title + " ----")
	fmt.Println("目录:", config.SSHDir())
	idx := Select("请选择 / Select: ", labels)
	if idx < 0 {
		return nil
	}
	return &groups[idx]
}

// PickGPGGroup 选择一组 GPG 密钥（公钥 + 私钥）
func PickGPGGroup(title string) *core.KeyGroup {
	groups, err := core.ListGPGGroups(config.GPGDir())
	if err != nil {
		fmt.Println("读取目录失败:", err)
		Pause()
		return nil
	}
	if len(groups) == 0 {
		fmt.Println()
		fmt.Println("还没有 GPG 密钥，请先生成。")
		Pause()
		return nil
	}

	labels := make([]string, len(groups))
	for i, g := range groups {
		labels[i] = g.Display()
	}

	fmt.Println()
	fmt.Println("---- " + title + " ----")
	fmt.Println("目录:", config.GPGDir())
	idx := Select("请选择 / Select: ", labels)
	if idx < 0 {
		return nil
	}
	return &groups[idx]
}
