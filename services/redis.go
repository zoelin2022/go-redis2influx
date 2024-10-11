package services

import (
	"go-redis2influx/databases"
	"go-redis2influx/global"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func ReadRedisData() {

	config := global.EnvConfig.Redis

	// 初始化 Redis 客戶端
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: config.Address,
		DB:   config.DB,
	})

	// 創建消費者群組（如果不存在）
	err := rdb.XGroupCreateMkStream(ctx, config.StreamKey, config.GroupName, "0").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP Consumer Group name already exists") {
		global.Logger.Error(fmt.Sprintf("Failed to create consumer group: %v", err),
			zap.Any(global.LogEvent.RedisGroupCreate.Name, global.LogEvent.RedisGroupCreate))
		return
	}

	// 設置阻塞時間和讀取數量
	blockDuration := time.Duration(config.BlockMs) * time.Millisecond

	for {
		// 檢查 InfluxDB 連線狀態
		for !databases.InfluxdbConnectionAvailable() {
			global.Logger.Warn("InfluxDB is unavailable, retrying...",
				zap.Any(global.LogEvent.ConnectInfluxDB.Name, global.LogEvent.ConnectInfluxDB))
			// 從環境參數中獲取重試延遲
			time.Sleep(time.Duration(global.EnvConfig.Redis.RetryDelay) * time.Second)
		}

		// 讀取 Stream 中的消息（使用配置中的 Count 和 Block 參數）
		streams, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    config.GroupName,
			Consumer: config.ConsumerName,
			Streams:  []string{config.StreamKey, ">"},
			Count:    int64(config.Count), // 每次讀取數據的數量，取決於配置
			Block:    blockDuration,
		}).Result()

		if err != nil && err != redis.Nil {
			global.Logger.Error(fmt.Sprintf("Error reading from Redis stream: %v", err),
				zap.Any(global.LogEvent.ReadRedisStream.Name, global.LogEvent.ReadRedisStream))
			// 重試前等待設置的重試延遲時間
			time.Sleep(time.Duration(global.EnvConfig.Redis.RetryDelay) * time.Second)
			continue
		}

		var batchData []string  // 存放這次讀取的所有數據
		var messageIDs []string // 存放所有成功處理的消息ID

		for _, stream := range streams {
			for _, message := range stream.Messages {
				data, ok := message.Values[config.MessageField].(string)
				if !ok {
					global.Logger.Error(fmt.Sprintf("Failed to parse data from message: %v", message),
						zap.Any(global.LogEvent.ReadRedisStream.Name, global.LogEvent.ReadRedisStream))
					continue
				}

				// 將數據加入到批量數據集中
				batchData = append(batchData, data)
				messageIDs = append(messageIDs, message.ID)
			}

			// 將數據寫入 InfluxDB 成功後刪除
			if len(batchData) > 0 {
				err = databases.WriteLineProtocol(batchData)
				if err == nil {
					// 成功寫入後，刪除這些已處理的消息
					for _, msgID := range messageIDs {
						rdb.XDel(ctx, config.StreamKey, msgID)
					}
				}
			}
		}

		// 將這次批量讀取的所有數據一次性寫入 InfluxDB
		if len(batchData) > 0 {
			err = databases.WriteLineProtocol(batchData) // 將所有讀取到的數據一次性寫入
			if err != nil {
				global.Logger.Error(fmt.Sprintf("Failed to write batch data to InfluxDB: %v", err),
					zap.Any(global.LogEvent.OutputInfluxDB.Name, global.LogEvent.OutputInfluxDB))
				// 如果寫入失敗，保留數據，下次重試
				continue
			}

			// 批次確認成功寫入的所有消息
			err = rdb.XAck(ctx, config.StreamKey, config.GroupName, messageIDs...).Err()
			if err != nil {
				global.Logger.Error(fmt.Sprintf("Failed to batch acknowledge messages: %v", err),
					zap.Any(global.LogEvent.AckRedisMessage.Name, global.LogEvent.AckRedisMessage))
			} else {
				global.Logger.Info(fmt.Sprintf("Successfully acknowledged %d records from Redis Stream", len(messageIDs)),
					zap.Any(global.LogEvent.AckRedisMessage.Name, global.LogEvent.AckRedisMessage))
			}

			global.Logger.Info(fmt.Sprintf("Successfully written %d records to InfluxDB", len(batchData)),
				zap.Any(global.LogEvent.OutputInfluxDB.Name, global.LogEvent.OutputInfluxDB))
		}

		// 在迴圈中等待 blockDuration 再進行下一次迴圈
		time.Sleep(blockDuration)
	}
}

func ProcessRemainingDataFromRedis() {
	config := global.EnvConfig.Redis
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: config.Address,
		DB:   config.DB,
	})

	// 讀取 Redis 中所有未處理的數據
	streams, err := rdb.XRange(ctx, config.StreamKey, "-", "+").Result()
	if err != nil {
		global.Logger.Error(fmt.Sprintf("Failed to read from Redis stream: %v", err),
			zap.Any(global.LogEvent.ReadRedisStream.Name, global.LogEvent.ReadRedisStream))
		return
	}

	var messageIDs []string
	var batchData []string

	for _, message := range streams {
		data, ok := message.Values[config.MessageField].(string)
		if !ok {
			global.Logger.Error(fmt.Sprintf("Failed to parse data from message: %v", message),
				zap.Any(global.LogEvent.ReadRedisStream.Name, global.LogEvent.ReadRedisStream))
			continue
		}

		// 收集數據重新寫入 InfluxDB
		batchData = append(batchData, data)
		messageIDs = append(messageIDs, message.ID)
	}

	// 批量重新寫入 InfluxDB
	if len(batchData) > 0 {
		err = databases.WriteLineProtocol(batchData)
		if err == nil {
			// 記錄成功寫入的筆數
			global.Logger.Info(fmt.Sprintf("Successfully re-written %d records to InfluxDB", len(batchData)),
				zap.Any(global.LogEvent.OutputInfluxDB.Name, global.LogEvent.OutputInfluxDB))

			// 批量刪除 Redis 中的這些消息
			if len(messageIDs) > 0 {
				err = rdb.XDel(ctx, config.StreamKey, messageIDs...).Err()
				if err != nil {
					global.Logger.Error(fmt.Sprintf("Failed to batch delete messages: %v", err),
						zap.Any(global.LogEvent.AckRedisMessage.Name, global.LogEvent.AckRedisMessage))
				} else {
					global.Logger.Info(fmt.Sprintf("Successfully deleted %d records from Redis Stream", len(messageIDs)),
						zap.Any(global.LogEvent.AckRedisMessage.Name, global.LogEvent.AckRedisMessage))
				}
			}
		} else {
			global.Logger.Error(fmt.Sprintf("Failed to re-write data to InfluxDB: %v", err),
				zap.Any(global.LogEvent.OutputInfluxDB.Name, global.LogEvent.OutputInfluxDB))
		}
	} else {
		global.Logger.Info("No residual data found in Redis to process.",
			zap.Any(global.LogEvent.ReadRedisStream.Name, global.LogEvent.ReadRedisStream))
	}
}
