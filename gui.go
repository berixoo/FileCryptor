// gui.go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func showGUI(filePath string, isEncrypted bool, exeDir string) {
	a := app.New()
	w := a.NewWindow("FileCryptor")
	w.Resize(fyne.NewSize(400, 300))
	w.CenterOnScreen()

	// File info
	fileName := widget.NewLabel("文件: " + filePath)
	fileName.Wrapping = fyne.TextWrapWord

	// Operation type
	operation := "加密"
	if isEncrypted {
		operation = "解密"
	}
	opLabel := widget.NewLabel("操作: " + operation)

	// Password input
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("请输入密码")

	// Confirm password (only for encryption)
	var confirmEntry *widget.Entry
	if !isEncrypted {
		confirmEntry = widget.NewPasswordEntry()
		confirmEntry.SetPlaceHolder("请再次输入密码")
	}

	// Progress bar
	progress := widget.NewProgressBar()
	progress.Hide()

	// Progress channel for thread-safe UI updates
	progressChan := make(chan float64, 10)
	doneChan := make(chan error, 1)

	// Buttons
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

			if isEncrypted {
				err = processDecryption(filePath, exeDir, password, progressFunc)
			} else {
				err = processEncryption(filePath, exeDir, password, progressFunc)
			}

			doneChan <- err
		}()
	})

	cancelBtn := widget.NewButton("取消", func() {
		w.Close()
	})

	// Goroutine to handle progress updates on main thread
	go func() {
		for {
			select {
			case p := <-progressChan:
				progress.SetValue(p)
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

	// Layout
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

func processEncryption(inputPath, exeDir string, password string, progress func(float64)) error {
	// Create encrypted directory
	encDir := filepath.Join(exeDir, "encrypted")
	if err := os.MkdirAll(encDir, 0755); err != nil {
		return fmt.Errorf("创建加密目录失败: %w", err)
	}

	// Output file path
	fileName := filepath.Base(inputPath)
	outputPath := filepath.Join(encDir, fileName+".enc")

	progress(0.1)

	// Read input
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}

	progress(0.3)

	// Encrypt
	ciphertext, err := encrypt(plaintext, password)
	if err != nil {
		return fmt.Errorf("加密失败: %w", err)
	}

	progress(0.8)

	// Write output
	if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	progress(1.0)
	return nil
}

func processDecryption(inputPath, exeDir string, password string, progress func(float64)) error {
	// Create decrypted directory
	decDir := filepath.Join(exeDir, "decrypted")
	if err := os.MkdirAll(decDir, 0755); err != nil {
		return fmt.Errorf("创建解密目录失败: %w", err)
	}

	// Output file path (remove .enc extension)
	fileName := filepath.Base(inputPath)
	outputName := strings.TrimSuffix(fileName, ".enc")
	outputPath := filepath.Join(decDir, outputName)

	progress(0.1)

	// Read input
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}

	progress(0.3)

	// Decrypt
	plaintext, err := decrypt(ciphertext, password)
	if err != nil {
		return fmt.Errorf("解密失败（密码错误或文件损坏）: %w", err)
	}

	progress(0.8)

	// Write output
	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	progress(1.0)
	return nil
}
