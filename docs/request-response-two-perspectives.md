# 请求/响应双向视角架构分析

> 目的：同时从 **用户输入 → 服务器处理**（自上而下）与 **服务器响应 → 用户接收**（自下而上）两个角度，梳理 OneGoServer 的调用链与各层职责，以便排查问题与性能优化。

---

## 1. 自上而下：用户输入路径

| 顺序 | 层级 | 关键职责 | 典型文件 |
| ---- | ---- | -------- | -------- |
| 1 | Client | 发送 HTTP/JSON 请求，携带 Header（Auth、Trace） | Postman / Web
| 2 | Router (Gin) | URL → Controller 映射 | `internal/router/router.go` |
| 3 | Middleware | Trace/RL/Auth/Recovery | `internal/middleware/*` |
| 4 | Controller | 解析请求头/路径，调用 **API DTO** 绑定 | `internal/controller/*.go` |
| 5 | **API (DTO)** | `ShouldBindJSON(&DTO)` 参数校验/映射 | `api/v1/*.go` |
| 6 | Service | 事务控制、聚合多 Biz、调用队列 | `internal/service/*.go` |
| 7 | Biz | 领域规则、状态机、计算 | `internal/biz/*` |
| 8 | Repository | GORM DAO，读写 DB | `internal/data/repository/*.go` |
| 9 | External | PostgreSQL / Redis / Kafka | Docker / Cloud |

### Mermaid 流程图（上行）

```mermaid
flowchart TD
    A[Client] --> B[Router]
    B --> C[Middleware]
    C --> D[Controller]
    D --> X[API (DTO)]
    X --> E[Service]
    E --> F[Biz]
    F --> G[Repository]
    G --> H[PostgreSQL / Redis / Kafka]
```

---

## 2. 自下而上：服务器响应路径

| 顺序 | 层级 | 关键职责 |
| ---- | ---- | -------- |
| 1 | External | 返回查询结果 / ACK |
| 2 | Repository | 将原始数据映射 Domain Model |
| 3 | Biz | 封装业务对象，执行业务后处理 |
| 4 | Service | 汇总数据，处理错误 → DTO |
| 5 | **API (DTO)** | 将内部对象映射响应 DTO |
| 6 | Controller | 封装统一 `Response{code,msg,data}` |
| 7 | Middleware | Trace 记录、错误转换、Header 填充 |
| 8 | Router | 交由 Gin 写 Socket |
| 9 | Client | 接收 JSON，渲染 UI |

### Mermaid 流程图（下行）

```mermaid
flowchart TD
    H[PostgreSQL / Redis / Kafka] --> G[Repository]
    G --> F[Biz]
    F --> E[Service]
    E --> X[API (DTO)]
    X --> D[Controller]
    D --> C[Middleware]
    C --> B[Router]
    B --> A[Client]
```

---

## 3. 双向视角下的优势

1. **单向思考不足**：双向视角有助于发现响应耗时瓶颈（DB→Repo→Biz）与请求过载点（Middleware→Controller）。
2. **Debug 路径清晰**：定位问题时可由上至下逐层断点，也可由下至上追溯数据来源。
3. **性能优化**：识别长耗时 SQL、冗余 Biz 计算、重复序列化等。
4. **可观测性**：TraceID 贯穿两条路径，配合日志与指标快速闭环。

---

## 4. 后续行动清单

| # | 文件 | 操作 |
|---|------|------|
| 1 | `docs/request-response-two-perspectives.md` | 更新 API 层并补全章节 |
| 2 | (可选) | 在 README 或开发者文档中引用本文件 |

如需深入某层代码示例，请随时提出。