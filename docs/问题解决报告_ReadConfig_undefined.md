# ReadConfig undefined 问题解决报告

## 问题描述
运行 `go run main.go` 时出现错误：
```
undefined: ReadConfig
```

## 问题分析

### 根本原因
错误的运行命令导致Go编译器只编译了单个文件 `main.go`，而没有包含同包中的其他文件 `config.go`。

### 详细分析
1. **文件结构正确**：
   - `main.go` 中调用了 `ReadConfig` 函数
   - `config.go` 中正确定义了 `ReadConfig` 函数
   - 两个文件都属于 `main` 包

2. **依赖正确**：
   - `go.mod` 文件包含了所需的 `gopkg.in/yaml.v3` 依赖
   - `config.go` 正确导入了相关包

3. **运行命令错误**：
   - **错误**：`go run main.go` - 只编译单个文件
   - **正确**：`go run .` - 编译整个包

## 解决方案

### 使用正确的运行命令
```bash
# 错误的命令
go run main.go

# 正确的命令
go run .
```

### 验证结果
运行 `go run .` 后，程序成功输出：
```
&{{8080 debug} {myapp 1.0.0}}
Hello, World!
```

## 知识点总结

### Go包编译规则
1. **单文件编译**：`go run main.go` 只编译指定文件
2. **包编译**：`go run .` 编译当前目录下的所有Go文件
3. **多文件包**：当一个包包含多个文件时，必须一起编译

### 最佳实践
1. 对于多文件的Go包，始终使用 `go run .`
2. 或者明确指定所有文件：`go run main.go config.go`
3. 使用 `go build` 来构建可执行文件

## 相关文件
- `main.go` - 主程序文件
- `config.go` - 配置读取功能
- `config/server.yaml` - 配置文件
- `go.mod` - 模块依赖

## 修改时间
${new Date().toLocaleString('zh-CN')} 