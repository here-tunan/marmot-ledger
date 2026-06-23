# Marmot Ledger SQL Developer Guide

> This is the single database guide for Marmot Ledger. For a fresh database, use only `sql/marmot_ledger.sql`.

## 1. Fresh database setup

```bash
mysql < sql/marmot_ledger.sql
```

The bootstrap file creates the database, creates all current MVP tables, and inserts required reference data. It does **not** create demo users, demo accounts, or demo records.

```text
sql/marmot_ledger.sql
├── CREATE DATABASE marmot_ledger
├── current schema
└── required reference data
```

## 2. Core rules

- `financial_event` records what happened.
- `ledger_entry` records how Bucket balances changed.
- `bucket` stores current balances.
- `account` is a platform/institution/grouping layer and does not store balances.
- `ledger_entry.amount` is signed: positive means Bucket balance increases, negative means it decreases.
- Do not introduce a new `transaction` core table.
- Ordinary statistics use original event currency; different currencies are not merged by default.
- Refunds offset expense and are not counted as ordinary income.
- Family statistics aggregate active family members by original record currency.

## 3. Schema overview

| Table | Purpose | Key notes |
| --- | --- | --- |
| `user` | Login user and profile | Has `role` for admin access; normal users default to `user`. |
| `currency` | Enabled currency dictionary | Seeded with CNY, USD, HKD, EUR, JPY, GBP, SGD, AUD, NZD. |
| `account_template` | System account templates | Used by admin/template management and quick creation. |
| `category_group` | System reporting groups | Used as reference grouping vocabulary. |
| `category_template` | System category templates | Used by category quick creation and admin template management. |
| `channel_template` | System channel templates | One channel such as WECHAT supports income/expense/refund/transfer by event type. |
| `personal_channel` | User-owned payment/receipt channels | Can be imported from templates or created manually. |
| `account` | User-owned platform/institution/group | Does not store balance. Buckets under the account store balances. |
| `bucket` | Asset, liability, or virtual funds container | Stores `balance`, `initial_balance`, `bucket_type`, and `bucket_nature`. |
| `category` | User-owned income/expense category | Flexible personal categorization. |
| `financial_event` | Ledger event | Stores event type, category, channel, time, original currency/amount, stats flag, and status. |
| `ledger_entry` | Balance entry | Stores signed Bucket balance delta and `balance_after`. |
| `family` | Family workspace | Owned by one user; no shared family account is created. |
| `family_member` | Family membership | Tracks role and invite/join/leave status. |
| `family_category_group` | Family custom reporting group | Lets a family define shared reporting groups. |
| `family_category_group_member` | Family group/category mapping | Maps user categories into family-level groups. |

## 4. Table details

### `user`

Important columns:

- `account`: unique login account.
- `password`: password credential field. Do not expose this field in API responses or token payloads.
- `name`: display name.
- `role`: `user` or `admin`.
- `is_deleted`: soft-delete flag.
- `gmt_create`, `gmt_modified`: legacy-style user timestamps.

Indexes:

- `uk_user_account(account)`
- `idx_user_role(role)`
- `idx_user_is_deleted(is_deleted)`

### `currency`

Important columns:

- `code`: ISO-like 3-letter currency code.
- `name`, `symbol`, `precision_digits`: display and amount formatting metadata.
- `enabled`, `sort`: frontend/admin ordering.

Seeded currencies:

| Code | Name | Symbol |
| --- | --- | --- |
| CNY | 人民币 | ¥ |
| USD | 美元 | $ |
| HKD | 港币 | HK$ |
| EUR | 欧元 | € |
| JPY | 日元 | ¥ |
| GBP | 英镑 | £ |
| SGD | 新加坡元 | S$ |
| AUD | 澳大利亚元 | A$ |
| NZD | 新西兰元 | NZ$ |

### `account_template`

Used for quick account creation and admin template management.

Important columns:

- `provider_code`: unique stable code such as `WECHAT`, `ALIPAY`, `CMB`, `CASH`.
- `type`: `cash`, `wallet`, `bank`, `credit`, `investment`, `liability`, or `other`.
- `icon`, `color`: frontend display metadata.
- `enabled`, `sort`: admin visibility and order.

### `category_group`

System-level reporting vocabulary. Personal statistics currently group by user category, while family statistics can use family category groups. This table remains useful as a standard reference set.

