package core

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

type SSHResult struct {
	PrivatePath string
	PublicPath  string
	Fingerprint string
}

func GenerateSSHKey(name, keyType string, bits int, outDir string) (*SSHResult, error) {
	if name == "" {
		return nil, fmt.Errorf("名称不能为空")
	}
	if keyType == "" {
		keyType = "ed25519"
	}
	if bits == 0 {
		bits = 4096
	}

	var priv interface{}
	switch keyType {
	case "ed25519":
		_, k, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, err
		}
		priv = k
	case "rsa":
		k, err := rsa.GenerateKey(rand.Reader, bits)
		if err != nil {
			return nil, err
		}
		priv = k
	default:
		return nil, fmt.Errorf("不支持的类型: %s", keyType)
	}

	signer, err := ssh.NewSignerFromKey(priv)
	if err != nil {
		return nil, err
	}

	privPath := filepath.Join(outDir, name+"_"+keyType)
	pubPath := privPath + ".pub"

	block, err := ssh.MarshalPrivateKey(priv, "")
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(privPath, pem.EncodeToMemory(block), 0600); err != nil {
		return nil, err
	}
	if err := os.WriteFile(pubPath, ssh.MarshalAuthorizedKey(signer.PublicKey()), 0644); err != nil {
		return nil, err
	}

	return &SSHResult{
		PrivatePath: privPath,
		PublicPath:  pubPath,
		Fingerprint: ssh.FingerprintSHA256(signer.PublicKey()),
	}, nil
}
