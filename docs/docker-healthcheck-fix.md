# Docker Compose 健康检查问题解决报告

## 🚨 问题发现

在启动kafka-ui时遇到错误：
```
dependency failed to start: container kafka has no healthcheck configured
```

## 🔍 问题分析

### 原因
- `kafka-ui` 服务配置了依赖检查：`condition: service_healthy`
- 但是 `kafka` 服务没有配置健康检查（healthcheck）
- Docker Compose 无法判断kafka是否就绪，导致依赖启动失败

### 依赖关系
```yaml
kafka-ui:
  depends_on:
    kafka:
      condition: service_healthy  # ❌ 要求kafka有健康检查
```

## ✅ 解决方案

### 1. 为Kafka添加健康检查

```yaml
kafka:
  # ... 其他配置
  healthcheck:
    test: ["CMD-SHELL", "timeout 10s bash -c '</dev/tcp/localhost/9092'"]
    interval: 15s
    timeout: 10s
    retries: 5
    start_period: 30s
```

**健康检查说明**:
- `test`: 通过TCP连接测试9092端口是否可达
- `interval`: 每15秒检查一次
- `timeout`: 单次检查超时10秒
- `retries`: 失败重试5次
- `start_period`: 启动后等待30秒再开始检查

### 2. 替代方案（如果不需要UI）

移除kafka-ui的依赖检查：
```yaml
kafka-ui:
  depends_on:
    - kafka  # ✅ 简单依赖，不要求健康检查
```

或者完全移除kafka-ui服务（仅启动kafka）：
```bash
docker-compose -f docker-compose.kafka.yml up -d kafka
```

## 🔧 最终配置

完整的工作配置（`docker-compose.kafka-fixed.yml`）：

```yaml
services:
  kafka:
    image: bitnami/kafka:4.0.0
    container_name: kafka
    ports:
      - "9092:9092"
    environment:
      # KRaft模式配置
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
      
      # 固定cluster ID
      KAFKA_CLUSTER_ID: "onego-kafka-cluster"
      
      # JVM优化
      KAFKA_HEAP_OPTS: "-Xmx256m -Xms256m"
      
    restart: unless-stopped
    healthcheck:  # ✅ 添加健康检查
      test: ["CMD-SHELL", "timeout 10s bash -c '</dev/tcp/localhost/9092'"]
      interval: 15s
      timeout: 10s
      retries: 5
      start_period: 30s

  kafka-ui:
    image: provectuslabs/kafka-ui:latest
    container_name: kafka-ui
    ports:
      - "8081:8080"
    environment:
      KAFKA_CLUSTERS_0_NAME: "onego-kafka"
      KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS: "kafka:9092"
    depends_on:
      kafka:
        condition: service_healthy  # ✅ 现在可以正常工作
    restart: unless-stopped
```

## 🎯 验证步骤

1. **启动服务**:
   ```bash
   docker-compose -f docker-compose.kafka-fixed.yml up -d
   ```

2. **检查健康状态**:
   ```bash
   docker ps
   # 应该看到 kafka (healthy) 状态
   ```

3. **测试应用连接**:
   ```bash
   go run main.go
   # 应该看到 Kafka initialized successfully
   ```

4. **访问Kafka UI** (可选):
   ```
   http://localhost:8081
   ```

## 📊 总结

| 问题 | 原因 | 解决方案 |
|------|------|----------|
| 依赖启动失败 | kafka缺少健康检查 | 添加TCP端口检测 |
| 服务状态不明 | 无法判断就绪状态 | 配置健康检查参数 |

**关键要点**:
- Docker Compose的`condition: service_healthy`要求目标服务有健康检查
- 健康检查应该测试服务的实际可用性（如端口连通性）
- 合理设置检查间隔和重试次数

---
*修复时间: 2025-07-01 11:10*  
*状态: ✅ 已解决* 