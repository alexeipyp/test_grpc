package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	EnvDev  = "dev"
	EnvProd = "prod"
)

type Config struct {
	Environment            string                       `json:"env"`
	Log                    LogConfig                    `json:"log"`
	ExternalHTTPConnection ExternalHTTPConnectionConfig `json:"external_http_connection"`
	Cache                  CacheConfig                  `json:"cache"`
	GRPC                   GRPCConfig                   `json:"grpc"`
	Scheduler              SchedulerConfig              `json:"scheduler"`
}

type SchedulerConfig struct {
	QueueSize      int `json:"queue_size"`
	WorkerPoolSize int `json:"worker_pool_size"`
}

type CacheConfig struct {
	Lifetime ParsedDuration `json:"lifetime"`
}

type ExternalHTTPConnectionConfig struct {
	Host     string         `json:"host"`
	Username string         `json:"username"`
	Password string         `json:"password"`
	Timeout  ParsedDuration `json:"timeout"`
	Port     int            `json:"port"`
}

type GRPCConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type LogConfig struct {
	LogsDir              string `env:"LOGS_PATH" env-required:"true"`
	LogFilename          string `json:"log_filename" env-default:"log.log"`
	DebugLogFilename     string `json:"debug_log_filename" env-default:"debuglog.log"`
	GRPCTraceLogFilename string `json:"grpc_trace_log_filename" env-default:"grpctracelog.log"`
	HTTPTraceLogFilename string `json:"http_trace_log_filename" env-default:"httptracelog.log"`
}

type ParsedDuration string

func (p ParsedDuration) TryDuration() (time.Duration, error) {
	return time.ParseDuration(string(p))
}

func (p ParsedDuration) Duration() time.Duration {
	dur, _ := time.ParseDuration(string(p))
	return dur
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist " + path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func (cfg *Config) Validate() error {
	errList := []error{
		cfg.Scheduler.Validate(),
		cfg.Cache.Validate(),
		cfg.GRPC.Validate(),
		cfg.ExternalHTTPConnection.Validate(),
		cfg.Log.Validate(),
		ValidateEnv(cfg.Environment),
	}
	if isErrSliceContainsNonNilValue(errList) {
		return constructAnError(
			"validation failed",
			"\n",
			errList...,
		)
	}
	return nil
}

func fetchConfigPath() string {
	res := os.Getenv("CONFIG_PATH")
	return res
}
