package core

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type RevokedEntry struct {
	Serial    string
	RevokedAt time.Time
	Reason    string
}

// RevokeCert 用 CA 签出一张 CRL，包含被吊销的证书序列号
// revokedCerts 是待吊销的证书文件路径列表
func RevokeCert(ca *CAInfo, revokedCerts []string, outDir string) (string, error) {
	parentCert, parentKey, err := loadCAKeyPair(ca)
	if err != nil {
		return "", fmt.Errorf("加载 CA 失败: %w", err)
	}

	var revoked []pkix.RevokedCertificate
	for _, path := range revokedCerts {
		data, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("读取 %s 失败: %w", path, err)
		}
		cert, err := parseFirstCert(data)
		if err != nil {
			return "", fmt.Errorf("解析 %s 失败: %w", path, err)
		}
		revoked = append(revoked, pkix.RevokedCertificate{
			SerialNumber:   cert.SerialNumber,
			RevocationTime: time.Now(),
		})
	}

	tmpl := &x509.RevocationList{
		SignatureAlgorithm:  parentCert.SignatureAlgorithm,
		RevokedCertificateEntries: []x509.RevocationListEntry{},
		Number:              big.NewInt(time.Now().Unix()),
		ThisUpdate:          time.Now(),
		NextUpdate:          time.Now().AddDate(0, 0, 30),
	}
	for _, r := range revoked {
		tmpl.RevokedCertificateEntries = append(tmpl.RevokedCertificateEntries, x509.RevocationListEntry{
			SerialNumber:   r.SerialNumber,
			RevocationTime: r.RevocationTime,
		})
	}

	der, err := x509.CreateRevocationList(rand.Reader, tmpl, parentCert, parentKey)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(outDir, 0700); err != nil {
		return "", err
	}
	name := ca.Subject + ".crl"
	outPath := filepath.Join(outDir, name)
	f, err := os.Create(outPath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if err := pem.Encode(f, &pem.Block{Type: "X509 CRL", Bytes: der}); err != nil {
		return "", err
	}
	return outPath, nil
}

// ParseCRL 读取 CRL，返回被吊销证书的序列号列表
func ParseCRL(path string) ([]RevokedEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("PEM 解析失败")
	}
	crl, err := x509.ParseRevocationList(block.Bytes)
	if err != nil {
		return nil, err
	}
	var out []RevokedEntry
	for _, r := range crl.RevokedCertificateEntries {
		out = append(out, RevokedEntry{
			Serial:    r.SerialNumber.String(),
			RevokedAt: r.RevocationTime,
			Reason:    fmt.Sprintf("%d", r.ReasonCode),
		})
	}
	return out, nil
}

// IsCertRevoked 判断某张证书是否在 CRL 中
func IsCertRevoked(certPath, crlPath string) (bool, error) {
	certData, err := os.ReadFile(certPath)
	if err != nil {
		return false, err
	}
	cert, err := parseFirstCert(certData)
	if err != nil {
		return false, err
	}

	entries, err := ParseCRL(crlPath)
	if err != nil {
		return false, err
	}
	target := strings.TrimSpace(cert.SerialNumber.String())
	for _, e := range entries {
		if strings.TrimSpace(e.Serial) == target {
			return true, nil
		}
	}
	return false, nil
}

// CountCRL 统计 CRL 中被吊销证书数量
func CountCRL(path string) (int, error) {
	entries, err := ParseCRL(path)
	if err != nil {
		return 0, err
	}
	return len(entries), nil
}
