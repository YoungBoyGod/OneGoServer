# =============================================================================
# OneGoServer PowerShell 构建脚本
# =============================================================================

param(
    [string]$Command = "help",
    [string]$Platform = "",
    [switch]$Verbose = $false
)

# 项目信息
$PROJECT_NAME = "OneGoServer"
$BINARY_NAME = "onego"
$PACKAGE = "github.com/YoungBoyGod/OneGoServer"

# 版本信息
try {
    $VERSION = & git describe --tags --always --dirty 2>$null
    if ($LASTEXITCODE -ne 0) { $VERSION = "v0.1.0" }
} catch {
    $VERSION = "v0.1.0"
}

$BUILD_TIME = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")

try {
    $GIT_COMMIT = & git rev-parse HEAD 2>$null
    if ($LASTEXITCODE -ne 0) { $GIT_COMMIT = "unknown" }
} catch {
    $GIT_COMMIT = "unknown"
}

try {
    $GIT_BRANCH = & git rev-parse --abbrev-ref HEAD 2>$null
    if ($LASTEXITCODE -ne 0) { $GIT_BRANCH = "unknown" }
} catch {
    $GIT_BRANCH = "unknown"
}

$GO_VERSION = (& go version).Split(' ')[2]
$BUILD_BY = "$env:USERNAME@$env:COMPUTERNAME"

# 目录
$BUILD_DIR = "bin"
$DOCKER_DIR = "deployments/docker"

# Docker信息
$DOCKER_REGISTRY = "youngboygod"
$DOCKER_IMAGE = "$DOCKER_REGISTRY/$BINARY_NAME"
$DOCKER_TAG = $VERSION

# 构建标志
$LDFLAGS = "-ldflags `"-X '$PACKAGE/cmd.Version=$VERSION' -X '$PACKAGE/cmd.BuildTime=$BUILD_TIME' -X '$PACKAGE/cmd.GitCommit=$GIT_COMMIT' -X '$PACKAGE/cmd.GitBranch=$GIT_BRANCH' -X '$PACKAGE/cmd.GoVersion=$GO_VERSION' -X '$PACKAGE/cmd.BuildBy=$BUILD_BY' -w -s`""

# =============================================================================
# 工具函数
# =============================================================================

function Write-Success($message) {
    Write-Host "✅ $message" -ForegroundColor Green
}

function Write-Info($message) {
    Write-Host "📋 $message" -ForegroundColor Cyan
}

function Write-Warning($message) {
    Write-Host "⚠️  $message" -ForegroundColor Yellow
}

function Write-Error($message) {
    Write-Host "❌ $message" -ForegroundColor Red
}

function Write-Progress($message) {
    Write-Host "🚀 $message" -ForegroundColor Blue
}

# =============================================================================
# 构建函数
# =============================================================================

function Show-Help {
    Write-Host "OneGoServer PowerShell 构建脚本" -ForegroundColor Yellow
    Write-Host "==============================" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "用法: .\build.ps1 [-Command <command>] [-Platform <platform>] [-Verbose]"
    Write-Host ""
    Write-Host "可用命令:"
    Write-Host "  help            显示此帮助信息"
    Write-Host "  info            显示构建信息"
    Write-Host "  build           构建本地版本"
    Write-Host "  build-windows   构建Windows版本"
    Write-Host "  build-linux     构建Linux版本"
    Write-Host "  build-darwin    构建MacOS版本"
    Write-Host "  build-all       交叉编译所有平台"
    Write-Host "  clean           清理构建文件"
    Write-Host "  deps            更新依赖"
    Write-Host "  test            运行测试"
    Write-Host "  dev             开发模式运行"
    Write-Host "  docker          构建Docker镜像"
    Write-Host "  version         显示版本信息"
    Write-Host ""
    Write-Host "示例:"
    Write-Host "  .\build.ps1 build"
    Write-Host "  .\build.ps1 build-all"
    Write-Host "  .\build.ps1 docker"
    Write-Host "  .\build.ps1 test -Verbose"
}

function Show-Info {
    Write-Info "构建信息"
    Write-Host "================"
    Write-Host "项目名称: $PROJECT_NAME"
    Write-Host "二进制名: $BINARY_NAME"
    Write-Host "版本号:   $VERSION"
    Write-Host "构建时间: $BUILD_TIME"
    Write-Host "Git提交:  $GIT_COMMIT"
    Write-Host "Git分支:  $GIT_BRANCH"
    Write-Host "Go版本:   $GO_VERSION"
    Write-Host "构建者:   $BUILD_BY"
    Write-Host "Docker镜像: $DOCKER_IMAGE`:$DOCKER_TAG"
}

function Update-Deps {
    Write-Progress "更新依赖..."
    & go mod tidy
    & go mod download
    if ($LASTEXITCODE -eq 0) {
        Write-Success "依赖更新完成"
    } else {
        Write-Error "依赖更新失败"
        exit 1
    }
}

function Clean-Build {
    Write-Progress "清理构建文件..."
    if (Test-Path $BUILD_DIR) {
        Remove-Item -Recurse -Force $BUILD_DIR
    }
    if (Test-Path "release") {
        Remove-Item -Recurse -Force "release"
    }
    if (Test-Path "coverage.out") {
        Remove-Item -Force "coverage.out"
    }
    if (Test-Path "coverage.html") {
        Remove-Item -Force "coverage.html"
    }
    Write-Success "清理完成"
}

