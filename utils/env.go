package utils

import (
	"go-redis2influx/global"
	"fmt"
	"log"

	"go-redis2influx/models"
	"os"

	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func LoadEnvironment() {
	loadEnvConfigFile()
	loadEventLogConfig()
}

func loadEnvConfigFile() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/go-redis2influx")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {

			log.Println("沒有發現 config.yml，改抓取環境變數")
			viper.AutomaticEnv()
			viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		} else {
			// 有找到 config.yml 但是發生了其他未知的錯誤
			panic(fmt.Sprintf("Fatal error config file: %v\n", err.Error()))
		}
	}

	// 创建配置结构体实例
	var config models.EnvironmentModel

	// 将配置文件内容解析到结构体中
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	// 檢查路徑是否為目錄
	if err := os.MkdirAll(config.Log.Path, os.ModePerm); err != nil {
		global.Logger.Error(fmt.Sprintf("Create %v Error: %v", config.Log.Path, err.Error()),
			zap.Any(global.LogEvent.LoadEnvConfig.Name, global.LogEvent.LoadEnvConfig))
	}

	// 加上預設檔案名稱
	defaultFileName := "bimap.log"
	config.Log.Path = filepath.Join(config.Log.Path, defaultFileName)

	global.EnvConfig = &config
}

func loadEventLogConfig() {
	// 直接初始化 LogEvent 配置
	global.LogEvent = &models.LogEvent{
		OutputInfluxDB: models.Event{
			Name:        "OutputInfluxDB",
			Code:        "INFLUX01",
			Category:    "InfluxDB",
			Level:       "",
			Threshold:   "",
			Description: "Logs related to InfluxDB output",
		},
		ConnectInfluxDB: models.Event{
			Name:        "ConnectInfluxDB",
			Code:        "INFLUX02",
			Category:    "InfluxDB",
			Level:       "",
			Threshold:   "",
			Description: "Logs related to InfluxDB connection",
		},
		LoggerWrite: models.Event{
			Name:        "LoggerWrite",
			Code:        "LOG01",
			Category:    "Logger",
			Level:       "",
			Threshold:   "",
			Description: "Logs related to logger writing",
		},
		LoadEnvConfig: models.Event{
			Name:        "LoadEnvConfig",
			Code:        "ENV01",
			Category:    "Config",
			Level:       "",
			Threshold:   "",
			Description: "Logs related to loading environment configuration",
		},
		ConnectRedis: models.Event{
			Name:        "ConnectRedis",
			Code:        "REDIS01",
			Category:    "Redis",
			Level:       "",
			Threshold:   "",
			Description: "Logs related to establishing connection to Redis",
		},
		ReadRedisStream: models.Event{
			Name:        "ReadRedisStream",
			Code:        "REDIS02",
			Category:    "Redis",
			Level:       "",
			Threshold:   "",
			Description: "Logs related to reading from Redis Stream",
		},
		AckRedisMessage: models.Event{
			Name:        "AckRedisMessage",
			Code:        "REDIS03",
			Category:    "Redis",
			Level:       "",
			Threshold:   "",
			Description: "Logs related to acknowledging Redis messages",
		},
		RedisCommandError: models.Event{
			Name:        "RedisCommandError",
			Code:        "REDIS04",
			Category:    "Redis",
			Level:       "Error",
			Threshold:   "",
			Description: "Logs related to errors when executing Redis commands",
		},
		RedisConnectionLost: models.Event{
			Name:        "RedisConnectionLost",
			Code:        "REDIS05",
			Category:    "Redis",
			Level:       "Error",
			Threshold:   "",
			Description: "Logs related to losing connection to Redis",
		},
		ReconnectRedis: models.Event{
			Name:        "ReconnectRedis",
			Code:        "REDIS06",
			Category:    "Redis",
			Level:       "",
			Threshold:   "",
			Description: "Logs related to reconnecting to Redis after losing connection",
		},
		RedisWrite: models.Event{
			Name:        "RedisWrite",
			Code:        "REDIS07",
			Category:    "Redis",
			Level:       "",
			Threshold:   "",
			Description: "Logs related to writing data to Redis",
		},
		RedisGroupCreate: models.Event{
			Name:        "RedisGroupCreate",
			Code:        "REDIS08",
			Category:    "Redis",
			Level:       "",
			Threshold:   "",
			Description: "Logs related to creating Redis Consumer Group",
		},
	}
}

// func loadEventLogConfigFile() {
// 	viper.SetConfigName("log")
// 	viper.SetConfigType("yml")
// 	viper.AddConfigPath(".")

// 	if err := viper.ReadInConfig(); err != nil {
// 		log.Fatalf("Error reading config file, %s", err)
// 	}

// 	var logEvent models.LogEvent
// 	err := viper.Sub("log_event").Unmarshal(&logEvent)
// 	if err != nil {
// 		log.Fatalf("unable to decode into struct, %v", err)
// 	}

// 	global.LogEvent = &logEvent

// }
