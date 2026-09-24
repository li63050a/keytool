package core

import (
	"fmt"
	"io"
	"os"

	"github.com/ProtonMail/go-crypto/openpgp"
)

func LoadPublicKeyRing(path string) (openpgp.EntityList, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return openpgp.ReadArmoredKeyRing(f)
}

func LoadPrivateKeyRing(path, password string) (openpgp.EntityList, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return openpgp.ReadArmoredKeyRing(f)
}

func GPGEncryptFile(pubKeyPath, inputPath, outPath string) error {
	entities, err := LoadPublicKeyRing(pubKeyPath)
	if err != nil {
		return fmt.Errorf("加载公钥失败: %w", err)
	}

	in, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer out.Close()

	w, err := openpgp.Encrypt(out, entities, nil, nil, nil)
	if err != nil {
		return err
	}
	if _, err := io.Copy(w, in); err != nil {
		return err
	}
	return w.Close()
}

func GPGDecryptFile(privKeyPath, inputPath, outPath string) error {
	entities, err := LoadPrivateKeyRing(privKeyPath, "")
	if err != nil {
		return fmt.Errorf("加载私钥失败: %w", err)
	}

	in, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer in.Close()

	md, err := openpgp.ReadMessage(in, entities, nil, nil)
	if err != nil {
		return fmt.Errorf("读取加密消息失败: %w", err)
	}

	out, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, md.UnverifiedBody); err != nil {
		return err
	}
	return nil
}

func GPGSignFile(privKeyPath, inputPath, outPath string) error {
	entities, err := LoadPrivateKeyRing(privKeyPath, "")
	if err != nil {
		return fmt.Errorf("加载私钥失败: %w", err)
	}
	if len(entities) == 0 {
		return fmt.Errorf("未找到私钥")
	}

	in, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer out.Close()

	err = openpgp.ArmoredDetachSign(out, entities[0], in, nil)
	if err != nil {
		return fmt.Errorf("签名失败: %w", err)
	}
	return nil
}

func GPGVerifyFile(pubKeyPath, sigPath, inputPath string) error {
	entities, err := LoadPublicKeyRing(pubKeyPath)
	if err != nil {
		return fmt.Errorf("加载公钥失败: %w", err)
	}

	sig, err := os.Open(sigPath)
	if err != nil {
		return err
	}
	defer sig.Close()

	in, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer in.Close()

	_, err = openpgp.CheckArmoredDetachedSignature(entities, in, sig, nil)
	if err != nil {
		return fmt.Errorf("签名验证失败: %w", err)
	}
	return nil
}
