package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

func GenerateSelfSignedCert(cn string, dnsNames []string, days int, outDir string) (string, string, error) {
	if cn == "" {
		return "", "", fmt.Errorf("CN 不能为空")
	}
	if days <= 0 {
		days = 365
	}
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", err
	}
	serial, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	tmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject:      pkix.Name{CommonName: cn},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(0, 0, days),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     dnsNames,
	}
	if ip := net.ParseIP(cn); ip != nil {
		tmpl.IPAddresses = []net.IP{ip}
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &priv.PublicKey, priv)
	if err != nil {
		return "", "", err
	}
	certPath := outDir + "/" + cn + ".crt"
	keyPath := outDir + "/" + cn + ".key"
	certOut, _ := os.Create(certPath)
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: der})
	certOut.Close()
	keyOut, _ := os.Create(keyPath)
	keyBytes, _ := x509.MarshalECPrivateKey(priv)
	pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes})
	keyOut.Close()
	os.Chmod(keyPath, 0600)
	return certPath, keyPath, nil
}

func GenerateCSR(cn string, dnsNames []string, outDir string) (string, string, error) {
	if cn == "" {
		return "", "", fmt.Errorf("CN 不能为空")
	}
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", err
	}
	tmpl := &x509.CertificateRequest{
		Subject:  pkix.Name{CommonName: cn},
		DNSNames: dnsNames,
	}
	der, err := x509.CreateCertificateRequest(rand.Reader, tmpl, priv)
	if err != nil {
		return "", "", err
	}
	csrPath := outDir + "/" + cn + ".csr"
	keyPath := outDir + "/" + cn + ".csr.key"
	csrOut, _ := os.Create(csrPath)
	pem.Encode(csrOut, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: der})
	csrOut.Close()
	keyOut, _ := os.Create(keyPath)
	keyBytes, _ := x509.MarshalECPrivateKey(priv)
	pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes})
	keyOut.Close()
	os.Chmod(keyPath, 0600)
	return csrPath, keyPath, nil
}

type SiteCertInfo struct {
	Subject   string
	Issuer    string
	NotBefore time.Time
	NotAfter  time.Time
	DNSNames  []string
}

func CheckSiteCert(hostport string) (*SiteCertInfo, error) {
	conn, err := tls.Dial("tcp", hostport, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return nil, fmt.Errorf("无证书")
	}
	c := certs[0]
	return &SiteCertInfo{
		Subject:   c.Subject.CommonName,
		Issuer:    c.Issuer.CommonName,
		NotBefore: c.NotBefore,
		NotAfter:  c.NotAfter,
		DNSNames:  c.DNSNames,
	}, nil
}

func VerifyCertChain(certPath, caPath string) error {
	certData, err := os.ReadFile(certPath)
	if err != nil {
		return err
	}
	caData, err := os.ReadFile(caPath)
	if err != nil {
		return err
	}
	block, _ := pem.Decode(certData)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caData) {
		return fmt.Errorf("CA 证书解析失败")
	}
	_, err = cert.Verify(x509.VerifyOptions{Roots: pool})
	return err
}
