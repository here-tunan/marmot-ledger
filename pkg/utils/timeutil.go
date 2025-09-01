package utils

import (
	"fmt"
	"go-my-life/pkg/model"
	"time"
)

const formatter = "2006-01-02 15:04:05"
const Loc = "Asia/Shanghai"

func Format(t time.Time) string {
	location, err := time.LoadLocation(Loc)
	if err != nil {
		fmt.Println("时区加载错误:", err)
		return ""
	}
	// 转换为指定时区的时间
	t = t.In(location)

	formattedTime := t.Format(formatter)

	return formattedTime
}

func LocalTimeFormat(lt model.LocalTime) string {
	t := time.Time(lt)
	location, err := time.LoadLocation(Loc)
	if err != nil {
		fmt.Println("时区加载错误:", err)
		return ""
	}
	// 转换为指定时区的时间
	t = t.In(location)

	formattedTime := t.Format(formatter)

	return formattedTime
}

func ParseTimeInLoc(formatter string, timeStr string) time.Time {
	location, err := time.LoadLocation(Loc)
	if err != nil {
		fmt.Println("时区加载错误:", err)
		return time.Time{}
	}
	curTime, err := time.ParseInLocation(formatter, timeStr, location)
	return curTime
}

// ParseTimeString 解析多种时间格式
func ParseTimeString(timeStr string) (time.Time, error) {
	// 常见的时间格式
	timeFormats := []string{
		"2006-01-02",          // 2025-08-31
		"2006/1/2 15:04",      // 2006/1/2 15:04
		"2006/01/02 15:04",    // 2006/01/02 15:04
		"2006-01-02 15:04",    // 2006-01-02 15:04
		"2006/1/2 15:04:05",   // 2006/1/2 15:04:05
		"2006/01/02 15:04:05", // 2006/01/02 15:04:05
		"2006-01-02 15:04:05", // 2006-01-02 15:04:05
		"2006/1/2",            // 2006/1/2
		"2006/01/02",          // 2006/01/02
	}

	location, err := time.LoadLocation(Loc)
	if err != nil {
		fmt.Println("时区加载错误:", err)
		location = time.Local // 使用本地时区作为后备
	}

	for _, format := range timeFormats {
		if parsedTime, err := time.ParseInLocation(format, timeStr, location); err == nil {
			return parsedTime, nil
		}
	}

	return time.Time{}, fmt.Errorf("不支持的时间格式: %s", timeStr)
}
