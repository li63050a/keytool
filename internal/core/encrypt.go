package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

const (
	saltSize   = 16
	keySize    = 32
	iterations = 100000
)

func EncryptFileAES(password, inputPath, outPath string) error {
	if password == "" {
		return fmt.Errorf("密码不能为空")
	}
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	key := pbkdf2.Key([]byte(password), salt, iterations, keySize, sha256.New)
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return err
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	out := append([]byte{}, salt...)
	out = append(out, nonce...)
	out = append(out, ciphertext...)

	return os.WriteFile(outPath, out, 0600)
}

func DecryptFileAES(password, inputPath, outPath string) error {
	if password == "" {
		return fmt.Errorf("密码不能为空")
	}
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}
	if len(data) < saltSize {
		return fmt.Errorf("文件格式错误")
	}

	salt := data[:saltSize]
	rest := data[saltSize:]

	key := pbkdf2.Key([]byte(password), salt, iterations, keySize, sha256.New)
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	if len(rest) < gcm.NonceSize() {
		return fmt.Errorf("文件损坏")
	}

	nonce := rest[:gcm.NonceSize()]
	ciphertext := rest[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("解密失败（密码可能错误）")
	}

	return os.WriteFile(outPath, plaintext, 0644)
}
