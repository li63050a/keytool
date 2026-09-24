package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAll(t *testing.T) {
	tmp := t.TempDir()

	t.Run("GPG生成", func(t *testing.T) {
		res, err := GenerateGPGKey("TestUser", "test@example.com", tmp)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(res.PublicPath); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(res.PrivatePath); err != nil {
			t.Fatal(err)
		}
		t.Logf("公钥: %s", res.PublicPath)
		t.Logf("指纹: %s", res.Fingerprint)
	})

	t.Run("SSH生成", func(t *testing.T) {
		res, err := GenerateSSHKey("testkey", "ed25519", 0, tmp)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(res.PrivatePath); err != nil {
			t.Fatal(err)
		}
		t.Logf("指纹: %s", res.Fingerprint)
	})

	t.Run("AES加解密", func(t *testing.T) {
		plain := filepath.Join(tmp, "plain.txt")
		enc := filepath.Join(tmp, "enc.bin")
		dec := filepath.Join(tmp, "dec.txt")

		os.WriteFile(plain, []byte("hello secret"), 0644)

		if err := EncryptFileAES("123456", plain, enc); err != nil {
			t.Fatal(err)
		}
		if err := DecryptFileAES("123456", enc, dec); err != nil {
			t.Fatal(err)
		}
		data, _ := os.ReadFile(dec)
		if string(data) != "hello secret" {
			t.Fatalf("解密内容不对: %q", string(data))
		}
		t.Log("加解密成功")
	})

	t.Run("GPG加密签名", func(t *testing.T) {
		gpgRes, err := GenerateGPGKey("Signer", "signer@example.com", tmp)
		if err != nil {
			t.Fatal(err)
		}

		plain := filepath.Join(tmp, "msg.txt")
		enc := filepath.Join(tmp, "msg.gpg")
		sig := filepath.Join(tmp, "msg.sig")

		os.WriteFile(plain, []byte("signed message"), 0644)

		if err := GPGEncryptFile(gpgRes.PublicPath, plain, enc); err != nil {
			t.Fatal(err)
		}
		t.Log("GPG加密成功")

		if err := GPGSignFile(gpgRes.PrivatePath, plain, sig); err != nil {
			t.Fatal(err)
		}
		t.Log("GPG签名成功")
	})
}
