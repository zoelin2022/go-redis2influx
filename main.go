package main

import (
	"go-redis2influx/services"
	"go-redis2influx/utils"
)

func main() {
	// 加載環境參數和初始化日誌系統
	utils.LoadEnvironment()
	utils.InitLogger()

	// 啟動 Redis 消費者處理數據
	go services.ReadRedisData()

	// 阻塞主線程
	select {}
}
