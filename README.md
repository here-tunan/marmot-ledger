# Marmot Ledger

Marmot Ledger（土拨鼠账本）是一个面向个人与家庭的财务账本系统，用于记录现金流、资产、负债、多币种账户、投资、报销、押金以及共享家庭财务。

<p align="center">
  <img src="img/marmot-ledger-1.png" alt="Marmot Ledger" width="180" />
  <img src="img/marmot-ledger-2.png" alt="Marmot Ledger" width="180" />
</p>

当前仓库已经从历史项目中清理出基础框架。历史业务模块已移除；后续财务能力将基于新的核心模型重新实现。

## 当前状态

当前保留的是最小可运行基础：

- Go + Fiber 后端服务
- MySQL / Redis 基础设施连接
- 用户登录、token 校验、用户信息接口
- Vue 3 + Vite 前端应用
- Element Plus UI
- Pinia 状态管理
- 登录页、系统首页、用户中心

## 领域模型方向

后续新财务模块应基于以下核心模型实现：

- `financial_event`：财务事件，表示发生了什么
- `ledger_entry`：余额分录，表示事件造成的 Bucket 余额变化
- `bucket`：资产 / 负债 / 虚拟资金容器
- `account`：平台 / 机构 / 分组，不保存余额
- `category`：用户自己的分类
- `category_group`：系统 / 家庭统计用的聚合分类
- `channel_template`：收款 / 支付渠道模板
- `currency`：全平台公共币种字典
- `exchange_rate`：汇率缓存
- `investment_snapshot`：投资估值 / 收益快照

> 新项目不要把 `transaction` 作为主表名。账单统计回答“发生了什么”，资产报表回答“我有什么”。

详细设计见：

- [`docs/MARMOT_LEDGER_HANDOFF.md`](docs/MARMOT_LEDGER_HANDOFF.md)
- [`docs/财务系统升级设计方案.md`](docs/财务系统升级设计方案.md)
- [`docs/导入识别与规则学习设计方案.md`](docs/导入识别与规则学习设计方案.md)

## 项目结构

```text
.
├── api/                         # Fiber API 路由
│   ├── user.go                  # 用户接口
│   └── web.go                   # Web 服务启动、CORS、token 中间件
├── env/                         # 环境配置
├── internal/
│   ├── domain/                  # 领域实体与仓储
│   │   ├── entity/user/
│   │   └── repository/userdb/
│   ├── infrastructure/          # MySQL / Redis 初始化
│   └── service/                 # 应用服务
├── pkg/                         # 通用模型、错误、返回结构、工具
├── ui/                          # Vue 3 前端
│   └── src/
│       ├── api/                 # 前端 API 封装
│       ├── components/          # 基础布局组件
│       ├── router/              # 路由
│       ├── stores/              # Pinia store
│       └── views/               # 登录、首页、用户中心
└── docs/                        # 设计与 handoff 文档
```

## 本地运行

### 后端

准备 MySQL 和 Redis，并根据本地环境修改：

```text
env/dev.yaml
```

启动后端：

```bash
go run main.go
```

默认读取 `env/dev.yaml`。如需切换环境：

```bash
export MARMOT_LEDGER_ENV=prod
```

### 前端

```bash
cd ui
npm install
npm run dev
```

## 构建与验证

后端：

```bash
go test ./...
go build ./...
```

前端：

```bash
npm --prefix ui run build
```

## 开发约定

- 后续财务模块从 `financial_event` / `ledger_entry` / `bucket` 开始建模。
- `account` 不保存余额，余额由其下 Bucket 汇总。
- `ledger_entry.amount` 使用正负数表示余额变化，不增加 direction 字段。
- `refund` 抵扣支出，不作为普通收入。
- 家庭资产报表展示家庭成员各自的账户和 Bucket，不创建家庭共享 Account。
