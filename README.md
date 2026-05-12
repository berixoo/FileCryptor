# FileCryptor

A portable file and folder encryption tool with drag-and-drop GUI, built with Go and Fyne.

## Features

- **AES-256-GCM** encryption with scrypt key derivation
- **Drag-and-drop** interface — just drop files/folders onto the exe
- **Folder support** — recursively encrypts all files with multi-threaded processing
- **Portable** — single executable, no installation required
- **Cross-platform** — built with Go and Fyne

## Download

Download the latest `FileCryptor.exe` from [Releases](../../releases).

## Usage

### Encrypt a file

1. Drag a file onto `FileCryptor.exe`
2. Enter password (twice for confirmation)
3. Click "Start"
4. Encrypted file (`.enc`) is created in the same directory

### Decrypt a file

1. Drag a `.enc` file onto `FileCryptor.exe`
2. Enter password
3. Click "Start"
4. Decrypted file is created in the same directory

### Encrypt a folder

1. Drag a folder onto `FileCryptor.exe`
2. Enter password (twice for confirmation)
3. Click "Start"
4. Encrypted folder (`foldername_encrypted/`) is created in the parent directory

### Decrypt a folder

1. Drag a `*_encrypted` folder onto `FileCryptor.exe`
2. Enter password
3. Click "Start"
4. Decrypted folder is created alongside the encrypted one

## Build from Source

### Prerequisites

- Go 1.21+
- Fyne CLI (optional, for icon packaging)

### Build

```bash
# Clone the repository
git clone https://github.com/yourusername/FileCryptor.git
cd FileCryptor

# Build
go build -ldflags="-s -w" -o FileCryptor.exe .

# Or with icon (requires fyne CLI)
fyne package -icon icon.png -name FileCryptor
```

### Run Tests

```bash
go test -v ./...
```

## Project Structure

```
FileCryptor/
├── main.go          # Entry point, CLI argument handling
├── crypto.go        # AES-256-GCM encryption/decryption core
├── crypto_test.go   # Unit tests for crypto functions
├── folder.go        # Folder traversal and batch processing
├── folder_test.go   # Unit tests for folder operations
├── gui.go           # Fyne GUI with progress bar
├── icon.png         # Application icon
└── go.mod           # Go module definition
```

## Security

- **AES-256-GCM** — authenticated encryption providing confidentiality and integrity
- **scrypt** — password-based key derivation (N=32768, r=8, p=1)
- **Random salt** — each file gets a unique 16-byte salt
- **Random nonce** — each encryption uses a unique 12-byte nonce
- **No password storage** — passwords are never saved to disk

## License

[MIT](LICENSE)
