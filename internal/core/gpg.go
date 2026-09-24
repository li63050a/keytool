package core

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
)

type GPGResult struct {
	PublicPath  string
	PrivatePath string
	Fingerprint string
}

func GenerateGPGKey(name, email, outDir string) (*GPGResult, error) {
	if name == "" || email == "" {
		return nil, fmt.Errorf("姓名和邮箱不能为空")
	}

	cfg := &packet.Config{
		DefaultHash:   10,
		DefaultCipher: packet.CipherAES256,
	}

	entity, err := openpgp.NewEntity(name, "", email, cfg)
	if err != nil {
		return nil, fmt.Errorf("生成失败: %w", err)
	}

	fp := fmt.Sprintf("%X", entity.PrimaryKey.Fingerprint)
	pubPath := filepath.Join(outDir, email+"-public.asc")
	privPath := filepath.Join(outDir, email+"-private.asc")

	if err := exportEntity(entity, pubPath, true); err != nil {
		return nil, err
	}
	if err := exportEntity(entity, privPath, false); err != nil {
		return nil, err
	}

	return &GPGResult{PublicPath: pubPath, PrivatePath: privPath, Fingerprint: fp}, nil
}

func exportEntity(entity *openpgp.Entity, path string, publicOnly bool) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	armorType := openpgp.PublicKeyType
	if !publicOnly {
		armorType = openpgp.PrivateKeyType
	}

	w, err := armor.Encode(f, armorType, map[string]string{"Version": "KeyTool"})
	if err != nil {
		return err
	}

	var werr error
	if publicOnly {
		werr = entity.Serialize(w)
	} else {
		werr = entity.SerializePrivate(w, nil)
	}
	if werr != nil {
		return werr
	}
	return w.Close()
}

var _ io.WriteCloser
