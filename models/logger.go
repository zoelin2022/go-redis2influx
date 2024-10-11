package models

type Event struct {
	Name        string `mapstructure:"name"`
	Code        string `mapstructure:"code"`
	Category    string `mapstructure:"category"`
	Level       string `mapstructure:"level"`
	Threshold   string `mapstructure:"threshold"`
	Description string `mapstructure:"description"`
}

type LogEvent struct {
	// InfluxDB Events
	OutputInfluxDB  Event `mapstructure:"output_influxdb"`
	ConnectInfluxDB Event `mapstructure:"connect_influxdb"`

	// Logger Event
	LoggerWrite Event `mapstructure:"logger_write"`

	// Configuration Event
	LoadEnvConfig Event `mapstructure:"load_env_config"`

	// Redis Events
	ConnectRedis        Event `mapstructure:"connect_redis"`
	ReadRedisStream     Event `mapstructure:"read_redis_stream"`
	AckRedisMessage     Event `mapstructure:"ack_redis_message"`
	RedisCommandError   Event `mapstructure:"redis_command_error"`
	RedisConnectionLost Event `mapstructure:"redis_connection_lost"`
	ReconnectRedis      Event `mapstructure:"reconnect_redis"`
	RedisWrite          Event `mapstructure:"redis_write"`
	RedisGroupCreate    Event `mapstructure:"redis_group_create"`
}