function Build-Local {
    Write-Progress "构建 $BINARY_NAME $VERSION..."
    Clean-Build
    Update-Deps
    
    if (!(Test-Path $BUILD_DIR)) {
        New-Item -ItemType Directory -Path $BUILD_DIR | Out-Null
    }
    
    $output = Join-Path $BUILD_DIR "$BINARY_NAME.exe"
    $buildCmd = "go build $LDFLAGS -o `"$output`" ./main.go"
    
    if ($Verbose) {
        Write-Host "执行命令: $buildCmd"
    }
    
    Invoke-Expression $buildCmd
    
    if ($LASTEXITCODE -eq 0) {
        Write-Success "构建完成: $output"
    } else {
        Write-Error "构建失败"
        exit 1
    }
}

function Build-Windows {
    Write-Progress "构建Windows版本..."
    Clean-Build
    Update-Deps
    
    if (!(Test-Path $BUILD_DIR)) {
        New-Item -ItemType Directory -Path $BUILD_DIR | Out-Null
    }
    
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    $output = Join-Path $BUILD_DIR "$BINARY_NAME-windows-amd64.exe"
    
    & go build $LDFLAGS.Split(' ') -o $output ./main.go
    
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Windows版本构建完成: $output"
    } else {
        Write-Error "Windows版本构建失败"
        exit 1
    }
    
    # 重置环境变量
    Remove-Item Env:GOOS -ErrorAction SilentlyContinue
    Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
}

function Build-Linux {
    Write-Progress "构建Linux版本..."
    
    if (!(Test-Path $BUILD_DIR)) {
        New-Item -ItemType Directory -Path $BUILD_DIR | Out-Null
    }
    
    $env:GOOS = "linux"
    $env:GOARCH = "amd64"
    $output = Join-Path $BUILD_DIR "$BINARY_NAME-linux-amd64"
    
    & go build $LDFLAGS.Split(' ') -o $output ./main.go
    
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Linux版本构建完成: $output"
    } else {
        Write-Error "Linux版本构建失败"
        exit 1
    }
    
    # 重置环境变量
    Remove-Item Env:GOOS -ErrorAction SilentlyContinue
    Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
}

function Build-Darwin {
    Write-Progress "构建MacOS版本..."
    
    if (!(Test-Path $BUILD_DIR)) {
        New-Item -ItemType Directory -Path $BUILD_DIR | Out-Null
    }
    
    # MacOS AMD64
    $env:GOOS = "darwin"
    $env:GOARCH = "amd64"
    $output = Join-Path $BUILD_DIR "$BINARY_NAME-darwin-amd64"
    & go build $LDFLAGS.Split(' ') -o $output ./main.go
    
    # MacOS ARM64
    $env:GOARCH = "arm64"
    $output = Join-Path $BUILD_DIR "$BINARY_NAME-darwin-arm64"
    & go build $LDFLAGS.Split(' ') -o $output ./main.go
    
    if ($LASTEXITCODE -eq 0) {
        Write-Success "MacOS版本构建完成"
    } else {
        Write-Error "MacOS版本构建失败"
        exit 1
    }
    
    # 重置环境变量
    Remove-Item Env:GOOS -ErrorAction SilentlyContinue
    Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
}

function Build-All {
    Write-Progress "交叉编译所有平台..."
    Build-Windows
    Build-Linux
    Build-Darwin
    Write-Success "所有平台构建完成!"
    Get-ChildItem $BUILD_DIR
}

function Run-Tests {
    Write-Progress "运行测试..."
    & go test -v ./...
    if ($LASTEXITCODE -eq 0) {
        Write-Success "测试完成"
    } else {
        Write-Error "测试失败"
        exit 1
    }
}

function Start-Dev {
    Write-Progress "开发模式启动..."
    Update-Deps
    & go run main.go server --verbose
}

function Build-Docker {
    Write-Progress "构建Docker镜像 $DOCKER_IMAGE`:$DOCKER_TAG..."
    
    $buildArgs = @(
        "--build-arg", "VERSION=$VERSION",
        "--build-arg", "BUILD_TIME=$BUILD_TIME",
        "--build-arg", "GIT_COMMIT=$GIT_COMMIT",
        "--build-arg", "GIT_BRANCH=$GIT_BRANCH",
        "--build-arg", "GO_VERSION=$GO_VERSION",
        "--build-arg", "BUILD_BY=$BUILD_BY",
        "-t", "$DOCKER_IMAGE`:$DOCKER_TAG",
        "-t", "$DOCKER_IMAGE`:latest",
        "-f", "$DOCKER_DIR/Dockerfile",
        "."
    )
    
    & docker build @buildArgs
    
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Docker镜像构建完成: $DOCKER_IMAGE`:$DOCKER_TAG"
    } else {
        Write-Error "Docker镜像构建失败"
        exit 1
    }
}

function Show-Version {
    Build-Local
    $binaryPath = Join-Path $BUILD_DIR "$BINARY_NAME.exe"
    & $binaryPath version --verbose
}

# =============================================================================
# 主函数
# =============================================================================

switch ($Command.ToLower()) {
    "help" { Show-Help }
    "info" { Show-Info }
    "build" { Build-Local }
    "build-windows" { Build-Windows }
    "build-linux" { Build-Linux }
    "build-darwin" { Build-Darwin }
    "build-all" { Build-All }
    "clean" { Clean-Build }
    "deps" { Update-Deps }
    "test" { Run-Tests }
    "dev" { Start-Dev }
    "docker" { Build-Docker }
    "version" { Show-Version }
    default { 
        Write-Error "未知命令: $Command"
        Show-Help
        exit 1
    }
} 