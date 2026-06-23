# Claude Code 项目配置与规范

## 项目概述

Marmot Ledger（土拨鼠账本）是面向个人与家庭的财务账本系统，用于记录现金流、资产、负债、多币种账户、报销、押金、投资与共享家庭财务。

当前仓库已从历史项目清理为新项目基础框架。后续开发必须基于 Marmot Ledger 新财务领域模型，不要恢复旧 money / transaction / family / health / train / OSS / Elasticsearch 等历史业务模块。

## 当前进度（截至 2026-06-11）

### 已完成

- 远程仓库已切换为 `https://github.com/here-tunan/marmot-ledger.git`。
- 数据库统一为 `marmot_ledger`，SQL 初始化脚本位于 `sql/`。
- Go module 已迁移为 `marmot-ledger`，历史 `go-my-life` / `my-life` / `life` 命名已清理。
- 保留基础框架：Go + Fiber、MySQL、Redis、用户登录/token、Vue 3 + Vite、Element Plus、Pinia、Vue Router、Axios。
- 设计文档整理至 `docs/`，UI 方向为 `Calm Marmot Finance`。
- 品牌图位于：`img/marmot-ledger-1.png`、`img/marmot-ledger-2.png`。
- 已实现后端核心模块：Account、Bucket、FinancialEvent、LedgerEntry、Category、CategoryGroup、Statistics、Family。
- Bucket 创建会自动生成 `balance_adjustment` financial_event 和 ledger_entry，初始化余额不进入统计。
- 已实现统一 `POST /api/record`，首批支持：`income`、`expense`、`transfer`、`refund`。
- Record 创建/编辑/删除使用事务，编辑/删除第一版采用 rollback-and-rebuild。
- 已加入轻量 Record Strategy：income / expense / transfer / refund 分策略构建事件和分录。
- 已实现 Records Center：查询、详情、编辑、删除、分页。
- 已实现统计：个人/家庭 summary、category group 按原始币种分别聚合；refund 作为支出抵扣，不作为收入。
- 已实现 Family：多家庭、成员邀请、接受/拒绝、家庭统计聚合。
- 已实现前端页面：Dashboard、Accounts、Buckets、Record、Records、Categories、Family、Login、Header、Sidebar、User Center。
- 已接入全局中英文 i18n；中文 Bucket 术语统一为“资金桶”，英文保留 “Bucket”。
- 已统一货币前端展示工具 `ui/src/utils/currency.js`，包含国旗、中文名、英文名、符号。
- 当前支持货币：CNY、USD、HKD、EUR、JPY、GBP、SGD、AUD、NZD。
- Family / Buckets / Record / Records 的货币下拉已改为统一 currencyOptions，并支持 AUD / NZD。
- 已移除普通账单统计的 base currency 口径：`financial_event` 不再保存 `base_currency` / `base_amount` / `exchange_rate`，`family` 不再保存 `base_currency`。
- 已验证过：`go test ./...`、`go build ./...`、`npm --prefix ui run build`；接口 smoke test 曾覆盖 Account/Bucket 初始化、income、expense、transfer、refund。

### 最近完成

- 修复 Records 分页 total 统计问题。
- 修复全局页面底部空间过窄问题。
- 修复 refund 统计语义：refund 进入统计，作为 expense offset。
- 修复 Family route i18n key 和部分 Bucket 术语不一致问题。
- SQL 初始化已整合为单文件：`sql/marmot_ledger.sql`。
- 前端新增统一货币展示映射：`ui/src/utils/currency.js`。

### 待继续

- 继续做全局 i18n/硬编码文案扫描：ThemePicker、UserCard、Header aria-label、动态 enum label 等。
- Dashboard 的 Family 模式需要进一步完善：家庭选择、成员卡片、家庭统计视图、无家庭引导。
- Records / Record 页面动态事件类型与 entry role 建议使用 `te()` fallback，避免显示 raw i18n key。
- 可考虑安全清理：`/api/user/info` 和 Redis token payload 不应包含 password。
- 稳定后再提交当前大批变更；除非用户明确要求，不要 push。

## 核心文档

- `README.md` / `README.zh-CN.md`
- `docs/SQL_DEVELOPER_GUIDE.md`
- `docs/DESIGN.md`

当前项目只保留必要文档；SQL 初始化以 `sql/marmot_ledger.sql` 为准。

## 品牌与 UI

