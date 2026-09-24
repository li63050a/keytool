package i18n

type Lang string

const (
	ZH Lang = "zh"
	EN Lang = "en"
)

var Current Lang = ZH

var dict = map[Lang]map[string]string{
	ZH: {
		"app.title":      "KeyTool 密钥管理",
		"menu.gpg":       "GPG 密钥管理",
		"menu.ssh":       "SSH 密钥管理",
		"menu.crypt":     "文件加解密",
		"menu.cert":      "证书工具",
		"menu.tools":     "通用工具",
		"menu.list":      "列出所有密钥",
		"menu.settings":  "设置",
		"menu.exit":      "退出",
		"menu.back":      "返回上一级",
		"menu.select":    "请选择: ",
		"menu.invalid":   "无效选项",
		"menu.pause":     "按回车继续...",
		"menu.bye":       "再见。",
		"gpg.title":      "GPG 密钥管理",
		"gpg.gen":        "生成 GPG 密钥",
		"gpg.enc":        "GPG 加密文件",
		"gpg.dec":        "GPG 解密文件",
		"gpg.sign":       "GPG 签名文件",
		"gpg.verify":     "GPG 验证签名",
		"gpg.import":     "导入 GPG 公钥",
		"gpg.delete":     "删除 GPG 密钥",
		"gpg.info":       "查看 GPG 密钥详情",
		"ssh.title":      "SSH 密钥管理",
		"ssh.gen":        "生成 SSH 密钥",
		"ssh.list":       "查看 SSH 密钥",
		"ssh.info":       "查看 SSH 密钥详情",
		"ssh.delete":     "删除 SSH 密钥",
		"ssh.copy":       "复制公钥到服务器",
		"ssh.config":     "生成 SSH config",
		"crypt.title":    "文件加解密",
		"crypt.enc":      "AES 加密",
		"crypt.dec":      "AES 解密",
		"cert.title":     "证书工具",
		"cert.view":      "查看 SSL 证书",
		"cert.selfgen":   "生成自签名证书",
		"cert.csr":       "生成 CSR",
		"cert.check":     "检查网站证书",
		"cert.chain":     "验证证书链",
		"tools.title":    "通用工具",
		"tools.hash":     "计算哈希",
		"tools.b64e":     "Base64 编码",
		"tools.b64d":     "Base64 解码",
		"tools.backup":   "备份所有密钥",
		"tools.restore":  "恢复密钥备份",
		"tools.qr":       "生成公钥二维码",
		"settings.title": "设置",
		"settings.lang":  "当前语言: %s",
		"settings.switch":"切换到 %s",
		"settings.saved": "语言已切换为: %s",
	},
	EN: {
		"app.title":      "KeyTool Key Manager",
		"menu.gpg":       "GPG Key Management",
		"menu.ssh":       "SSH Key Management",
		"menu.crypt":     "File Encryption",
		"menu.cert":      "Certificate Tools",
		"menu.tools":     "Utilities",
		"menu.list":      "List All Keys",
		"menu.settings":  "Settings",
		"menu.exit":      "Exit",
		"menu.back":      "Back",
		"menu.select":    "Select: ",
		"menu.invalid":   "Invalid option",
		"menu.pause":     "Press Enter to continue...",
		"menu.bye":       "Bye.",
		"gpg.title":      "GPG Key Management",
		"gpg.gen":        "Generate GPG Key",
		"gpg.enc":        "GPG Encrypt File",
		"gpg.dec":        "GPG Decrypt File",
		"gpg.sign":       "GPG Sign File",
		"gpg.verify":     "GPG Verify Signature",
		"gpg.import":     "Import GPG Public Key",
		"gpg.delete":     "Delete GPG Key",
		"gpg.info":       "View GPG Key Info",
		"ssh.title":      "SSH Key Management",
		"ssh.gen":        "Generate SSH Key",
		"ssh.list":       "List SSH Keys",
		"ssh.info":       "View SSH Key Info",
		"ssh.delete":     "Delete SSH Key",
		"ssh.copy":       "Copy Public Key to Server",
		"ssh.config":     "Generate SSH config",
		"crypt.title":    "File Encryption",
		"crypt.enc":      "AES Encrypt",
		"crypt.dec":      "AES Decrypt",
		"cert.title":     "Certificate Tools",
		"cert.view":      "View SSL Certificate",
		"cert.selfgen":   "Generate Self-signed Cert",
		"cert.csr":       "Generate CSR",
		"cert.check":     "Check Website Certificate",
		"cert.chain":     "Verify Certificate Chain",
		"tools.title":    "Utilities",
		"tools.hash":     "Compute Hash",
		"tools.b64e":     "Base64 Encode",
		"tools.b64d":     "Base64 Decode",
		"tools.backup":   "Backup All Keys",
		"tools.restore":  "Restore Backup",
		"tools.qr":       "QR Code for Public Key",
		"settings.title": "Settings",
		"settings.lang":  "Current language: %s",
		"settings.switch":"Switch to %s",
		"settings.saved": "Language switched to: %s",
	},
}

func T(key string) string {
	if m, ok := dict[Current]; ok {
		if v, ok := m[key]; ok {
			return v
		}
	}
	return key
}

func Set(l string) {
	if l == "en" {
		Current = EN
	} else {
		Current = ZH
	}
}

func CurrentCode() string {
	if Current == EN {
		return "en"
	}
	return "zh"
}

func CurrentName() string {
	if Current == EN {
		return "English"
	}
	return "中文"
}

func OtherLangName() string {
	if Current == EN {
		return "中文"
	}
	return "English"
}
