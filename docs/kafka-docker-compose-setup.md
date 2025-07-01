# OneGo服务器 Kafka Docker Compose 配置指南

## 概述

为OneGo服务器配置了基于Docker Compose的Kafka服务，使用KRaft模式（无需Zookeeper），适用于开发环境。

## Docker Compose 配置

### 文件：`docker-compose.kafka.yml`

```yaml
services:
  kafka:
    image: bitnami/kafka:4.0.0
    container_name: kafka
    ports:
      - "9092:9092"
    environment:
      # KRaft 模式配置（简化版）
      KAFKA_CFG_PROCESS_ROLES: "controller,broker"
      KAFKA_CFG_NODE_ID: "1"
      KAFKA_CFG_CONTROLLER_QUORUM_VOTERS: "1@localhost:9093"
      KAFKA_CFG_CONTROLLER_LISTENER_NAMES: "CONTROLLER"
      
      # 监听器配置
      KAFKA_CFG_LISTENERS: "PLAINTEXT://:9092,CONTROLLER://:9093"
      KAFKA_CFG_ADVERTISED_LISTENERS: "PLAINTEXT://localhost:9092"
      KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP: "CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT"
      KAFKA_CFG_INTER_BROKER_LISTENER_NAME: "PLAINTEXT"
      
      # 开发环境设置
      KAFKA_CFG_AUTO_CREATE_TOPICS_ENABLE: "true"
      KAFKA_CFG_NUM_PARTITIONS: "1"
      KAFKA_CFG_DEFAULT_REPLICATION_FACTOR: "1"
      KAFKA_CFG_MIN_INSYNC_REPLICAS: "1"
      
      # 固定cluster ID避免重复初始化
      KAFKA_CLUSTER_ID: "onego-kafka-cluster"
      
      # JVM优化
      KAFKA_HEAP_OPTS: "-Xmx256m -Xms256m"
      
    restart: unless-stopped
```

## 配置优化说明

### 1. KRaft 模式优势
- **无需Zookeeper**: 简化架构，减少资源消耗
- **更快启动**: 减少组件依赖，启动速度更快
- **更好稳定性**: 内置元数据管理，更稳定

### 2. 开发环境优化
- **固定Cluster ID**: 避免重复初始化数据
- **低内存配置**: JVM堆内存设置为256MB
- **自动创建Topic**: 开发时无需手动创建Topic
- **单副本配置**: 适合单节点开发环境

### 3. 网络配置
- **端口映射**: 9092端口用于客户端连接
- **监听器设置**: 支持本地和容器内访问

## 启动命令

```bash
# 启动Kafka服务
docker-compose -f docker-compose.kafka.yml up -d kafka

# 查看服务状态
docker-compose -f docker-compose.kafka.yml ps

# 查看日志
docker logs kafka

# 停止服务
docker-compose -f docker-compose.kafka.yml down
```

## 应用程序配置

### config.yaml 配置
```yaml
kafka:
  brokers: ["localhost:9092"]
  client_id: "onego-server"
  version: "2.8.1"
  enable_sasl: false
  enable_tls: false
  producer_return_successes: true
  producer_return_errors: true
  producer_required_acks: 1
  producer_retry_max: 3
  consumer_group_id: "onego-server"
  consumer_offset_initial: "earliest"
```

## 验证连接

### 1. 容器状态检查
```bash
docker logs kafka --tail 20
```

成功启动的标志：
```
[BrokerServer id=1] Transition from STARTING to STARTED
Kafka version: 4.0.0
Kafka Server started
```

### 2. 应用程序连接测试
```bash
go run main.go
```

成功连接的日志：
```json
{"level":"info","msg":"Connecting to Kafka brokers","brokers":["localhost:9092"],"client_id":"onego-server"}
{"level":"info","msg":"Kafka operation completed","type":"kafka_operation","operation":"connect"}
{"level":"info","msg":"Kafka initialized successfully","brokers":["localhost:9092"]}
```

## 故障排除

### 常见问题及解决方案

1. **容器重复重启**
   - 原因：数据卷冲突或配置错误
   - 解决：`docker-compose down -v` 清理数据卷后重新启动

2. **连接被拒绝**
   - 原因：端口未正确映射或容器未完全启动
   - 解决：等待30秒让容器完全启动，检查端口映射

3. **内存不足**
   - 原因：JVM堆内存设置过高
   - 解决：调整 `KAFKA_HEAP_OPTS` 参数

## 生产环境建议

⚠️ **当前配置仅适用于开发环境**

生产环境建议：
1. 增加副本数量 (`KAFKA_CFG_DEFAULT_REPLICATION_FACTOR`)
2. 配置持久化存储卷
3. 启用安全认证 (SASL/SSL)
4. 增加JVM内存分配
5. 配置监控和日志收集

## 容器资源使用

- **CPU**: 低负载下 < 10%
- **内存**: ~256MB JVM堆内存
- **存储**: 临时数据，重启后清空
- **网络**: 端口9092对外开放

---
*配置时间: 2025-07-01 11:03*  
*状态: ✅ 运行正常*  
*版本: Kafka 4.0.0 (KRaft模式)* 