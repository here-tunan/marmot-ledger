package main

import (
	"log"
	"marmot-ledger/api"
	_ "marmot-ledger/internal/infrastructure"
)

func main() {
	log.Println("Marmot Ledger is starting!")
	// 启用web服务
	api.Start()
}
