# ratelimiter

限流器服务，纯 Go 标准库实现。

## 运行

```bash
go run ./cmd/server
```

## API 概览

- `/api/rules` — 限流规则 CRUD
- `/api/clients` — 客户端 CRUD + 封禁/解封
- `/api/access-records` — 访问记录（只读追加）
- `/api/buckets` — 令牌桶状态 + refill
- `/api/alerts` — 告警 CRUD
- `/api/check` — 限流检查（核心业务）
- `/api/stats/overview` — 统计概览
- `/api/stats/top-clients` — Top N 客户端

## 模块

- `pkg/httpx` — HTTP 通用工具
- `pkg/idgen` — ID 生成
- `pkg/logger` — 日志
- `internal/config` — 配置
- `internal/model` — 领域模型
- `internal/store` — 数据访问（内存实现）
- `internal/service` — 业务逻辑
- `internal/handler` — HTTP 处理器
