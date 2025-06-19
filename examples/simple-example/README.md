# 新手入门示例：Cobra + Gin + YAML

这是一个最简单的入门示例，适合刚开始学习Go Web开发的新手。

## 📁 文件说明

```
simple-example/
├── main.go      (主程序文件，包含所有代码)
├── config.yaml  (配置文件)
└── README.md    (本说明文件)
```

## 🚀 快速开始

### 1. 安装依赖

```bash
# 进入示例目录
cd examples/simple-example

# 初始化Go模块
go mod init simple-example

# 安装依赖包
go get github.com/gin-gonic/gin
go get github.com/spf13/cobra
go get gopkg.in/yaml.v3
```

### 2. 编译程序

```bash
go build -o simple-app main.go
```

### 3. 运行程序

```bash
# 启动服务器
./simple-app server

# 查看帮助
./simple-app --help

# 查看版本
./simple-app version

# 使用自定义配置文件
./simple-app server --config my-config.yaml
```

## 🧪 测试API

启动服务器后，你可以在浏览器中访问或使用curl测试：

```bash
# 测试首页
curl http://localhost:8080/

# 测试ping接口
curl http://localhost:8080/ping

# 查看配置信息
curl http://localhost:8080/config

# 查看帮助信息
curl http://localhost:8080/help
```

## 📝 代码结构详解

### 1. 配置结构体 (第1步)
```go
type Config struct {
    Server ServerConfig `yaml:"server"`
    App    AppConfig    `yaml:"app"`
}
```
- 定义程序需要的配置项
- `yaml:"server"` 标签告诉程序从YAML的哪个部分读取数据

### 2. 配置加载 (第2步)
```go
func LoadConfig(filename string) (*Config, error) {
    data, err := os.ReadFile(filename)
    // ...
    err = yaml.Unmarshal(data, &config)
    // ...
}
```
- 读取YAML文件内容
- 解析YAML数据到Go结构体

### 3. Gin服务器 (第3步)
```go
func CreateServer(config *Config) *gin.Engine {
    router := gin.Default()
    router.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })
    return router
}
```
- 创建HTTP路由
- 定义API接口的处理逻辑

### 4. Cobra命令 (第5步)
```go
var serverCmd = &cobra.Command{
    Use:   "server",
    Short: "启动Web服务器",
    Run: func(cmd *cobra.Command, args []string) {
        startServer()
    },
}
```
- 定义命令行子命令
- 当用户执行命令时调用相应函数

## 🔄 程序执行流程

```
用户输入命令
    ↓
./simple-app server
    ↓
Cobra解析命令 → 识别"server"子命令 → 调用startServer()
    ↓
LoadConfig() → 读取config.yaml → 解析到Config结构体
    ↓
CreateServer() → 创建Gin路由器 → 注册API路由
    ↓
router.Run() → 启动HTTP服务器 → 监听端口8080
    ↓
等待用户请求...
```

## 🎯 学习要点

### 作为新手，重点理解这几个概念：

1. **结构体标签**: `yaml:"server"` 告诉程序如何映射数据
2. **错误处理**: 每个可能出错的操作都要检查 `err`
3. **路由注册**: `router.GET()` 就是告诉程序"访问这个网址时做什么"
4. **JSON响应**: `c.JSON()` 把Go数据转换成网页能理解的格式

### 常见问题：

**Q: 为什么端口要用字符串而不是数字？**
A: 因为在网络编程中，端口通常作为字符串处理，方便拼接到地址中。

**Q: gin.H 是什么？**
A: 它是 `map[string]interface{}` 的简写，用来创建JSON响应数据。

**Q: 为什么要用 Cobra？**
A: Cobra让程序更专业，支持子命令、帮助信息、参数验证等功能。

## 🔧 实验建议

1. **修改配置文件**: 改变端口号，看服务器是否在新端口启动
2. **添加新路由**: 在`CreateServer()`中添加新的API接口
3. **修改响应内容**: 改变JSON响应的内容
4. **添加新命令**: 仿照`versionCmd`添加其他命令

## 📚 下一步学习

掌握这个示例后，你可以学习：

1. **数据库集成**: 学习如何连接MySQL、PostgreSQL等数据库
2. **中间件使用**: 学习CORS、日志记录、身份验证等中间件
3. **项目结构**: 学习如何组织更大的项目（参考主项目结构）
4. **部署发布**: 学习如何将程序部署到服务器

记住：编程是一个循序渐进的过程，先把基础打牢，再慢慢学习高级功能！ 