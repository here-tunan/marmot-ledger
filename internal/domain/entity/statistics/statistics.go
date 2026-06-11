package statistics

import "github.com/shopspring/decimal"

type Summary struct {
	Currency     string          `json:"currency"`
	Income       decimal.Decimal `json:"income"`
	GrossExpense decimal.Decimal `json:"grossExpense"`
	Refund       decimal.Decimal `json:"refund"`
	NetExpense   decimal.Decimal `json:"netExpense"`
	Expense      decimal.Decimal `json:"expense"`
	Net          decimal.Decimal `json:"net"`
	EventCount   int64           `json:"eventCount"`
}

type CategoryGroupItem struct {
	CategoryGroupId   int64           `json:"categoryGroupId"`
	CategoryGroupCode string          `json:"categoryGroupCode"`
	CategoryGroupName string          `json:"categoryGroupName"`
	Type              string          `json:"type"`
	Amount            decimal.Decimal `json:"amount"`
	Count             int64           `json:"count"`
	Color             string          `json:"color"`
	Icon              string          `json:"icon"`
}

type CategoryGroupStats struct {
	Currency string              `json:"currency"`
	Type     string              `json:"type"`
	Items    []CategoryGroupItem `json:"items"`
}
