# 文件夹加密功能实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 扩展 FileCryptor 支持文件夹加密/解密，使用多线程并行处理，保持目录结构。

**Architecture:** 新增 `folder.go` 处理文件夹遍历和批量加密，使用 worker pool 模式并发处理文件，通过 channel 串行化进度更新。

**Tech Stack:** Go, fyne.io/fyne/v2, golang.org/x/crypto, sync, runtime

---

## File Structure

```
FileCryptor/
├── main.go              # 修改：检测文件/文件夹类型
├── crypto.go            # 不变：AES-256-GCM 加密/解密核心
├── folder.go            # 新增：文件夹遍历和批量处理
├── folder_test.go       # 新增：文件夹操作测试
├── gui.go               # 修改：支持文件夹操作的 GUI
└── ...
```

---

### Task 1: 实现文件夹遍历函数

**Files:**
- Create: `folder.go`
- Create: `folder_test.go`

- [ ] **Step 1: 写失败测试**

```go
// folder_test.go
package main

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestCollectFiles(t *testing.T) {
	// 创建临时目录结构
	tmpDir := t.TempDir()
	os.MkdirAll(filepath.Join(tmpDir, "sub"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("hello"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file2.txt"), []byte("world"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "sub", "file3.txt"), []byte("sub"), 0644)

	files, err := collectFiles(tmpDir)
	if err != nil {
		t.Fatalf("collectFiles failed: %v", err)
	}

	// 排序以便比较
	sort.Strings(files)

	expected := []string{
		filepath.Join(tmpDir, "file1.txt"),
		filepath.Join(tmpDir, "file2.txt"),
		filepath.Join(tmpDir, "sub", "file3.txt"),
	}

	if len(files) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(files))
	}

	for i, f := range files {
		if f != expected[i] {
			t.Errorf("file[%d] = %q, want %q", i, f, expected[i])
		}
	}
}
```

- [ ] **Step 2: 运行测试验证失败**

```bash
go test -v -run TestCollectFiles
```

Expected: FAIL with "undefined: collectFiles"

- [ ] **Step 3: 写最小实现**

```go
// folder.go
package main

import (
	"path/filepath"
	"os"
)

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
```

- [ ] **Step 4: 运行测试验证通过**

```bash
go test -v -run TestCollectFiles
```

Expected: PASS

- [ ] **Step 5: 提交**

```bash
git add folder.go folder_test.go
git commit -m "feat: add folder file collection"
```

---

### Task 2: 实现路径转换函数

**Files:**
- Modify: `folder.go`
- Modify: `folder_test.go`

- [ ] **Step 1: 写失败测试**

```go
func TestConvertPaths(t *testing.T) {
	tests := []struct {
		inputDir   string
		inputFile  string
		isEncrypt  bool
		wantOutput string
	}{
		{
			inputDir:   "C:\\docs\\myfolder",
			inputFile:  "C:\\docs\\myfolder\\file1.txt",
			isEncrypt:  true,
			wantOutput: "C:\\docs\\myfolder_encrypted\\file1.txt.enc",
		},
		{
			inputDir:   "C:\\docs\\myfolder",
			inputFile:  "C:\\docs\\myfolder\\sub\\file2.txt",
			isEncrypt:  true,
			wantOutput: "C:\\docs\\myfolder_encrypted\\sub\\file2.txt.enc",
		},
		{
			inputDir:   "C:\\docs\\myfolder_encrypted",
			inputFile:  "C:\\docs\\myfolder_encrypted\\file1.txt.enc",
			isEncrypt:  false,
			wantOutput: "C:\\docs\\myfolder\\file1.txt",
		},
	}

	for _, tt := range tests {
		got := convertPath(tt.inputDir, tt.inputFile, tt.isEncrypt)
		if got != tt.wantOutput {
			t.Errorf("convertPath(%q, %q, %v) = %q, want %q",
				tt.inputDir, tt.inputFile, tt.isEncrypt, got, tt.wantOutput)
		}
	}
}
```

- [ ] **Step 2: 运行测试验证失败**

```bash
go test -v -run TestConvertPaths
```

Expected: FAIL with "undefined: convertPath"