Seeded groups include expense groups such as `FOOD`, `GROCERY`, `TRANSPORT`, `HOUSING`, `SHOPPING`, `TRAVEL`, and income groups such as `SALARY`, `BONUS`, `INVESTMENT_INCOME`, `REIMBURSEMENT`.

### `category_template`

Used to create user categories quickly.

Important columns:

- `template_code`: unique stable code.
- `name`: default display name.
- `type`: `income` or `expense`.
- `icon`, `color`: frontend display metadata.
- `enabled`, `sort`: admin visibility and order.

Seeded templates include rich expense categories such as dining, grocery, transport, utilities, shopping, healthcare, education, travel, personal care, pets, and family/children; income templates include salary, bonus, dividends, interest, reimbursement, side income, rental income, and gifts.

### `channel_template`

System channel templates for payment and receipt channels.

Important columns:

- `channel_code`: unique stable code.
- `channel_type`: `wallet`, `bank`, `card`, `cash`, `system`, or `other`.
- `provider_code`: optional relation-like provider code.
- `supported_event_types`: comma-separated event types.

Important rule:

```text
income + WECHAT  = WeChat receipt / 微信收款
expense + WECHAT = WeChat payment / 微信支付
```

Do not split channel codes into `WECHAT_PAY` and `WECHAT_RECEIVE`.

### `personal_channel`

User-owned channel created from templates or manually.

Important columns:

- `user_id`: owner.
- `channel_template_id`: optional source template.
- `name`, `channel_type`, `provider_code`, `supported_event_types`, `icon`.
- `is_active`, `is_deleted`: lifecycle flags.

Uniqueness:

- `uk_personal_channel_user_name_deleted(user_id, name, is_deleted)` prevents duplicate active names while allowing soft-deleted history.

### `account`

Account means platform, institution, or grouping layer.

Important columns:

- `user_id`: owner.
- `name`: user-visible account name.
- `type`: account type.
- `icon`, `color`: display metadata.
- `is_active`, `is_deleted`: lifecycle flags.

Account does not store balance. Use Bucket aggregation for balances.

### `bucket`

Bucket / 资金桶 stores the actual balance.

Important columns:

- `user_id`: owner.
- `account_id`: parent account.
- `currency`: references `currency.code`.
- `balance`: current balance.
- `initial_balance`: user-entered opening balance.
- `bucket_type`: cash, wallet, bank, credit, investment asset, receivable, deposit, liability, etc.
- `bucket_nature`: `asset` or `liability`.
- `is_active`, `is_deleted`: lifecycle flags.

Net worth rule:

```text
net_worth = sum(asset bucket balances) - sum(liability bucket balances)
```

### `category`

User-owned category for income/expense records.

Important columns:

- `user_id`: owner.
- `name`: display name.
- `type`: `income` or `expense`.
- `icon`, `color`: display metadata.
- `is_active`, `is_deleted`: lifecycle flags.

Uniqueness:

- `uk_category_user_type_name_deleted(user_id, type, name, is_deleted)`.

### `financial_event`

Financial event records the business fact.

Important columns:

- `event_group_id`: groups related events for split, exchange, paired, or investment scenarios.
- `related_financial_event_id`: links settlement/collection/refund-like events.
- `event_type`: event type.
- `category_id`: optional user category.
- `channel_id`: optional personal channel.
- `event_time`: business time.
- `currency`, `amount`: original currency and amount.
- `include_in_statistics`: whether ordinary income/expense stats include this event.
- `source`: usually `manual`.
- `status`: usually `active`.
- `is_deleted`: soft-delete flag.

Primary MVP event types:

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

Statistics rules:

- `income` and `expense` normally enter ordinary statistics.
- `refund` offsets expense.
- `transfer`, `exchange`, receivable/deposit/loan principal flows, investment buy/sell, and `balance_adjustment` do not enter ordinary statistics by default.
- `investment_income` belongs to investment reporting and should not be treated as ordinary salary-like income.

### `ledger_entry`

Ledger entries record Bucket balance changes.

Important columns:

- `financial_event_id`: parent event.
- `bucket_id`: affected Bucket.
- `currency`: entry currency.
- `amount`: signed balance delta.
- `balance_after`: Bucket balance after applying this entry.
- `entry_role`: semantic role in multi-entry events.

