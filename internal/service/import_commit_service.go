package service

import (
	"errors"
	"marmot-ledger/internal/domain/entity/record"
	"marmot-ledger/internal/domain/repository/currencydb"
	"marmot-ledger/internal/infrastructure"
	"strings"

	"xorm.io/xorm"
)

const importCommitChunkSize = 100

// 只支持这几种 scenario 走批量入账；其他类型（split/investment_sell/family_transfer）自开事务，与本流程冲突。
var supportedImportScenarios = map[string]bool{
	EventTypeIncome:   true,
	EventTypeExpense:  true,
	EventTypeRefund:   true,
	EventTypeTransfer: true,
}

// ImportCommitRow 表示前端从预览页提交的一行入账请求。RowIndex 让后端能把结果映射回预览行。
type ImportCommitRow struct {
	RowIndex int                  `json:"rowIndex"`
	Request  record.RecordRequest `json:"request"`
}

type ImportCommitFailure struct {
	RowIndex int    `json:"rowIndex"`
	Error    string `json:"error"`
}

type ImportCommitResult struct {
	SuccessCount      int                   `json:"successCount"`
	FailedCount       int                   `json:"failedCount"`
	Failures          []ImportCommitFailure `json:"failures"`
	SuccessRowIndexes []int                 `json:"successRowIndexes"`
}

// CommitImportRecords 把批量行按 importCommitChunkSize 分组，每组一个事务；
// 组内任一行失败则回滚整组，同组行统一标 failed（附首个失败原因）。
func CommitImportRecords(userId int64, rows []ImportCommitRow) (*ImportCommitResult, error) {
	result := &ImportCommitResult{
		Failures:          []ImportCommitFailure{},
		SuccessRowIndexes: []int{},
	}
	if len(rows) == 0 {
		return result, nil
	}

	for start := 0; start < len(rows); start += importCommitChunkSize {
		end := start + importCommitChunkSize
		if end > len(rows) {
			end = len(rows)
		}
		chunk := rows[start:end]
		outcomes := commitImportChunk(userId, chunk)
		for i, outcome := range outcomes {
			rowIndex := chunk[i].RowIndex
			if outcome.err != nil {
				result.Failures = append(result.Failures, ImportCommitFailure{RowIndex: rowIndex, Error: outcome.err.Error()})
				result.FailedCount++
			} else {
				result.SuccessRowIndexes = append(result.SuccessRowIndexes, rowIndex)
				result.SuccessCount++
			}
		}
	}
	return result, nil
}

// rowOutcome 记录单行入账结果；err 为 nil 表示成功。
type rowOutcome struct {
	err error
}

// commitImportChunk 每行一个独立事务：坏行只影响自己，同片其他行照常 commit。
// 前端保持 100 行/片是为了 UX（进度条 + 避开 axios timeout），后端不再把片内所有行绑一个事务。
func commitImportChunk(userId int64, chunk []ImportCommitRow) []rowOutcome {
	outcomes := make([]rowOutcome, len(chunk))
	for i := range chunk {
		req := &chunk[i].Request
		req.Source = EventSourceImport
		outcomes[i] = rowOutcome{err: commitOneImportRowInOwnTx(userId, req)}
	}
	return outcomes
}

// commitOneImportRowInOwnTx 每行自开 session 与事务，成功即 commit，失败即 rollback。
func commitOneImportRowInOwnTx(userId int64, req *record.RecordRequest) error {
	session := infrastructure.Mysql.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}
	committed := false
	defer func() {
		if !committed {
			_ = session.Rollback()
		}
	}()
	if err := commitOneImportRow(session, userId, req); err != nil {
		return err
	}
	if err := session.Commit(); err != nil {
		return err
	}
	committed = true
	return nil
}

// commitOneImportRow 复用现有 record pipeline，但共享外部 session、不 commit。
func commitOneImportRow(session *xorm.Session, userId int64, req *record.RecordRequest) error {
	scenario := strings.TrimSpace(req.Scenario)
	if !supportedImportScenarios[scenario] {
		return errors.New("scenario not supported for import: " + scenario)
	}
	if err := validateRecordRequest(req); err != nil {
		return err
	}

	currency := strings.ToUpper(strings.TrimSpace(req.Currency))
	if _, err := currencydb.GetEnabledCurrency(session, currency); err != nil {
		return err
	}

	ctx, err := buildRecordContext(session, userId, req, currency)
	if err != nil {
		return err
	}
	strategy, ok := recordStrategies[scenario]
	if !ok {
		return errors.New("record scenario is unsupported")
	}
	buildResult, err := strategy.Build(ctx)
	if err != nil {
		return err
	}
	return persistRecordBuildResult(session, userId, buildResult)
}
