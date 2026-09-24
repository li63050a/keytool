<div align="center">

# 🔐 KeyTool

**一个用 Go 编写的命令行密钥管理工具**

支持 GPG、SSH、AES 加解密、SSL 证书查看

在 Termux、Linux、macOS 上均可运行

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-Apache--2.0-blue?style=for-the-badge)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Termux%20%7C%20Linux%20%7C%20macOS-lightgrey?style=for-the-badge)]()
[![Go Report Card](https://goreportcard.com/badge/github.com/li63050a/keytool?style=for-the-badge)](https://goreportcard.com/report/github.com/li63050a/keytool)

[功能特性](#-功能特性) • [安装](#-安装) • [快速开始](#-快速开始) • [功能详解](#-功能详解) • [常见问题](#-常见问题) • [安全说明](#-安全说明)

</div>

---

## 📖 项目简介

**KeyTool** 是一个用 Go 语言编写的命令行密钥管理工具，目标是把日常用到的密码学操作集中到一个简单的中文菜单里，避免每次都要回忆各种命令的参数。

它解决了什么问题？

- 日常需要 GPG 签名、SSH 登录、加密文件、查看证书，但每个工具的语法都不一样
- GPG 的命令行参数复杂，英文菜单难懂
- 手机（Termux）上不方便敲长命令
- 希望有一个统一的地方存放所有密钥

KeyTool 用一个简单的中文数字菜单，把这些操作全部包起来。

---

## ✨ 功能特性

### 已实现

| 模块 | 功能 | 状态 |
|:---|:---|:---:|
| **GPG** | 生成 RSA 4096 位 OpenPGP 密钥对 | ✅ |
| **GPG** | 用公钥加密文件 | ✅ |
| **GPG** | 用私钥生成分离签名 | ✅ |
| **SSH** | 生成 ed25519 密钥 | ✅ |
| **SSH** | 生成 RSA 4096 密钥 | ✅ |
| **AES** | AES-256-GCM 文件加密 | ✅ |
| **AES** | AES-256-GCM 文件解密 | ✅ |
| **SSL** | 解析 PEM 证书信息 | ✅ |
| **统计** | 查看密钥数量 | ✅ |

### 计划中

| 模块 | 功能 | 状态 |
|:---|:---|:---:|
| GPG | 私钥解密文件 | 🚧 |
| GPG | 验证分离签名 | 🚧 |
| GPG | 导出公钥到 GitHub | 🚧 |
| GPG | 列出所有 GPG 密钥 | 🚧 |
| SSH | 列出所有 SSH 密钥 | 🚧 |
| SSH | 一键复制公钥到服务器 | 🚧 |
| SSL | 生成自签名证书 | 🚧 |
| SSL | 检查网站证书有效期 | 🚧 |
| 通用 | 配置文件支持 | 🚧 |
| 通用 | 多语言支持 | 🚧 |

---

## 🚀 安装

### 前置条件

- Go 1.21 或更高版本
- Termux、Linux 或 macOS

检查 Go 版本：

```bash
go version
```

方式一：go install（最简单）

```bash
go install github.com/li63050a/keytool@latest
```

安装后运行：

```bash
keytool
```

如果提示 command not found，把 $GOPATH/bin 加到 PATH：

```bash
export PATH=$PATH:$(go env GOPATH)/bin
echo 'export PATH=$PATH:'$(go env GOPATH)'/bin' >> ~/.bashrc
```

方式二：从源码编译

```bash
git clone https://github.com/li63050a/keytool.git
cd keytool
go build -o keytool .
./keytool
```

方式三：Termux 用户

```bash
# 更新包列表
pkg update && pkg upgrade

# 安装 Go
pkg install golang git

# 下载并安装
go install github.com/li63050a/keytool@latest

# 加入 PATH
export PATH=$PATH:$(go env GOPATH)/bin

# 运行
keytool
```

方式四：交叉编译到其他平台

在任意平台编译到其他系统：

```bash
# 编译到 Linux amd64
GOOS=linux GOARCH=amd64 go build -o keytool-linux .

# 编译到 macOS arm64（M 系列芯片）
GOOS=darwin GOARCH=arm64 go build -o keytool-mac .

# 编译到 Windows
GOOS=windows GOARCH=amd64 go build -o keytool.exe .

# 编译到 Android ARM（Termux 内可直接运行）
GOOS=linux GOARCH=arm64 go build -o keytool-arm64 .
```

---

🎯 快速开始

运行程序

```bash
./keytool
```

你会看到中文菜单：

```
============================
      KeyTool 密钥管理
============================
  1. 生成 GPG 密钥
  2. 生成 SSH 密钥
  3. 查看密钥统计
  4. 文件加密 (AES-256)
  5. 文件解密 (AES-256)
  6. GPG 加密文件
  7. GPG 签名文件
  8. 查看 SSL 证书
  0. 退出
============================
```

操作方式

1. 输入对应数字（如 1）
2. 按回车
3. 按提示输入参数
4. 完成后按回车返回主菜单

---

📚 功能详解

1️⃣ 生成 GPG 密钥

选择 1，输入姓名和邮箱。

示例：

```
---- 生成 GPG 密钥 ----
姓名: 张三
邮箱: zhangsan@example.com
生成中...
完成！
  公钥: /root/.keytool/keys/gpg/zhangsan@example.com-public.asc
  私钥: /root/.keytool/keys/gpg/zhangsan@example.com-private.asc
  指纹: ABCD1234EFGH5678...
```

生成的文件：

文件 说明 能不能给别人
xxx-public.asc 公钥 ✅ 可以公开
xxx-private.asc 私钥 ❌ 绝对不能给别人

公钥的用途：

· 别人用它给你加密文件
· 别人用它验证你的签名

私钥的用途：

· 解密别人发给你的文件
· 给文件签名

2️⃣ 生成 SSH 密钥

选择 2，输入名称，选择类型。

示例：

```
---- 生成 SSH 密钥 ----
名称（用作文件名）: mykey
类型:
  1. ed25519（推荐）
  2. rsa
请选择 [1-2]: 1
生成中...
完成！
  私钥: /root/.keytool/keys/ssh/mykey_ed25519
  公钥: /root/.keytool/keys/ssh/mykey_ed25519.pub
  指纹: SHA256:xxxxx
```

ed25519 vs RSA：

类型 密钥长度 速度 安全性 推荐度
ed25519 256 位 快 高 ⭐⭐⭐⭐⭐
RSA 4096 位 慢 高 ⭐⭐⭐

公钥怎么用？

放到服务器的 ~/.ssh/authorized_keys：

```bash
cat mykey_ed25519.pub | ssh user@server "cat >> ~/.ssh/authorized_keys"
```

或者粘贴到 GitHub：

1. 打开 https://github.com/settings/keys
2. New SSH key
3. 粘贴公钥内容
4. 保存

3️⃣ 查看密钥统计

选择 3，显示：

```
---- 密钥统计 ----
  GPG 密钥: 2 个
  SSH 密钥: 3 个
  存放目录: /root/.keytool/keys
```

4️⃣ 文件加密（AES-256）

选择 4，输入文件路径、输出路径、密码。

加密流程：

1. 生成 16 字节随机盐值
2. 用 PBKDF2-SHA256（10 万次迭代）从密码派生 32 字节密钥
3. 生成 12 字节随机 nonce
4. 用 AES-256-GCM 加密文件内容
5. 把 盐值 + nonce + 密文 写入输出文件

特点：

· 对称加密，速度快
· 只需密码，不需要密钥对
· GCM 模式自带完整性校验
· 密码错误无法解密

示例：

```
输入文件路径: /root/secret.txt
输出文件路径: /root/secret.txt.enc
密码: MyP@ssw0rd123
完成！加密文件: /root/secret.txt.enc
```

5️⃣ 文件解密（AES-256）

选择 5，输入加密文件路径、输出路径、密码。

示例：

```
输入加密文件路径: /root/secret.txt.enc
输出文件路径: /root/secret-dec.txt
密码: MyP@ssw0rd123
完成！解密文件: /root/secret-dec.txt
```

注意：

· 密码必须和加密时完全一致
· 密码错误会提示 解密失败（密码可能错误）
· 加密文件被篡改也会解密失败

6️⃣ GPG 加密文件

选择 6，用对方的公钥加密文件。

示例：

```
---- GPG 加密文件 ----
公钥路径 (.asc): /root/.keytool/keys/gpg/lisi@example.com-public.asc
输入文件路径: /root/msg.txt
输出文件路径: /root/msg.gpg
完成！加密文件: /root/msg.gpg
```

原理：

· 文件先用随机对称密钥加密
· 对称密钥再用对方公钥加密
· 只有对方的私钥能解出对称密钥，再解密文件

注意： 只能加密，不能解密。解密需要对方的私钥。

7️⃣ GPG 签名文件

选择 7，用自己的私钥生成签名。

示例：

```
---- GPG 签名文件 ----
私钥路径 (.asc): /root/.keytool/keys/gpg/zhangsan@example.com-private.asc
输入文件路径: /root/document.pdf
输出签名路径 (.asc): /root/document.pdf.asc
完成！签名文件: /root/document.pdf.asc
```

用途：

· 把原文件和签名文件一起发给别人
· 别人用你的公钥验证签名
· 确认文件是你发的，且未被篡改

分离签名的好处：

· 原文件保持原样
· 签名单独一个文件
· 适合给二进制文件（如 PDF、ZIP）签名

8️⃣ 查看 SSL 证书

选择 8，输入 PEM 格式证书路径。

示例：

```
---- 查看 SSL 证书 ----
证书文件路径 (.pem/.crt): /root/cert.pem

证书信息:
  主体: example.com
  颁发者: Let's Encrypt Authority X3
  生效: 2024-01-01 00:00:00
  过期: 2024-04-01 00:00:00
  是否 CA: false
  域名: example.com, www.example.com
```

支持的格式：

· .pem
· .crt
· .cer（PEM 编码的）

用途：

· 检查证书什么时候过期
· 确认证书覆盖哪些域名
· 排查证书链问题

---

📂 密钥存放位置

所有密钥统一存放在用户主目录下的 .keytool/ 目录：

```
~/.keytool/
├── keys/
│   ├── gpg/                          # GPG 密钥
│   │   ├── zhangsan@example.com-public.asc
│   │   └── zhangsan@example.com-private.asc
│   └── ssh/                          # SSH 密钥
│       ├── mykey_ed25519
│       ├── mykey_ed25519.pub
│       ├── server_rsa
│       └── server_rsa.pub
└── config.toml                       # 预留配置文件
```

各系统下的路径：

系统 路径
Linux /home/用户名/.keytool/
Termux /data/data/com.termux/files/home/.keytool/
macOS /Users/用户名/.keytool/
Windows C:\Users\用户名\AppData\Roaming\keytool\

---

🏗️ 项目结构

```
keytool/
├── main.go                           # 程序入口
├── go.mod                            # Go 模块定义
├── go.sum                            # 依赖校验
├── README.md                         # 本文档
├── LICENSE                           # Apache-2.0
└── internal/                         # 内部包
    ├── config/
    │   └── config.go                 # 路径和目录管理
    ├── core/                         # 核心逻辑
    │   ├── gpg.go                    # GPG 密钥生成
    │   ├── ssh.go                    # SSH 密钥生成
    │   ├── encrypt.go                # AES 文件加解密
    │   ├── gpgenc.go                 # GPG 加密与签名
    │   ├── ssl.go                    # SSL 证书解析
    │   ├── helpers.go                # 辅助函数
    │   └── core_test.go              # 单元测试
    └── menu/
        └── menu.go                   # 中文菜单界面
```

分层设计：

```
        main.go
           │
           ▼
    ┌──────────────┐
    │    menu      │  只负责和用户交互
    └──────┬───────┘
           │
           ▼
    ┌──────────────┐
    │    core      │  只负责具体功能
    └──────┬───────┘
           │
           ▼
    ┌──────────────┐
    │   config     │  只负责路径管理
    └──────────────┘
```

为什么这么分？

· core 不打印、不读键盘，只接收参数、返回结果
· 以后如果想加 TUI、GUI 或 Web 界面，直接调用 core 即可
· 单元测试只测 core，UI 不用测

---

🧪 运行测试

```bash
go test ./...
```

带详细输出：

```bash
go test ./internal/core/ -v
```

测试覆盖：

· GPG 密钥生成
· SSH 密钥生成
· AES 加解密往返
· GPG 加密
· GPG 签名

预期输出：

```
=== RUN   TestAll
=== RUN   TestAll/GPG生成
=== RUN   TestAll/SSH生成
=== RUN   TestAll/AES加解密
=== RUN   TestAll/GPG加密签名
--- PASS: TestAll (0.84s)
PASS
ok      github.com/li63050a/keytool/internal/core       0.872s
```

---

❓ 常见问题

生成 GPG 密钥时报 unsupported preferred hash function？

原因： DefaultHash 配置错误。

修复： 本项目已在 internal/core/gpg.go 中使用 crypto.SHA256，不会有这个问题。如果你用的是旧版本，请更新：

```bash
go install github.com/li63050a/keytool@latest
```

Termux 里中文显示有空格，比如「密 钥 管 理」？

原因： Termux 终端的字体对 CJK 字符宽度渲染问题，不是程序 bug。

解决方案：

1. 长按屏幕 → More → Style → 选择支持 CJK 等宽的字体（如 Noto Sans Mono CJK）
2. 或者忽略，功能完全正常

密钥文件放哪了？

默认在 ~/.keytool/keys/：

· GPG 密钥在 gpg/ 子目录
· SSH 密钥在 ssh/ 子目录

私钥丢了怎么办？

GPG 和 SSH 私钥一旦丢失就无法恢复。

务必提前备份 ~/.keytool/ 目录到加密的离线介质。

备份方式：

```bash
# 打包加密备份
tar czf keytool-backup.tar.gz ~/.keytool/
./keytool  # 选择 4，加密 keytool-backup.tar.gz
```

能加密目录吗？

目前只支持单个文件。加密目录可以先打包：

```bash
tar czf mydir.tar.gz mydir/
./keytool  # 选择 4，加密 mydir.tar.gz
```

支持 Windows 吗？

理论上可以编译，但未在 Windows 上测试。密钥路径会使用 %APPDATA%。

为什么 go install 之后找不到命令？

把 $GOPATH/bin 加到 PATH：

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

永久生效：

```bash
echo 'export PATH=$PATH:'$(go env GOPATH)'/bin' >> ~/.bashrc
source ~/.bashrc
```

能同时管理多个身份吗？

可以。生成多个 GPG 密钥即可，每个用不同的邮箱。使用时指定对应的路径。

AES 密码忘了怎么办？

无法恢复。 AES 加密是不可逆的，密码丢失意味着文件永远无法解密。

建议： 把密码记录在密码管理器里（如 Bitwarden、KeePass）。

---

🔒 安全说明

私钥保护

· 私钥绝对不能泄露。 一旦泄露，别人可以冒充你签名、解密你的文件。
· 建议把私钥备份到加密的离线介质（如加密 U 盘）。
· 不要把私钥上传到网盘、GitHub、微信等公开平台。
· GPG 私钥文件 .asc 本身没有密码保护，建议把整个 ~/.keytool/ 打包加密后再备份。

密码选择

· AES 加密的密码建议至少 16 位。
· 包含大小写字母、数字、符号。
· 不要在不同文件间重复使用同一个密码。
· 不要用生日、手机号、常见单词。
· 推荐用密码管理器生成随机密码。

算法说明

用途 算法 说明
GPG 密钥 RSA 4096 兼容性最好
GPG 加密 AES-256 OpenPGP 默认对称算法
GPG 签名 SHA-256 哈希算法
SSH 密钥 ed25519 现代、快速、安全
SSH 密钥 RSA 4096 兼容老系统
文件加密 AES-256-GCM 带认证的对称加密
密钥派生 PBKDF2-SHA256 10 万次迭代

威胁模型

能防住：

· 本地文件被其他人看到（加密后看不到内容）
· 文件传输过程中被篡改（GCM 会检测）
· 文件被中间人替换（GPG 签名可验证）

防不住：

· 电脑被完全控制（键盘记录器可以窃取密码）
· 私钥被复制（如果备份介质不安全）
· 弱密码被暴力破解（AES 强度依赖密码强度）

---

🛠️ 开发指南

环境准备

```bash
git clone https://github.com/li63050a/keytool.git
cd keytool
go mod download
```

目录说明

目录 职责
internal/config 路径管理，不依赖其他包
internal/core 核心功能，不依赖 menu
internal/menu 用户交互，依赖 core 和 config

添加新功能

1. 在 internal/core/ 下写功能函数
2. 在 internal/core/core_test.go 加测试
3. 在 internal/menu/menu.go 加菜单项
4. 运行 go test ./... 确认通过
5. 运行 go build -o keytool . 确认能编译
6. 提交 PR

代码规范

· 用 gofmt 格式化：gofmt -w .
· 用 go vet 检查：go vet ./...
· 函数名用英文，注释可以中文
· 错误信息用中文，方便用户理解

提交规范

```
feat: 添加 GPG 解密功能
fix: 修复 hash 配置错误
docs: 更新 README
test: 添加 AES 测试
refactor: 重构 config 模块
```

---

🗺️ 路线图

v0.1.x（当前）

☑ GPG 密钥生成
☑ SSH 密钥生成
☑ AES 文件加解密
☑ GPG 加密和签名
☑ SSL 证书解析
☑ 单元测试

v0.2.0（计划）

☐ GPG 私钥解密文件
☐ GPG 验证分离签名
☐ 列出所有 GPG 密钥
☐ 列出所有 SSH 密钥
☐ 导出公钥到 GitHub

v0.3.0（计划）

☐ 配置文件支持
☐ 生成自签名证书
☐ 检查网站证书有效期
☐ 一键复制公钥到服务器

v1.0.0（目标）

☐ 稳定的 API
☐ 完整的文档
☐ 多语言支持
☐ GUI 版本

---

🤝 贡献

欢迎提交 Issue 和 Pull Request！

提交 PR 前请确保：

```bash
go test ./...
go build -o keytool .
gofmt -l .
```

都通过。

贡献方式：

1. Fork 本仓库
2. 创建分支：git checkout -b feature/xxx
3. 提交改动：git commit -m "feat: xxx"
4. 推送分支：git push origin feature/xxx
5. 打开 Pull Request

---

📄 更新日志

v0.1.0 — 2026-09-24

新增：

· GPG 密钥生成（RSA 4096）
· SSH 密钥生成（ed25519、RSA 4096）
· AES-256-GCM 文件加密
· AES-256-GCM 文件解密
· GPG 加密文件
· GPG 签名文件
· SSL 证书解析
· 密钥统计
· 单元测试

修复：

· 修复 unsupported preferred hash function 错误

---

🙏 致谢

· ProtonMail/go-crypto — OpenPGP 实现
· golang.org/x/crypto — SSH 和 PBKDF2
· Shields.io — 徽章

---

📜 许可证

Apache-2.0

```
Copyright 2026 li63050a

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

---

<div align="center">

⭐ 如果这个项目对你有帮助，请给个 Star ⭐

Made with ❤️ by li63050a

</div>
