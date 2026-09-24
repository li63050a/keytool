package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type CAInfo struct {
	Name      string
	CertPath  string
	KeyPath   string
	IsRoot    bool
	Subject   string
	Issuer    string
	NotBefore time.Time
	NotAfter  time.Time
}

type CAResult struct {
	Name     string
	CertPath string
	KeyPath  string
}

func ListCAs(baseDir string) ([]CAInfo, error) {
	dir := filepath.Join(baseDir, "ca")
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var cas []CAInfo
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".crt") {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".crt")
		certPath := filepath.Join(dir, e.Name())
		keyPath := filepath.Join(dir, name+".key")
		if _, err := os.Stat(keyPath); err != nil {
			continue
		}
		info, err := readCertBrief(certPath)
		if err != nil {
			continue
		}
		cas = append(cas, CAInfo{
			Name:      name,
			CertPath:  certPath,
			KeyPath:   keyPath,
			IsRoot:    info.Subject == info.Issuer,
			Subject:   info.Subject,
			Issuer:    info.Issuer,
			NotBefore: info.NotBefore,
			NotAfter:  info.NotAfter,
		})
	}
	sort.Slice(cas, func(i, j int) bool { return cas[i].Name < cas[j].Name })
	return cas, nil
}

func readCertBrief(path string) (*CertInfo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	rest := data
	for {
		block, r := pem.Decode(rest)
		if block == nil {
			return nil, fmt.Errorf("PEM 解析失败")
		}
		if block.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err == nil {
				return &CertInfo{
					Subject:   cert.Subject.CommonName,
					Issuer:    cert.Issuer.CommonName,
					NotBefore: cert.NotBefore,
					NotAfter:  cert.NotAfter,
					DNSNames:  cert.DNSNames,
					IsCA:      cert.IsCA,
				}, nil
			}
		}
		rest = r
	}
}

func GenerateRootCA(name string, days int, baseDir string) (*CAResult, error) {
	if name == "" {
		return nil, fmt.Errorf("CA 名称不能为空")
	}
	if days <= 0 {
		days = 3650
	}
	dir := filepath.Join(baseDir, "ca")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}

	priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, err
	}
	serial, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	tmpl := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: name, Organization: []string{"KeyTool Root CA"}},
		NotBefore:             time.Now().Add(-time.Minute),
		NotAfter:              time.Now().AddDate(0, 0, days),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            2,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &priv.PublicKey, priv)
	if err != nil {
		return nil, err
	}

	certPath := filepath.Join(dir, name+".crt")
	keyPath := filepath.Join(dir, name+".key")

	if err := writeCert(certPath, der); err != nil {
		return nil, err
	}
	if err := writeECKey(keyPath, priv); err != nil {
		return nil, err
	}
	return &CAResult{Name: name, CertPath: certPath, KeyPath: keyPath}, nil
}

func GenerateIntermediateCA(name string, parent *CAInfo, days int, baseDir string) (*CAResult, error) {
	if name == "" {
		return nil, fmt.Errorf("CA 名称不能为空")
	}
	if parent == nil {
		return nil, fmt.Errorf("必须指定父 CA")
	}
	if days <= 0 {
		days = 1825
	}
	dir := filepath.Join(baseDir, "ca")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}

	parentCert, parentKey, err := loadCAKeyPair(parent)
	if err != nil {
		return nil, fmt.Errorf("加载父 CA 失败: %w", err)
	}

	priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, err
	}
	serial, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	tmpl := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: name, Organization: []string{"KeyTool Intermediate CA"}},
		NotBefore:             time.Now().Add(-time.Minute),
		NotAfter:              time.Now().AddDate(0, 0, days),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            0,
		MaxPathLenZero:        true,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, parentCert, &priv.PublicKey, parentKey)
	if err != nil {
		return nil, err
	}

	certPath := filepath.Join(dir, name+".crt")
	keyPath := filepath.Join(dir, name+".key")

	if err := writeCertChain(certPath, [][]byte{der, parentCert.Raw}); err != nil {
		return nil, err
	}
	if err := writeECKey(keyPath, priv); err != nil {
		return nil, err
	}
	return &CAResult{Name: name, CertPath: certPath, KeyPath: keyPath}, nil
}

func loadCAKeyPair(ca *CAInfo) (*x509.Certificate, *ecdsa.PrivateKey, error) {
	certData, err := os.ReadFile(ca.CertPath)
	if err != nil {
		return nil, nil, err
	}
	var cert *x509.Certificate
	rest := certData
	for {
		block, r := pem.Decode(rest)
		if block == nil {
			break
		}
		if block.Type == "CERTIFICATE" {
			c, err := x509.ParseCertificate(block.Bytes)
			if err == nil {
				cert = c
				break
			}
		}
		rest = r
	}
	if cert == nil {
		return nil, nil, fmt.Errorf("证书解析失败")
	}

	keyData, err := os.ReadFile(ca.KeyPath)
	if err != nil {
		return nil, nil, err
	}
	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, nil, fmt.Errorf("密钥解析失败")
	}

	var key *ecdsa.PrivateKey
	switch block.Type {
	case "EC PRIVATE KEY":
		key, err = x509.ParseECPrivateKey(block.Bytes)
	case "PRIVATE KEY":
		k, err2 := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err2 != nil {
			return nil, nil, err2
		}
		var ok bool
		key, ok = k.(*ecdsa.PrivateKey)
		if !ok {
			return nil, nil, fmt.Errorf("不是 ECDSA 私钥")
		}
	default:
		return nil, nil, fmt.Errorf("未知密钥类型: %s", block.Type)
	}
	if err != nil {
		return nil, nil, err
	}
	return cert, key, nil
}

func writeCert(path string, der []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return pem.Encode(f, &pem.Block{Type: "CERTIFICATE", Bytes: der})
}

func writeCertChain(path string, ders [][]byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, der := range ders {
		if err := pem.Encode(f, &pem.Block{Type: "CERTIFICATE", Bytes: der}); err != nil {
			return err
		}
	}
	return nil
}

func writeECKey(path string, key *ecdsa.PrivateKey) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	der, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return err
	}
	if err := pem.Encode(f, &pem.Block{Type: "EC PRIVATE KEY", Bytes: der}); err != nil {
		return err
	}
	return os.Chmod(path, 0600)
}
