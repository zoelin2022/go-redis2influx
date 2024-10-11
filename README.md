# go-redis2influx

本專案是一個用於從 Redis Stream 消費數據並將其寫入 InfluxDB 的 Go 應用程序。它包含連線檢查和重試機制，當 InfluxDB 無法連線時，會暫停消費數據，直到 InfluxDB 恢復連線為止。

## 功能

- 從 Redis Stream 中消費數據，數據以 InfluxDB Line Protocol 格式存儲。
- 將數據寫入 InfluxDB v2，可以指定 org 和 bucket。
- 若 InfluxDB 連線失敗則暫停消費，並定期重試，直到連線恢復。

## start

sudo systemctl start go-redis2influx.service

## status

sudo systemctl status go-redis2influx.service

## enable

sudo systemctl enable go-redis2influx.service

## log file

/var/log/go-redis2influx/bimap.log

## journalctl 看 log

journalctl -f -u go-redis2influx.service

```log
2024-10-07 08:22:31	INFO	Successfully written 119 records to InfluxDB	{"OutputInfluxDB": {"Name":"OutputInfluxDB","Code":"INFLUX01","Category":"InfluxDB","Level":"","Threshold":"","Description":"Logs related to InfluxDB output"}}
2024-10-07 08:22:33	INFO	Successfully written 187 records to InfluxDB	{"OutputInfluxDB": {"Name":"OutputInfluxDB","Code":"INFLUX01","Category":"InfluxDB","Level":"","Threshold":"","Description":"Logs related to InfluxDB output"}}
2024-10-07 08:22:34	INFO	Successfully written 175 records to InfluxDB	{"OutputInfluxDB": {"Name":"OutputInfluxDB","Code":"INFLUX01","Category":"InfluxDB","Level":"","Threshold":"","Description":"Logs related to InfluxDB output"}}
```
