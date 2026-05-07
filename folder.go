package main

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
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

// BatchResult holds the results of a batch encrypt/decrypt operation.
type BatchResult struct {
	TotalFiles   int
	SuccessFiles int
	FailedFiles  int
	Errors       []error
}

func batchEncrypt(dir string, password string, workerCount int, progress func()) BatchResult {
	// 收集所有文件
	files, err := collectFiles(dir)
	if err != nil {
		return BatchResult{Errors: []error{err}}
	}

	result := BatchResult{TotalFiles: len(files)}

	// 预创建所有输出目录
	outputDir := filepath.Join(filepath.Dir(dir), filepath.Base(dir)+"_encrypted")
	for _, f := range files {
		relPath, _ := filepath.Rel(dir, f)
		outPath := filepath.Join(outputDir, relPath+".enc")
		os.MkdirAll(filepath.Dir(outPath), 0755)
	}

	// 创建任务通道
	taskChan := make(chan string, len(files))
	var wg sync.WaitGroup
	var mu sync.Mutex

	// 启动 workers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filePath := range taskChan {
				relPath, _ := filepath.Rel(dir, filePath)
				outPath := filepath.Join(outputDir, relPath+".enc")

				err := encryptFile(filePath, outPath, password)

				mu.Lock()
				if err != nil {
					result.FailedFiles++
					result.Errors = append(result.Errors, err)
				} else {
					result.SuccessFiles++
				}
				mu.Unlock()

				progress()
			}
		}()
	}

	// 分发任务
	for _, f := range files {
		taskChan <- f
	}
	close(taskChan)

	// 等待完成
	wg.Wait()

	return result
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
