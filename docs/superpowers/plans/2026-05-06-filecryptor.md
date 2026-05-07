# FileCryptor Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a portable Windows file encryption/decryption tool with drag-and-drop GUI using AES-256-GCM.

**Architecture:** Single Go binary with Fyne GUI. Files are dragged onto the exe, which opens a GUI window for password input. Encrypted files go to `encrypted/`, decrypted files go to `decrypted/`. Uses scrypt for key derivation and AES-256-GCM for authenticated encryption.

**Tech Stack:** Go 1.21+, fyne.io/fyne/v2 (GUI), golang.org/x/crypto (scrypt)

---

## File Structure

```
FileCryptor/
├── main.go              # Entry point, CLI arg parsing, GUI launch
├── crypto.go            # AES-256-GCM encrypt/decrypt, scrypt key derivation
├── crypto_test.go       # Unit tests for crypto functions
├── gui.go               # Fyne GUI: password window, progress bar
├── go.mod               # Module definition
├── go.sum               # Dependency checksums
├── encrypted/           # Output for encrypted files (auto-created)
└── decrypted/           # Output for decrypted files (auto-created)
```

---

### Task 1: Initialize Go Module

**Files:**
- Create: `go.mod`

- [ ] **Step 1: Initialize Go module**

```bash
cd D:\Desktop\workspace\FileCryptor
go mod init filecryptor
```

- [ ] **Step 2: Add dependencies**

```bash
go get fyne.io/fyne/v2
go get golang.org/x/crypto
```

- [ ] **Step 3: Verify module**

```bash
cat go.mod
```

Expected output should contain module name and dependencies.

- [ ] **Step 4: Commit**

```bash
git init
git add go.mod go.sum
git commit -m "chore: initialize Go module with dependencies"
```

---

### Task 2: Implement Crypto Core - Key Derivation

**Files:**
- Create: `crypto.go`
- Create: `crypto_test.go`

- [ ] **Step 1: Write failing test for key derivation**

```go
// crypto_test.go
package main

import (
	"testing"
)

func TestDeriveKey(t *testing.T) {
	password := "testpassword123"
	salt := make([]byte, 16)
	for i := range salt {
		salt[i] = byte(i)
	}

	key1 := deriveKey(password, salt)
	key2 := deriveKey(password, salt)

	if len(key1) != 32 {
		t.Errorf("expected key length 32, got %d", len(key1))
	}

	// Same password + salt should produce same key
	for i := range key1 {
		if key1[i] != key2[i] {
			t.Error("same password+salt should produce same key")
			break
		}
	}

	// Different salt should produce different key
	salt2 := make([]byte, 16)
	for i := range salt2 {
		salt2[i] = byte(i + 1)
	}
	key3 := deriveKey(password, salt2)

	same := true
	for i := range key1 {
		if key1[i] != key3[i] {
			same = false
			break
		}
	}
	if same {
		t.Error("different salt should produce different key")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test -v -run TestDeriveKey
```

Expected: FAIL with "undefined: deriveKey"

- [ ] **Step 3: Write minimal implementation**

