package menu

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/li63050a/keytool/internal/config"
	"github.com/li63050a/keytool/internal/core"
)

// SelectFile 列出目录下的文件，方向键选择
func SelectFile(title, dir string, include, exclude []string, defaultIdx int) string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("读取目录失败:", err)
		Pause()
		return ""
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()

		if len(include) > 0 {
			ok := false
			for _, s := range include {
				if strings.Contains(name, s) {
					ok = true
					break
				}
			}
			if !ok {
				continue
			}
		}
		if len(exclude) > 0 {
			skip := false
			for _, s := range exclude {
				if strings.Contains(name, s) {
					skip = true
					break
				}
			}
			if skip {
				continue
			}
		}
		files = append(files, name)
	}

	if len(files) == 0 {
		fmt.Println()
		fmt.Println("目录下没有匹配的文件:", dir)
		Pause()
		return ""
	}

	sort.Strings(files)

	fmt.Println()
	fmt.Println("---- " + title + " ----")
	fmt.Println("目录:", dir)

	idx := SelectWithDefault("请选择 / Select: ", files, defaultIdx)
	if idx < 0 {
		return ""
	}
	return filepath.Join(dir, files[idx])
}

func PickGPGPublic() string {
	return SelectFile("选择 GPG 公钥", config.GPGDir(),
		[]string{"-public.asc"}, nil, 0)
}

func PickGPGPrivate() string {
	return SelectFile("选择 GPG 私钥", config.GPGDir(),
		[]string{"-private.asc"}, nil, 0)
}

func PickGPGAny() string {
	return SelectFile("选择 GPG 密钥", config.GPGDir(),
		[]string{".asc"}, nil, 0)
}

func PickSSHPrivate() string {
	return SelectFile("选择 SSH 私钥", config.SSHDir(),
		nil, []string{".pub"}, 0)
}

func PickSSHPublic() string {
	return SelectFile("选择 SSH 公钥", config.SSHDir(),
		[]string{".pub"}, nil, 0)
}

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
