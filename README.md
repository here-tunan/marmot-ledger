# Marmot Ledger

<p align="center">
  <img src="img/marmot-ledger-1.png" alt="Marmot Ledger" width="180" />
  <img src="img/marmot-ledger-2.png" alt="Marmot Ledger" width="180" />
</p>

<p align="center">
  <strong>A calm personal and family finance ledger.</strong><br />
  <a href="README.zh-CN.md">简体中文</a> · English
</p>

Marmot Ledger is a finance ledger for people who want to understand everyday cash flow, assets, liabilities, multi-currency accounts, investments, reimbursements, deposits, and shared household finance without turning their personal life into an accounting system.

---

## Highlights

- **Accounts** — organize cash, wallets, banks, cards, investment platforms, and other places where money lives.
- **Buckets** — track real balances for assets, liabilities, receivables, deposits, investment cash, and investment assets.
- **Records** — record income, expense, transfer, refund, exchange, receivable/deposit/loan flows, split bills, family transfers, and investment activities.
- **Ledger entries** — every balance-changing event creates signed entries so Bucket balances stay traceable.
- **Categories** — manage flexible personal categories and family-level reporting groups.
- **Channels** — manage payment and receipt channels such as WeChat, Alipay, cash, bank transfer, cards, and custom channels.
- **Dashboard** — view net worth, assets, liabilities, income/expense, refunds, trends, category statistics, and investment P/L by original currency.
- **Family workspace** — invite members, aggregate active-member records, and keep family statistics without creating shared family accounts.
- **Admin templates** — manage account, channel, and category system templates.
- **Bilingual UI** — Chinese and English i18n are built in.

---

## Product principles

1. **Ledger first, not transaction first.** The core model is `financial_event` + `ledger_entry` + `bucket`, not a generic `transaction` table.
2. **Accounts group; Buckets hold balances.** Accounts are platforms or institutions. Buckets store current balances.
3. **Signed entries.** `ledger_entry.amount > 0` increases a Bucket; `amount < 0` decreases it.
4. **Original currency reporting.** Ordinary reports do not merge different currencies by default.
5. **Refund offsets expense.** Refunds reduce expense instead of becoming ordinary income.
6. **Family data stays member-owned.** Family reports aggregate active members’ own records and Buckets; Marmot Ledger does not create shared family Accounts.

---

## Tech stack

### Backend

- Go
- Fiber
- Xorm
- MySQL
- Redis

### Frontend

- Vue 3
- Vite
- Element Plus
- Pinia
- Vue Router
- Axios
- vue-i18n
- ECharts

---

## Quick start

### 1. Prepare services

Install and start:

- MySQL 8.x recommended
- Redis 6.x+ recommended

Update local configuration if needed:

```text
env/dev.yaml
```

### 2. Create database

For a fresh local database, run the single bootstrap file:

```bash
mysql < sql/marmot_ledger.sql
```

This file creates the database, current schema, indexes, foreign keys, and required reference data. It does **not** create demo users or demo records.

See [`docs/SQL_DEVELOPER_GUIDE.md`](docs/SQL_DEVELOPER_GUIDE.md) for the full table guide and data policy.

### 3. Start backend

```bash
go run main.go
```

By default the backend reads `env/dev.yaml`. To switch environment:

```bash
export MARMOT_LEDGER_ENV=prod
```

### 4. Start frontend

```bash
npm --prefix ui install
npm --prefix ui run dev
```

Then open the Vite dev URL shown in your terminal.

---

## Validation

Backend:

```bash
go test ./...
go build ./...
```

Frontend:

```bash
npm --prefix ui run build
```

Fresh database smoke check:

```bash
mysql < sql/marmot_ledger.sql
```

Then verify in the app:

1. Register and log in.
2. Create an Account.
3. Create a Bucket.
4. Create income, expense, refund, and transfer records.
5. Check Dashboard statistics by original currency.
6. Create a Family workspace and invite/accept a member if family testing is needed.

---

## Project structure

```text
.
├── api/                         # Fiber API routes
├── env/                         # Environment configuration
├── internal/
│   ├── domain/                  # Domain entities and repositories
│   ├── infrastructure/          # MySQL / Redis initialization
│   └── service/                 # Application services and record strategies
├── pkg/                         # Shared models, errors, responses, utilities
├── sql/
│   └── marmot_ledger.sql        # Single fresh database bootstrap
├── ui/                          # Vue 3 frontend
│   └── src/
│       ├── api/                 # Frontend API clients
│       ├── components/          # Layout and reusable components
│       ├── i18n/                # Chinese / English locales
│       ├── router/              # Vue Router
│       ├── stores/              # Pinia stores
│       ├── styles/              # Global styles
│       └── views/               # Product pages
├── docs/
│   ├── DESIGN.md                # Calm Marmot Finance UI direction
│   └── SQL_DEVELOPER_GUIDE.md   # Database guide
└── img/                         # Brand images
```

---

## Domain model

```text
User
  ├── Account -> Bucket -> LedgerEntry
  ├── FinancialEvent -> LedgerEntry
  ├── Category
  ├── PersonalChannel
  ├── Family / FamilyMember
  └── Templates / reference data
```

| Model | Meaning |
| --- | --- |
| `financial_event` | What happened: income, expense, refund, transfer, exchange, investment, etc. |
| `ledger_entry` | How Bucket balances changed because of the event. |
| `bucket` | Asset/liability/funds container that stores current balance. |
| `account` | Platform, institution, or grouping layer. It does not store balance. |
| `category` | User-owned income/expense category. |
| `category_group` | System reporting-group reference vocabulary. |
| `family_category_group` | Family-owned reporting group for shared family statistics. |
| `personal_channel` | User-owned payment/receipt channel. |
| `channel_template` | System channel template. |
| `currency` | Shared currency dictionary. |

---

## Documentation

- [`docs/SQL_DEVELOPER_GUIDE.md`](docs/SQL_DEVELOPER_GUIDE.md) — schema, reference data, domain flows, and fresh database setup.
- [`docs/DESIGN.md`](docs/DESIGN.md) — Calm Marmot Finance brand and UI direction.

---

## Development conventions

- Keep new finance features based on `financial_event`, `ledger_entry`, and `bucket`.
- Do not restore old unrelated business modules.
- Do not introduce `transaction` as the primary domain table.
- All balance-changing operations must run in a database transaction.
- Chinese UI uses **资金桶** for Bucket; English UI keeps **Bucket**.
- New visible UI copy should go through i18n.
- Ordinary statistics should stay original-currency based unless an explicit FX feature is designed.
- The SQL bootstrap must not contain demo users, predictable admin credentials, personal sample data, or hardcoded `user_id = 1` records.

---

## Post-MVP roadmap

- Improve family-mode dashboard selection and member views.
- Add import recognition and rule-learning workflows when the product needs import automation.
- Continue polishing i18n and remove hardcoded UI copy.
- Harden auth/token payloads so password fields are never exposed.
- Introduce production migration tooling when the project has deployed databases that need incremental upgrades.

---

## License

No license has been selected yet.
