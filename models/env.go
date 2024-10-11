package models

type EnvironmentModel struct {
	Redis struct {
		Address      string `mapstructure:"address"`
		DB           int    `mapstructure:"db"`
		GroupName    string `mapstructure:"group_name"`
		ConsumerName string `mapstructure:"consumer_name"`
		StreamKey    string `mapstructure:"stream_key"`
		MessageField string `mapstructure:"message_field"`
		Count        int    `mapstructure:"count"`
		BlockMs      int    `mapstructure:"block_ms"`
		RetryDelay   int    `mapstructure:"retry_delay"`
	}

	Log struct {
		Level   string `mapstructure:"level"`
		Path    string `mapstructure:"path"`
		MaxSize int    `mapstructure:"maxsize"`
		MaxAge  int    `mapstructure:"maxage"`
	}

	Influxdb struct {
		URL     string `mapstructure:"url"`
		Token   string `mapstructure:"token"`
		Org     string `mapstructure:"org"`
		Bucket  string `mapstructure:"bucket"`
		Options struct {
			SetBatchSize          int    `mapstructure:"set_batch_size"`
			SetLogLevel           int    `mapstructure:"set_log_level"`
			SetUseGzip            bool   `mapstructure:"set_use_gzip"`
			SetFlushInterval      int    `mapstructure:"set_flush_interval"`
			SetMaxRetries         int    `mapstructure:"set_max_retries"`
			SetMaxRetryTime       int    `mapstructure:"set_max_retry_time"`
			SetRetryInterval      int    `mapstructure:"set_retry_interval"`
			SetHTTPRequestTimeout int    `mapstructure:"set_http_request_timeout"`
			SetApplicationName    string `mapstructure:"set_application_name"`
		} `mapstructure:"options"`
	}
}
