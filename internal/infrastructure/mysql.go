package infrastructure

import (
	"fmt"
	"marmot-ledger/env"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

var Mysql *xorm.Engine

func init() {
	dataSourName := fmt.Sprintf("%s:%s@tcp(%s:%d)/marmot_ledger?charset=utf8mb4&parseTime=true",
		env.Prop.Mysql.Username, env.Prop.Mysql.Password, env.Prop.Mysql.Host, env.Prop.Mysql.Port)

	newEngine, _ := xorm.NewEngine("mysql", dataSourName)

	// 设置驼峰转换
	newEngine.SetMapper(names.GonicMapper{})

	Mysql = newEngine
}
