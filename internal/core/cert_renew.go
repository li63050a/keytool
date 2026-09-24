package core

import (
	"crypto/elliptic"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

type RenewOptions struct {
	Days int
}

// RenewCert 用父 CA 续期一张叶子证书
// certPath 是旧证书，keyPath 是旧私钥，parent 是签发 CA
// 如果 reuseKey 为 true，复用旧私钥；否则生成新私钥
func RenewCert(certPath, keyPath string, parent *CAInfo, opts RenewOptions, reuseKey bool, outDir string) (*LeafResult, error) {
	if opts.Days <= 0 {
		opts.Days = 365
	}

	oldData, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	oldCert, err := parseFirstCert(oldData)
	if err != nil {
		return nil, fmt.Errorf("解析旧证书失败: %w", err)
	}

	var priv *ecdsa.PrivateKey
	if reuseKey {
		k, err := loadECKey(keyPath)
		if err != nil {
			return nil, fmt.Errorf("加载旧私钥失败: %w", err)
		}
		priv = k
	} else {
		k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, err
		}
		priv = k
	}

	parentCert, parentKey, err := loadCAKeyPair(parent)
	if err != nil {
		return nil, fmt.Errorf("加载父 CA 失败: %w", err)
	}

	serial, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	tmpl := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               oldCert.Subject,
		NotBefore:             time.Now().Add(-time.Minute),
		NotAfter:              time.Now().AddDate(0, 0, opts.Days),
		KeyUsage:              oldCert.KeyUsage,
		ExtKeyUsage:           oldCert.ExtKeyUsage,
		DNSNames:              oldCert.DNSNames,
		IPAddresses:           oldCert.IPAddresses,
		BasicConstraintsValid: true,
	}
	if len(tmpl.DNSNames) == 0 && len(tmpl.IPAddresses) == 0 {
		if ip := net.ParseIP(oldCert.Subject.CommonName); ip != nil {
			tmpl.IPAddresses = []net.IP{ip}
		} else {
			tmpl.DNSNames = []string{oldCert.Subject.CommonName}
		}
	}

	der, err := x509.CreateCertificate(rand.Reader, tmpl, parentCert, &priv.PublicKey, parentKey)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(outDir, 0700); err != nil {
		return nil, err
	}
	cn := oldCert.Subject.CommonName
	certOut := filepath.Join(outDir, cn+".crt")
	keyOut := filepath.Join(outDir, cn+".key")

	if err := writeCertChain(certOut, [][]byte{der, parentCert.Raw}); err != nil {
		return nil, err
	}
	if err := writeECKey(keyOut, priv); err != nil {
		return nil, err
	}

	return &LeafResult{
		CertPath: certOut,
		KeyPath:  keyOut,
		Subject:  cn,
		Issuer:   parentCert.Subject.CommonName,
	}, nil
}

func parseFirstCert(data []byte) (*x509.Certificate, error) {
	rest := data
	for {
		block, r := pem.Decode(rest)
		if block == nil {
			return nil, fmt.Errorf("没有找到 CERTIFICATE 块")
		}
		if block.Type == "CERTIFICATE" {
			c, err := x509.ParseCertificate(block.Bytes)
			if err == nil {
				return c, nil
			}
		}
		rest = r
	}
}

func loadECKey(path string) (*ecdsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("PEM 解析失败")
	}
	switch block.Type {
	case "EC PRIVATE KEY":
		return x509.ParseECPrivateKey(block.Bytes)
	case "PRIVATE KEY":
		k, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		ec, ok := k.(*ecdsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("不是 ECDSA 密钥")
		}
		return ec, nil
	}
	return nil, fmt.Errorf("未知密钥类型: %s", block.Type)
}

var _ = pkix.Name{}