- [ ] **Step 3: 写最小实现**

```go
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
```

- [ ] **Step 4: 运行测试验证通过**

```bash
go test -v -run TestConvertPaths
```

Expected: PASS

- [ ] **Step 5: 提交**

```bash
git add folder.go folder_test.go
git commit -m "feat: add path conversion for folder encryption"
```

---

### Task 3: 实现 Worker Pool 批量处理

**Files:**
- Modify: `folder.go`
- Modify: `folder_test.go`

- [ ] **Step 1: 写失败测试**

```go
func TestBatchEncrypt(t *testing.T) {
	// 创建临时目录结构
	tmpDir := t.TempDir()
	os.MkdirAll(filepath.Join(tmpDir, "sub"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("hello"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file2.txt"), []byte("world"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "sub", "file3.txt"), []byte("sub"), 0644)

	// 进度计数
	progressCount := 0
	progressFunc := func() {
		progressCount++
	}

	// 执行批量加密
	result := batchEncrypt(tmpDir, "password123", 2, progressFunc)

	// 验证结果
	if result.TotalFiles != 3 {
		t.Errorf("expected 3 total files, got %d", result.TotalFiles)
	}
	if result.SuccessFiles != 3 {
		t.Errorf("expected 3 success files, got %d", result.SuccessFiles)
	}
	if result.FailedFiles != 0 {
		t.Errorf("expected 0 failed files, got %d", result.FailedFiles)
	}

	// 验证输出文件存在
	outputDir := filepath.Join(filepath.Dir(tmpDir), filepath.Base(tmpDir)+"_encrypted")
	if _, err := os.Stat(filepath.Join(outputDir, "file1.txt.enc")); os.IsNotExist(err) {
		t.Error("file1.txt.enc not created")
	}
	if _, err := os.Stat(filepath.Join(outputDir, "sub", "file3.txt.enc")); os.IsNotExist(err) {
		t.Error("sub/file3.txt.enc not created")
	}

	// 验证进度调用次数
	if progressCount != 3 {
		t.Errorf("expected 3 progress calls, got %d", progressCount)
	}
}
```

- [ ] **Step 2: 运行测试验证失败**

```bash
go test -v -run TestBatchEncrypt
```

Expected: FAIL with "undefined: batchEncrypt"

- [ ] **Step 3: 写最小实现**

```go
import (
	"runtime"
	"sync"
)

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
```

- [ ] **Step 4: 运行测试验证通过**

```bash
go test -v -run TestBatchEncrypt
```

Expected: PASS

- [ ] **Step 5: 提交**

```bash
git add folder.go folder_test.go
git commit -m "feat: add batch encryption with worker pool"
```

---

### Task 4: 实现批量解密

**Files:**
- Modify: `folder.go`
- Modify: `folder_test.go`

- [ ] **Step 1: 写失败测试**

```go
func TestBatchDecrypt(t *testing.T) {
	// 创建临时目录并加密
	tmpDir := t.TempDir()
	os.MkdirAll(filepath.Join(tmpDir, "sub"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("hello"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "sub", "file2.txt"), []byte("world"), 0644)

	// 先加密
	encResult := batchEncrypt(tmpDir, "password123", 2, func() {})
	if encResult.SuccessFiles != 2 {
		t.Fatalf("encryption failed: %v", encResult.Errors)
	}

	// 解密
	encDir := filepath.Join(filepath.Dir(tmpDir), filepath.Base(tmpDir)+"_encrypted")
	progressCount := 0
	decResult := batchDecrypt(encDir, "password123", 2, func() { progressCount++ })

	if decResult.SuccessFiles != 2 {
		t.Errorf("expected 2 success files, got %d", decResult.SuccessFiles)
	}

	// 验证解密后的文件内容
	decDir := filepath.Join(filepath.Dir(encDir), filepath.Base(encDir))
	decDir = strings.TrimSuffix(decDir, "_encrypted")
	content, _ := os.ReadFile(filepath.Join(decDir, "file1.txt"))
	if string(content) != "hello" {
		t.Errorf("decrypted content = %q, want %q", content, "hello")
	}
}
```

