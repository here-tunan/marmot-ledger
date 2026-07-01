# Marmot Ledger · 土拨鼠账本

<p align="center">
  <img src="img/marmot-ledger-1.png" alt="Marmot Ledger" width="180" />
  <img src="img/marmot-ledger-2.png" alt="Marmot Ledger" width="180" />
</p>

<p align="center">
  <strong>一个温和、清晰、可信的个人与家庭财务账本。</strong><br />
  简体中文 · <a href="README.md">English</a>
</p>

土拨鼠账本面向个人与家庭财务管理：记录日常收支、资产、负债、多币种账户、投资、报销、押金与家庭共享财务。它关注“钱发生了什么”和“我现在拥有什么”，而不是把生活变成复杂会计系统。

---

## 核心能力

- **账户** — 管理现金、钱包、银行、卡片、投资平台等资金所在位置。
- **资金桶** — 追踪资产、负债、应收款、押金、投资现金、投资资产等真实余额。
- **记账记录** — 支持收入、支出、转账、退款、换汇、应收/押金/借出、分摊、家庭转账与投资场景。
- **余额分录** — 每个影响余额的财务事件都会生成正负分录，资金桶余额可追溯。
- **分类** — 管理个人收支分类，并支持家庭维度的统计分组。
- **渠道** — 管理微信、支付宝、现金、银行转账、银行卡和自定义收付款渠道。
- **总览** — 按原始币种查看净资产、资产、负债、收支、退款、趋势、分类统计与投资盈亏。
- **家庭工作台** — 邀请成员，聚合 active 成员记录，不创建家庭共享账户。
- **模板管理** — 管理账户、渠道、分类系统模板。
- **表格导入** — 为每个平台创建导入配置，把 xlsx / csv 账单的列映射到 record 字段；支持默认值 + 条件规则（`contains` / `equals` / `notContains` / `notEquals` / `containsAny` / `notContainsAny` / `equalsAny` / `notEqualsAny`）解析出分类、渠道、收支类型、资金桶；行过滤规则（drop / keep）可跳过明细行；预览页可编辑、异常行即时标红，"仅异常"过滤 + "跳到下一个异常"帮助定位；分片批量入账（100 行/HTTP，每行独立事务）+ 实时进度条，坏行不拖累其他行。
- **双语界面** — 内置中文与英文 i18n。

---

## 产品原则

1. **账本优先，不以 transaction 为核心。** 核心模型是 `financial_event` + `ledger_entry` + `bucket`。
2. **账户负责分组，资金桶承载余额。** 账户是平台或机构，资金桶保存当前余额。
3. **分录使用正负数。** `ledger_entry.amount > 0` 表示资金桶余额增加；`amount < 0` 表示减少。
4. **按原始币种统计。** 普通报表默认不合并不同币种。
5. **退款抵扣支出。** 退款减少支出，不作为普通收入。
6. **家庭数据仍归成员所有。** 家庭报表聚合 active 成员自己的记录和资金桶，不创建家庭共享账户。

---

## 技术栈

### 后端

- Go
- Fiber
- Xorm
- MySQL
- Redis

### 前端

- Vue 3
- Vite
- Element Plus
- Pinia
- Vue Router
- Axios
- vue-i18n
- ECharts

---

## 快速开始

### 1. 准备服务

请先安装并启动：

- 推荐 MySQL 8.x
- 推荐 Redis 6.x+

如需调整本地配置，请修改：

```text
env/dev.yaml
```

### 2. 创建数据库

新环境只需要执行一个 SQL 文件：

```bash
mysql < sql/marmot_ledger.sql
```

这个文件会创建数据库、当前表结构、索引、外键和必要参考数据。它**不会**创建演示用户或演示账单。

完整表结构和数据策略见 [`docs/SQL_DEVELOPER_GUIDE.md`](docs/SQL_DEVELOPER_GUIDE.md)。

### 3. 启动后端

```bash
go run main.go
```

默认读取 `env/dev.yaml`。如需切换环境：

