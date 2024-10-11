package global

import (
	"go-redis2influx/models"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var (
	EnvConfig *models.EnvironmentModel
	InfluxDB  influxdb2.Client
	Crontab   *cron.Cron
	Logger    *zap.Logger
	LogEvent  *models.LogEvent
)