- [ ] **Step 2: 运行测试验证失败**

```bash
go test -v -run TestBatchDecrypt
```

Expected: FAIL with "undefined: batchDecrypt"

- [ ] **Step 3: 写最小实现**

```go
func batchDecrypt(dir string, password string, workerCount int, progress func()) BatchResult {
	// 收集所有 .enc 文件
	files, err := collectFiles(dir)
	if err != nil {
		return BatchResult{Errors: []error{err}}
	}

	result := BatchResult{TotalFiles: len(files)}

	// 预创建所有输出目录
	outputDir := strings.TrimSuffix(dir, "_encrypted")
	if outputDir == dir {
		outputDir = dir + "_decrypted"
	}
	for _, f := range files {
		relPath, _ := filepath.Rel(dir, f)
		outPath := filepath.Join(outputDir, strings.TrimSuffix(relPath, ".enc"))
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
				outPath := filepath.Join(outputDir, strings.TrimSuffix(relPath, ".enc"))

				err := decryptFile(filePath, outPath, password)

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
```

- [ ] **Step 4: 运行测试验证通过**

```bash
go test -v -run TestBatchDecrypt
```

Expected: PASS

- [ ] **Step 5: 提交**

```bash
git add folder.go folder_test.go
git commit -m "feat: add batch decryption with worker pool"
```

---

### Task 5: 修改 main.go 支持文件夹

**Files:**
- Modify: `main.go`

- [ ] **Step 1: 修改 main.go 检测文件/文件夹**

```go
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
```

- [ ] **Step 2: 验证编译**

```bash
go build -o FileCryptor.exe .
```

Expected: 编译成功（可能有未使用变量警告）

- [ ] **Step 3: 提交**

```bash
git add main.go
git commit -m "feat: detect file/folder type in main"
```

---

### Task 6: 修改 GUI 支持文件夹操作

**Files:**
- Modify: `gui.go`

- [ ] **Step 1: 修改 showGUI 函数签名和逻辑**

