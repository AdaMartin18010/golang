@echo off
chcp 65001 >nul

echo ==========================================
echo     Go语言AI-Agent架构演示程序
echo ==========================================
echo.

REM 检查Go环境
echo 1. 检查Go环境...
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo 错误: 未找到Go环境，请先安装Go 1.24+
    pause
    exit /b 1
)

for /f "tokens=3" %%i in ('go version') do set GO_VERSION=%%i
echo    Go版本: %GO_VERSION%
echo.

REM 检查依赖
echo 2. 检查项目依赖...
go mod tidy
echo    依赖检查完成
echo.

REM 运行测试
echo 3. 运行测试套件...
go test -v
if %errorlevel% neq 0 (
    echo 错误: 测试失败
    pause
    exit /b 1
)
echo    所有测试通过 ✓
echo.

REM 运行基准测试
echo 4. 运行基准测试...
go test -bench=. -benchmem
echo    基准测试完成 ✓
echo.

REM 构建项目
echo 5. 构建项目...
go build -o ai-agent-demo.exe agent.go learning.go coordinator.go router.go specialized_agents.go main.go
if %errorlevel% neq 0 (
    echo 错误: 构建失败
    pause
    exit /b 1
)
echo    构建完成 ✓
echo.

REM 运行演示
echo 6. 运行AI-Agent架构演示...
echo ==========================================
ai-agent-demo.exe
echo ==========================================
echo.

REM 清理
echo 7. 清理临时文件...
del ai-agent-demo.exe >nul 2>&1
echo    清理完成 ✓
echo.

echo ==========================================
echo     演示完成！
echo ==========================================
echo.
echo 项目特点:
echo   ✓ 完整的AI-Agent架构实现
echo   ✓ 多种专业代理类型
echo   ✓ 智能协调和负载均衡
echo   ✓ 学习和决策引擎
echo   ✓ 完整的测试覆盖
echo   ✓ 高性能并发处理
echo.
echo 技术栈:
echo   - Go 1.24+ 标准库
echo   - 无外部依赖
echo   - 企业级架构设计
echo   - 云原生就绪
echo.
pause
