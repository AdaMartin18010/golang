# CI/CD Integration Guide

**Go Formal Verifier (FV)** 的 CI/CD 集成指南

---

## 📋 目录

1. [概述](#概述)
2. [GitHub Actions](#github-actions)
3. [GitLab CI](#gitlab-ci)
4. [Jenkins](#jenkins)
5. [配置最佳实践](#配置最佳实践)
6. [报告发布](#报告发布)
7. [故障排查](#故障排查)

---

## 概述

FV工具专为CI/CD环境设计，提供以下特性：

- ✅ **非交互模式**: 完全支持自动化运行
- ✅ **配置文件**: 通过 `.fv.yaml` 统一配置
- ✅ **退出码控制**: 根据分析结果返回适当的退出码
- ✅ **多格式报告**: JSON、HTML、Markdown 等格式
- ✅ **质量门槛**: 设置最低质量分数
- ✅ **严格模式**: CI专用的严格检查配置

---

## GitHub Actions

### 基础配置

在项目根目录创建 `.github/workflows/fv-analysis.yml`:

```yaml
name: Formal Verification

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  analyze:
    runs-on: ubuntu-latest
    
    steps:
    - name: Checkout code
      uses: actions/checkout@v3
    
    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Install FV Tool
      run: |
        git clone https://github.com/your-org/formal-verifier.git
        cd formal-verifier
        go build -o fv ./cmd/fv
        sudo mv fv /usr/local/bin/
    
    - name: Run FV Analysis
      run: |
        fv analyze \
          --dir=. \
          --format=json \
          --output=fv-report.json \
          --fail-on-error \
          --no-color
    
    - name: Upload Report
      if: always()
      uses: actions/upload-artifact@v3
      with:
        name: fv-report
        path: fv-report.json
```

### 使用配置文件

```yaml
    - name: Generate Strict Config
      run: fv init-config --output=.fv-ci.yaml --strict
    
    - name: Run FV Analysis with Config
      run: fv analyze --config=.fv-ci.yaml --no-color
```

### HTML报告发布

```yaml
    - name: Generate HTML Report
      run: |
        fv analyze \
          --dir=. \
          --format=html \
          --output=fv-report.html \
          --no-color
    
    - name: Deploy to GitHub Pages
      if: github.ref == 'refs/heads/main'
      uses: peaceiris/actions-gh-pages@v3
      with:
        github_token: ${{ secrets.GITHUB_TOKEN }}
        publish_dir: ./
        publish_branch: gh-pages
        destination_dir: reports
```

### 质量门槛

```yaml
    - name: Quality Gate Check
      run: |
        fv analyze \
          --dir=. \
          --format=json \
          --output=fv-report.json \
          --no-color
        
        # 检查质量分数
        QUALITY=$(jq -r '.stats.quality_score' fv-report.json)
        echo "Quality Score: $QUALITY"
        
        if [ $QUALITY -lt 80 ]; then
          echo "❌ Quality score $QUALITY is below threshold 80"
          exit 1
        fi
```

### PR 评论集成

```yaml
    - name: Generate Markdown Report
      run: |
        fv analyze \
          --dir=. \
          --format=markdown \
          --output=fv-report.md \
          --no-color
    
    - name: Comment PR
      if: github.event_name == 'pull_request'
      uses: actions/github-script@v6
      with:
        github-token: ${{secrets.GITHUB_TOKEN}}
        script: |
          const fs = require('fs');
          const report = fs.readFileSync('fv-report.md', 'utf8');
          
          github.rest.issues.createComment({
            issue_number: context.issue.number,
            owner: context.repo.owner,
            repo: context.repo.repo,
            body: report
          });
```

---

## GitLab CI

### 基础配置1

在项目根目录创建 `.gitlab-ci.yml`:

```yaml
image: golang:1.26.2

stages:
  - setup
  - analyze
  - report

variables:
  FV_VERSION: "v1.0.0"
  GO111MODULE: "on"

before_script:
  - apt-get update -qq
  - apt-get install -y -qq git

install_fv:
  stage: setup
  script:
    - git clone https://github.com/your-org/formal-verifier.git
    - cd formal-verifier
    - go build -o fv ./cmd/fv
    - cp fv /usr/local/bin/
  artifacts:
    paths:
      - formal-verifier/fv
    expire_in: 1 hour

fv_analysis:
  stage: analyze
  dependencies:
    - install_fv
  script:
    - cp formal-verifier/fv /usr/local/bin/
    - fv init-config --output=.fv-ci.yaml --strict
    - fv analyze --config=.fv-ci.yaml --no-color
  artifacts:
    reports:
      # GitLab 支持的测试报告格式
      junit: fv-report.json
    paths:
      - fv-report.json
      - fv-report.html
    expire_in: 30 days
  allow_failure: false

generate_report:
  stage: report
  dependencies:
    - fv_analysis
  script:
    - cp formal-verifier/fv /usr/local/bin/
    - fv analyze --dir=. --format=html --output=public/index.html --no-color
  artifacts:
    paths:
      - public
  only:
    - main
```

### Merge Request 注释

```yaml
fv_mr_comment:
  stage: report
  dependencies:
    - fv_analysis
  script:
    - cp formal-verifier/fv /usr/local/bin/
    - fv analyze --dir=. --format=markdown --output=report.md --no-color
    - |
      curl --request POST \
        --header "PRIVATE-TOKEN: ${CI_JOB_TOKEN}" \
        --data-urlencode "body@report.md" \
        "${CI_API_V4_URL}/projects/${CI_PROJECT_ID}/merge_requests/${CI_MERGE_REQUEST_IID}/notes"
  only:
    - merge_requests
```

### Pages 部署

```yaml
pages:
  stage: report
  dependencies:
    - fv_analysis
  script:
    - cp formal-verifier/fv /usr/local/bin/
    - mkdir -p public
    - fv analyze --dir=. --format=html --output=public/index.html --no-color
  artifacts:
    paths:
      - public
  only:
    - main
```

---

## Jenkins

### Pipeline 配置

创建 `Jenkinsfile`:

```groovy
pipeline {
    agent any
    
    environment {
        FV_VERSION = 'v1.0.0'
        GO_VERSION = '1.21'
    }
    
    stages {
        stage('Setup') {
            steps {
                script {
                    // 安装 Go
                    sh '''
                        wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
                        tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
                        export PATH=$PATH:/usr/local/go/bin
                    '''
                    
                    // 安装 FV
                    sh '''
                        git clone https://github.com/your-org/formal-verifier.git
                        cd formal-verifier
                        /usr/local/go/bin/go build -o fv ./cmd/fv
                        cp fv /usr/local/bin/
                    '''
                }
            }
        }
        
        stage('Analysis') {
            steps {
                script {
                    sh '''
                        fv init-config --output=.fv-ci.yaml --strict
                        fv analyze --config=.fv-ci.yaml --no-color
                    '''
                }
            }
        }
        
        stage('Quality Gate') {
            steps {
                script {
                    def qualityScore = sh(
                        script: 'jq -r .stats.quality_score fv-report.json',
                        returnStdout: true
                    ).trim().toInteger()
                    
                    echo "Quality Score: ${qualityScore}"
                    
                    if (qualityScore < 80) {
                        error("Quality score ${qualityScore} is below threshold 80")
                    }
                }
            }
        }
        
        stage('Publish Report') {
            steps {
                script {
                    sh '''
                        fv analyze \
                          --dir=. \
                          --format=html \
                          --output=fv-report.html \
                          --no-color
                    '''
                    
                    publishHTML([
                        reportDir: '.',
                        reportFiles: 'fv-report.html',
                        reportName: 'FV Analysis Report'
                    ])
                }
            }
        }
    }
    
    post {
        always {
            archiveArtifacts artifacts: 'fv-report.*', fingerprint: true
        }
        success {
            echo 'FV Analysis passed!'
        }
        failure {
            echo 'FV Analysis failed!'
        }
    }
}
```

### Declarative Pipeline with Docker

```groovy
pipeline {
    agent {
        docker {
            image 'golang:1.26.2'
        }
    }
    
    stages {
        stage('FV Analysis') {
            steps {
                sh '''
                    # 安装 FV
                    git clone https://github.com/your-org/formal-verifier.git
                    cd formal-verifier
                    go build -o ../fv ./cmd/fv
                    cd ..
                    
                    # 运行分析
                    ./fv analyze \
                      --dir=. \
                      --format=json \
                      --output=fv-report.json \
                      --fail-on-error \
                      --no-color
                '''
            }
        }
    }
}
```

---

## 配置最佳实践

### 1. CI专用配置文件

创建 `.fv-ci.yaml`:

```yaml
project:
  root_dir: .
  recursive: true
  include_tests: true
  exclude_patterns:
    - vendor
    - .git
    - testdata

analysis:
  concurrent: true
  workers: 4
  timeout: 600

report:
  format: json
  output_path: fv-report.json

rules:
  complexity:
    cyclomatic_threshold: 5
    max_function_lines: 30
    max_parameters: 3

output:
  color_output: false    # CI环境禁用颜色
  show_progress: false   # CI环境禁用进度条
  fail_on_error: true
  min_quality_score: 80
```

### 2. 分级检查策略

#### 开发环境

```bash
fv analyze --dir=. --format=text
```

#### Pull Request

```bash
fv analyze --dir=. --format=markdown --output=pr-report.md
```

#### 主分支

```bash
fv analyze --config=.fv-ci.yaml --fail-on-error
```

#### 发布前

```bash
fv analyze --config=.fv-strict.yaml --fail-on-error
```

### 3. 缓存优化

#### GitHub Actions3

```yaml
- name: Cache FV Tool
  uses: actions/cache@v3
  with:
    path: ~/bin/fv
    key: fv-${{ runner.os }}-${{ env.FV_VERSION }}
```

#### GitLab CI3

```yaml
cache:
  key: fv-${CI_COMMIT_REF_SLUG}
  paths:
    - .fv-cache/
```

### 4. 并行化分析

```yaml
# GitHub Actions Matrix Strategy
strategy:
  matrix:
    module: [api, core, utils]

steps:
  - name: Analyze ${{ matrix.module }}
    run: fv analyze --dir=./${{ matrix.module }} --no-color
```

---

## 报告发布

### 1. GitHub Pages

```bash
# 自动发布到 gh-pages 分支
fv analyze --format=html --output=public/index.html --no-color
git checkout gh-pages
cp -r public/* .
git add .
git commit -m "Update FV report"
git push origin gh-pages
```

### 2. Markdown Badge

在 `README.md` 中添加质量徽章:

```markdown
[![FV Quality](https://img.shields.io/badge/FV%20Quality-85%25-green)](./fv-report.html)
```

### 3. Slack 通知

```bash
# 获取质量分数
QUALITY=$(jq -r '.stats.quality_score' fv-report.json)

# 发送 Slack 消息
curl -X POST -H 'Content-type: application/json' \
  --data "{\"text\":\"FV Analysis Complete: Quality Score ${QUALITY}\"}" \
  ${SLACK_WEBHOOK_URL}
```

### 4. 邮件报告

```bash
# 生成HTML报告
fv analyze --format=html --output=report.html --no-color

# 发送邮件
mail -s "FV Analysis Report" \
  -a "Content-Type: text/html" \
  team@example.com < report.html
```

---

## 故障排查

### 常见问题

#### 1. 退出码非零

**问题**: CI任务因FV返回非零退出码而失败

**解决方案**:

```bash
# 允许失败但记录
fv analyze --dir=. --no-color || echo "FV analysis found issues"

# 或者移除 --fail-on-error
fv analyze --dir=. --no-color
```

#### 2. 内存不足

**问题**: 分析大型项目时内存耗尽

**解决方案**:

```yaml
# 限制并发worker数量
analysis:
  workers: 2
  max_file_size: 512  # KB
```

#### 3. 超时

**问题**: 分析时间过长

**解决方案**:

```yaml
analysis:
  timeout: 1200  # 20分钟
  concurrent: true
  workers: 8
```

#### 4. 权限问题

**问题**: 无法写入报告文件

**解决方案**:

```bash
# 确保输出目录存在且有写权限
mkdir -p reports
chmod 755 reports
fv analyze --output=reports/fv-report.json --no-color
```

### 调试模式

```bash
# 启用详细输出
export FV_DEBUG=1
fv analyze --dir=. --no-color

# 或使用配置
analysis:
  verbose: true
```

---

## 完整示例

### GitHub Actions 完整工作流

```yaml
name: Complete FV Pipeline

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  fv-analysis:
    runs-on: ubuntu-latest
    
    steps:
    - name: Checkout
      uses: actions/checkout@v3
    
    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Cache FV
      uses: actions/cache@v3
      with:
        path: ~/bin
        key: fv-${{ runner.os }}-v1.0.0
    
    - name: Install FV
      if: steps.cache.outputs.cache-hit != 'true'
      run: |
        git clone https://github.com/your-org/formal-verifier.git /tmp/fv
        cd /tmp/fv
        go build -o ~/bin/fv ./cmd/fv
        chmod +x ~/bin/fv
    
    - name: Setup Config
      run: |
        export PATH=$PATH:~/bin
        fv init-config --output=.fv-ci.yaml --strict
    
    - name: Run Analysis
      run: |
        export PATH=$PATH:~/bin
        fv analyze --config=.fv-ci.yaml --no-color
    
    - name: Generate Reports
      if: always()
      run: |
        export PATH=$PATH:~/bin
        fv analyze --format=html --output=fv-report.html --no-color
        fv analyze --format=markdown --output=fv-report.md --no-color
    
    - name: Upload Artifacts
      if: always()
      uses: actions/upload-artifact@v3
      with:
        name: fv-reports
        path: |
          fv-report.json
          fv-report.html
          fv-report.md
    
    - name: Quality Gate
      run: |
        QUALITY=$(jq -r '.stats.quality_score' fv-report.json)
        echo "::notice::Quality Score: $QUALITY"
        if [ $QUALITY -lt 80 ]; then
          echo "::error::Quality score $QUALITY is below threshold"
          exit 1
        fi
    
    - name: Comment PR
      if: github.event_name == 'pull_request'
      uses: actions/github-script@v6
      with:
        github-token: ${{secrets.GITHUB_TOKEN}}
        script: |
          const fs = require('fs');
          const report = fs.readFileSync('fv-report.md', 'utf8');
          
          github.rest.issues.createComment({
            issue_number: context.issue.number,
            owner: context.repo.owner,
            repo: context.repo.repo,
            body: `## 🔍 FV Analysis Report\n\n${report}`
          });
    
    - name: Deploy Report
      if: github.ref == 'refs/heads/main'
      uses: peaceiris/actions-gh-pages@v3
      with:
        github_token: ${{ secrets.GITHUB_TOKEN }}
        publish_dir: ./
        destination_dir: reports/${{ github.sha }}
```

---

## 总结

FV工具提供了完善的CI/CD集成支持，通过配置文件和命令行参数的组合，可以灵活适应各种CI/CD环境和需求。

**关键要点**:

1. 使用专门的CI配置文件（`.fv-ci.yaml`）
2. 在CI环境中禁用颜色和进度条
3. 合理设置质量门槛
4. 充分利用缓存机制
5. 自动化报告发布和通知

**下一步**: 查看[实战教程](Tutorial.md)了解更多使用场景！