Examples:

```text
Expense from CNY cash:
financial_event.event_type = expense
ledger_entry.amount = -88.00
ledger_entry.entry_role = expense

Transfer CNY cash to CNY bank:
ledger_entry #1 amount = -1000.00, entry_role = transfer_out
ledger_entry #2 amount = +1000.00, entry_role = transfer_in
```

### `family`

Family workspace for household finance collaboration.

Important columns:

- `name`: family name.
- `owner_user_id`: creator/owner.
- `is_deleted`: soft-delete flag.

Family reports aggregate member data. A family does not own shared accounts or shared Buckets.

### `family_member`

Tracks family membership and invitation state.

Important columns:

- `family_id`, `user_id`.
- `role`: `owner`, `admin`, or `member`.
- `status`: `invited`, `active`, `rejected`, or `left`.
- `display_name`: name inside the family workspace.
- `invited_by_user_id`, `invited_at`, `joined_at`, `left_at`.

### `family_category_group`

Family-owned custom reporting group.

Important columns:

- `family_id`: owner family.
- `name`: group name.
- `type`: `income` or `expense`.
- `icon`, `color`.
- `created_by_user_id`.
- `sort`, `is_active`, `is_deleted`.

### `family_category_group_member`

Many-to-many relation between family category groups and user categories.

Important columns:

- `family_group_id`: family category group.
- `category_id`: user category.
- `added_by_user_id`: who added the mapping.
- `added_at`: mapping time.

## 5. Common domain flows

### Bucket creation

When a user creates a Bucket:

1. Insert `bucket` with `initial_balance` and `balance`.
2. Insert `financial_event` with:
   - `event_type = balance_adjustment`
   - `include_in_statistics = false`
3. Insert one `ledger_entry` with:
   - `amount = initial_balance`
   - `balance_after = initial_balance`
   - `entry_role = adjustment`

All steps must run in one database transaction.

### Income

- One `financial_event` with `event_type = income`.
- One positive `ledger_entry` into the selected Bucket.
- Usually included in ordinary statistics.

### Expense

- One `financial_event` with `event_type = expense`.
- One negative `ledger_entry` from the selected Bucket.
- Usually included in ordinary statistics.

### Refund

- One `financial_event` with `event_type = refund`.
- One positive `ledger_entry` into the selected Bucket.
- Counts as expense offset, not ordinary income.

### Transfer

- One `financial_event` with `event_type = transfer`.
- Two `ledger_entry` rows:
  - negative from source Bucket
  - positive into target Bucket
- Excluded from ordinary income/expense statistics.

### Family statistics

- Uses active members from `family_member`.
- Aggregates members’ own events and Buckets.
- Does not create family-owned shared accounts.
- Does not merge different currencies by default.

## 6. Reference data policy

Reference data in `sql/marmot_ledger.sql` is required and safe to rerun.

Do:

- Keep required dictionaries/templates in the unified SQL file.
- Use stable codes such as `WECHAT`, `CASH`, `FOOD_DINING`.
- Use `ON DUPLICATE KEY UPDATE` for reference rows.

Do not:

- Add demo users to the normal bootstrap.
- Add demo accounts, Buckets, or records to the normal bootstrap.
- Hardcode `user_id = 1` in seed logic.
- Ship predictable default admin credentials.

## 7. Fresh database smoke checklist

After running:

```bash
mysql < sql/marmot_ledger.sql
```

Verify:

1. Backend starts with `go run main.go`.
2. Frontend starts with `npm --prefix ui run dev`.
3. A new user can register and log in.
4. Currencies are available in dropdowns.
5. Account templates, channel templates, and category templates load.
6. A user can create an Account.
7. A user can create a Bucket.
8. Bucket creation creates a balance-adjustment event and ledger entry.
9. A user can create income, expense, refund, and transfer records.
10. Dashboard/statistics display amounts by original currency.
11. Family workspace can invite/accept members and show active-member statistics.

## 8. Future migration policy

The MVP uses one fresh bootstrap file for clarity. Once Marmot Ledger has production deployments that need incremental upgrades, introduce a migration tool and keep migrations separate from the fresh bootstrap.

Suggested future tools:

- Goose
- Atlas
- Flyway
- Liquibase

Until then, keep `sql/marmot_ledger.sql` as the single current schema and reference-data source.
