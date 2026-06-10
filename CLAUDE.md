# Claude Code 项目配置与规范

## 项目概述

Marmot Ledger（土拨鼠账本）是一个面向个人与家庭的财务账本系统，用于记录现金流、资产、负债、多币种账户、投资、报销、押金以及共享家庭财务。

当前仓库已经从旧 `my-life` 项目中清理出基础框架。旧的 `transaction` 记账、家庭、健康、训练、OSS 上传和 Elasticsearch 业务模块已移除。后续财务模块应基于新的 Marmot Ledger 领域模型重新实现。

## 当前保留的基础能力

### 后端

- Go + Fiber Web 服务
- MySQL / Redis 基础设施连接
- Redis token 中间件
- 用户登录、token 校验、token 刷新、用户信息更新
- 通用返回结构、错误码、时间模型和工具函数

### 前端

- Vue 3 + Vite
- Element Plus
- Pinia
- Vue Router
- Axios 请求封装
- 登录页、系统首页、用户中心
- Header、Sidebar、Tags 基础布局组件

## 核心设计文档

- `docs/MARMOT_LEDGER_HANDOFF.md`
- `docs/财务系统升级设计方案.md`
- `docs/导入识别与规则学习设计方案.md`

这些文档是后续新财务系统设计与实现的依据，不要删除。

## 新财务领域模型方向

核心模型：

```text
financial_event：财务事件，表示发生了什么。
ledger_entry：余额分录，表示这件事导致哪些 Bucket 余额变化。
bucket：资产 / 负债 / 虚拟资金池。
account：平台 / 机构 / 分组，不保存余额。
category：用户自己的分类。
category_group：系统 / 家庭统计用的聚合分类。
channel_template：收款 / 支付渠道模板。
currency：全平台公共币种字典。
exchange_rate：汇率缓存。
investment_snapshot：投资估值 / 收益快照。
```

推荐关系：

```text
User
  ├── Account
  │     └── Bucket
  │           └── LedgerEntry
  ├── FinancialEvent
  │     ├── LedgerEntry
  │     └── ExchangeDetail
  ├── Category
  │     └── CategoryGroup
  ├── InvestmentAsset
  │     └── InvestmentSnapshot
  └── ExchangeRate
```

## 关键领域决策

### 命名

使用：

```text
financial_event
ledger_entry
bucket
```

不要在新项目中把 `transaction` 作为主表名。

字段命名：

```text
event_type
event_group_id
event_time
financial_event_id
```

### Account

```text
Account = platform / institution / grouping.
Account does not store balance.
Account total balance is computed from child Buckets.
Account always belongs to a user.
No family-shared Account is needed.
```

### Bucket

```text
Bucket = asset / liability / virtual funds container.
Bucket is a core model, not optional.
Bucket balance is stored in bucket.balance.
Bucket balance can be recalculated from ledger_entry.
```

Bucket nature：

```text
asset
liability
```

净资产：

```text
net_worth = sum(asset bucket balances) - sum(liability bucket balances)
```

### LedgerEntry

`ledger_entry.amount` 使用正负数表示余额变化，不要增加 direction 字段：

```text
amount > 0：这个 Bucket 余额增加
amount < 0：这个 Bucket 余额减少
```

`ledger_entry.amount` 只表达余额变化。收入/支出的业务含义来自 `financial_event.event_type`。

### 初始余额

创建 Bucket 时：

```text
bucket.initial_balance = user input
bucket.balance = initial_balance
create financial_event event_type = balance_adjustment
create ledger_entry.amount = initial_balance
remark = 初始化余额
include_in_statistics = false
```

### 余额一致性

所有影响 Bucket 余额的操作必须在数据库事务中完成。

同一事务内：

```text
create/update financial_event
create/delete/rebuild ledger_entry
update bucket.balance
```

编辑/删除事件第一版使用 rollback-and-rebuild，不做 reversal-ledger audit mode。

### FinancialEvent 类型

第一版确认类型：