- 品牌图优先使用：`img/marmot-ledger-1.png`、`img/marmot-ledger-2.png`。
- UI 方向：Calm Marmot Finance。
- 风格参考：Wise 50% + Linear 35% + Notion 15%。
- 背景：warm cream；卡片：16px radius、轻阴影；按钮 active scale `0.96`。
- 品牌/focus 颜色：苔绿 `#2f7d5c`；收入颜色：红色 `#ef4444`；支出颜色：橙色 `#f97316`。
- 中文 UI 中 Bucket 统一叫“资金桶”；英文 UI 保留 “Bucket”。

## 核心领域模型

```text
financial_event：财务事件，表示发生了什么。
ledger_entry：余额分录，表示事件导致哪些 Bucket 余额变化。
bucket：资产 / 负债 / 虚拟资金池，保存余额。
account：平台 / 机构 / 分组，不保存余额。
category：用户自己的分类。
category_group：系统 / 家庭统计用聚合分类。
channel_template：收款 / 支付渠道模板。
currency：公共币种字典。
exchange_rate：汇率缓存。
investment_snapshot：投资估值 / 收益快照。
```

推荐关系：

```text
User
  ├── Account -> Bucket -> LedgerEntry
  ├── FinancialEvent -> LedgerEntry / ExchangeDetail
  ├── Category -> CategoryGroup
  ├── InvestmentAsset -> InvestmentSnapshot
  └── ExchangeRate
```

## 关键领域规则

### 命名

- 使用 `financial_event`、`ledger_entry`、`bucket`。
- 不要把 `transaction` 作为新主领域模型或主表名。
- 核心字段：`event_type`、`event_group_id`、`event_time`、`financial_event_id`。

### Account

- Account = platform / institution / grouping。
- Account 不保存余额；账户总余额从子 Bucket 汇总。
- Account 始终属于用户。
- 不创建家庭共享 Account。

### Bucket

- Bucket = asset / liability / virtual funds container。
- `bucket.balance` 保存当前余额，可由 ledger_entry 重算。
- Bucket nature：`asset` / `liability`。
- 净资产：`sum(asset bucket balances) - sum(liability bucket balances)`。

### LedgerEntry

- `ledger_entry.amount` 使用正负数表示余额变化，不增加 direction 字段。
- `amount > 0`：Bucket 余额增加。
- `amount < 0`：Bucket 余额减少。
- 业务含义来自 `financial_event.event_type`，不是来自 entry direction。

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

所有影响 Bucket 余额的操作必须在数据库事务中完成：

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

- income / expense 默认进入普通收支统计。
- refund 用于抵扣支出，不算普通收入。
- transfer / exchange / receivable / deposit / loan / investment_buy / investment_sell / balance_adjustment 默认不进入普通收支统计。
- investment_income 只进入投资统计，不进入普通收入统计。

### Category / CategoryGroup

- 第一版不要引入 OperationType。
- Category = user-owned category。
- CategoryGroup = system/family reporting aggregation category。
- 家庭报表按 `category_group_id` 聚合，不按 `category_id` 聚合。

### ChannelTemplate

渠道使用模板表，不拆分微信收款/微信支付：

```text
income + WECHAT = 微信收款
expense + WECHAT = 微信支付
income + ALIPAY = 支付宝收款
expense + ALIPAY = 支付宝支付
```

不要创建：`WECHAT_PAY`、`WECHAT_RECEIVE`、`ALIPAY_PAY`、`ALIPAY_RECEIVE`。

## 报表边界

- 账单：看我干了啥。
- 资产：看我有啥。
- 账单统计使用：`event_type + category/category_group + time range`。
- 资产报表使用：`bucket.balance + bucket_nature + currency`。
- 普通账单统计按 `financial_event.currency + amount` 分原始币种聚合，不使用本位币折算口径。
- 不同币种金额不合并、不默认换算；UI 按币种分块展示。
- `financial_event` 不保存 `base_currency` / `base_amount` / `exchange_rate`；`family` 不保存 `base_currency`。
- 家庭账单统计不按 Account / Bucket / Channel 分组，也不做家庭本位币折算；按 active 成员事件原始币种聚合。
- 家庭资产报表展示每个家庭成员自己的 Account / Bucket，不创建家庭共享 Account。

## 常用验证命令

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

- 新财务模块必须基于 `financial_event` / `ledger_entry` / `bucket` 扩展。
- 不恢复旧业务模块，除非明确是为了查阅历史实现。
- 余额变化必须使用数据库事务。
- `docs/` 下设计文档优先级高于旧 README 或历史代码习惯。
- 代码风格保持现有 Go / Vue 写法；新增 UI 文案优先接入 i18n。
