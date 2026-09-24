package core

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ProtonMail/go-crypto/openpgp"
)

type GPGKeyInfo struct {
	Path        string
	Fingerprint string
	Identities  []string
	HasPrivate  bool
	Created     string
}

func ImportGPGKey(srcPath, destDir string) (string, error) {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return "", err
	}
	entities, err := openpgp.ReadArmoredKeyRing(strings.NewReader(string(data)))
	if err != nil {
		return "", fmt.Errorf("解析密钥失败: %w", err)
	}
	if len(entities) == 0 {
		return "", fmt.Errorf("未找到有效密钥")
	}
	email := ""
	if e := entities[0]; e != nil {
		for _, id := range e.Identities {
			email = id.UserId.Email
			break
		}
	}
	if email == "" {
		email = fmt.Sprintf("imported-%d", os.Getpid())
	}
	dest := filepath.Join(destDir, email+"-imported.asc")
	if err := os.WriteFile(dest, data, 0600); err != nil {
		return "", err
	}
	return dest, nil
}

func InspectGPGKey(path string) (*GPGKeyInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	entities, err := openpgp.ReadArmoredKeyRing(f)
	if err != nil {
		return nil, err
	}
	if len(entities) == 0 {
		return nil, fmt.Errorf("空密钥")
	}
	e := entities[0]
	info := &GPGKeyInfo{Path: path, HasPrivate: e.PrivateKey != nil}
	if e.PrimaryKey != nil {
		info.Fingerprint = fmt.Sprintf("%X", e.PrimaryKey.Fingerprint)
		info.Created = e.PrimaryKey.CreationTime.Format("2006-01-02")
	}
	for _, id := range e.Identities {
		info.Identities = append(info.Identities, id.Name+" <"+id.UserId.Email+">")
	}
	return info, nil
}

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
