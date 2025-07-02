@echo off
setlocal

echo ======================================================================
echo ==      Go Profile-Guided Optimization (PGO) Demo Flow            ==
echo ======================================================================
echo.

REM Step 0: Check for dependencies (Go and hey)
go version > nul 2>&1
if %errorlevel% neq 0 (
    echo Error: Go is not installed or not in PATH.
    goto:eof
)
hey -h > nul 2>&1
if %errorlevel% neq 0 (
    echo Error: 'hey' is not installed or not in PATH.
    echo Please run: go install github.com/rakyll/hey@latest
    goto:eof
)
echo Dependencies check: OK
echo.


REM Step 1: Run benchmark WITHOUT PGO to get a baseline
echo --- Step 1: Running benchmark WITHOUT PGO (Baseline) ---
go test -bench=. -benchmem > baseline.txt
echo Baseline benchmark results saved to baseline.txt
type baseline.txt
echo.


REM Step 2: Build and run the application in the background
echo --- Step 2: Building and running the web server... ---
go build -o pgo_demo.exe .
start "PGO_Server" /B .\pgo_demo.exe
echo Server started in the background. PID: %!
timeout /t 3 > nul
echo.


REM Step 3: Apply load and collect the profile
echo --- Step 3: Applying load and collecting CPU profile... ---
echo Running 'hey' for 15 seconds to generate load on the hot path...
hey -z 15s "http://localhost:8080/?n=30" > nul 2>&1

echo Fetching profile from pprof endpoint...
curl -s -o default.pgo "http://localhost:6060/debug/pprof/cpu?seconds=10"
if %errorlevel% neq 0 (
    echo Error: Failed to curl pprof endpoint. Is the server running?
    taskkill /IM pgo_demo.exe /F > nul 2>&1
    goto:eof
)
echo Profile saved to default.pgo
echo.


REM Step 4: Stop the server
echo --- Step 4: Stopping the web server... ---
taskkill /IM pgo_demo.exe /F > nul 2>&1
echo Server stopped.
echo.


REM Step 5: Run benchmark WITH PGO
echo --- Step 5: Running benchmark WITH PGO... ---
echo The Go tool will now automatically use 'default.pgo'
go test -bench=. -benchmem > pgo_enabled.txt
echo PGO-enabled benchmark results saved to pgo_enabled.txt
type pgo_enabled.txt
echo.


REM Step 6: Compare results
echo --- Step 6: Comparison ---
echo.
echo Baseline (without PGO):
type baseline.txt
echo.
echo With PGO:
type pgo_enabled.txt
echo.
echo Notice the ns/op (nanoseconds per operation) value. It should be lower with PGO.
echo.


REM Cleanup
del default.pgo
del pgo_demo.exe
del baseline.txt
del pgo_enabled.txt

endlocal 