package menu

import (
	"fmt"
	"strings"

	"github.com/li63050a/keytool/internal/config"
	"github.com/li63050a/keytool/internal/core"
	"github.com/li63050a/keytool/internal/i18n"
)

func MenuCert() {
	for {
		Section(i18n.T("cert.title"))
		fmt.Printf("  [1] %s\n", i18n.T("cert.view"))
		fmt.Printf("  [2] %s\n", i18n.T("cert.selfgen"))
		fmt.Printf("  [3] %s\n", i18n.T("cert.csr"))
		fmt.Printf("  [4] %s\n", i18n.T("cert.check"))
		fmt.Printf("  [5] %s\n", i18n.T("cert.chain"))
		fmt.Printf("  [0] %s\n", i18n.T("menu.back"))

		switch ReadLine(i18n.T("menu.select")) {
		case "1":
			certView()
		case "2":
			certSelfGen()
		case "3":
			certCSR()
		case "4":
			certCheck()
		case "5":
			certChain()
		case "0":
			return
		default:
			fmt.Println(i18n.T("menu.invalid"))
			Pause()
		}
	}
}

func certView() {
	Section(i18n.T("cert.view"))
	path := ReadLine("证书 / Cert: ")
	info, err := core.ParseCertFile(path)
	if err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("主体 / Subject :", info.Subject)
	fmt.Println("颁发 / Issuer  :", info.Issuer)
	fmt.Println("生效 / From    :", info.NotBefore.Format("2006-01-02 15:04:05"))
	fmt.Println("过期 / To      :", info.NotAfter.Format("2006-01-02 15:04:05"))
	fmt.Println("是否CA / IsCA  :", info.IsCA)
	if len(info.DNSNames) > 0 {
		fmt.Println("域名 / DNS     :", strings.Join(info.DNSNames, ", "))
	}
	Pause()
}

func certSelfGen() {
	Section(i18n.T("cert.selfgen"))
	cn := ReadLine("CN (域名/IP): ")
	dns := ReadLine("额外 DNS（逗号分隔，可空）: ")
	days := ReadLine("有效天数 / Days (365): ")
	d := 365
	fmt.Sscanf(days, "%d", &d)
	var dnsNames []string
	if dns != "" {
		dnsNames = strings.Split(dns, ",")
	}
	cert, key, err := core.GenerateSelfSignedCert(cn, dnsNames, d, config.CertDir())
	if err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("完成 / Done")
	fmt.Println("  证书 / Cert:", cert)
	fmt.Println("  私钥 / Key :", key)
	Pause()
}

func certCSR() {
	Section(i18n.T("cert.csr"))
	cn := ReadLine("CN: ")
	dns := ReadLine("额外 DNS（逗号分隔，可空）: ")
	var dnsNames []string
	if dns != "" {
		dnsNames = strings.Split(dns, ",")
	}
	csr, key, err := core.GenerateCSR(cn, dnsNames, config.CsrDir())
	if err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("完成 / Done")
	fmt.Println("  CSR:", csr)
	fmt.Println("  Key:", key)
	Pause()
}

func certCheck() {
	Section(i18n.T("cert.check"))
	host := ReadLine("主机:端口 (如 example.com:443): ")
	if !strings.Contains(host, ":") {
		host = host + ":443"
	}
	info, err := core.CheckSiteCert(host)
	if err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("主体 / Subject :", info.Subject)
	fmt.Println("颁发 / Issuer  :", info.Issuer)
	fmt.Println("生效 / From    :", info.NotBefore.Format("2006-01-02 15:04:05"))
	fmt.Println("过期 / To      :", info.NotAfter.Format("2006-01-02 15:04:05"))
	fmt.Println("域名 / DNS     :", strings.Join(info.DNSNames, ", "))
	Pause()
}

func certChain() {
	Section(i18n.T("cert.chain"))
	cert := ReadLine("证书路径 / Cert: ")
	ca := ReadLine("CA 路径 / CA: ")
	if err := core.VerifyCertChain(cert, ca); err != nil {
		fmt.Println("验证失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("验证通过 / OK")
	Pause()
}
