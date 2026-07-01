package service

const (
	EventTypeIncome            = "income"
	EventTypeExpense           = "expense"
	EventTypeRefund            = "refund"
	EventTypeTransfer          = "transfer"
	EventTypeExchange          = "exchange"
	EventTypeBalanceAdjustment = "balance_adjustment"
	EventTypeReceivableCreate  = "receivable_create"
	EventTypeReceivableCollect = "receivable_collect"
	EventTypeDepositCreate     = "deposit_create"
	EventTypeDepositRefund     = "deposit_refund"
	EventTypeLoanOut           = "loan_out"
	EventTypeLoanCollect       = "loan_collect"
	EventTypeInvestmentBuy     = "investment_buy"
	EventTypeInvestmentSell    = "investment_sell"
	EventTypeInvestmentIncome  = "investment_income"
	EventTypeInvestmentRevalue = "investment_revalue"
	EventTypeFamilyTransfer    = "family_transfer"

	EntryRoleIncome            = "income"
	EntryRoleExpense           = "expense"
	EntryRoleRefund            = "refund"
	EntryRoleTransferOut       = "transfer_out"
	EntryRoleTransferIn        = "transfer_in"
	EntryRoleExchangeOut       = "exchange_out"
	EntryRoleExchangeIn        = "exchange_in"
	EntryRoleAdjustment        = "adjustment"
	EntryRoleCashLeg           = "cash_leg"
	EntryRoleReceivableCreate  = "receivable_create"
	EntryRoleReceivableCollect = "receivable_collect"
	EntryRoleDepositCreate     = "deposit_create"
	EntryRoleDepositRefund     = "deposit_refund"
	EntryRoleLoanOut           = "loan_out"
	EntryRoleLoanCollect       = "loan_collect"
	EntryRoleInvestmentBuy     = "investment_buy"
	EntryRoleInvestmentSell    = "investment_sell"
	EntryRoleInvestmentIncome  = "investment_income"
	EntryRoleRevaluationGain   = "revaluation_gain"
	EntryRoleRevaluationLoss   = "revaluation_loss"
	EntryRoleFamilyTransferOut = "family_transfer_out"
	EntryRoleFamilyTransferIn  = "family_transfer_in"

	EventSourceManual = "manual"
	EventSourceImport = "import"
	EventStatusActive = "active"

	BucketNatureAsset     = "asset"
	BucketNatureLiability = "liability"
)

var validBucketTypes = map[string]bool{
	"cash":             true,
	"wallet":           true,
	"bank":             true,
	"credit":           true,
	"investment_cash":  true,
	"investment_asset": true,
	"receivable":       true,
	"deposit":          true,
	"loan_out":         true,
	"liability":        true,
	"virtual":          true,
}

var assetLikeBucketTypes = map[string]bool{
	"cash":             true,
	"wallet":           true,
	"bank":             true,
	"receivable":       true,
	"deposit":          true,
	"loan_out":         true,
	"investment_cash":  true,
	"investment_asset": true,
	"virtual":          true,
}

var cashSideBucketTypes = map[string]bool{
	"cash":    true,
	"wallet":  true,
	"bank":    true,
	"virtual": true,
}

var investmentBucketTypes = map[string]bool{
	"investment_cash":  true,
	"investment_asset": true,
}

var investmentCashSideBucketTypes = map[string]bool{
	"cash":            true,
	"wallet":          true,
	"bank":            true,
	"virtual":         true,
	"investment_cash": true,
}
