package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"go-my-life/internal/domain/repository/moneydb"
	"go-my-life/pkg/model"
	"go-my-life/pkg/utils"
	"io"
	"log"
	"mime/multipart"
	"strconv"
	"strings"
	"sync"
)

// ProcessAlipayCSV 处理支付宝CSV账单文件
func ProcessAlipayCSV(userId int64, csvFile multipart.File) (*ProcessResult, error) {
	result := &ProcessResult{
		Transactions: []moneydb.Transaction{},
		Warnings:     []string{},
	}

	// 读取文件内容
	fileContent, err := io.ReadAll(csvFile)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	// 检测并转换编码
	convertedContent, encoding, encodingErr := utils.DetectAndConvertEncoding(fileContent)
	result.Encoding = encoding
	if encodingErr != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("编码检测警告: %v", encodingErr))
		log.Printf("编码检测警告: %v", encodingErr)
	}
	log.Printf("检测到文件编码: %s", encoding)

	// 创建 CSV Reader
	reader := csv.NewReader(bytes.NewReader(convertedContent))
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true

	// 读取 CSV 数据
	records, err := reader.ReadAll()
	if err != nil {
		log.Print("Failed to read Alipay CSV record:", err)
		return nil, err
	}
	if len(records) == 0 {
		return result, nil
	}

	// 验证是支付宝账单格式
	if len(records) < 4 || !strings.Contains(records[3][0], "支付宝账户") {
		return nil, fmt.Errorf("不是有效的支付宝CSV账单格式")
	}

	fmt.Print("处理支付宝账单!")
	after, _ := strings.CutPrefix(records[2][0], "姓名：")
	fmt.Printf("支付宝用户：%s\n", after)

	// 处理支付宝账单
	transactions := processAlipayBill(userId, records)
	result.Transactions = transactions

	return result, nil
}

func processAlipayBill(userId int64, records [][]string) []moneydb.Transaction {
	var transactions []moneydb.Transaction
	var mu sync.Mutex
	var wg sync.WaitGroup
	isBillStart := false

	for _, record := range records {
		if len(record) == 0 {
			continue
		}

		if strings.HasPrefix(record[0], "交易时间") {
			isBillStart = true
			continue
		}

		if !isBillStart {
			continue
		}

		wg.Add(1)
		go func(record []string) {
			defer wg.Done()

			// 添加边界检查
			if len(record) < 7 {
				return
			}

			// 开始获取真正的账单行数据
			// 第二列 交易分类（筛选掉投资理财的）
			if record[1] == "投资理财" {
				return
			}

			// 第一列：交易时间
			timeStr := record[0]
			localTime, _ := utils.ParseTimeString(timeStr)

			// 第三列：交易对方; 第五列：商品
			description := record[2] + "_" + record[4]

			// 先暂且不管
			category := AnalysisCategory(description)

			// 第六列：收入/支出
			typeId := 1
			if record[5] == "收入" {
				typeId = 1
			} else {
				typeId = 2
			}

			// 第七列：金额(去除¥符号)
			amountStr := strings.TrimPrefix(record[6], "¥")
			amount, _ := strconv.ParseFloat(amountStr, 64)

			// 第八列收付款方式 有一个要特殊处理一下
			account := 2 // 默认支付宝
			if strings.Contains(record[7], "关爱通") {
				account = 12 // 关爱通
			}

			transaction := moneydb.Transaction{
				Amount:      amount,
				Description: description,
				UserId:      userId,
				Type:        typeId,
				Category:    category,
				Account:     account,
				Time:        model.LocalTime(localTime),
			}

			mu.Lock()
			transactions = append(transactions, transaction)
			mu.Unlock()
		}(record)
	}

	wg.Wait()
	return transactions
}
