package currencydb

import (
	"errors"
	"marmot-ledger/internal/infrastructure"

	"xorm.io/xorm"
)

type Currency struct {
	Id              int64  `json:"id" xorm:"pk autoincr 'id'"`
	Code            string `json:"code" xorm:"'code'"`
	Name            string `json:"name" xorm:"'name'"`
	Symbol          string `json:"symbol" xorm:"'symbol'"`
	PrecisionDigits int    `json:"precisionDigits" xorm:"'precision_digits'"`
	Enabled         bool   `json:"enabled" xorm:"'enabled'"`
	Sort            int    `json:"sort" xorm:"'sort'"`
}

func (Currency) TableName() string {
	return "currency"
}

func GetEnabledCurrency(session *xorm.Session, code string) (*Currency, error) {
	currency := &Currency{}
	has, err := session.Where("code = ? AND enabled = ?", code, 1).Get(currency)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("currency not found")
	}
	return currency, nil
}

func ListEnabledCurrencies() ([]Currency, error) {
	currencies := make([]Currency, 0)
	err := infrastructure.Mysql.Where("enabled = ?", 1).Asc("sort").Find(&currencies)
	return currencies, err
}
