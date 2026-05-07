// main.go
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法: FileCryptor.exe <文件或文件夹路径>")
		fmt.Println("请将文件或文件夹拖拽到本程序图标上")
		os.Exit(1)
	}

	targetPath := os.Args[1]

	// 检查路径是否存在
	info, err := os.Stat(targetPath)
	if os.IsNotExist(err) {
		fmt.Printf("错误: 路径不存在: %s\n", targetPath)
		os.Exit(1)
	}

	// 判断是文件还是文件夹
	isDir := info.IsDir()

	// 判断操作类型
	var isEncrypted bool
	if isDir {
		// 文件夹：检查名称是否以 _encrypted 结尾
		isEncrypted = strings.HasSuffix(targetPath, "_encrypted")
	} else {
		// 文件：检查扩展名是否为 .enc
		isEncrypted = strings.HasSuffix(strings.ToLower(targetPath), ".enc")
	}

	// 启动 GUI
	showGUI(targetPath, isEncrypted, isDir)
}
