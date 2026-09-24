package menu

import (
	"path/filepath"
	"fmt"
	"strings"

	"github.com/li63050a/keytool/internal/config"
	"github.com/li63050a/keytool/internal/core"
	"github.com/li63050a/keytool/internal/i18n"
)

func MenuCert() {
	last := 0
	for {
		fmt.Println()
		fmt.Println("---- " + i18n.T("cert.title") + " ----")
		opts := []string{
			"CA 管理",
			"签发证书（用 CA）",
			"证书续期",
			"证书吊销（生成 CRL）",
			"查看 CRL 内容",
			"导出证书链",
			"生成自签名证书",
			"查看证书信息",
			"生成 CSR",
			"检查网站证书",
			"验证证书链",
			"返回上一级",
		}
		choice := SelectWithDefault("请选择 / Select: ", opts, last)
		if choice < 0 {
			return
		}
		last = choice
		switch choice {
		case 0:
			menuCA()
		case 1:
			menuIssue()
		case 2:
			certRenew()
		case 3:
			certRevoke()
		case 4:
			certCRLView()
		case 5:
			certExportChain()
		case 6:
			certSelfGen()
		case 7:
			certView()
		case 8:
			certCSR()
		case 9:
			certCheck()
		case 10:
			certChain()
		case 11:
			return
		}
	}
}

// ---------- CA 管理 ----------

func menuCA() {
	last := 0
	for {
		fmt.Println()
		fmt.Println("---- CA 管理 ----")
		opts := []string{
			"生成根 CA",
			"生成中间 CA",
			"列出所有 CA",
			"返回上一级",
		}
		choice := SelectWithDefault("请选择 / Select: ", opts, last)
		if choice < 0 {
			return
		}
		last = choice
		switch choice {
		case 0:
			caGenRoot()
		case 1:
			caGenIntermediate()
		case 2:
			caList()
		case 3:
			return
		}
	}
}

func caGenRoot() {
	Section("生成根 CA")
	name := ReadLine("CA 名称（如 My Root CA）: ")
	if name == "" {
		fmt.Println("名称不能为空")
		Pause()
		return
	}
	days := ReadLine("有效天数 / Days (默认3650): ")
	d := 3650
	fmt.Sscanf(days, "%d", &d)
	fmt.Println("生成中... / Generating...")
	res, err := core.GenerateRootCA(name, d, config.HomeDir())
	if err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("完成 / Done")
	fmt.Println("  证书 / Cert:", res.CertPath)
	fmt.Println("  私钥 / Key :", res.KeyPath)
	fmt.Println("  提示：私钥请离线保管。")
	Pause()
}

func caGenIntermediate() {
	Section("生成中间 CA")
	parent := PickCA("选择父 CA")
	if parent == nil {
		return
	}
	name := ReadLine("中间 CA 名称: ")
	if name == "" {
		fmt.Println("名称不能为空")
		Pause()
		return
	}
	days := ReadLine("有效天数 / Days (默认1825): ")
	d := 1825
	fmt.Sscanf(days, "%d", &d)
	fmt.Println("生成中... / Generating...")
	res, err := core.GenerateIntermediateCA(name, parent, d, config.HomeDir())
	if err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("完成 / Done")
	fmt.Println("  证书 / Cert:", res.CertPath)
	fmt.Println("  私钥 / Key :", res.KeyPath)
	fmt.Println("  父 CA     :", parent.Subject)
	Pause()
}

func caList() {
	Section("所有 CA")
	cas, err := core.ListCAs(config.HomeDir())
	if err != nil {
		fmt.Println("读取失败 / Failed:", err)
		Pause()
		return
	}
	if len(cas) == 0 {
		fmt.Println("  （还没有 CA）")
		Pause()
		return
	}
	for i, c := range cas {
		kind := "中间"
		if c.IsRoot {
			kind = "根"
		}
		fmt.Printf("  %d. [%s] %s\n", i+1, kind, c.Subject)
		fmt.Printf("      有效期: %s → %s\n",
			c.NotBefore.Format("2006-01-02"),
			c.NotAfter.Format("2006-01-02"))
		fmt.Printf("      证书: %s\n", c.CertPath)
		fmt.Println()
	}
	Pause()
}

// ---------- 签发证书 ----------

func menuIssue() {
	last := 0
	for {
		fmt.Println()
		fmt.Println("---- 签发证书 ----")
		opts := []string{
			"服务器证书（用 CA 签发）",
			"客户端证书（用 CA 签发）",
			"返回上一级",
		}
		choice := SelectWithDefault("请选择 / Select: ", opts, last)
		if choice < 0 {
			return
		}
		last = choice
		switch choice {
		case 0:
			issueLeaf(true, false, "服务器证书")
		case 1:
			issueLeaf(false, true, "客户端证书")
		case 2:
			return
		}
	}
}

func issueLeaf(isServer, isClient bool, label string) {
	Section("签发" + label)
	parent := PickCA("选择签发用的 CA")
	if parent == nil {
		return
	}
	cn := ReadLine("CN (域名/用户名): ")
	if cn == "" {
		fmt.Println("CN 不能为空")
		Pause()
		return
	}
	dns := ReadLine("额外 DNS（逗号分隔，可空）: ")
	ips := ReadLine("IP 地址（逗号分隔，可空）: ")
	days := ReadLine("有效天数 / Days (默认365): ")
	d := 365
	fmt.Sscanf(days, "%d", &d)

	var dnsNames, ipList []string
	if dns != "" {
		for _, s := range strings.Split(dns, ",") {
			if t := strings.TrimSpace(s); t != "" {
				dnsNames = append(dnsNames, t)
			}
		}
	}
	if ips != "" {
		for _, s := range strings.Split(ips, ",") {
			if t := strings.TrimSpace(s); t != "" {
				ipList = append(ipList, t)
			}
		}
	}

	fmt.Println("签发中... / Signing...")
	res, err := core.IssueLeafCert(core.LeafOptions{
		Name:     cn,
		DNSNames: dnsNames,
		IPs:      ipList,
		Days:     d,
		IsServer: isServer,
		IsClient: isClient,
	}, parent, config.CertDir())
	if err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("完成 / Done")
	fmt.Println("  证书 / Cert  :", res.CertPath)
	fmt.Println("  私钥 / Key   :", res.KeyPath)
	fmt.Println("  主体 / Subject:", res.Subject)
	fmt.Println("  签发 / Issuer :", res.Issuer)
	Pause()
}

