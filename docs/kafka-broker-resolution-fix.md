# Kafka Broker地址解析问题修复报告

## 问题描述

在运行Kafka集成测试时，出现了 `dial tcp: lookup kafka: no such host` 错误，但是使用命令行工具可以正常连接到 `localhost:9092`。

## 问题分析

### 错误现象
```
dial tcp: lookup kafka: no such host
```

### 调试发现
通过添加调试日志发现：

1. **初始连接成功**：使用 `localhost:9092` 成功建立连接
2. **元数据获取成功**：能够获取主题列表，说明连接正常
3. **消息发送失败**：尝试连接到 `kafka` 主机名失败

### 根本原因

这是一个典型的 **Kafka Advertised Listeners 配置问题**：

1. 客户端首次连接到 `localhost:9092`
2. Kafka服务器返回集群元数据，包含所有broker的advertised地址
3. 如果Kafka配置中使用了 `kafka` 作为advertised hostname
4. 后续操作会尝试连接到 `kafka:9092`，导致DNS解析失败

## Kafka Advertised Listeners 机制

Kafka使用 `advertised.listeners` 配置告诉客户端如何连接到broker：

```properties
# Kafka server.properties
listeners=PLAINTEXT://0.0.0.0:9092
advertised.listeners=PLAINTEXT://kafka:9092  # 问题所在！
```

当客户端连接时：
1. 连接到配置的初始broker地址
2. 获取集群元数据，包含所有broker的advertised地址
3. 后续所有操作使用advertised地址

## 解决方案

### 方案1：修改Kafka服务器配置（推荐）

#### Docker Compose环境
```yaml
version: '3'
services:
  kafka:
    image: confluentinc/cp-kafka:latest
    environment:
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://localhost:9092  # 使用localhost
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
    ports:
      - "9092:9092"
```

#### 原生Kafka配置
```properties
# server.properties
listeners=PLAINTEXT://0.0.0.0:9092
advertised.listeners=PLAINTEXT://localhost:9092  # 使用localhost而不是kafka
```

### 方案2：系统hosts映射

在 `/etc/hosts` 文件中添加：
```
127.0.0.1 kafka
```

### 方案3：代码层面容错处理（已实现）

修改测试配置，支持多个broker地址：

```go
func getTestKafkaConfig() *config.KafkaConfig {
    kafkaHost := os.Getenv("KAFKA_HOST")
    if kafkaHost == "" {
        kafkaHost = "localhost"
    }
    
    kafkaPort := os.Getenv("KAFKA_PORT")
    if kafkaPort == "" {
        kafkaPort = "9092"
    }
    
    // 支持多个broker地址，用于容错
    brokers := []string{
        fmt.Sprintf("%s:%s", kafkaHost, kafkaPort),
    }
    
    // 如果是localhost，也尝试kafka主机名（用于Docker环境）
    if kafkaHost == "localhost" {
        brokers = append(brokers, fmt.Sprintf("kafka:%s", kafkaPort))
    }
    
    return &config.KafkaConfig{
        Brokers: brokers,
        // ...
    }
}
```

## 验证命令

### 检查Kafka配置
```bash
# 检查Kafka进程配置
ps aux | grep kafka

# 检查Docker容器环境变量
docker inspect <kafka_container_id> | grep -A 20 Env
```

### 测试连接
```bash
# 使用localhost连接
kafka-console-producer.sh --bootstrap-server localhost:9092 --topic test

# 使用kafka主机名连接（应该失败）
kafka-console-producer.sh --bootstrap-server kafka:9092 --topic test
```

### 查看元数据
```bash
kafka-metadata-shell.sh --snapshot /path/to/metadata.log
```

## 最佳实践

### 1. 生产环境配置
```properties
# 使用FQDN或IP地址
advertised.listeners=PLAINTEXT://kafka.example.com:9092
```

### 2. 开发环境配置
```properties
# 使用localhost便于本地开发
advertised.listeners=PLAINTEXT://localhost:9092
```

### 3. Docker环境配置
```yaml
environment:
  # 根据访问方式选择
  KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://localhost:9092  # 主机访问
  # 或
  KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092      # 容器间访问
```

### 4. 混合环境配置
```yaml
environment:
  # 支持多种访问方式
  KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://localhost:9092,INTERNAL://kafka:9093
  KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,INTERNAL:PLAINTEXT
  KAFKA_INTER_BROKER_LISTENER_NAME: INTERNAL
```

## 监控建议

### 1. 添加连接日志
```go
// 在connect方法中添加
pkglog.LogInfo("Connecting to Kafka brokers",
    zap.Strings("brokers", km.config.Brokers),
    zap.String("client_id", km.config.ClientID))
```

### 2. 监控DNS解析
```bash
# 监控DNS查询
sudo tcpdump -i any port 53 | grep kafka
```

### 3. 网络连接跟踪
```bash
# 跟踪网络连接
netstat -an | grep 9092
```

## 总结

这个问题的核心是Kafka的advertised listeners配置与客户端期望的地址不匹配。通过正确配置Kafka服务器的advertised listeners，可以从根本上解决这个问题。代码层面的容错处理可以作为补充方案，提高系统的健壮性。 