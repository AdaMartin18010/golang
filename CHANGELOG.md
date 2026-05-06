# Changelog

All notable changes to this project's Go 1.26.2 modernization effort are documented in this file.

## [1.26.2] - 2026-05-06

### Emergency Fixes (P0)
- Removed fictional API `runtime.SetGoroutineLeakCallback` from 3 documentation files
- Corrected `new(expr)` compiler semantics documentation and examples
- Added working `crypto/hpke` example (RFC 9180)
- Fixed `iter`/`unique` version attribution to Go 1.23 (was incorrectly labeled Go 1.25)
- Added CVE-2026-27143/27144 security advisory documentation
- Removed 3 empty `cmd/` stubs (cli, graphql-server, mqtt-client)

### Core Modernization (P1)
- Migrated `math/rand` → `math/rand/v2` in `pkg/logger/`, `pkg/loadbalancer/`, `pkg/sampling/`
- Realigned `pkg/memory/weak_pointer.go` from `runtime.SetFinalizer` simulation to `weak.Make[T]` API
- Implemented `pkg/validator.Struct()` with `validate` tags
- Completed `cmd/grpc-server/main.go` service registration (User + Health)
- Added `go fix` CI gate
- Added `runtime/secret`, `bytes.Buffer.Peek`, `net.Dialer.DialTCP` examples
- Added PGO build targets to Makefile (`pgo-profile`, `pgo-build`, `pgo-auto-build`)

### Knowledge Base Sync (P2)
- Added docs for `new(expr)` compiler semantics
- Added Green Tea GC formal model documentation
- Added `crypto/hpke` deep dive
- Added post-quantum TLS deployment guide
- Added `go fix` modernization guide
- Updated master index

### Architecture Evolution (P3)
- Established quarterly authority audit template (`docs/tracking/quarterly-authority-audit.md`)
- Created `scripts/quality-check.sh`
- Created `examples/experimental/` sandbox (simd-archsimd, goroutineleak-profile, runtime-secret)

### Bug Fixes
- Fixed `pkg/registry` deadlock (`notifyWatchers` nested RLock inside Lock)
- Fixed `pkg/registry` `CleanupExpiredServices` test (LastSeen overwrite)
- Fixed `pkg/utils/crypto` AES key size bug (base64 decode)
- Fixed `pkg/utils/id` `SequentialID` race (atomic counter)
- Fixed `pkg/utils/strings` `Mask` test expectation
- Fixed `pkg/logger` sample rate logic (0.0 should mean no logs for non-error levels)

### Version Attribution Cleanup
- Fixed numerous "Go 1.25" mislabels across `examples/modern-features/`, `docs/`, `go-knowledge-base/`
- Updated all README version badges from `Go 1.25.3` → `Go 1.26.2`
- Updated all submodule `go.mod` directives to `go 1.26` (28+ files)
- Updated all production Dockerfiles to `golang:1.26.2-alpine`
- Updated CI workflow Go version matrices to `['1.26.x']`
- Updated `.golangci.yml` with `go: "1.26"`

### Code Quality
- Formatted all Go source files with `gofmt`
- Integrated `scripts/quality-check.sh` into CI pipeline
- All `go build`, `go vet`, `go test` pass across `pkg/`, `internal/`, `cmd/`, `examples/`

## [1.26.0] - 2026-04-07

### Security
- Upgraded Docker base image from `golang:1.21-alpine` to `golang:1.26.2-alpine`
- Documented CVE-2026-27143 (prove/loopbce induction variable wrap bug)
- Documented CVE-2026-27144 (SSA lowering no-op conversion bypass)

## Migration Notes

- Go 1.26.2 is the **minimum supported version** for this codebase
- Previous Go versions (1.21–1.25) are no longer tested in CI
- `GOEXPERIMENT=goroutineleakprofile` replaces the fictional `runtime.SetGoroutineLeakCallback`
