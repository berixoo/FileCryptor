# FileCryptor

**[English](README_EN.md)** | 中文

便携式文件/文件夹加密工具，支持拖拽操作，基于 Go + Fyne 构建。

## 功能特性

- **AES-256-GCM** 加密，scrypt 密钥派生
- **拖拽操作** — 直接拖放文件/文件夹到程序图标
- **文件夹支持** — 递归加密所有文件，多线程并行处理
- **绿色便携** — 单个可执行文件，无需安装
- **覆盖确认** — 输出文件已存在时会提示

## 下载

从 [Releases](../../releases) 下载最新的 `FileCryptor.exe`。

## 使用方法

### 加密文件

1. 将文件拖拽到 `FileCryptor.exe` 上
2. 输入密码（需输入两次确认）
3. 点击"开始"
4. 加密文件（`.enc`）生成在同目录下

### 解密文件

1. 将 `.enc` 文件拖拽到 `FileCryptor.exe` 上
2. 输入密码
3. 点击"开始"
4. 解密文件生成在同目录下

### 加密文件夹

1. 将文件夹拖拽到 `FileCryptor.exe` 上
2. 输入密码（需输入两次确认）
3. 点击"开始"
4. 加密后的文件夹（`原文件夹名_encrypted/`）生成在父目录下

### 解密文件夹

1. 将 `*_encrypted` 文件夹拖拽到 `FileCryptor.exe` 上
2. 输入密码
3. 点击"开始"
4. 解密后的文件夹生成在加密文件夹旁边

## 从源码构建

### 环境要求

- Go 1.21+
- Fyne CLI（可选，用于打包图标）

### 构建

```bash
# 克隆仓库
git clone https://github.com/berixoo/FileCryptor.git
cd FileCryptor

# 构建
go build -ldflags="-s -w" -o FileCryptor.exe .

# 或带图标打包（需要 fyne CLI）
fyne package -icon icon.png -name FileCryptor
```

### 运行测试

```bash
go test -v ./...
```

## 项目结构

```
FileCryptor/
├── main.go          # 入口，命令行参数处理
├── crypto.go        # AES-256-GCM 加密/解密核心
├── crypto_test.go   # 加密函数单元测试
├── folder.go        # 文件夹遍历和批量处理
├── folder_test.go   # 文件夹操作单元测试
├── gui.go           # Fyne GUI 界面
├── icon.png         # 应用图标
└── go.mod           # Go 模块定义
```

## 安全说明

- **AES-256-GCM** — 认证加密，提供机密性和完整性保护
- **scrypt** — 密码派生密钥（N=32768, r=8, p=1）
- **随机盐值** — 每个文件使用独立的 16 字节盐
- **随机 Nonce** — 每次加密使用独立的 12 字节 Nonce
- **不存储密码** — 密码永远不会写入磁盘

## 许可证

[MIT](LICENSE)
