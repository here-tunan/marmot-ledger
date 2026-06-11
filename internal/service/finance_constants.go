package service

const (
	EventTypeIncome            = "income"
	EventTypeExpense           = "expense"
	EventTypeRefund            = "refund"
	EventTypeTransfer          = "transfer"
	EventTypeBalanceAdjustment = "balance_adjustment"

	EntryRoleIncome      = "income"
	EntryRoleExpense     = "expense"
	EntryRoleRefund      = "refund"
	EntryRoleTransferOut = "transfer_out"
	EntryRoleTransferIn  = "transfer_in"
	EntryRoleAdjustment  = "adjustment"

	EventSourceManual = "manual"
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
