package main

import (
	"go-my-life/api"
	_ "go-my-life/internal/infrastructure"
	"log"
)

func main() {
	log.Println("Marmot Ledger is starting!")
	// 启用web服务
	api.Start()
}
