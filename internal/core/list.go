package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type KeyEntry struct {
	Filename string
	Path     string
	Size     int64
	IsPublic bool
}

func ListKeys(dir string) ([]KeyEntry, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var result []KeyEntry
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		name := e.Name()
		isPub := strings.Contains(name, "public") || strings.HasSuffix(name, ".pub") || strings.HasSuffix(name, ".crt")
		result = append(result, KeyEntry{
			Filename: name,
			Path:     filepath.Join(dir, name),
			Size:     info.Size(),
			IsPublic: isPub,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Filename < result[j].Filename
	})
	return result, nil
}

func FormatKeyList(entries []KeyEntry) string {
	if len(entries) == 0 {
		return "  （空 / Empty）\n"
	}
	var b strings.Builder
	for i, e := range entries {
		kind := "私钥"
		if e.IsPublic {
			kind = "公钥"
		}
		b.WriteString(fmt.Sprintf("  %d. [%s] %s (%d B)\n", i+1, kind, e.Filename, e.Size))
	}
	return b.String()
}

func ExportPublicKey(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func CountFiles(dir string) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}
	n := 0
	for _, e := range entries {
		if !e.IsDir() {
			n++
		}
	}
	return n, nil
}
