#!/bin/bash
# Quality Check Script for Go Knowledge Base Project
# 
# Usage: ./scripts/quality-check.sh
# 
# Checks:
# 1. Go code compilation
# 2. Go tests (short mode)
# 3. go vet
# 4. go fix modernization check
# 5. API existence verification for documented features
# 6. Markdown link check (basic)

set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

PASS=0
FAIL=0

log_pass() {
    echo -e "${GREEN}✓${NC} $1"
    ((PASS++))
}

log_fail() {
    echo -e "${RED}✗${NC} $1"
    ((FAIL++))
}

log_info() {
    echo -e "${YELLOW}ℹ${NC} $1"
}

echo "========================================"
echo "Quality Check - $(date)"
echo "Go Version: $(go version)"
echo "========================================"
echo

# 1. Go Compilation
log_info "Checking Go compilation..."
if go build ./pkg/... ./internal/... ./cmd/... 2>/dev/null; then
    log_pass "Go compilation (pkg/internal/cmd)"
else
    log_fail "Go compilation (pkg/internal/cmd)"
fi

# 2. Go Tests (short mode)
log_info "Running short tests..."
if go test -short ./pkg/... ./internal/... 2>/dev/null; then
    log_pass "Go short tests"
else
    log_fail "Go short tests"
fi

# 3. go vet
log_info "Running go vet..."
if go vet ./pkg/... ./internal/... ./cmd/... 2>/dev/null; then
    log_pass "go vet"
else
    log_fail "go vet"
fi

# 4. go fix modernization check
log_info "Checking go fix modernization..."
if go fix ./... 2>/dev/null; then
    if git diff --quiet 2>/dev/null; then
        log_pass "go fix (no modernization needed)"
    else
        log_fail "go fix found unmodernized code. Run 'go fix ./...' and commit."
        git diff --stat
        git checkout -- . 2>/dev/null || true
    fi
else
    log_fail "go fix execution failed"
fi

# 5. API Existence Verification
log_info "Verifying documented API existence..."

# List of APIs that MUST exist in Go 1.26
declare -a REQUIRED_APIS=(
    "errors.AsType"
    "log/slog.NewMultiHandler"
    "crypto/hpke.Seal"
    "crypto/hpke.Open"
    "bytes.Buffer.Peek"
    "math/rand/v2.Float64"
)

API_FAIL=0
for api in "${REQUIRED_APIS[@]}"; do
    if go doc "$api" >/dev/null 2>&1; then
        : # API exists
    else
        echo "  Missing API: $api"
        API_FAIL=1
    fi
done

if [ "$API_FAIL" -eq 0 ]; then
    log_pass "All documented APIs exist"
else
    log_fail "Some documented APIs do not exist"
fi

# 6. Fictional API Check (must NOT exist)
log_info "Checking for fictional APIs..."
if grep -r "runtime.SetGoroutineLeakCallback" --include="*.md" --include="*.go" . 2>/dev/null; then
    log_fail "Fictional API 'runtime.SetGoroutineLeakCallback' found in codebase"
else
    log_pass "No fictional APIs found"
fi

# Summary
echo
echo "========================================"
echo "Results: $PASS passed, $FAIL failed"
echo "========================================"

if [ "$FAIL" -gt 0 ]; then
    exit 1
fi
