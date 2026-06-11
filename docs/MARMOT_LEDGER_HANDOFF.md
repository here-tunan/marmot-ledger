# Marmot Ledger Handoff

## Project

Repository name:

```text
marmot-ledger
```

Product name:

```text
Marmot Ledger
```

Chinese name:

```text
土拨鼠账本
```

Repo description:

```text
A personal and family financial ledger for tracking cash flow, assets, liabilities, multi-currency accounts, investments, reimbursements, deposits, and shared household finances.
```

Positioning:

```text
A personal and family financial ledger for assets, cash flow, liabilities, investments, and shared household finances.
```

---

## Core documents

Current design docs:

```text
marmot-ledger/docs/财务系统升级设计方案.md
marmot-ledger/docs/导入识别与规则学习设计方案.md
```

---

## Core domain model

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

Recommended relation:

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

---

## Naming decisions

Use:

```text
financial_event
ledger_entry
bucket
```

Do not use `transaction` as the primary table name in the new project.

Field naming:

```text
event_type
event_group_id
event_time
financial_event_id
```

---

## Key decisions

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

Confirmed bucket types:

```text
cash              现金
wallet            电子钱包
bank              银行账户
credit            信用卡
investment_cash   投资账户现金
investment_asset  投资资产
receivable        应收款
deposit           押金 / 保证金
loan_out          借出款
liability         负债
virtual           虚拟桶
```

Bucket nature:

```text
asset
liability
```

Net worth:

```text
net_worth = sum(asset bucket balances) - sum(liability bucket balances)
```

Negative balance policy:

```text
cash / wallet / bank：不允许负数
receivable / deposit / loan_out：不允许负数
investment_cash / investment_asset：不允许负数
credit / liability：balance 表示欠款金额，按正数保存
virtual：按具体用途判断，第一阶段尽量少用
```

### LedgerEntry

Use signed amount, no direction field:

```text
amount > 0：这个 Bucket 余额增加
amount < 0：这个 Bucket 余额减少
```

`ledger_entry.amount` expresses balance change only. Income/expense meaning comes from `financial_event.event_type`.

Recommended ledger_entry fields:

```text
ledger_entry
  id
  financial_event_id
  user_id
  bucket_id
  currency
  amount
  balance_after
  entry_role optional
  created_at
```

Optional `entry_role` examples:

```text
income
expense
transfer_out
transfer_in
exchange_out
exchange_in
fee
receivable_create
receivable_collect
investment_buy
investment_sell
adjustment
```

### Initial balance

When creating a Bucket:

```text
bucket.initial_balance = user input
bucket.balance = initial_balance
create financial_event event_type = balance_adjustment
create ledger_entry.amount = initial_balance
remark = 初始化余额
include_in_statistics = false
```

This makes initial balance traceable and recalculable.

### Balance consistency

All operations that affect Bucket balance must run in a DB transaction.

Within the same transaction:

```text
create/update financial_event
create/delete/rebuild ledger_entry
update bucket.balance
```

Before updating buckets, lock relevant bucket rows, e.g. MySQL:

```sql
SELECT * FROM bucket WHERE id IN (...) FOR UPDATE;
```

Update by increment:

```sql
UPDATE bucket
SET balance = balance + ?
WHERE id = ?;
```

`?` is `ledger_entry.amount`.

### Edit/delete financial event

First version uses rollback-and-rebuild:

```text
1. Open DB transaction.
2. Read old ledger_entry rows.
3. Roll back bucket.balance using reverse amount.
4. Delete or soft-delete old ledger_entry rows.
5. Update/delete financial_event.
6. Generate new ledger_entry rows if editing.
7. Update bucket.balance.
8. Commit.
```

No reversal-ledger audit mode in v1.

---

## FinancialEvent types

Confirmed first version:

```text
income                收入
expense               支出
refund                退款
transfer              转账
exchange              换汇
receivable_create     形成应收
receivable_collect    收回应收
deposit_create        支付押金
deposit_refund        退回押金
loan_out              借出
loan_collect          收回借出款
investment_buy        投资买入
investment_sell       投资卖出
investment_income     投资收入
balance_adjustment    余额调整
```

Statistics:

```text
income / expense 默认进入普通收入支出统计
refund 用于抵扣支出，不算普通收入
transfer / exchange / receivable / deposit / loan / investment_buy / investment_sell / balance_adjustment 默认不进入普通收入支出统计
investment_income 只进入投资统计，不进入普通收入统计
```

### Refund

Refund uses `refund`, not `income`.

Example:

