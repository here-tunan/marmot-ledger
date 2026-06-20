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

type TrendPoint struct {
	Period       string          `json:"period"`
	Income       decimal.Decimal `json:"income"`
	GrossExpense decimal.Decimal `json:"grossExpense"`
	Refund       decimal.Decimal `json:"refund"`
	NetExpense   decimal.Decimal `json:"netExpense"`
}

type InvestmentSummary struct {
	Currency      string                  `json:"currency"`
	MarketValue   decimal.Decimal         `json:"marketValue"`
	BuyAmount     decimal.Decimal         `json:"buyAmount"`
	RevalueAmount decimal.Decimal         `json:"revalueAmount"`
	IncomeAmount  decimal.Decimal         `json:"incomeAmount"`
	ProfitLoss    decimal.Decimal         `json:"profitLoss"`
	ReturnRate    decimal.Decimal         `json:"returnRate"`
	Buckets       []InvestmentBucketStats `json:"buckets"`
}

type InvestmentBucketStats struct {
	BucketId      int64           `json:"bucketId"`
	BucketName    string          `json:"bucketName"`
	BucketType    string          `json:"bucketType"`
	Currency      string          `json:"currency"`
	MarketValue   decimal.Decimal `json:"marketValue"`
	BuyAmount     decimal.Decimal `json:"buyAmount"`
	RevalueAmount decimal.Decimal `json:"revalueAmount"`
	ProfitLoss    decimal.Decimal `json:"profitLoss"`
	ReturnRate    decimal.Decimal `json:"returnRate"`
}

type NetWorthTrendPoint struct {
	Period    string          `json:"period"`
	Currency  string          `json:"currency"`
	Asset     decimal.Decimal `json:"asset"`
	Liability decimal.Decimal `json:"liability"`
	NetWorth  decimal.Decimal `json:"netWorth"`
}

type FamilyAssets struct {
	Members []MemberAsset `json:"members"`
}

type MemberAsset struct {
	UserId      int64                `json:"userId"`
	Account     string               `json:"account"`
	Name        string               `json:"name"`
	DisplayName string               `json:"displayName"`
	Role        string               `json:"role"`
	Totals      []MemberAssetTotal   `json:"totals"`
	Accounts    []MemberAssetAccount `json:"accounts"`
}

type MemberAssetTotal struct {
	Currency    string          `json:"currency"`
	Asset       decimal.Decimal `json:"asset"`
	Liability   decimal.Decimal `json:"liability"`
	NetWorth    decimal.Decimal `json:"netWorth"`
	BucketCount int64           `json:"bucketCount"`
}

type MemberAssetAccount struct {
	Id      int64               `json:"id"`
	Name    string              `json:"name"`
	Type    string              `json:"type"`
	Buckets []MemberAssetBucket `json:"buckets"`
}

type MemberAssetBucket struct {
	Id           int64           `json:"id"`
	Name         string          `json:"name"`
	Currency     string          `json:"currency"`
	Balance      decimal.Decimal `json:"balance"`
	BucketType   string          `json:"bucketType"`
	BucketNature string          `json:"bucketNature"`
}
