// gui.go
package main

import (
	"fmt"
	"os"
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
	w.Resize(fyne.NewSize(400, 300))
	w.CenterOnScreen()

	// File info
	fileName := widget.NewLabel("文件: " + targetPath)
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
	// Channel for overwrite confirmation
	confirmChan := make(chan string, 1)
	confirmResultChan := make(chan bool, 1)

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
				err = processDecryption(targetPath, password, progressFunc, confirmChan, confirmResultChan)
			} else {
				err = processEncryption(targetPath, password, progressFunc, confirmChan, confirmResultChan)
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

func processEncryption(inputPath, password string, progress func(float64), confirmChan chan<- string, confirmResultChan <-chan bool) error {
	// Output file path (same directory as input)
	outputPath := inputPath + ".enc"

	// Check if output file exists
	if _, err := os.Stat(outputPath); err == nil {
		confirmChan <- outputPath
		if !<-confirmResultChan {
			return fmt.Errorf("操作已取消")
		}
	}

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

func processDecryption(inputPath, password string, progress func(float64), confirmChan chan<- string, confirmResultChan <-chan bool) error {
	// Output file path (remove .enc extension)
	outputPath := strings.TrimSuffix(inputPath, ".enc")

	// Check if output file exists
	if _, err := os.Stat(outputPath); err == nil {
		confirmChan <- outputPath
		if !<-confirmResultChan {
			return fmt.Errorf("操作已取消")
		}
	}

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