```text
expense 100 -> bucket -100
refund 30 -> bucket +30
report expense net = 70
income does not increase
```

### Fees

All fees are separate `expense` financial events and linked by `event_group_id`.

Applies to:

```text
transfer fee
exchange fee
investment fee
withdrawal fee
```

### Transfer fee

```text
transfer principal does not affect income/expense stats
fee is optional
no fee: transfer financial_event only
with fee: transfer financial_event + expense financial_event
same event_group_id
```

### Exchange fee

Same as transfer fee.

### Investment

```text
investment_buy / investment_sell only record movement between cash buckets and investment buckets.
investment_snapshot records market value, profit/loss, return rate.
investment_income only enters investment statistics.
```

---

## Category and CategoryGroup

No OperationType in v1.

```text
Category = user-owned category
CategoryGroup = system/family reporting aggregation category
```

Example:

```text
用户A: 外卖 -> FOOD, 买菜 -> GROCERY
用户B: 吃饭 -> FOOD, 超市 -> GROCERY
```

Family reports group by `category_group_id`, not `category_id`.

---

## ChannelTemplate

Channels use a template table. Do not split WeChat receive/pay.

```text
income + WECHAT = 微信收款
expense + WECHAT = 微信支付
income + ALIPAY = 支付宝收款
expense + ALIPAY = 支付宝支付
```

Do not create:

```text
WECHAT_PAY
WECHAT_RECEIVE
ALIPAY_PAY
ALIPAY_RECEIVE
```

Channel examples:

```text
WECHAT
ALIPAY
BANK_TRANSFER
BANK_CARD
CREDIT_CARD
UNIONPAY
CASH
AUTO_DEBIT
REFUND_ORIGINAL
POS
QR_CODE
OTHER
```

---

## Currency and exchange rate

```text
currency = global platform dictionary
base_currency = reporting currency
```

Both user and family have changeable base currency:

```text
user.base_currency
family.base_currency
```

Default:

```text
CNY
```

Exchange rate strategy:

```text
Use free daily FX API first.
Prefer local exchange_rate cache.
Allow manual input if API fails.
Allow user to override rate.
financial_event stores final exchange_rate snapshot.
```

Suggested first API:

```text
Frankfurter
```

---

## Personal and family reports

Product definition:

```text
账单：看我干了啥
资产：看我有啥
```

Separate:

```text
账单统计
资产报表
```

Bill statistics use:

```text
event_type + category/category_group + time range
```

Asset reports use:

```text
bucket.balance + bucket_nature + currency
```

Family bill statistics do not group by Account/Bucket/Channel.

Family bill dimensions:

```text
member
income/expense
CategoryGroup
trend
time range
```

Family asset report:

```text
Show each family member's own Account/Bucket.
No family-shared Account.
No forced aggregation of all users' WeChat/Alipay accounts.
```

Family asset totals:

```text
family_total_assets = all members' asset buckets converted to family.base_currency
family_total_liabilities = all members' liability buckets converted to family.base_currency
family_net_worth = assets - liabilities
```

Transaction/financial event details must show Account/Bucket/Channel.

---

## Import recognition design

Separate design doc:

```text
导入识别与规则学习设计方案.md
```

Core idea:

```text
ImportRecord
FinancialEventDraft
ImportMappingRule
MerchantAlias
MerchantCategoryStats
```

Layered recognition:

```text
rules
merchant normalization
history statistics
AI suggestion later
```

This is intentionally parked for later, after core model stabilizes.

---

## Suggested new repository CLAUDE.md seed

```markdown
# Marmot Ledger Project Instructions

Marmot Ledger is a personal and family financial ledger for tracking cash flow, assets, liabilities, multi-currency accounts, investments, reimbursements, deposits, and shared household finances.

Core models:
- financial_event: a financial event, what happened.
- ledger_entry: balance changes caused by a financial event.
- bucket: asset/liability/virtual funds container.
- account: platform/institution/grouping; never stores balance.
- category: user-owned category.
- category_group: reporting aggregation category.
- channel_template: payment/receive channel template.

Key decisions:
- Use financial_event, not transaction, as the main event table.
- ledger_entry uses signed amount; no direction field.
- Account does not store balance.
- Bucket is a core model.
- Bucket initial balance creates a balance_adjustment financial_event and ledger_entry.
- Bill reports answer “what happened”; asset reports answer “what I have”.
- Bill statistics do not group by Account/Bucket.
- Family does not need shared accounts; family asset reports show members' own accounts/buckets.
- investment_income only enters investment statistics.
- refund offsets expense and is not income.
```
