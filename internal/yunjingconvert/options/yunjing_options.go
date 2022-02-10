package options

import (
	"github.com/spf13/pflag"
)

type YunjingKafkaOptions struct {
	AlertTopic    string   `json:"alert-topic"  mapstructure:"alert-topic"`
	HostTopics    []string `json:"host-topics"  mapstructure:"host-topics"`
	ConsumerGroup string   `json:"consumer-group"  mapstructure:"consumer-group"`
	OffsetOldest  bool     `json:"offset-oldest"  mapstructure:"offset-oldest"`
	Version       string   `json:"version"  mapstructure:"version"`
	Assignor      string   `json:"assignor"  mapstructure:"assignor"`
	Brokers       []string `json:"brokers"  mapstructure:"brokers"`
}

func NewYunjingOptions() *YunjingKafkaOptions {
	return &YunjingKafkaOptions{
		AlertTopic:    "event_msg",
		HostTopics:    []string{"fast_msg", "host_login", "bash_scan", "asset_accout", "asset_port", "asset_process"},
		ConsumerGroup: "cg_soc_yunjing",
	}
}

// Validate verifies flags passed to MySQLOptions.
func (o *YunjingKafkaOptions) Validate() []error {
	errs := []error{}

	return errs
}

// AddFlags adds flags related to mysql storage for a specific APIServer to the specified FlagSet.
func (o *YunjingKafkaOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.AlertTopic, "yunjing-kafka.alert-topic", o.AlertTopic, ""+
		"Yujing alert log topic.")

	fs.StringSliceVar(&o.HostTopics, "yunjing-kafka.host-topics", o.HostTopics, ""+
		"Yujing host log topics.")

	fs.StringVar(&o.ConsumerGroup, "yunjing-kafka.consumer-groupid", o.ConsumerGroup, ""+
		"The consumer groupid for consumer log from yunjing kafka.")
}
