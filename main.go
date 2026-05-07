// main.go
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法: FileCryptor.exe <文件路径>")
		fmt.Println("请将文件拖拽到本程序图标上")
		os.Exit(1)
	}

	filePath := os.Args[1]

	// Verify file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("错误: 文件不存在: %s\n", filePath)
		os.Exit(1)
	}

	// Determine operation based on file extension
	isEncrypted := strings.HasSuffix(strings.ToLower(filePath), ".enc")

	// Launch GUI
	showGUI(filePath, isEncrypted)
}
