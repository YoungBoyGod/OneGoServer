# 配置文件读取迁移到Viper

## 修改概述

本次修改将OneGoClient的配置文件读取方式从原生的`yaml.NewDecoder`迁移到使用`viper`库，提供更强大的配置管理功能。

## 修改内容

### 1. 依赖变更

**移除的依赖：**
- `os` 包（不再需要直接操作文件）
- `gopkg.in/yaml.v3` 包（由viper内部处理）

**新增的依赖：**
- `github.com/spf13/viper` 包（配置管理库）

### 2. 结构体标签更新

为所有配置结构体字段添加了`mapstructure`标签，以支持viper的自动解析：

```go
// 示例：ClientConfig结构体
type ClientConfig struct {
    Host        string `yaml:"host" mapstructure:"host"`
    Port        int    `yaml:"port" mapstructure:"port"`
    Name        string `yaml:"name" mapstructure:"name"`
    // ... 其他字段
}
```

### 3. LoadConfig函数重构

**原有实现：**
```go
func LoadConfig(path string) (*Config, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    var cfg Config
    decoder := yaml.NewDecoder(f)
    if err := decoder.Decode(&cfg); err != nil {
        return nil, err
    }
    return &cfg, nil
}
```

**新实现：**
```go
func LoadConfig(path string) (*Config, error) {
    viper.SetConfigFile(path)
    viper.SetConfigType("yaml")

    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }

    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}
```

## 优势

1. **功能增强**: viper提供更多配置来源支持（环境变量、命令行参数、远程配置等）
2. **自动重载**: 支持配置文件热重载
3. **类型安全**: 更好的类型转换和验证
4. **标准化**: 使用Go社区广泛采用的配置管理库
5. **扩展性**: 为后续添加更多配置功能奠定基础

## 影响分析

- **兼容性**: 保持原有API不变，现有调用代码无需修改
- **性能**: viper内部优化了配置读取性能
- **维护性**: 代码更简洁，易于理解和维护

## 测试建议

1. 验证现有配置文件格式兼容性
2. 测试配置文件不存在或格式错误的错误处理
3. 验证所有配置字段正确映射到结构体

## 修改时间

2024年12月19日

## 相关文件

- `OneGoClient/internal/config/config.go` - 主要修改文件
- `OneGoClient/go.mod` - 已包含viper依赖 