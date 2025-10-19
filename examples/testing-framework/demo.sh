#!/bin/bash

# Go语言现代化 - 完整测试体系演示脚本

set -e

echo "=========================================="
echo "Go语言现代化 - 完整测试体系演示"
echo "=========================================="
echo

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_step() {
    echo -e "${PURPLE}[STEP]${NC} $1"
}

print_header() {
    echo -e "${CYAN}=========================================="
    echo "$1"
    echo "==========================================${NC}"
}

# 检查Go环境
check_go_environment() {
    print_step "检查Go环境..."
    
    if ! command -v go &> /dev/null; then
        print_error "Go未安装，请先安装Go 1.24+"
        exit 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    print_info "Go版本: $GO_VERSION"
    
    if [[ "$GO_VERSION" < "1.24" ]]; then
        print_warning "建议使用Go 1.24+版本"
    fi
    
    print_success "Go环境检查完成"
}

# 安装依赖
install_dependencies() {
    print_step "安装项目依赖..."
    
    if [ -f "go.mod" ]; then
        go mod download
        go mod tidy
        print_success "依赖安装完成"
    else
        print_error "go.mod文件不存在"
        exit 1
    fi
}

# 构建项目
build_project() {
    print_step "构建项目..."
    
    if go build -o testing-system main.go; then
        print_success "项目构建成功"
    else
        print_error "项目构建失败"
        exit 1
    fi
}

# 运行基本演示
run_basic_demo() {
    print_header "运行基本演示"
    
    print_step "执行集成测试演示..."
    ./testing-system 2>&1 | grep -A 20 "1. 集成测试演示" || true
    
    print_step "执行性能测试演示..."
    ./testing-system 2>&1 | grep -A 15 "2. 性能测试演示" || true
    
    print_step "执行质量监控演示..."
    ./testing-system 2>&1 | grep -A 10 "3. 质量监控演示" || true
    
    print_success "基本演示完成"
}

# 运行完整演示
run_full_demo() {
    print_header "运行完整演示"
    
    print_step "启动完整测试工作流..."
    timeout 30s ./testing-system 2>&1 || true
    
    print_success "完整演示完成"
}

# 运行单元测试
run_unit_tests() {
    print_header "运行单元测试"
    
    print_step "执行单元测试..."
    if go test -v ./...; then
        print_success "单元测试通过"
    else
        print_error "单元测试失败"
        exit 1
    fi
}

# 运行基准测试
run_benchmark_tests() {
    print_header "运行基准测试"
    
    print_step "执行基准测试..."
    if go test -bench=. -benchmem ./...; then
        print_success "基准测试完成"
    else
        print_warning "基准测试部分失败"
    fi
}

# 检查代码质量
check_code_quality() {
    print_header "检查代码质量"
    
    print_step "格式化代码..."
    go fmt ./...
    print_success "代码格式化完成"
    
    print_step "检查代码..."
    if command -v golangci-lint &> /dev/null; then
        golangci-lint run || print_warning "代码检查发现问题"
    else
        print_warning "golangci-lint未安装，跳过代码检查"
    fi
}

# 生成测试报告
generate_test_report() {
    print_header "生成测试报告"
    
    print_step "生成测试覆盖率报告..."
    mkdir -p test-results
    go test -coverprofile=test-results/coverage.out ./...
    go tool cover -html=test-results/coverage.out -o test-results/coverage.html
    
    if [ -f "test-results/coverage.html" ]; then
        print_success "测试覆盖率报告生成完成: test-results/coverage.html"
    else
        print_warning "测试覆盖率报告生成失败"
    fi
}

# 清理构建文件
cleanup() {
    print_step "清理构建文件..."
    rm -f testing-system
    print_success "清理完成"
}

# 显示帮助信息
show_help() {
    echo "用法: $0 [选项]"
    echo
    echo "选项:"
    echo "  -h, --help     显示此帮助信息"
    echo "  -b, --build    仅构建项目"
    echo "  -t, --test     运行单元测试"
    echo "  -p, --perf     运行性能测试"
    echo "  -d, --demo     运行基本演示"
    echo "  -f, --full     运行完整演示"
    echo "  -q, --quality  检查代码质量"
    echo "  -r, --report   生成测试报告"
    echo "  -a, --all      运行所有检查"
    echo "  -c, --clean    清理构建文件"
    echo
    echo "示例:"
    echo "  $0 --all       运行完整的演示和测试"
    echo "  $0 --demo      仅运行基本演示"
    echo "  $0 --test      仅运行单元测试"
}

# 主函数
main() {
    case "${1:-}" in
        -h|--help)
            show_help
            exit 0
            ;;
        -b|--build)
            check_go_environment
            install_dependencies
            build_project
            ;;
        -t|--test)
            check_go_environment
            install_dependencies
            run_unit_tests
            ;;
        -p|--perf)
            check_go_environment
            install_dependencies
            run_benchmark_tests
            ;;
        -d|--demo)
            check_go_environment
            install_dependencies
            build_project
            run_basic_demo
            cleanup
            ;;
        -f|--full)
            check_go_environment
            install_dependencies
            build_project
            run_full_demo
            cleanup
            ;;
        -q|--quality)
            check_go_environment
            check_code_quality
            ;;
        -r|--report)
            check_go_environment
            install_dependencies
            generate_test_report
            ;;
        -a|--all)
            check_go_environment
            install_dependencies
            build_project
            run_unit_tests
            run_benchmark_tests
            check_code_quality
            generate_test_report
            run_basic_demo
            cleanup
            ;;
        -c|--clean)
            cleanup
            ;;
        "")
            print_header "默认演示模式"
            check_go_environment
            install_dependencies
            build_project
            run_basic_demo
            cleanup
            ;;
        *)
            print_error "未知选项: $1"
            show_help
            exit 1
            ;;
    esac
}

# 捕获中断信号
trap 'print_warning "演示被中断"; cleanup; exit 1' INT TERM

# 运行主函数
main "$@"

print_header "演示完成"
print_success "感谢使用Go语言现代化测试体系！"
