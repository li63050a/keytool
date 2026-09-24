package core

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
)

type SSHKeyInfo struct {
	Path        string
	PublicPath  string
	Type        string
	Fingerprint string
	BitSize     int
}

func InspectSSHKey(privPath string) (*SSHKeyInfo, error) {
	data, err := os.ReadFile(privPath)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(data)
	if err != nil {
		return nil, fmt.Errorf("解析私钥失败: %w", err)
	}
	pub := signer.PublicKey()
	return &SSHKeyInfo{
		Path:        privPath,
		PublicPath:  privPath + ".pub",
		Type:        pub.Type(),
		Fingerprint: ssh.FingerprintSHA256(pub),
		BitSize:     pub.(ssh.CryptoPublicKey).CryptoPublicKey().(interface{ Size() int }).Size() * 8,
	}, nil
}

func DeleteSSHKey(privPath string) error {
	if err := os.Remove(privPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	os.Remove(privPath + ".pub")
	return nil
}

func CopyPubKeyToServer(pubPath, user, host, port string) error {
	pubData, err := os.ReadFile(pubPath)
	if err != nil {
		return err
	}
	if port == "" {
		port = "22"
	}
	remote := fmt.Sprintf("%s@%s", user, host)
	cmd := exec.Command("ssh", "-p", port, remote,
		"mkdir -p ~/.ssh && chmod 700 ~/.ssh && cat >> ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys")
	cmd.Stdin = strings.NewReader(string(pubData))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func GenerateSSHConfig(name, host, user, port, keyPath string) (string, error) {
	if port == "" {
		port = "22"
	}
	abs, _ := filepath.Abs(keyPath)
	config := fmt.Sprintf(`
Host %s
    HostName %s
    User %s
    Port %s
    IdentityFile %s
    IdentitiesOnly yes
`, name, host, user, port, abs)
	return config, nil
}