```go
// crypto.go
package main

import (
	"golang.org/x/crypto/scrypt"
)

func deriveKey(password string, salt []byte) []byte {
	key, _ := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	return key
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test -v -run TestDeriveKey
```

Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add crypto.go crypto_test.go
git commit -m "feat: add scrypt key derivation"
```

---

### Task 3: Implement Crypto Core - Encrypt Function

**Files:**
- Modify: `crypto.go`
- Modify: `crypto_test.go`

- [ ] **Step 1: Write failing test for encryption**

Add to `crypto_test.go`:

```go
func TestEncrypt(t *testing.T) {
	plaintext := []byte("Hello, World! This is a test message.")
	password := "strongpassword"

	ciphertext, err := encrypt(plaintext, password)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	// ciphertext should be longer than plaintext (salt + nonce + tag)
	if len(ciphertext) <= len(plaintext) {
		t.Error("ciphertext should be longer than plaintext")
	}

	// First 16 bytes should be salt
	// Next 12 bytes should be nonce
	// Rest should be encrypted data + tag
	if len(ciphertext) != 16+12+len(plaintext)+16 {
		t.Errorf("unexpected ciphertext length: %d", len(ciphertext))
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test -v -run TestEncrypt
```

Expected: FAIL with "undefined: encrypt"

- [ ] **Step 3: Write minimal implementation**

Add to `crypto.go`:

```go
import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

func encrypt(plaintext []byte, password string) ([]byte, error) {
	// Generate random salt
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("generating salt: %w", err)
	}

	// Derive key
	key := deriveKey(password, salt)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("creating cipher: %w", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("creating GCM: %w", err)
	}

	// Generate random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("generating nonce: %w", err)
	}

	// Encrypt
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	// Combine: salt + nonce + ciphertext
	result := make([]byte, 0, len(salt)+len(nonce)+len(ciphertext))
	result = append(result, salt...)
	result = append(result, nonce...)
	result = append(result, ciphertext...)

	return result, nil
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test -v -run TestEncrypt
```

Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add crypto.go crypto_test.go
git commit -m "feat: add AES-256-GCM encryption"
```

---

### Task 4: Implement Crypto Core - Decrypt Function

**Files:**
- Modify: `crypto.go`
- Modify: `crypto_test.go`

- [ ] **Step 1: Write failing test for decryption**

Add to `crypto_test.go`:

```go
func TestDecrypt(t *testing.T) {
	plaintext := []byte("Hello, World! This is a test message.")
	password := "strongpassword"

	// Encrypt first
	ciphertext, err := encrypt(plaintext, password)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	// Decrypt
	decrypted, err := decrypt(ciphertext, password)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}

	// Should match original
	if string(decrypted) != string(plaintext) {
		t.Errorf("decrypted text doesn't match: got %q, want %q", decrypted, plaintext)
	}
}

func TestDecryptWrongPassword(t *testing.T) {
	plaintext := []byte("Secret message")
	password := "correctpassword"

	ciphertext, err := encrypt(plaintext, password)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	// Try decrypt with wrong password
	_, err = decrypt(ciphertext, "wrongpassword")
	if err == nil {
		t.Error("expected error with wrong password, got nil")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test -v -run TestDecrypt
```

Expected: FAIL with "undefined: decrypt"

- [ ] **Step 3: Write minimal implementation**

Add to `crypto.go`:

```go
func decrypt(ciphertext []byte, password string) ([]byte, error) {
	if len(ciphertext) < 16+12+16 {
		return nil, fmt.Errorf("ciphertext too short")
	}

	// Extract salt, nonce, and encrypted data
	salt := ciphertext[:16]
	nonce := ciphertext[16:28]
	encrypted := ciphertext[28:]

	// Derive key
	key := deriveKey(password, salt)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("creating cipher: %w", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("creating GCM: %w", err)
	}

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed (wrong password or corrupted file): %w", err)
	}

	return plaintext, nil
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test -v -run TestDecrypt
```

Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add crypto.go crypto_test.go
git commit -m "feat: add AES-256-GCM decryption"
```

---

### Task 5: Implement Crypto Core - File Operations

**Files:**
- Modify: `crypto.go`
- Modify: `crypto_test.go`

- [ ] **Step 1: Write failing test for file encryption**

Add to `crypto_test.go`:

```go
import (
	"os"
	"path/filepath"
)

func TestEncryptFile(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()

	// Create test file
	inputFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(inputFile, []byte("Hello, World!"), 0644)

	// Encrypt
	outputFile := filepath.Join(tmpDir, "test.txt.enc")
	err := encryptFile(inputFile, outputFile, "password123")
	if err != nil {
		t.Fatalf("encryptFile failed: %v", err)
	}

	// Verify output exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("encrypted file not created")
	}

	// Verify output is different from input
	input, _ := os.ReadFile(inputFile)
	output, _ := os.ReadFile(outputFile)
	if string(input) == string(output) {
		t.Error("encrypted file should be different from input")
	}
}

func TestDecryptFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create and encrypt test file
	inputFile := filepath.Join(tmpDir, "test.txt")
	original := []byte("Hello, World!")
	os.WriteFile(inputFile, original, 0644)

	encryptedFile := filepath.Join(tmpDir, "test.txt.enc")
	encryptFile(inputFile, encryptedFile, "password123")

	// Decrypt
	decryptedFile := filepath.Join(tmpDir, "test.txt.dec")
	err := decryptFile(encryptedFile, decryptedFile, "password123")
	if err != nil {
		t.Fatalf("decryptFile failed: %v", err)
	}

	// Verify content matches
	decrypted, _ := os.ReadFile(decryptedFile)
	if string(decrypted) != string(original) {
		t.Errorf("decrypted content doesn't match: got %q, want %q", decrypted, original)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test -v -run TestEncryptFile
```

Expected: FAIL with "undefined: encryptFile"

- [ ] **Step 3: Write minimal implementation**

Add to `crypto.go`:

```go
func encryptFile(inputPath, outputPath, password string) error {
	// Read input file
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("reading input file: %w", err)
	}

	// Encrypt
	ciphertext, err := encrypt(plaintext, password)
	if err != nil {
		return fmt.Errorf("encrypting: %w", err)
	}

	// Write output file
	if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
		return fmt.Errorf("writing output file: %w", err)
	}

	return nil
}

