package main

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

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
