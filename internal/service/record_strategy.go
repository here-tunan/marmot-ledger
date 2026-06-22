package service

import (
	"marmot-ledger/internal/domain/entity/record"
	"marmot-ledger/internal/domain/repository/bucketdb"
	"marmot-ledger/internal/domain/repository/categorydb"
	"marmot-ledger/internal/domain/repository/financialeventdb"
	"marmot-ledger/internal/domain/repository/ledgerentrydb"
	"marmot-ledger/pkg/model"

	"github.com/shopspring/decimal"
	"xorm.io/xorm"
)

type RecordStrategy interface {
	Build(ctx *RecordContext) (*RecordBuildResult, error)
}

type RecordContext struct {
	UserId    int64
	Request   *record.RecordRequest
	Session   *xorm.Session
	Amount    decimal.Decimal
	Currency  string
	EventTime model.LocalTime
	Buckets   map[int64]*bucketdb.Bucket
	Category  *categorydb.CategoryView
	ChannelId *int64
}

type RecordBuildResult struct {
	Event          *financialeventdb.FinancialEvent
	Entries        []*ledgerentrydb.LedgerEntry
	BucketBalances map[int64]decimal.Decimal
}

var recordStrategies = map[string]RecordStrategy{
	EventTypeIncome:            IncomeRecordStrategy{},
	EventTypeExpense:           ExpenseRecordStrategy{},
	EventTypeRefund:            RefundRecordStrategy{},
	EventTypeTransfer:          TransferRecordStrategy{},
	EventTypeExchange:          ExchangeRecordStrategy{},
	EventTypeBalanceAdjustment: BalanceAdjustmentRecordStrategy{},
	EventTypeReceivableCreate:  ReceivableCreateRecordStrategy{},
	EventTypeReceivableCollect: ReceivableCollectRecordStrategy{},
	EventTypeDepositCreate:     DepositCreateRecordStrategy{},
	EventTypeDepositRefund:     DepositRefundRecordStrategy{},
	EventTypeLoanOut:           LoanOutRecordStrategy{},
	EventTypeLoanCollect:       LoanCollectRecordStrategy{},
	EventTypeInvestmentBuy:     InvestmentBuyRecordStrategy{},
	EventTypeInvestmentIncome:  InvestmentIncomeRecordStrategy{},
	EventTypeInvestmentRevalue: InvestmentRevalueRecordStrategy{},
}

func singleBalance(bucketId int64, balance decimal.Decimal) map[int64]decimal.Decimal {
	return map[int64]decimal.Decimal{bucketId: balance}
}

func mapBalance(firstBucketId int64, firstBalance decimal.Decimal, secondBucketId int64, secondBalance decimal.Decimal) map[int64]decimal.Decimal {
	return map[int64]decimal.Decimal{
		firstBucketId:  firstBalance,
		secondBucketId: secondBalance,
	}
}