// ---------- 续期 ----------

func certRenew() {
	Section("证书续期")
	certPath := SelectFile("选择要续期的证书", config.CertDir(), []string{".crt"}, nil, 0)
	if certPath == "" {
		return
	}
	base := strings.TrimSuffix(filepath.Base(certPath), ".crt")
	keyPath := filepath.Join(filepath.Dir(certPath), base+".key")

	parent := PickCA("选择签发用 CA")
	if parent == nil {
		return
	}
	days := ReadLine("新有效天数 / Days (默认365): ")
	d := 365
	fmt.Sscanf(days, "%d", &d)

	reuse := Select("是否复用旧私钥？", []string{"否 / No（生成新私钥）", "是 / Yes（复用旧私钥）"})
	if reuse < 0 {
		return
	}

	fmt.Println("续期中... / Renewing...")
	res, err := core.RenewCert(certPath, keyPath, parent, core.RenewOptions{Days: d}, reuse == 1, config.CertDir())
	if err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("完成 / Done")
	fmt.Println("  证书 / Cert  :", res.CertPath)
	fmt.Println("  私钥 / Key   :", res.KeyPath)
	fmt.Println("  主体 / Subject:", res.Subject)
	fmt.Println("  签发 / Issuer :", res.Issuer)
	Pause()
}

// ---------- 吊销 ----------

func certRevoke() {
	Section("证书吊销（生成 CRL）")
	ca := PickCA("选择签发 CRL 的 CA")
	if ca == nil {
		return
	}

	var revoked []string
	for {
		c := SelectFile("选择要吊销的证书（可多次）", config.CertDir(), []string{".crt"}, nil, 0)
		if c == "" {
			break
		}
		revoked = append(revoked, c)
		more := Select("继续添加？", []string{"否 / No", "是 / Yes"})
		if more != 1 {
			break
		}
	}
	if len(revoked) == 0 {
		fmt.Println("没有选择任何证书")
		Pause()
		return
	}

	fmt.Println("生成中... / Generating...")
	out, err := core.RevokeCert(ca, revoked, config.CADir())
	if err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("完成 / Done")
	fmt.Println("  CRL:", out)
	fmt.Printf("  已吊销 %d 张证书\n", len(revoked))
	Pause()
}

func certCRLView() {
	Section("查看 CRL 内容")
	path := SelectFile("选择 CRL 文件", config.CADir(), []string{".crl"}, nil, 0)
	if path == "" {
		return
	}
	entries, err := core.ParseCRL(path)
	if err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	if len(entries) == 0 {
		fmt.Println("  （没有吊销记录）")
		Pause()
		return
	}
	for i, e := range entries {
		fmt.Printf("  %d. 序列号: %s\n", i+1, e.Serial)
		fmt.Printf("     吊销时间: %s\n", e.RevokedAt.Format("2006-01-02 15:04:05"))
	}
	Pause()
}

// ---------- 导出链 ----------

func certExportChain() {
	Section("导出证书链")
	certPath := SelectFile("选择证书", config.CertDir(), []string{".crt"}, nil, 0)
	if certPath == "" {
		return
	}
	base := strings.TrimSuffix(filepath.Base(certPath), ".crt")
	outPath := filepath.Join(config.CertDir(), base+"-chain.pem")

	fmt.Println("导出中... / Exporting...")
	n, err := core.ExportChain(certPath, outPath)
	if err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("完成 / Done")
	fmt.Println("  输出:", outPath)
	fmt.Printf("  包含 %d 张证书（从叶子到根）\n", n)
	Pause()
}

// ---------- 自签名 ----------

func certSelfGen() {
	Section("生成自签名证书")
	cn := ReadLine("CN (域名/IP): ")
	dns := ReadLine("额外 DNS（逗号分隔，可空）: ")
	days := ReadLine("有效天数 / Days (365): ")
	d := 365
	fmt.Sscanf(days, "%d", &d)
	var dnsNames []string
	if dns != "" {
		dnsNames = strings.Split(dns, ",")
	}
	res, err := core.IssueLeafCert(core.LeafOptions{
		Name:     cn,
		DNSNames: dnsNames,
		Days:     d,
		IsServer: true,
		IsClient: true,
	}, nil, config.CertDir())
	if err != nil {
		fmt.Println("失败 / Failed:", err)
		Pause()
		return
	}
	fmt.Println("完成 / Done")
	fmt.Println("  证书 / Cert:", res.CertPath)
	fmt.Println("  私钥 / Key :", res.KeyPath)
	Pause()
}

// ---------- 其他 ----------

func certView() {
	Section("查看证书信息")
	path := SelectFile("选择证书", config.CertDir(), []string{".crt", ".pem"}, nil, 0)
	if path == "" {
		return
	}
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

func certCSR() {
	Section("生成 CSR")
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
	Section("检查网站证书")
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
	Section("验证证书链")
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