```bash
export MARMOT_LEDGER_ENV=prod
```

### 4. 启动前端

```bash
npm --prefix ui install
npm --prefix ui run dev
```

然后打开终端中显示的 Vite 开发地址。

---

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

数据库初始化 smoke check：

```bash
mysql < sql/marmot_ledger.sql
```

然后在应用中验证：

1. 注册并登录。
2. 创建账户。
3. 创建资金桶。
4. 创建收入、支出、退款、转账记录。
5. 查看 Dashboard 按原始币种展示统计。
6. 如需测试家庭功能，创建家庭并邀请/接受成员。

---

## 项目结构

```text
.
├── api/                         # Fiber API 路由
├── env/                         # 环境配置
├── internal/
│   ├── domain/                  # 领域实体与仓储
│   ├── infrastructure/          # MySQL / Redis 初始化
│   └── service/                 # 应用服务与记账策略
├── pkg/                         # 通用模型、错误、返回结构、工具
├── sql/
│   └── marmot_ledger.sql        # 单文件数据库初始化脚本
├── ui/                          # Vue 3 前端
│   └── src/
│       ├── api/                 # 前端 API 封装
│       ├── components/          # 布局与复用组件
│       ├── i18n/                # 中文 / 英文语言包
│       ├── router/              # Vue Router
│       ├── stores/              # Pinia stores
│       ├── styles/              # 全局样式
│       └── views/               # 产品页面
├── docs/
│   ├── DESIGN.md                # Calm Marmot Finance UI 方向
│   └── SQL_DEVELOPER_GUIDE.md   # 数据库说明
└── img/                         # 品牌图片
```

---

## 领域模型

```text
User
  ├── Account -> Bucket -> LedgerEntry
  ├── FinancialEvent -> LedgerEntry
  ├── Category
  ├── PersonalChannel
  ├── Family / FamilyMember
  └── Templates / reference data
```

| 模型 | 含义 |
| --- | --- |
| `financial_event` | 发生了什么：收入、支出、退款、转账、换汇、投资等。 |
| `ledger_entry` | 事件如何影响资金桶余额。 |
| `bucket` | 资产 / 负债 / 资金容器，保存当前余额。 |
| `account` | 平台、机构或分组，不保存余额。 |
| `category` | 用户自己的收支分类。 |
| `category_group` | 系统统计分组参考字典。 |
| `family_category_group` | 家庭自己的统计分组。 |
| `personal_channel` | 用户自己的收付款渠道。 |
| `channel_template` | 系统渠道模板。 |
| `currency` | 公共币种字典。 |

---

## 文档

- [`docs/SQL_DEVELOPER_GUIDE.md`](docs/SQL_DEVELOPER_GUIDE.md) — 表结构、参考数据、领域流程和数据库初始化。
- [`docs/DESIGN.md`](docs/DESIGN.md) — Calm Marmot Finance 品牌与 UI 方向。

---

## 开发约定

- 新财务能力基于 `financial_event`、`ledger_entry`、`bucket` 扩展。
- 不恢复旧的无关业务模块。
- 不把 `transaction` 作为主领域表。
- 所有影响余额的操作必须在数据库事务中完成。
- 中文 UI 使用 **资金桶**；英文 UI 保留 **Bucket**。
- 新增可见 UI 文案优先接入 i18n。
- 普通统计保持原始币种口径，除非单独设计汇率能力。
- SQL 初始化脚本不能包含演示用户、默认管理员密码、个人样例数据或硬编码 `user_id = 1` 的记录。

---

## Post-MVP 方向

- 完善家庭模式下的 Dashboard 选择与成员视图。
- 表格导入的幂等提交与"整批撤销"（当前未做去重，重复提交会重复入账；后续通过 request-id + 服务端去重表 + `event_group_id` 分组撤销来解决）。
- 持续清理 i18n 和硬编码 UI 文案。
- 加强 auth/token payload，确保 password 字段不会暴露。
- 当项目出现需要增量升级的生产数据库后，引入正式 migration 工具。

---

## License

暂未选择 License。
