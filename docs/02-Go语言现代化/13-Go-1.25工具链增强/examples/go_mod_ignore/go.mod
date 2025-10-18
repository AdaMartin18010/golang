module example.com/go_mod_ignore_demo

go 1.25

require (
    github.com/gin-gonic/gin v1.9.1
)

// 忽略非 Go 代码目录
ignore (
    // 文档和示例
    ./docs/...
    ./examples/...
    
    // 临时和输出文件
    ./tmp/...
    ./_output/...
    
    // 前端代码
    ./web/...
    
    // 脚本和工具
    ./scripts/...
    ./hack/...
)

