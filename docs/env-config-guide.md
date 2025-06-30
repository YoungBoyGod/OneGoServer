# 环境变量配置指南

## 概述
本指南说明如何配置OneGoServer以从.env文件读取数据库配置。

## 安装依赖
首先需要安装godotenv包：
```bash
go get github.com/joho/godotenv
```

## 创建.env文件
在项目根目录创建.env文件，包含以下配置：

```env
# 数据库配置环境变量
DB_TYPE=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=admin
DB_PASSWORD=onegoserver@123
DB_NAME=onegoserver
DB_CHARSET=utf8
DB_MAX_IDLE_CONNS=10
DB_MAX_OPEN_CONNS=100
DB_CONN_MAX_LIFETIME=1h
DB_AUTO_MIGRATE=true

# 服务器配置
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_MODE=debug
```

## 配置优先级
1. 环境变量（从.env文件加载）
2. config.yaml文件中的默认值

## 使用方法
环境变量会覆盖yaml配置文件中的相应值，使得在不同环境中部署时更加灵活。

## 修改内容
- 更新了配置结构体以匹配config.yaml格式
- 添加了从.env文件读取配置的功能
- 增强了环境变量绑定支持更多数据库配置项 