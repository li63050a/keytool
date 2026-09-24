package i18n

type Lang string

const (
	ZH Lang = "zh"
	EN Lang = "en"
)

var Current Lang = ZH

var dict = map[Lang]map[string]string{
	ZH: {
		"app.title":       "KeyTool 密钥管理",
		"main.gpg":        "GPG 密钥管理",
		"main.ssh":        "SSH 密钥管理",
		"main.crypt":      "文件加解密",
		"main.cert":       "证书工具",
		"main.list":       "列出所有密钥",
		"main.settings":   "设置",
		"main.exit":       "退出",
		"menu.back":       "返回主菜单",
		"menu.select":     "请选择: ",
		"menu.invalid":    "无效选项，请重试。",
		"menu.pause":      "按回车返回...",
		"menu.bye":        "再见。",

		"gpg.title":       "GPG 密钥管理",
		"gpg.gen":         "生成 GPG 密钥",
		"gpg.enc":         "GPG 加密文件",
		"gpg.dec":         "GPG 解密文件",
		"gpg.sign":        "GPG 签名文件",
		"gpg.verify":      "GPG 验证签名",

		"ssh.title":       "SSH 密钥管理",
		"ssh.gen":         "生成 SSH 密钥",
		"ssh.list":        "查看 SSH 密钥",

		"crypt.title":     "文件加解密",
		"crypt.enc":       "AES 加密",
		"crypt.dec":       "AES 解密",

		"cert.title":      "证书工具",
		"cert.view":       "查看 SSL 证书",

		"list.title":      "所有密钥",

		"settings.title":  "设置",
		"settings.lang":   "当前语言: %s",
		"settings.switch": "切换到 %s",
		"settings.saved":  "语言已切换为: %s",
	},
	EN: {
		"app.title":       "KeyTool Key Manager",
		"main.gpg":        "GPG Key Management",
		"main.ssh":        "SSH Key Management",
		"main.crypt":      "File Encryption",
		"main.cert":       "Certificate Tools",
		"main.list":       "List All Keys",
		"main.settings":   "Settings",
		"main.exit":       "Exit",
		"menu.back":       "Back to Main Menu",
		"menu.select":     "Select: ",
		"menu.invalid":    "Invalid option, try again.",
		"menu.pause":      "Press Enter to continue...",
		"menu.bye":        "Bye.",

		"gpg.title":       "GPG Key Management",
		"gpg.gen":         "Generate GPG Key",
		"gpg.enc":         "GPG Encrypt File",
		"gpg.dec":         "GPG Decrypt File",
		"gpg.sign":        "GPG Sign File",
		"gpg.verify":      "GPG Verify Signature",

		"ssh.title":       "SSH Key Management",
		"ssh.gen":         "Generate SSH Key",
		"ssh.list":        "List SSH Keys",

		"crypt.title":     "File Encryption",
		"crypt.enc":       "AES Encrypt",
		"crypt.dec":       "AES Decrypt",

		"cert.title":      "Certificate Tools",
		"cert.view":       "View SSL Certificate",

		"list.title":      "All Keys",

		"settings.title":  "Settings",
		"settings.lang":   "Current language: %s",
		"settings.switch": "Switch to %s",
		"settings.saved":  "Language switched to: %s",
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
