package databases

import (
	"go-redis2influx/global"
	"go-redis2influx/models"
	"context"
	"fmt"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"go.uber.org/zap"
)

func InfluxdbConnectionAvailable() bool {
	client := NewInfluxDBClient(time.Second)
	ctx := context.Background()

	_, err := client.Health(ctx)
	if err != nil {
		global.Logger.Error(err.Error(),
			zap.Any(global.LogEvent.ConnectInfluxDB.Name, global.LogEvent.ConnectInfluxDB))
		global.InfluxDB = client
		return false
	}
	return true
}

func LoadInfluxDB() {
	client := NewInfluxDBClient(time.Second)
	ctx := context.Background()

	d, err := client.Health(ctx)
	if err != nil {
		global.Logger.Error(err.Error(),
			zap.Any(global.LogEvent.ConnectInfluxDB.Name, global.LogEvent.ConnectInfluxDB))
		global.InfluxDB = client
		return
	}
	global.InfluxDB = client
	global.Logger.Info(fmt.Sprintf("influxdb connection success: %v", *d.Message),
		zap.Any(global.LogEvent.ConnectInfluxDB.Name, global.LogEvent.ConnectInfluxDB))

}

func NewInfluxDBClient(precision time.Duration) influxdb2.Client {
	url := global.EnvConfig.Influxdb.URL
	token := global.EnvConfig.Influxdb.Token
	env := global.EnvConfig.Influxdb.Options

	influxdb := influxdb2.NewClientWithOptions(url, token,
		influxdb2.DefaultOptions().
			//*** 指定最粗時間戳的精度
			SetPrecision(precision).
			//*** 最佳批量大小是 5000 行
			SetBatchSize(uint(env.SetBatchSize)).
			//*** 0=error, 1=warning, 2=info, 3=debug, nil 禁用
			SetLogLevel(uint(env.SetLogLevel)).
			//*** 壓縮傳輸可提高5倍速
			SetUseGZip(env.SetUseGzip).
			//*** 如果緩衝區3s尚未寫入，則刷新緩衝區
			SetFlushInterval(uint(env.SetFlushInterval)).
			//*** 失敗寫入的最大重試次數
			SetMaxRetries(uint(env.SetMaxRetries)).
			//*** 失敗寫入的最大重試時間
			SetMaxRetryTime(uint(env.SetMaxRetryTime)).
			//*** 重試之間的最大延遲
			SetRetryInterval(uint(env.SetFlushInterval)).
			//*** 長時間寫資料設定
			SetHTTPRequestTimeout(uint(env.SetHTTPRequestTimeout)).
			SetApplicationName(env.SetApplicationName))

	return influxdb
}

// * 寫入 InfluxDB
func WriteLineProtocol(data []string) error {

	db := global.EnvConfig.Influxdb
	client := NewInfluxDBClient(time.Second)
	writeAPI := client.WriteAPIBlocking(db.Org, db.Bucket)

	if err := writeAPI.WriteRecord(context.Background(), strings.Join(data, "\n")); err != nil {
		global.Logger.Error(fmt.Sprintf("WriteToInfluxDB Error: %v", err),
			zap.Any(global.LogEvent.OutputInfluxDB.Name, global.LogEvent.OutputInfluxDB))
		return err
	}

	return nil
}

// * 寫入 InfluxDB
func WriteToInfluxDB(data []models.Point) error {
	points := []string{}

	db := global.EnvConfig.Influxdb
	client := NewInfluxDBClient(time.Second)
	writeAPI := client.WriteAPIBlocking(db.Org, db.Bucket)

	for _, line := range data {
		points = append(points, toLineProtocol(line))
		if global.EnvConfig.Log.Level == "error" {
			fmt.Printf("%v\n", toLineProtocol(line))
		}
	}

	if err := writeAPI.WriteRecord(context.Background(), points...); err != nil {
		global.Logger.Error(fmt.Sprintf("WriteToInfluxDB Error: %v", err),
			zap.Any(global.LogEvent.OutputInfluxDB.Name, global.LogEvent.OutputInfluxDB))
		return err
	}

	return nil
}

func toLineProtocol(p models.Point) string {
	var builder strings.Builder

	builder.WriteString(p.Name)

	var tagStrings []string
	for key, value := range p.Tags {
		tagStrings = append(tagStrings, fmt.Sprintf("%s=%s", key, value))
	}
	builder.WriteString("," + strings.Join(tagStrings, ","))

	var fieldStrings []string
	for key, value := range p.Fields {
		fieldStrings = append(fieldStrings, fmt.Sprintf("%s=%v", key, value))
	}
	builder.WriteString(" " + strings.Join(fieldStrings, ","))

	builder.WriteString(" " + fmt.Sprintf("%d", p.Time))

	return builder.String()
}