func decryptFile(inputPath, outputPath, password string) error {
	// Read input file
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("reading input file: %w", err)
	}

	// Decrypt
	plaintext, err := decrypt(ciphertext, password)
	if err != nil {
		return fmt.Errorf("decrypting: %w", err)
	}

	// Write output file
	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		return fmt.Errorf("writing output file: %w", err)
	}

	return nil
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test -v -run "TestEncryptFile|TestDecryptFile"
```

Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add crypto.go crypto_test.go
git commit -m "feat: add file encrypt/decrypt operations"
```

---

### Task 6: Implement GUI - Window Structure

**Files:**
- Create: `gui.go`
- Modify: `main.go`

- [ ] **Step 1: Create main.go entry point**

```go
// main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"
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

	// Get exe directory for output
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("错误: 无法获取程序路径: %v\n", err)
		os.Exit(1)
	}
	exeDir := filepath.Dir(exePath)

	// Launch GUI
	showGUI(filePath, isEncrypted, exeDir)
}
```

- [ ] **Step 2: Create GUI skeleton**

```go
// gui.go
package main

import (
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

	// Status label
	statusLabel := widget.NewLabel("")

	// Buttons
	startBtn := widget.NewButton("开始", func() {
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
			if isEncrypted {
				err = processDecryption(filePath, exeDir, password, func(p float64) {
					progress.SetValue(p)
				})
			} else {
				err = processEncryption(filePath, exeDir, password, func(p float64) {
					progress.SetValue(p)
				})
			}

			if err != nil {
				dialog.ShowError(err, w)
				startBtn.Enable()
				progress.Hide()
			} else {
				dialog.ShowInformation("成功", operation+"完成！", w)
				w.Close()
			}
		}()
	})

	cancelBtn := widget.NewButton("取消", func() {
		w.Close()
	})

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
	content.Add(statusLabel)
	content.Add(container.NewHBox(startBtn, cancelBtn))

	w.SetContent(content)
	w.ShowAndRun()
}
```

- [ ] **Step 3: Add missing imports to main.go**

Update `main.go` imports:

```go
import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)
```

- [ ] **Step 4: Verify compilation**

```bash
go build -o /dev/null .
```

Expected: Compilation succeeds (may have unused function warnings)

- [ ] **Step 5: Commit**

```bash
git add main.go gui.go
git commit -m "feat: add GUI window structure"
```

---

### Task 7: Implement GUI - Processing Functions

**Files:**
- Modify: `gui.go`

- [ ] **Step 1: Add processEncryption function**

Add to `gui.go`:

```go
import (
	"fmt"
	"os"
	"path/filepath"
)

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
```

- [ ] **Step 2: Fix imports in gui.go**

Ensure `gui.go` has these imports:

```go
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
```

- [ ] **Step 3: Verify compilation**

```bash
go build -o FileCryptor.exe .
```

Expected: Compilation succeeds, FileCryptor.exe created

- [ ] **Step 4: Commit**

```bash
git add gui.go
git commit -m "feat: add encryption/decryption processing"
```

---

### Task 8: Integration Test - Manual Testing

**Files:**
- None (manual testing)

- [ ] **Step 1: Create test file**

```bash
echo "Hello, World! This is a test file for encryption." > test.txt
```

- [ ] **Step 2: Test drag-and-drop encryption**

Drag `test.txt` onto `FileCryptor.exe`:
- Enter password when prompted
- Verify `encrypted/test.txt.enc` is created

- [ ] **Step 3: Test drag-and-drop decryption**

Drag `encrypted/test.txt.enc` onto `FileCryptor.exe`:
- Enter same password
- Verify `decrypted/test.txt` is created with original content

- [ ] **Step 4: Test wrong password**

Drag `encrypted/test.txt.enc` onto `FileCryptor.exe`:
- Enter wrong password
- Verify error message appears

- [ ] **Step 5: Clean up test files**

```bash
rm test.txt
rm -rf encrypted decrypted
```

- [ ] **Step 6: Commit final version**

```bash
git add -A
git commit -m "feat: complete FileCryptor with GUI encryption/decryption"
```

---

### Task 9: Build Portable Executable

**Files:**
- None (build step)

- [ ] **Step 1: Build optimized executable**

```bash
go build -ldflags="-s -w" -o FileCryptor.exe .
```

- [ ] **Step 2: Verify executable size**

```bash
ls -lh FileCryptor.exe
```

Expected: ~10-15MB executable

- [ ] **Step 3: Test standalone execution**

Copy `FileCryptor.exe` to a different directory and test drag-and-drop.

- [ ] **Step 4: Final commit**

```bash
git add -A
git commit -m "build: create portable executable"
```
