# 文件夹加密功能设计

## 概述

扩展 FileCryptor 支持文件夹加密/解密，使用多线程并行处理，保持目录结构。

## 核心需求

- 支持拖拽文件夹进行加密/解密
- 逐文件加密，每个文件生成 .enc 文件
- 递归处理所有子目录，保持完整目录结构
- 输出到源文件夹同目录（如 `myfolder` → `myfolder_encrypted`）
- 多线程并行加密，并发数 = CPU 核心数
- 解决并发写入问题

## 架构设计

### 文件结构

```
FileCryptor/
├── main.go          # 检测文件/文件夹类型
├── crypto.go        # AES-256-GCM 加密/解密核心
├── folder.go        # 新增：文件夹遍历和批量处理
├── gui.go           # 支持文件夹操作的 GUI
└── ...
```

### 核心流程

1. 拖入文件/文件夹 → `main.go` 检测类型
2. 如果是文件 → 走现有逻辑
3. 如果是文件夹 → 调用 `folder.go` 的批量处理函数
4. 批量处理函数遍历所有文件，用 worker pool 并发加密
5. 进度通过 channel 传回 GUI

## 并发写入问题解决方案

### 问题分析

多个 goroutine 同时写入不同文件时，可能遇到：
1. 同一目录下并发创建文件 → OS 层面安全，无需特殊处理
2. 多个 goroutine 同时创建子目录 → `os.MkdirAll` 并发安全
3. 进度更新竞争 → 使用 channel 串行化

### 解决方案：Worker Pool 模式

```go
type Task struct {
    InputPath  string
    OutputPath string
}

// 使用 buffered channel 作为任务队列
taskChan := make(chan Task, workerCount)

// 每个 worker 独立处理，无共享状态
for i := 0; i < workerCount; i++ {
    go func() {
        for task := range taskChan {
            encryptFile(task.InputPath, task.OutputPath, password)
            progressChan <- 1  // 通知完成一个文件
        }
    }()
}
```

### 关键点

- 每个 worker 只操作自己的输入/输出文件，无竞争
- 目录创建在分发任务前完成（单线程）
- 进度通过 channel 串行更新到 GUI

## 文件夹遍历

### 遍历逻辑

```go
func collectFiles(dir string) ([]string, error) {
    var files []string
    filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if !info.IsDir() {
            files = append(files, path)
        }
        return nil
    })
    return files, nil
}
```

### 输出目录结构

```
源: D:\docs\myfolder\
  ├── file1.txt
  ├── file2.txt
  └── sub\
      └── file3.txt

输出: D:\docs\myfolder_encrypted\
  ├── file1.txt.enc
  ├── file2.txt.enc
  └── sub\
      └── file3.txt.enc
```

## 进度显示

- 预先计算总文件数
- 进度条显示：`已处理文件数 / 总文件数`
- 状态标签显示当前处理的文件名（可选）

## 错误处理

- 单个文件加密失败 → 记录错误，继续处理其他文件
- 全部完成后显示结果：成功 X 个，失败 Y 个
- 失败文件列表可展开查看

## 依赖

- 复用现有 `crypto.go` 的 `encrypt`/`decrypt` 函数
- 新增 `folder.go` 处理文件夹逻辑
- 修改 `main.go` 和 `gui.go` 支持文件夹操作
