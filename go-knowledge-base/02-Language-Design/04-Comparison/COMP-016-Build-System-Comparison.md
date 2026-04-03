# Build System and Dependency Management Comparison

## Executive Summary

Build systems and package managers are critical for software development workflow. This document compares Go Modules, Cargo, npm, Maven/Gradle, pip, and other build tools across features, performance, and developer experience.

---

## Table of Contents

- [Build System and Dependency Management Comparison](#build-system-and-dependency-management-comparison)
  - [Executive Summary](#executive-summary)
  - [Table of Contents](#table-of-contents)
  - [Go Modules](#go-modules)
  - [Rust Cargo](#rust-cargo)
  - [JavaScript npm/yarn/pnpm](#javascript-npmyarnpnpm)
  - [Java Maven/Gradle](#java-mavengradle)
  - [Python pip/poetry](#python-pippoetry)
  - [C++ CMake](#c-cmake)
  - [Comparison Matrix](#comparison-matrix)

---

## Go Modules

Go's official dependency management since Go 1.11:

```go
// go.mod - Module definition
module github.com/example/myproject

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/stretchr/testify v1.8.4
    golang.org/x/crypto v0.15.0
)

require (
    github.com/bytedance/sonic v1.9.1 // indirect
    github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
    // ... transitive dependencies
)

// go.sum - Cryptographic checksums
github.com/gin-gonic/gin v1.9.1 h1:qU9TgL/...=
github.com/gin-gonic/gin v1.9.1/go.mod h1:hV...=
```

```bash
# Common commands
go mod init github.com/example/myproject    # Initialize module
go get github.com/gin-gonic/gin@latest      # Add dependency
go get github.com/gin-gonic/gin@v1.9.1      # Specific version
go get -u ./...                              # Update all dependencies
go mod tidy                                  # Remove unused, add missing
go mod download                              # Download dependencies
go mod vendor                                # Create vendor directory
go mod verify                                # Verify checksums
go list -m all                               # List all dependencies

# Build commands
go build ./...                              # Build all packages
go build -o myapp cmd/app/main.go          # Build specific output
go test ./...                               # Run tests
go test -race ./...                         # Run with race detector
go test -cover ./...                        # With coverage

# Environment variables
export GOPROXY=https://proxy.golang.org,direct  # Module proxy
export GOSUMDB=sum.golang.org                    # Checksum database
export GOPRIVATE=*.internal.com                  # Private repos
export GONOPROXY=*.internal.com
export GONOSUMDB=*.internal.com
```

**Key Features:**

- Minimal configuration (just go.mod)
- Automatic version selection (MVS - Minimal Version Selection)
- Cryptographic verification
- Module proxy support
- Vendor support for reproducible builds
- No central registry (uses VCS)

---

## Rust Cargo

Cargo is Rust's build system and package manager:

```toml
# Cargo.toml
[package]
name = "myproject"
version = "1.0.0"
edition = "2021"
authors = ["Your Name <you@example.com>"]
description = "A sample project"
license = "MIT"
repository = "https://github.com/example/myproject"
rust-version = "1.70"

[dependencies]
serde = { version = "1.0", features = ["derive"] }
tokio = { version = "1", features = ["full"] }
reqwest = { version = "0.11", default-features = false, features = ["json"] }

[dependencies.sqlx]
version = "0.7"
optional = true

[dev-dependencies]
tokio-test = "0.4"
mockall = "0.12"

[features]
default = ["database"]
database = ["sqlx"]
async = []

[profile.release]
opt-level = 3
lto = true
codegen-units = 1
panic = "abort"

[profile.dev]
opt-level = 0
debug = true
```

```bash
# Common commands
cargo init                    # Initialize new project
cargo new myproject           # Create new project directory
cargo add serde               # Add dependency
cargo add --dev tokio-test    # Add dev dependency
cargo add sqlx --features runtime-tokio
cargo update                  # Update dependencies
cargo update -p serde         # Update specific package
cargo build                   # Build debug
cargo build --release         # Build optimized release
cargo test                    # Run tests
cargo test --features database # Run with feature
cargo check                   # Fast syntax/type check
cargo clippy                  # Run linter
cargo fmt                     # Format code
cargo doc --open              # Generate and open docs
cargo run                     # Build and run
cargo run --release           # Run optimized
cargo bench                   # Run benchmarks
cargo publish                 # Publish to crates.io
cargo install ripgrep         # Install binary crate
cargo tree                    # Show dependency tree
cargo audit                   # Check for security advisories
```

**Key Features:**

- Comprehensive build configuration
- Feature flags for conditional compilation
- Built-in testing and benchmarking
- Documentation generation
- Workspace support
- Cargo.lock for reproducible builds
- crates.io registry

---

## JavaScript npm/yarn/pnpm

```json
{
  "name": "myproject",
  "version": "1.0.0",
  "description": "A sample project",
  "main": "index.js",
  "type": "module",
  "scripts": {
    "build": "tsc && vite build",
    "dev": "vite",
    "test": "vitest",
    "test:coverage": "vitest --coverage",
    "lint": "eslint src --ext .ts,.tsx",
    "format": "prettier --write src"
  },
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "axios": "~1.6.0"
  },
  "devDependencies": {
    "@types/react": "^18.2.0",
    "@vitejs/plugin-react": "^4.2.0",
    "typescript": "^5.3.0",
    "vite": "^5.0.0",
    "vitest": "^1.0.0",
    "eslint": "^8.0.0",
    "prettier": "^3.0.0"
  },
  "engines": {
    "node": ">=18.0.0"
  },
  "packageManager": "pnpm@8.0.0"
}
```

```bash
# npm commands
npm init                      # Initialize project
npm install                   # Install dependencies
npm install react             # Add dependency
npm install --save-dev jest   # Add dev dependency
npm install react@18.2.0      # Specific version
npm update                    # Update dependencies
npm update react              # Update specific package
npm outdated                  # Check for updates
npm run build                 # Run build script
npm test                      # Run tests
npm publish                   # Publish to npm
npm link                      # Link local package
npm audit                     # Security audit
npm audit fix                 # Fix security issues

# yarn commands
yarn init
yarn install
yarn add react
yarn add --dev jest
yarn upgrade
yarn build
yarn test

# pnpm commands (faster, disk efficient)
pnpm init
pnpm install
pnpm add react
pnpm add -D jest
pnpm update
pnpm run build
pnpm test
```

**Key Features:**

- Mature ecosystem with 2M+ packages
- Lock files for reproducibility
- Semantic versioning with `^` and `~`
- npm scripts for task automation
- Workspaces for monorepos
- Private registry support

---

## Java Maven/Gradle

```xml
<!-- pom.xml (Maven) -->
<?xml version="1.0" encoding="UTF-8"?>
<project>
    <modelVersion>4.0.0</modelVersion>

    <groupId>com.example</groupId>
    <artifactId>myproject</artifactId>
    <version>1.0.0-SNAPSHOT</version>
    <packaging>jar</packaging>

    <properties>
        <maven.compiler.source>21</maven.compiler.source>
        <maven.compiler.target>21</maven.compiler.target>
        <project.build.sourceEncoding>UTF-8</project.build.sourceEncoding>
    </properties>

    <dependencies>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-web</artifactId>
            <version>3.2.0</version>
        </dependency>

        <dependency>
            <groupId>org.junit.jupiter</groupId>
            <artifactId>junit-jupiter</artifactId>
            <version>5.10.0</version>
            <scope>test</scope>
        </dependency>
    </dependencies>

    <build>
        <plugins>
            <plugin>
                <groupId>org.springframework.boot</groupId>
                <artifactId>spring-boot-maven-plugin</artifactId>
            </plugin>
        </plugins>
    </build>
</project>
```

```groovy
// build.gradle (Gradle)
plugins {
    id 'java'
    id 'org.springframework.boot' version '3.2.0'
    id 'io.spring.dependency-management' version '1.1.4'
}

group = 'com.example'
version = '1.0.0-SNAPSHOT'

java {
    sourceCompatibility = '21'
}

repositories {
    mavenCentral()
}

dependencies {
    implementation 'org.springframework.boot:spring-boot-starter-web'
    implementation 'org.springframework.boot:spring-boot-starter-data-jpa'

    runtimeOnly 'com.h2database:h2'

    testImplementation 'org.springframework.boot:spring-boot-starter-test'
    testImplementation 'org.junit.jupiter:junit-jupiter'
}

tasks.named('test') {
    useJUnitPlatform()
}

// Custom tasks
task customTask(type: JavaExec) {
    main = 'com.example.Main'
    classpath = sourceSets.main.runtimeClasspath
}
```

```bash
# Maven commands
mvn archetype:generate          # Create new project
mvn compile                     # Compile source
mvn test                        # Run tests
mvn package                     # Build JAR/WAR
mvn install                     # Install to local repo
mvn deploy                      # Deploy to remote repo
mvn clean                       # Clean build artifacts
mvn dependency:tree             # Show dependency tree
mvn dependency:resolve          # Download dependencies
mvn versions:display-updates    # Check for updates
mvn spring-boot:run             # Run Spring Boot app

# Gradle commands
gradle init                     # Initialize project
gradle build                    # Full build
gradle test                     # Run tests
gradle run                      # Run application
gradle bootRun                  # Run Spring Boot
gradle dependencies             # Show dependency tree
gradle dependencyUpdates        # Check for updates
gradle clean                    # Clean build
gradle wrapper --version        # Update wrapper
```

---

## Python pip/poetry

```txt
# requirements.txt (pip)
flask==2.3.0
requests>=2.31.0,<3.0.0
numpy>=1.24.0
pytest~=7.4.0
black>=23.0.0; extra == "dev"
```

```toml
# pyproject.toml (Poetry)
[tool.poetry]
name = "myproject"
version = "1.0.0"
description = "A sample project"
authors = ["Your Name <you@example.com>"]
readme = "README.md"
license = "MIT"
python = "^3.11"

[tool.poetry.dependencies]
fastapi = "^0.104.0"
uvicorn = {extras = ["standard"], version = "^0.24.0"}
pydantic = "^2.5.0"
sqlalchemy = {version = "^2.0.0", optional = true}

[tool.poetry.group.dev.dependencies]
pytest = "^7.4.0"
pytest-cov = "^4.1.0"
black = "^23.0.0"
mypy = "^1.7.0"
ruff = "^0.1.0"

[tool.poetry.extras]
database = ["sqlalchemy"]

[build-system]
requires = ["poetry-core"]
build-backend = "poetry.core.masonry.api"
```

```bash
# pip commands
pip install -r requirements.txt         # Install from file
pip install flask                       # Install package
pip install flask==2.3.0                # Specific version
pip install -e .                        # Editable install
pip list                                # List packages
pip freeze > requirements.txt           # Generate requirements
pip check                               # Check for conflicts
pip audit                               # Security audit

# Poetry commands
poetry init                             # Initialize project
poetry add fastapi                      # Add dependency
poetry add --group dev pytest           # Add dev dependency
poetry add sqlalchemy --optional        # Add optional
poetry install                          # Install dependencies
poetry install --with dev               # With dev group
poetry install --extras database        # With extras
poetry update                           # Update all
poetry update fastapi                   # Update specific
poetry lock                             # Update lock file
poetry run pytest                       # Run in venv
poetry shell                            # Enter venv shell
poetry build                            # Build wheel/sdist
poetry publish                          # Publish to PyPI
```

---

## C++ CMake

```cmake
# CMakeLists.txt
cmake_minimum_required(VERSION 3.20)
project(myproject VERSION 1.0.0 LANGUAGES CXX)

# C++ standard
set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# Find packages
find_package(Boost 1.80 REQUIRED COMPONENTS system filesystem)
find_package(Threads REQUIRED)

# FetchContent for dependencies
include(FetchContent)
FetchContent_Declare(
    googletest
    GIT_REPOSITORY https://github.com/google/googletest.git
    GIT_TAG v1.14.0
)
FetchContent_MakeAvailable(googletest)

# Library target
add_library(mylib STATIC
    src/mylib.cpp
    src/utils.cpp
)

target_include_directories(mylib
    PUBLIC
        $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
        $<INSTALL_INTERFACE:include>
    PRIVATE
        ${CMAKE_CURRENT_SOURCE_DIR}/src
)

target_link_libraries(mylib
    PUBLIC Boost::system
    PRIVATE Threads::Threads
)

target_compile_options(mylib PRIVATE
    -Wall -Wextra -Wpedantic
    $<$<CONFIG:Release>:-O3>
    $<$<CONFIG:Debug>:-O0 -g>
)

# Executable target
add_executable(myapp src/main.cpp)
target_link_libraries(myapp PRIVATE mylib)

# Testing
enable_testing()
add_executable(mytest test/test_main.cpp)
target_link_libraries(mytest PRIVATE mylib gtest_main)
add_test(NAME mytest COMMAND mytest)

# Installation
install(TARGETS mylib myapp
    LIBRARY DESTINATION lib
    ARCHIVE DESTINATION lib
    RUNTIME DESTINATION bin
    INCLUDES DESTINATION include
)

install(DIRECTORY include/ DESTINATION include)
```

```bash
# CMake workflow
mkdir build && cd build
cmake ..                                    # Configure
cmake -DCMAKE_BUILD_TYPE=Release ..        # Release build
cmake --build .                             # Build
cmake --build . --parallel 4                # Parallel build
cmake --build . --target myapp             # Build specific target
cmake --install . --prefix /usr/local      # Install
ctest                                       # Run tests
cpack                                       # Create package

# With presets (CMake 3.19+)
cmake --preset=release                      # Use preset
cmake --build --preset=release
cmake --test --preset=release
```

---

## Comparison Matrix

| Feature | Go | Cargo | npm | Maven | Gradle | Poetry | CMake |
|---------|-----|-------|-----|-------|--------|--------|-------|
| Configuration | go.mod | Cargo.toml | package.json | pom.xml | build.gradle | pyproject.toml | CMakeLists.txt |
| Lock File | go.sum | Cargo.lock | package-lock.json | N/A | gradle.lockfile | poetry.lock | N/A |
| Central Registry | No (VCS) | crates.io | npmjs.com | Maven Central | Maven Central | PyPI | N/A |
| Version Resolution | MVS | SemVer | SemVer | Nearest | Dynamic | SemVer | N/A |
| Build Tool | Built-in | Built-in | External | Built-in | Built-in | External | Built-in |
| Test Runner | Built-in | Built-in | External | Built-in | Built-in | External | CTest |
| Workspace Support | Yes | Yes | Yes | Yes | Yes | Limited | Yes |
| Compile Time | Fast | Moderate | N/A | Slow | Slow | N/A | Slow |
| Binary Caching | Module cache | Target dir | node_modules | Local repo | Build cache | Virtual env | Build dir |
| Reproducibility | High | High | High | Medium | Medium | High | Medium |

---

*Document Version: 1.0*
*Last Updated: 2026-04-03*
*Size: ~20KB*
