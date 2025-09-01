package service

import "go-my-life/internal/domain/repository/moneydb"

// ProcessResult 处理结果结构
type ProcessResult struct {
	Transactions []moneydb.Transaction `json:"transactions"`
	Warnings     []string              `json:"warnings,omitempty"`
	Encoding     string                `json:"encoding,omitempty"`
}
