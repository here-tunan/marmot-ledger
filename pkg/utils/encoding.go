package utils

import (
	"bytes"
	"fmt"
	"io"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// DetectAndConvertEncoding 检测并转换文件编码
func DetectAndConvertEncoding(data []byte) ([]byte, string, error) {
	// 尝试UTF-8（包含BOM）
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		return data[3:], "UTF-8 with BOM", nil
	}

	// 尝试直接当作UTF-8解析
	if IsValidUTF8(data) {
		return data, "UTF-8", nil
	}

	// 尝试GBK转UTF-8
	if convertedData, err := ConvertGBKToUTF8(data); err == nil {
		if IsValidUTF8(convertedData) {
			return convertedData, "GBK", nil
		}
	}

	// 返回原始数据，让程序继续尝试处理
	return data, "Unknown", fmt.Errorf("无法识别文件编码，可能导致解析错误")
}

// IsValidUTF8 检查数据是否为有效的UTF-8
func IsValidUTF8(data []byte) bool {
	return utf8.Valid(data)
}

// ConvertGBKToUTF8 将GBK编码转换为UTF-8
func ConvertGBKToUTF8(data []byte) ([]byte, error) {
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
