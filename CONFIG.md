# OneGoServer 配置说明

## 配置文件

项目使用 GoFrame 框架，支持多种配置文件格式。推荐使用 YAML 格式。

### 配置文件位置

GoFrame 会按以下顺序查找配置文件：

1. `config.yaml` (项目根目录)
2. `config.dev.yaml` (开发环境)
3. `config.prod.yaml` (生产环境)
4. `manifest/config/config.yaml`

### 快速开始

1. **复制配置模板**
   ```bash
   cp config.template.yaml config.yaml
   ```

2. **修改数据库配置**
   编辑 `config.yaml` 文件，修改数据库连接信息：
   ```yaml
   database:
     default:
       link: "pgsql:用户名:密码@tcp(主机:端口)/数据库名"
   ```

3. **运行项目**
   ```bash
   go run main.go
   ```

## 配置项说明

### 数据库配置
```yaml
database:
  default:
    link: "pgsql:admin:onegoserver@123@tcp(192.168.100.128:5432)/onegoserver"
    debug: true          # 开启调试模式
    charset: "utf8"      # 字符集
    maxIdle: 10          # 最大空闲连接数
    maxOpen: 100         # 最大打开连接数
    maxLifetime: "30s"   # 连接最大生命周期
```

### 服务器配置
```yaml
server:
  address: ":8000"           # 服务器监听地址
  openapiPath: "/api.json"   # OpenAPI 文档路径
  swaggerPath: "/swagger"    # Swagger UI 路径
  dumpRouterMap: false       # 是否显示路由映射
  graceful: true             # 优雅关闭
  gracefulTimeout: 30        # 优雅关闭超时时间
```

### 日志配置
```yaml
logger:
  level: "all"               # 日志级别
  stdout: true               # 输出到控制台
  file: "logs/server.log"    # 日志文件路径
  rotateSize: "100MB"        # 日志文件大小限制
  rotateExpire: "7d"         # 日志文件保留时间
```

## 环境变量

可以通过环境变量覆盖配置：

```bash
# 设置配置文件
export GF_GCFG_FILE=config.dev.yaml

# 设置数据库连接
export DB_HOST=192.168.100.128
export DB_PORT=5432
export DB_NAME=onegoserver
export DB_USER=admin
export DB_PASSWORD=onegoserver@123

# 设置服务器端口
export SERVER_PORT=8000
```

## 开发环境

开发环境推荐使用 `config.dev.yaml`：

```bash
# 设置开发环境配置
export GF_GCFG_FILE=config.dev.yaml
go run main.go
```

## 生产环境

生产环境建议：

1. 使用环境变量设置敏感信息
2. 关闭调试模式
3. 配置适当的日志级别
4. 设置合理的数据库连接池参数

```yaml
# config.prod.yaml
database:
  default:
    link: "${DB_LINK}"
    debug: false
    maxIdle: 20
    maxOpen: 200

server:
  address: ":8080"
  dumpRouterMap: false

logger:
  level: "error"
  stdout: false
  file: "/var/log/onegoserver/server.log"

system:
  debug: false
``` 