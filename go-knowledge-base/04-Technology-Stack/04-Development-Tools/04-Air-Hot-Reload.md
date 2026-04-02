# TS-DT-004: Air - Hot Reload for Go

> **维度**: Technology Stack > Development Tools
> **级别**: S (16+ KB)
> **标签**: #air #hot-reload #development #golang #live-reload
> **权威来源**:
>
> - [Air Documentation](https://github.com/cosmtrek/air) - GitHub

---

## 1. Air Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          Air Hot Reload Flow                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Developer                                                                  │
│     │                                                                       │
│     │ Save file                                                            │
│     ▼                                                                       │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                          Air Process                                 │   │
│  │                                                                      │   │
│  │  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐             │   │
│  │  │   Watch     │───►│   Build     │───►│    Run      │             │   │
│  │  │  File       │    │  (go build) │    │  Binary     │             │   │
│  │  │  Changes    │    │             │    │             │             │   │
│  │  └─────────────┘    └─────────────┘    └─────────────┘             │   │
│  │         ▲                                    │                       │   │
│  │         │           ┌─────────────┐         │                       │   │
│  │         └───────────│  Cleanup    │◄────────┘                       │   │
│  │                     │ (kill proc) │                                 │   │
│  │                     └─────────────┘                                 │   │
│  │                                                                      │   │
│  │  Configuration: .air.toml                                           │   │
│  │  - Watches .go files                                                 │   │
│  │  - Excludes vendor, test files                                       │   │
│  │  - Builds on change                                                  │   │
│  │  - Restarts process                                                  │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                         │
│                                    ▼                                         │
│                              Application                                     │
│                              Running                                         │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Installation and Configuration

```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Initialize Air configuration
air init

# This creates .air.toml
```

```toml
# .air.toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ."
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  kill_delay = "0s"
  log = "build-errors.log"
  send_interrupt = false
  stop_on_error = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  time_only = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
```

---

## 3. Usage

```bash
# Run with default config
air

# Run with specific config
air -c .air.toml

# Show version
air -v

# Build only (don't run)
air build

# Run with custom args
air -- arg1 arg2
```

---

## 4. Advanced Configuration

```toml
# .air.toml for web server
root = "."
tmp_dir = "tmp"

[build]
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/server"
  delay = 100
  exclude_dir = ["assets", "tmp", "vendor", "web/node_modules"]
  exclude_regex = ["_test.go"]
  include_ext = ["go", "html", "css", "js"]

[proxy]
  # Enable live reload for frontend
  enabled = true
  proxy_port = 8090
  app_port = 8080
```

---

## 5. Checklist

```
Air Configuration Checklist:
□ .air.toml in project root
□ Correct build command
□ Proper exclusions configured
□ Tmp directory created
□ Color coding enabled
□ Kill delay appropriate
□ Exclude test files
□ Include all needed extensions
```