```go
// gui.go
package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func showGUI(targetPath string, isEncrypted bool, isDir bool) {
	a := app.New()
	w := a.NewWindow("FileCryptor")
	w.Resize(fyne.NewSize(400, 350))
	w.CenterOnScreen()

	// 文件信息
	typeLabel := "文件"
	if isDir {
		typeLabel = "文件夹"
	}
	fileName := widget.NewLabel(typeLabel + ": " + targetPath)
	fileName.Wrapping = fyne.TextWrapWord

	// 操作类型
	operation := "加密"
	if isEncrypted {
		operation = "解密"
	}
	opLabel := widget.NewLabel("操作: " + operation)

	// 密码输入
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("请输入密码")

	// 确认密码（仅加密时）
	var confirmEntry *widget.Entry
	if !isEncrypted {
		confirmEntry = widget.NewPasswordEntry()
		confirmEntry.SetPlaceHolder("请再次输入密码")
	}

	// 进度条
	progress := widget.NewProgressBar()
	progress.Hide()

	// 进度通道
	progressChan := make(chan float64, 100)
	doneChan := make(chan error, 1)
	confirmChan := make(chan string, 1)
	confirmResultChan := make(chan bool, 1)

	// 按钮
	var startBtn *widget.Button
	startBtn = widget.NewButton("开始", func() {
		password := passwordEntry.Text
		if password == "" {
			dialog.ShowError(fmt.Errorf("密码不能为空"), w)
			return
		}

		if !isEncrypted && confirmEntry != nil {
			if password != confirmEntry.Text {
				dialog.ShowError(fmt.Errorf("两次密码不一致"), w)
				return
			}
		}

		progress.Show()
		startBtn.Disable()

		go func() {
			var err error
			progressFunc := func(p float64) {
				progressChan <- p
			}

			if isDir {
				// 文件夹批量处理
				workerCount := runtime.NumCPU()
				if isEncrypted {
					result := batchDecrypt(targetPath, password, workerCount, func() {
						progressChan <- -1 // -1 表示完成一个文件
					})
					if result.FailedFiles > 0 {
						err = fmt.Errorf("完成：%d 成功，%d 失败", result.SuccessFiles, result.FailedFiles)
					}
				} else {
					result := batchEncrypt(targetPath, password, workerCount, func() {
						progressChan <- -1
					})
					if result.FailedFiles > 0 {
						err = fmt.Errorf("完成：%d 成功，%d 失败", result.SuccessFiles, result.FailedFiles)
					}
				}
			} else {
				// 单文件处理
				if isEncrypted {
					outputPath := strings.TrimSuffix(targetPath, ".enc")
					if _, statErr := os.Stat(outputPath); statErr == nil {
						confirmChan <- outputPath
						if !<-confirmResultChan {
							err = fmt.Errorf("操作已取消")
						} else {
							err = decryptFile(targetPath, outputPath, password)
						}
					} else {
						err = decryptFile(targetPath, outputPath, password)
					}
				} else {
					outputPath := targetPath + ".enc"
					if _, statErr := os.Stat(outputPath); statErr == nil {
						confirmChan <- outputPath
						if !<-confirmResultChan {
							err = fmt.Errorf("操作已取消")
						} else {
							err = encryptFile(targetPath, outputPath, password)
						}
					} else {
						err = encryptFile(targetPath, outputPath, password)
					}
				}
				progressChan <- 1.0
			}

			doneChan <- err
		}()
	})

	cancelBtn := widget.NewButton("取消", func() {
		w.Close()
	})

	// 处理进度更新的 goroutine
	go func() {
		var fileCount float64
		var totalFiles float64

		for {
			select {
			case p := <-progressChan:
				if p == -1 {
					// 文件夹模式：完成一个文件
					fileCount++
					if totalFiles > 0 {
						progress.SetValue(fileCount / totalFiles)
					}
				} else {
					// 单文件模式
					progress.SetValue(p)
				}
			case path := <-confirmChan:
				dialog.ShowConfirm("文件已存在", "文件已存在，是否覆盖？\n"+path, func(ok bool) {
					confirmResultChan <- ok
				}, w)
			case err := <-doneChan:
				if err != nil {
					dialog.ShowError(err, w)
					startBtn.Enable()
					progress.Hide()
				} else {
					dialog.ShowInformation("成功", operation+"完成！", w)
					w.Close()
				}
				return
			}
		}
	}()

	// 布局
	content := container.NewVBox(
		fileName,
		opLabel,
		widget.NewSeparator(),
		passwordEntry,
	)
	if confirmEntry != nil {
		content.Add(confirmEntry)
	}
	content.Add(progress)
	content.Add(container.NewHBox(startBtn, cancelBtn))

	w.SetContent(content)
	w.ShowAndRun()
}
```

- [ ] **Step 2: 验证编译**

```bash
go build -o FileCryptor.exe .
```

Expected: 编译成功

- [ ] **Step 3: 运行所有测试**

```bash
go test -v ./...
```

Expected: 所有测试通过

- [ ] **Step 4: 提交**

```bash
git add gui.go
git commit -m "feat: add folder encryption/decryption support to GUI"
```

---

### Task 7: 集成测试

**Files:**
- None (手动测试)

- [ ] **Step 1: 创建测试文件夹**

```powershell
mkdir test_folder
echo "file1" > test_folder\a.txt
mkdir test_folder\sub
echo "file2" > test_folder\sub\b.txt
```

- [ ] **Step 2: 测试文件夹加密**

拖拽 `test_folder` 到 `FileCryptor.exe`：
- 输入密码
- 验证 `test_folder_encrypted` 目录创建
- 验证 `.enc` 文件存在

- [ ] **Step 3: 测试文件夹解密**

拖拽 `test_folder_encrypted` 到 `FileCryptor.exe`：
- 输入相同密码
- 验证解密后的文件内容正确

- [ ] **Step 4: 测试单文件仍可用**

拖拽单个文件测试，确保原有功能不受影响。

- [ ] **Step 5: 清理测试文件**

```powershell
Remove-Item -Recurse test_folder, test_folder_encrypted
```

- [ ] **Step 6: 提交最终版本**

```bash
git add -A
git commit -m "feat: complete folder encryption/decryption feature"
```