```text
income
expense
refund
transfer
exchange
receivable_create
receivable_collect
deposit_create
deposit_refund
loan_out
loan_collect
investment_buy
investment_sell
investment_income
balance_adjustment
```

统计口径：

```text
income / expense 默认进入普通收入支出统计
refund 用于抵扣支出，不算普通收入
transfer / exchange / receivable / deposit / loan / investment_buy / investment_sell / balance_adjustment 默认不进入普通收入支出统计
investment_income 只进入投资统计，不进入普通收入统计
```

### Category 和 CategoryGroup

第一版不要引入 OperationType。

```text
Category = user-owned category
CategoryGroup = system/family reporting aggregation category
```

家庭报表按 `category_group_id` 聚合，不按 `category_id` 聚合。

### ChannelTemplate

渠道使用模板表，不要拆分微信收款/微信支付。

使用：

```text
income + WECHAT = 微信收款
expense + WECHAT = 微信支付
income + ALIPAY = 支付宝收款
expense + ALIPAY = 支付宝支付
```

不要创建：

```text
WECHAT_PAY
WECHAT_RECEIVE
ALIPAY_PAY
ALIPAY_RECEIVE
```

## 报表边界

产品定义：

```text
账单：看我干了啥
资产：看我有啥
```

分离：

```text
账单统计
资产报表
```

账单统计使用：

```text
event_type + category/category_group + time range
```

资产报表使用：

```text
bucket.balance + bucket_nature + currency
```

家庭账单统计不按 Account / Bucket / Channel 分组。

家庭资产报表展示每个家庭成员自己的 Account / Bucket，不创建家庭共享 Account。

## 当前项目结构

### 前端

- `ui/src/main.js` - 前端入口
- `ui/src/App.vue` - 根组件
- `ui/src/router/index.js` - 路由定义
- `ui/src/permission.js` - 路由鉴权
- `ui/src/api/request.js` - Axios 请求封装
- `ui/src/api/auth/auth.js` - 登录和 token API
- `ui/src/api/user/user.js` - 用户 API
- `ui/src/views/login/Login.vue` - 登录页
- `ui/src/views/HomeView.vue` - 后台布局页
- `ui/src/views/Dashboard.vue` - 系统首页
- `ui/src/views/User.vue` - 用户中心
- `ui/src/components/Header.vue` - 顶部栏
- `ui/src/components/Sidebar.vue` - 侧边栏
- `ui/src/components/Tags.vue` - 页面标签栏
- `ui/src/stores/` - Pinia store

### 后端

- `main.go` - 后端入口
- `api/web.go` - Fiber 服务启动、CORS、token 中间件、路由挂载
- `api/user.go` - 用户接口
- `env/properties.go` - 配置读取
- `env/dev.yaml` - 开发环境配置
- `internal/infrastructure/mysql.go` - MySQL 初始化
- `internal/infrastructure/redis.go` - Redis 初始化
- `internal/service/user_service.go` - 用户服务
- `internal/domain/entity/user/user.go` - 用户实体
- `internal/domain/repository/userdb/user.go` - 用户仓储
- `pkg/myresult/` - API 返回结构
- `pkg/myerror/` - 错误码
- `pkg/model/` - 日期时间模型
- `pkg/utils/` - 通用工具

## 技术栈

### 前端

- Vue 3 + Composition API
- Vite
- Element Plus
- Pinia
- Vue Router
- Axios

### 后端

- Go
- Fiber
- XORM
- MySQL
- Redis

## 验证命令

后端：

```bash
go test ./...
go build ./...
```

前端：

```bash
npm --prefix ui run build
```

## 开发注意事项

- 新财务模块不要复用旧 `transaction` 作为主领域模型。
- 不要恢复旧 money/family/health/train/OSS/ES 模块，除非明确是为了查阅历史实现。
- 涉及余额变化的后端操作必须设计数据库事务。
- 新 UI 中如果涉及收入/支出颜色，沿用项目约定：收入红色 `#ef4444`，支出绿色 `#10b981`。
- `docs/` 下设计文档优先级高于旧 README 或历史代码习惯。
