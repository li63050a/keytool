package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

type LeafOptions struct {
	Name     string
	DNSNames []string
	IPs      []string
	Days     int
	IsServer bool
	IsClient bool
}

type LeafResult struct {
	CertPath string
	KeyPath  string
	Subject  string
	Issuer   string
}

// IssueLeafCert 用父 CA 签发叶子证书。parent 为 nil 时生成自签名证书。
func IssueLeafCert(opts LeafOptions, parent *CAInfo, outDir string) (*LeafResult, error) {
	if opts.Name == "" {
		return nil, fmt.Errorf("CN 不能为空")
	}
	if opts.Days <= 0 {
		opts.Days = 365
	}
	if err := os.MkdirAll(outDir, 0700); err != nil {
		return nil, err
	}

	var parentCert *x509.Certificate
	var parentKey *ecdsa.PrivateKey
	if parent != nil {
		c, k, err := loadCAKeyPair(parent)
		if err != nil {
			return nil, fmt.Errorf("加载父 CA 失败: %w", err)
		}
		parentCert = c
		parentKey = k
	}

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	serial, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))

	var extUsage []x509.ExtKeyUsage
	if opts.IsServer {
		extUsage = append(extUsage, x509.ExtKeyUsageServerAuth)
	}
	if opts.IsClient {
		extUsage = append(extUsage, x509.ExtKeyUsageClientAuth)
	}
	if len(extUsage) == 0 {
		extUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}
	}

	tmpl := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: opts.Name},
		NotBefore:             time.Now().Add(-time.Minute),
		NotAfter:              time.Now().AddDate(0, 0, opts.Days),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           extUsage,
		DNSNames:              opts.DNSNames,
		BasicConstraintsValid: true,
		IsCA:                  false,
	}
	for _, ipStr := range opts.IPs {
		if ip := net.ParseIP(ipStr); ip != nil {
			tmpl.IPAddresses = append(tmpl.IPAddresses, ip)
		}
	}

	var der []byte
	if parentCert == nil {
		der, err = x509.CreateCertificate(rand.Reader, tmpl, tmpl, &priv.PublicKey, priv)
	} else {
		der, err = x509.CreateCertificate(rand.Reader, tmpl, parentCert, &priv.PublicKey, parentKey)
	}
	if err != nil {
		return nil, err
	}

	certPath := filepath.Join(outDir, opts.Name+".crt")
	keyPath := filepath.Join(outDir, opts.Name+".key")

	if parentCert != nil {
		if err := writeCertChain(certPath, [][]byte{der, parentCert.Raw}); err != nil {
			return nil, err
		}
	} else {
		if err := writeCert(certPath, der); err != nil {
			return nil, err
		}
	}
	if err := writeECKey(keyPath, priv); err != nil {
		return nil, err
	}

	issuer := opts.Name
	if parentCert != nil {
		issuer = parentCert.Subject.CommonName
	}

	return &LeafResult{
		CertPath: certPath,
		KeyPath:  keyPath,
		Subject:  opts.Name,
		Issuer:   issuer,
	}, nil
}
