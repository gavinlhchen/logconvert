package soc

import (
	"github.com/gavinlhchen/logconvert/log"
	"github.com/spf13/viper"
)

type Component struct {
	KafkaNodes []string `mapstructure:"kafka_nodes"`
}

type IsaGlobal struct {
	Component Component `mapstructure:"component"`
}

func NewConfig(path string) *IsaGlobal {
	v := viper.New()
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		log.Panicf("Read isaglobal config file err:%v", err)
	}

	var c IsaGlobal

	if err := v.Unmarshal(&c); err != nil {
		log.Panicf("Unmarshal isaglobal config file err:%v", err)
	}
	return &c
}

// RawEvent
type RawEvent struct {
	LogsourceIp        string `json:"logsource_ip"`
	LogsourceHost      string `json:"logsource_host"`
	LogsourceName      string `json:"logsource_name"`
	RawLogCharset      string `json:"raw_log_charset"`
	RawLog             string `json:"raw_log"`
	LogsourceTimestamp int64  `json:"logsource_timestamp"`
	EventUuid          string `json:"event_uuid"`
	LogsourceCategory  string `json:"logsource_category"`
}
