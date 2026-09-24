package core

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
)

// ExportChain 从叶子证书提取完整链，导出为 PEM 单文件
func ExportChain(certPath, outPath string) (int, error) {
	data, err := os.ReadFile(certPath)
	if err != nil {
		return 0, err
	}

	rest := data
	count := 0
	var blocks []*pem.Block
	for {
		block, r := pem.Decode(rest)
		if block == nil {
			break
		}
		if block.Type == "CERTIFICATE" {
			if _, err := x509.ParseCertificate(block.Bytes); err == nil {
				blocks = append(blocks, block)
				count++
			}
		}
		rest = r
	}
	if count == 0 {
		return 0, fmt.Errorf("没有找到有效证书")
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0700); err != nil {
		return 0, err
	}
	f, err := os.Create(outPath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	for _, b := range blocks {
		if err := pem.Encode(f, b); err != nil {
			return 0, err
		}
	}
	return count, nil
}

// SplitChain 把证书链拆成多个文件
func SplitChain(certPath, outDir, prefix string) (int, error) {
	data, err := os.ReadFile(certPath)
	if err != nil {
		return 0, err
	}
	if err := os.MkdirAll(outDir, 0700); err != nil {
		return 0, err
	}

	rest := data
	count := 0
	for {
		block, r := pem.Decode(rest)
		if block == nil {
			break
		}
		if block.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err == nil {
				count++
				label := fmt.Sprintf("%s-%d-%s.crt", prefix, count, cert.Subject.CommonName)
				out := filepath.Join(outDir, label)
				f, err := os.Create(out)
				if err != nil {
					return count, err
				}
				pem.Encode(f, block)
				f.Close()
			}
		}
		rest = r
	}
	return count, nil
}
