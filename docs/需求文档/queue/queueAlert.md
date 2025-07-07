# 队列告警需求文档

## 1. 功能描述

### 1.1 功能概述
队列告警功能用于监控队列运行状态、任务异常、资源瓶颈等情况，自动触发告警并通知相关人员，支持告警查询、筛选、处理和统计。

### 1.2 主要功能列表
- 队列异常自动告警
- 告警信息查询与筛选
- 告警通知与处理
- 告警统计与导出

### 1.3 支持的功能特性
- 多类型告警支持
- 告警自动触发与通知
- 告警处理流程可追溯

## 2. 功能目标

### 2.1 业务目标
- 提高队列异常响应速度
- 降低运维风险
- 满足合规与审计需求

### 2.2 技术目标
- 实时高效告警检测
- 告警数据一致性保障
- 支持大规模队列监控

### 2.3 安全目标
- 告警信息访问权限控制
- 敏感信息保护
- 告警操作可审计

## 3. 输入输出说明

### 3.1 输入参数

#### 3.1.1 必需参数
- `queue_id` (string): 队列ID

#### 3.1.2 可选参数
- `alert_type` (string): 告警类型（status/error/threshold等）
- `status` (string): 告警状态（active/handled/closed等）
- `start_time` (string): 查询起始时间
- `end_time` (string): 查询结束时间
- `keyword` (string): 关键字
- `page` (int): 页码
- `page_size` (int): 每页数量

### 3.2 输出参数

#### 3.2.1 成功响应
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "alerts": [
      {
        "alert_id": "alert_001",
        "queue_id": "queue_001",
        "alert_type": "error",
        "status": "active",
        "message": "队列阻塞异常",
        "created_at": "2024-01-16T18:20:00Z",
        "handled_at": null
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 20
  }
}
```

#### 3.2.2 错误响应
```json
{
  "code": 400,
  "message": "参数错误或无权限",
  "data": null
}
```

### 3.3 参数格式和约束
- 队列ID：1-64字符，字母、数字、下划线
- 告警类型：status、error、threshold等
- 分页：page>=1, page_size<=100

## 4. 涉及接口以及接口设计

### 4.1 API接口定义

#### 4.1.1 查询队列告警
```go
// 查询队列告警
GET /api/v1/queue/{queue_id}/alerts
```

**请求参数：**
- Path参数：queue_id
- Query参数：alert_type, status, start_time, end_time, keyword, page, page_size

**响应结构：**
```go
type QueueAlertListResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    struct {
        Alerts   []QueueAlert `json:"alerts"`
        Total    int          `json:"total"`
        Page     int          `json:"page"`
        PageSize int          `json:"page_size"`
    } `json:"data"`
}

