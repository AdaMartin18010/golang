# TS-DT-009: Go Build Modes and Cross-Compilation

> **维度**: Technology Stack > Development Tools
> **级别**: S (16+ KB)
> **标签**: #go-build #cross-compilation #cgo #build-tags #ldflags
> **权威来源**:
>
> - [go build documentation](https://golang.org/cmd/go/#hdr-Build_modes) - Go team
> - [Cross Compilation](https://dave.cheney.net/2015/08/22/cross-compilation-with-go) - Dave Cheney

---

## 1. Build Modes

### 1.1 Default Build Mode

```bash
# Default: executable binary
go build -o myapp

# Output:
# - Linux: ELF binary
# - Windows: PE binary (.exe)
# - macOS: Mach-O binary
```

### 1.2 Available Build Modes

```bash
# Build as archive (static library)
go build -buildmode=archive -o libmylib.a

# Build as shared library (C-shared)
go build -buildmode=c-shared -o libmylib.so

# Build as shared library (C-archive)
go build -buildmode=c-archive -o libmylib.a

# Build as plugin
go build -buildmode=plugin -o myplugin.so

# Build as PIE (Position Independent Executable)
go build -buildmode=pie -o myapp

# Build with race detector
go build -race -o myapp

# Build with coverage
go build -cover -o myapp
```

---

## 2. Cross-Compilation

### 2.1 Cross-Compilation Basics

```bash
# List available targets
go tool dist list

# Common cross-compilation examples:

# Linux AMD64
go build -o myapp-linux-amd64

# Linux ARM64
go build -o myapp-linux-arm64

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -o myapp-windows-amd64.exe

# macOS AMD64
GOOS=darwin GOARCH=amd64 go build -o myapp-darwin-amd64

# macOS ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o myapp-darwin-arm64

# FreeBSD
GOOS=freebsd GOARCH=amd64 go build -o myapp-freebsd-amd64

# WebAssembly
GOOS=js GOARCH=wasm go build -o myapp.wasm
```

### 2.2 Cross-Compilation Script

```bash
#!/bin/bash
# build-all.sh - Build for multiple platforms

PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "linux/arm"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "freebsd/amd64"
)

VERSION=$(git describe --tags --always)
LDFLAGS="-s -w -X main.Version=$VERSION"

for platform in "${PLATFORMS[@]}"; do
    GOOS=${platform%/*}
    GOARCH=${platform#*/}
    output="myapp-$GOOS-$GOARCH"

    if [ "$GOOS" = "windows" ]; then
        output="${output}.exe"
    fi

    echo "Building for $GOOS/$GOARCH..."
    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "$LDFLAGS" -o "dist/$output"
done
```

---

## 3. Build Tags and Constraints

### 3.1 File Constraints

```go
// +build linux

// This file only builds on Linux
package main

import "fmt"

func Platform() string {
    return "Linux"
}
```

```go
// +build windows

// This file only builds on Windows
package main

import "fmt"

func Platform() string {
    return "Windows"
}
```

### 3.2 New Build Tags (Go 1.17+)

```go
//go:build linux && amd64

// This file only builds on Linux AMD64
package main
```

### 3.3 Using Build Tags

```bash
# Build with specific tags
go build -tags "production"
go build -tags "debug"
go build -tags "linux"
go build -tags "production linux"

# Build without specific tags
go build -tags "!windows"
```

---

## 4. Linker Flags

### 4.1 Common ldflags

```bash
# Strip debug information (reduce binary size)
go build -ldflags "-s -w"

# Set version at build time
go build -ldflags "-X main.Version=1.0.0 -X main.BuildTime=$(date -u +%Y%m%d%H%M%S)"

# Disable CGO
go build -ldflags "-linkmode external -extldflags -static"

# Full static binary
go build -a -installsuffix cgo -ldflags "-s -w -extldflags '-static'"
```

### 4.2 Build for Production

```bash
# Optimized production build
go build -ldflags "-s -w \
    -X main.Version=$(git describe --tags) \
    -X main.Commit=$(git rev-parse --short HEAD) \
    -X main.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o myapp

# UPX compression (optional)
upx --best myapp
```

---

## 5. CGO and Cross-Compilation

```bash
# Disable CGO for easier cross-compilation
CGO_ENABLED=0 go build

# Enable CGO for specific platform
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build

# Cross-compile with CGO (requires cross-compiler)
CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ \
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -ldflags '-linkmode external -extldflags -static'
```

---

## 6. Checklist

```
Build Checklist:
□ Cross-compilation tested
□ Build tags used appropriately
□ Version information embedded
□ Debug symbols stripped for production
□ CGO disabled if not needed
□ Binary size optimized
□ Platform-specific code separated
□ CI/CD handles multi-platform builds
```
