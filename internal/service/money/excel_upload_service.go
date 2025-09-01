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
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// detectAndConvertEncoding 检测并转换文件编码
func detectAndConvertEncoding(data []byte) ([]byte, string, error) {
	// 尝试UTF-8（包含BOM）
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		return data[3:], "UTF-8 with BOM", nil
	}

	// 尝试直接当作UTF-8解析
	if isValidUTF8(data) {
		return data, "UTF-8", nil
	}

	// 尝试GBK转UTF-8
	if convertedData, err := convertGBKToUTF8(data); err == nil {
		if isValidUTF8(convertedData) {
			return convertedData, "GBK", nil
		}
	}

	// 返回原始数据，让程序继续尝试处理
	return data, "Unknown", fmt.Errorf("无法识别文件编码，可能导致解析错误")
}

// isValidUTF8 检查数据是否为有效的UTF-8
func isValidUTF8(data []byte) bool {
	return utf8.Valid(data)
}

// convertGBKToUTF8 将GBK编码转换为UTF-8
func convertGBKToUTF8(data []byte) ([]byte, error) {
	// 尝试GBK解码
	reader := transform.NewReader(bytes.NewReader(data), simplifiedchinese.GBK.NewDecoder())
	converted, err := io.ReadAll(reader)
	if err != nil {
		// 尝试GB2312解码
		reader = transform.NewReader(bytes.NewReader(data), simplifiedchinese.HZGB2312.NewDecoder())
		converted, err = io.ReadAll(reader)
		if err != nil {
			// 尝试GB18030解码
			reader = transform.NewReader(bytes.NewReader(data), simplifiedchinese.GB18030.NewDecoder())
			converted, err = io.ReadAll(reader)
		}
	}
	return converted, err
}

// ProcessResult 处理结果结构
type ProcessResult struct {
	Transactions []moneydb.Transaction `json:"transactions"`
	Warnings     []string              `json:"warnings,omitempty"`
	Encoding     string                `json:"encoding,omitempty"`
}

func ProcessTransactionExcel(userId int64, csvFile multipart.File) (*ProcessResult, error) {
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
	convertedContent, encoding, encodingErr := detectAndConvertEncoding(fileContent)
	result.Encoding = encoding
	if encodingErr != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("编码检测警告: %v", encodingErr))
		log.Printf("编码检测警告: %v", encodingErr)
	}
	log.Printf("检测到文件编码: %s", encoding)

	// 创建 CSV Reader
	reader := csv.NewReader(bytes.NewReader(convertedContent))
	// 这些csv文件都不规范，行与行之间的fields num不一样，会导致错误，加上这个就好了
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true // 忽略未被双引号括起来的字段

	// 读取 CSV 数据
	records, err := reader.ReadAll()
	if err != nil {
		log.Print("Failed to read transaction CSV record:", err)
		return nil, err
	}
	if len(records) == 0 {
		return result, nil
	}

	// ZWNBSP 微信的第一个单元格里面有这么个玩意，没法进行匹配, 用第二行进行匹配
	if len(records) > 1 && strings.HasPrefix(records[1][0], "微信昵称") {
		fmt.Print("处理微信账单!")
		fmt.Println(records[1][0])
		result.Transactions = processWeChatBill(userId, records)
		return result, nil
	}

	if len(records) > 3 && strings.Contains(records[3][0], "支付宝账户") {
		fmt.Print("处理支付宝账单!")
		after, _ := strings.CutPrefix(records[2][0], "姓名：")
		fmt.Printf("支付宝用户：%s\n", after)
		result.Transactions = processAlipayBill(userId, records)
		return result, nil
	}

	transactions, processErr := processMyExcel(userId, records)
	result.Transactions = transactions
	if processErr != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("处理CSV时发生错误: %v", processErr))
	}
	return result, nil
}