type QueueAlert struct {
    AlertID    string  `json:"alert_id"`
    QueueID    string  `json:"queue_id"`
    AlertType  string  `json:"alert_type"`
    Status     string  `json:"status"`
    Message    string  `json:"message"`
    CreatedAt  string  `json:"created_at"`
    HandledAt  *string `json:"handled_at"`
}
```

### 4.2 内部接口设计

#### 4.2.1 告警服务接口
```go
type QueueAlertService interface {
    ListAlerts(ctx context.Context, queueID string, filter QueueAlertListFilter) ([]QueueAlert, int, error)
}
```

#### 4.2.2 告警仓储接口
```go
type QueueAlertRepository interface {
    Query(ctx context.Context, queueID string, filter QueueAlertListFilter) ([]QueueAlert, int, error)
}
```

## 5. 数据结构设计

### 5.1 数据库表结构
- 新增 queue_alert 表

#### 5.1.1 队列告警表 (queue_alert)
```sql
CREATE TABLE queue_alert (
    alert_id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '告警ID',
    queue_id VARCHAR(64) NOT NULL COMMENT '队列ID',
    alert_type VARCHAR(32) NOT NULL COMMENT '告警类型',
    status VARCHAR(16) NOT NULL COMMENT '告警状态',
    message TEXT NOT NULL COMMENT '告警内容',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    handled_at TIMESTAMP NULL COMMENT '处理时间',
    INDEX idx_queue_id (queue_id),
    INDEX idx_alert_type (alert_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='队列告警表';
```

### 5.2 模型结构定义
```go
type QueueAlertListFilter struct {
    AlertType  string
    Status     string
    StartTime  string
    EndTime    string
    Keyword    string
    Page       int
    PageSize   int
}

type QueueAlert struct {
    AlertID    string  `json:"alert_id" db:"alert_id"`
    QueueID    string  `json:"queue_id" db:"queue_id"`
    AlertType  string  `json:"alert_type" db:"alert_type"`
    Status     string  `json:"status" db:"status"`
    Message    string  `json:"message" db:"message"`
    CreatedAt  string  `json:"created_at" db:"created_at"`
    HandledAt  *string `json:"handled_at" db:"handled_at"`
}
```

### 5.3 数据关系说明
- 告警与队列通过queue_id关联

## 6. 异常处理

### 6.1 输入验证异常
- 队列ID为空或格式错误
- 分页参数非法

### 6.2 业务逻辑异常
- 队列不存在
- 权限不足

### 6.3 系统异常
- 数据库连接失败
- 查询超时
- 系统资源不足

## 7. 交互流程图

```mermaid
flowchart TD
    A[用户请求队列告警] --> B{验证请求参数}
    B -->|参数错误| C[返回参数错误]
    B -->|参数正确| D[验证用户权限]
    D -->|权限不足| E[返回权限错误]
    D -->|权限正确| F[查询告警]
    F --> G{队列是否存在}
    G -->|不存在| H[返回队列不存在]
    G -->|存在| I[返回告警列表]
    
    style A fill:#e1f5fe
    style I fill:#c8e6c9
    style C fill:#ffcdd2
    style E fill:#ffcdd2
    style H fill:#ffcdd2
```

## 8. 交互时序图

```mermaid
sequenceDiagram
    participant Client as 客户端
    participant API as API网关
    participant Controller as 控制器
    participant Service as 告警服务
    participant Repository as 数据仓储
    participant DB as 数据库
    
    Client->>API: GET /api/v1/queue/{queue_id}/alerts
    API->>Controller: 路由到告警控制器
    Controller->>Controller: 验证请求参数
    Controller->>Service: 调用告警服务
    Service->>Repository: 查询告警
    Repository->>DB: 执行SQL查询
    DB-->>Repository: 返回告警数据
    Repository-->>Service: 返回告警列表
    Service-->>Controller: 返回处理结果
    Controller-->>API: 返回响应数据
    API-->>Client: 返回告警列表
```

## 9. 安全与权限

### 9.1 访问控制策略
- 需要用户身份验证
- 验证用户对告警的访问权限
- 支持基于角色的访问控制(RBAC)
- 记录所有告警查询操作

### 9.2 数据安全要求
- 敏感信息脱敏
- 告警数据加密存储
- 支持数据访问审计
- 防止SQL注入

### 9.3 身份验证机制
- JWT令牌
- API密钥
- 请求频率限制
- IP白名单

## 10. 日志与审计要求

### 10.1 操作日志要求
- 记录所有告警查询操作
- 记录参数和结果
- 记录操作人员身份
- 记录操作时间和IP

### 10.2 审计日志要求
- 记录告警访问轨迹
- 记录身份和权限
- 记录访问目的
- 支持长期保存

### 10.3 日志格式规范
```json
{
  "timestamp": "2024-01-16T18:20:00Z",
  "level": "INFO",
  "service": "queue-alert",
  "operation": "list_alerts",
  "user_id": "user_001",
  "queue_id": "queue_001",
  "parameters": {
    "alert_type": "error",
    "page": 1
  },
  "result": {
    "total": 1
  }
}
```

## 11. 测试用例

### 11.1 功能测试用例

#### 11.1.1 正常查询测试
**测试场景：** 查询队列告警
**输入数据：**
```json
{
  "queue_id": "queue_001"
}
```
**预期结果：**
- 返回状态码：200
- 返回告警列表

#### 11.1.2 条件筛选测试
**测试场景：** 按类型筛选
**输入数据：**
```json
{
  "queue_id": "queue_001",
  "alert_type": "error"
}
```
**预期结果：**
- 返回状态码：200
- 只包含error类型告警

### 11.2 性能测试用例

#### 11.2.1 大量告警查询测试
**测试场景：** 并发查询1000个队列告警
**预期结果：**
- 响应时间 < 2秒
- 错误率 < 1%

### 11.3 安全测试用例

#### 11.3.1 权限验证测试
**测试场景：** 验证用户告警访问权限
**测试数据：**
- 用户A：有权限
- 用户B：无权限
**预期结果：**
- 用户A：成功查询
- 用户B：返回权限错误

#### 11.3.2 敏感信息脱敏测试
**测试场景：** 敏感信息脱敏
**输入数据：**
```json
{
  "queue_id": "queue_with_sensitive"
}
```
**预期结果：**
- 敏感字段脱敏
- 不影响可用性

### 11.4 异常测试用例

#### 11.4.1 队列不存在测试
**测试场景：** 查询不存在的队列告警
**输入数据：**
```json
{
  "queue_id": "non_existent_queue"
}
```
**预期结果：**
- 返回404
- 错误信息友好

#### 11.4.2 参数错误测试
**测试场景：** 参数错误
**输入数据：**
```json
{
  "queue_id": "",
  "page": 0
}
```
**预期结果：**
- 返回参数验证错误
- 状态码400

#### 11.4.3 系统异常测试
**测试场景：** 数据库连接失败
**测试方法：** 临时关闭数据库
**预期结果：**
- 返回500
- 记录详细错误日志 