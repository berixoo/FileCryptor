package main

import (
	"os"
	"path/filepath"
	"strings"
)

func convertPath(baseDir, filePath string, isEncrypt bool) string {
	// 获取相对路径
	relPath, _ := filepath.Rel(baseDir, filePath)

	// 确定输出目录
	var outputDir string
	if isEncrypt {
		outputDir = filepath.Dir(baseDir) + string(filepath.Separator) + filepath.Base(baseDir) + "_encrypted"
	} else {
		// 解密时去掉 _encrypted 后缀
		base := filepath.Base(baseDir)
		if len(base) > 10 && base[len(base)-10:] == "_encrypted" {
			base = base[:len(base)-10]
		}
		outputDir = filepath.Dir(baseDir) + string(filepath.Separator) + base
	}

	// 处理文件名
	if isEncrypt {
		return filepath.Join(outputDir, relPath+".enc")
	}
	// 解密时去掉 .enc 后缀
	return filepath.Join(outputDir, strings.TrimSuffix(relPath, ".enc"))
}

func collectFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