func processWeChatBill(userId int64, records [][]string) []moneydb.Transaction {
	var transactions []moneydb.Transaction
	var mu sync.Mutex // 用于保证对 transactions 的并发安全
	var wg sync.WaitGroup
	isBillStart := false
	for _, record := range records {
		if isBillStart {
			// 每次启动一个新的goroutine时，计数器加1
			wg.Add(1)
			go func(record []string) {
				defer wg.Done() // 当前 goroutine 执行完后，计数器减 1

				// 开始获取真正的账单行数据
				// 第一列：交易时间
				timeStr := record[0]
				localTime, _ := utils.ParseTimeString(timeStr)
				// 第三列：交易对方; 第四列：商品
				description := record[2] + "_" + record[3]
				// 第五列：收入/支出
				typeId := 1
				if record[4] == "收入" {
					typeId = 1
				} else {
					typeId = 2
				}

				// 先暂且不管
				category := AnalysisCategory(description)

				// 第六列
				amount, _ := strconv.ParseFloat(strings.Trim(record[5], "¥ "), 64)
				transaction := moneydb.Transaction{
					Amount:      amount,
					Description: description,
					UserId:      userId,
					Type:        typeId,
					Category:    category,
					Account:     1,
					Time:        model.LocalTime(localTime),
				}

				// 使用锁确保 transactions 是线程安全的
				mu.Lock()
				transactions = append(transactions, transaction)
				mu.Unlock()
			}(record)

			continue
		}
		if record[0] == "交易时间" {
			isBillStart = true
		}
	}

	// 等待所有 goroutine 完成
	wg.Wait()
	return transactions
}

func processAlipayBill(userId int64, records [][]string) []moneydb.Transaction {
	var transactions []moneydb.Transaction
	isBillStart := false
	for _, record := range records {
		if isBillStart {
			// 开始获取真正的账单行数据
			// 第二列 交易分类（筛选掉投资理财的）
			if record[1] == "投资理财" {
				continue
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
			// 第七列 金额
			amount, _ := strconv.ParseFloat(record[6], 64)

			// 第八列收付款方式 有一个要特殊处理一下
			account := 2
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
			transactions = append(transactions, transaction)
			continue
		}
		if record[0] == "交易时间" {
			isBillStart = true
		}
	}
	return transactions
}

func processMyExcel(userId int64, records [][]string) ([]moneydb.Transaction, error) {
	var transactions []moneydb.Transaction
	for _, record := range records {
		// 添加边界检查防止index out of range
		if len(record) < 7 {
			return nil, fmt.Errorf("记录列数不足，期望至少7列，实际只有%d列", len(record))
		}

		if record[6] == "是" {
			continue
		}

		// 第1列 金额钱
		amount, _ := strconv.ParseFloat(record[0], 64)
		// 第2列 描述
		description := record[1]
		// 第3列 人
		//var userId int64 = 1
		//if record[2] == "Z" {
		//	userId = 2
		//}

		accountMap := map[string]int{
			"微信":   1,
			"支付宝":  2,
			"建设银行": 8,
			"中国银行": 6,
			"工商银行": 7,
			"农业银行": 9,
			"招商银行": 10,
			"医保":   5,
			"美团":   4,
			"给到":   12,
		}

		// 第四列 支出账户
		account := accountMap[record[3]]

		// 第5列 时间 - 支持多种格式
		localTime, err := utils.ParseTimeString(record[4])
		if err != nil {
			return nil, fmt.Errorf("无法解析时间格式 '%s': %v", record[4], err)
		}

		// 第6列 收支
		typeId := 1
		if record[5] == "收入" {
			typeId = 1
		} else {
			typeId = 2
		}

		// 第8列 类型
		categoryMap := map[string]int{
			"饮食":        1,
			"医疗":        3,
			"聚会":        5,
			"购物（含生活用品）": 2,
			"电子产品":      9,
			"房租":        10,
			"其他":        101,
			"交通":        11,
			"通讯":        8,
			"运动健身":      4,
			"教育":        6,
			"娱乐":        5,
			"水电煤气费":     7,
			"收入":        14,
		}
		category := categoryMap[record[7]]

		transaction := moneydb.Transaction{
			Amount:      amount,
			Description: description,
			UserId:      userId,
			Type:        typeId,
			Category:    category,
			Account:     account,
			Time:        model.LocalTime(localTime),
		}
		transactions = append(transactions, transaction)
		continue
	}
	return transactions, nil
}
