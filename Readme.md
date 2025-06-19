# 目录
├── cmd
│   ├── config.go # 配置管理命令
│   ├── root.go # 根命令
│   ├── server.go # 服务器命令
│   └── version.go # 版本命令
├── config
│   ├── server.dev.yaml # 开发环境配置
│   ├── server.prod.yaml # 生产环境配置
│   ├── server.test.yaml # 测试环境配置
│   └── server.yaml # 默认配置文件
├── docs
│   ├── 如何触发服务器关闭_演示指南.md
│   ├── 问题解决报告_ReadConfig_undefined.md
│   ├── 代码重构和清理报告.md
│   ├── 优雅重启功能实现报告.md
│   ├── 多环境配置使用指南.md
│   ├── 路由分离重构报告.md
│   ├── 配置管理架构建议.md
│   ├── 项目结构标准化重组报告.md
│   ├── Cobra_CLI框架集成报告.md
│   └── main.go代码详细解释.md
├── go.mod # 依赖管理文件
├── go.sum # 依赖管理文件
├── internal
│   ├── common # 公共模块
│   │   ├── errors # 错误处理
│   │   ├── logger # 日志处理
│   │   └── utils # 工具函数
│   ├── config
│   │   └── config.go # 配置管理
│   └── server 
│       ├── api
│       │   ├── handlers # 处理器
│       │   │   ├── api.go # API处理器
│       │   │   └── base.go # 基础处理器
│       │   ├── middleware # 中间件
│       │   └── routes # 路由
│       │       └── routes.go # 路由配置
│       ├── models # 模型
│       ├── repository # 仓库
│       ├── security # 安全
│       ├── server.go # 服务器
│       ├── services # 服务
│       └── tcp # TCP服务
├── learngo0619 # 项目根目录
├── main.go # 主函数
├── Makefile # 构建脚本
├── routes # 路由
└── tmp # 临时文件
    ├── build-errors.log # 构建错误日志
    └── main # 主函数