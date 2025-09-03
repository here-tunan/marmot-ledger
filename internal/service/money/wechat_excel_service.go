package service

import (
	"bytes"
	"fmt"
	"go-my-life/internal/domain/repository/moneydb"
	"go-my-life/pkg/model"
	"go-my-life/pkg/utils"
	"io"
	"log"
	"mime/multipart"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/xuri/excelize/v2"
)

// ProcessWeChatXLSX 处理微信XLSX账单文件
func ProcessWeChatXLSX(userId int64, xlsxFile multipart.File) (*ProcessResult, error) {
	result := &ProcessResult{
		Transactions: []moneydb.Transaction{},
		Warnings:     []string{},
		Encoding:     "XLSX",
	}

	// 读取XLSX文件
	fileContent, err := io.ReadAll(xlsxFile)
	if err != nil {
		return nil, fmt.Errorf("读取XLSX文件失败: %v", err)
	}

	// 打开Excel文件
	f, err := excelize.OpenReader(bytes.NewReader(fileContent))
	if err != nil {
		return nil, fmt.Errorf("打开XLSX文件失败: %v", err)
	}
	defer f.Close()

	// 获取工作表名称
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return result, fmt.Errorf("XLSX文件中没有工作表")
	}

	// 使用第一个工作表
	sheetName := sheets[0]
	log.Printf("处理微信XLSX账单，工作表: %s", sheetName)

	// 读取所有行
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("读取工作表数据失败: %v", err)
	}

	if len(rows) == 0 {
		return result, nil
	}

	// 处理微信XLSX账单
	transactions, processErr := processWeChatXLSXRows(userId, rows)
	result.Transactions = transactions
	if processErr != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("处理XLSX时发生错误: %v", processErr))
	}

	return result, nil
}

func processWeChatXLSXRows(userId int64, rows [][]string) ([]moneydb.Transaction, error) {
	var transactions []moneydb.Transaction
	var mu sync.Mutex
	var wg sync.WaitGroup
	isBillStart := false

	for _, row := range rows {
		if len(row) == 0 {
			continue
		}

		// 找到账单开始标志
		if strings.HasPrefix(row[0], "微信昵称") || strings.Contains(row[0], "交易时间") {
			isBillStart = true
			continue
		}

		if !isBillStart || len(row) < 8 {
			continue
		}

		wg.Add(1)
		go func(record []string) {
			defer wg.Done()

			// 添加边界检查
			if len(record) < 8 {
				return
			}

			// 第一列：交易时间
			timeStr := record[0]
			localTime, timeErr := utils.ParseTimeString(timeStr)
			if timeErr != nil {
				log.Printf("时间解析失败: %s, 错误: %v", timeStr, timeErr)
				return
			}

			// 第三列：交易对方; 第四列：商品
			description := record[2] + "_" + record[3]

			// 第五列：收入/支出
			typeId := 1
			if record[4] == "收入" {
				typeId = 1
			} else {
				typeId = 2
			}

			// 第六列：金额(去除¥符号)
			amountStr := strings.TrimPrefix(record[5], "¥")
			amount, amountErr := strconv.ParseFloat(amountStr, 64)
			if amountErr != nil {
				log.Printf("金额解析失败: %s, 错误: %v", amountStr, amountErr)
				return
			}

			// 第八列：当前状态 (筛选掉已全额退款的)
			if record[7] == "已全额退款" {
				return
			}

			category := AnalysisCategory(description, typeId)

			// 查找或创建微信账户
			wechatAccountId, err := FindOrCreateAccount(userId, "微信支付")
			if err != nil {
				wechatAccountId = 0 // 使用默认值
			}

			transaction := moneydb.Transaction{
				Amount:      amount,
				Description: description,
				UserId:      userId,
				Type:        typeId,
				Category:    category,
				Account:     int(wechatAccountId),
				Time:        model.LocalTime(localTime),
			}

			mu.Lock()
			transactions = append(transactions, transaction)
			mu.Unlock()
		}(row)
	}

	wg.Wait()

	// 按支付时间排序（从小到大）
	sort.Slice(transactions, func(i, j int) bool {
		return time.Time(transactions[i].Time).Before(time.Time(transactions[j].Time))
	})

	return transactions, nil
}
