package options

import (
	"github.com/spf13/pflag"
)

type YunjingOptions struct {
	AlertTopic      string   `json:"alert-topic"  mapstructure:"alert-topic"`
	HostTopics      []string `json:"host-topics"  mapstructure:"host-topics"`
	ConsumerGroupid string   `json:"consumer-groupid"  mapstructure:"consumer-groupid"`
}

func NewYunjingOptions() *YunjingOptions {
	return &YunjingOptions{
		AlertTopic:      "event_msg",
		HostTopics:      []string{"fast_msg", "host_login", "bash_scan", "asset_accout", "asset_port", "asset_process"},
		ConsumerGroupid: "cg_soc_yunjing",
	}
}

// Validate verifies flags passed to MySQLOptions.
func (o *YunjingOptions) Validate() []error {
	errs := []error{}

	return errs
}

// AddFlags adds flags related to mysql storage for a specific APIServer to the specified FlagSet.
func (o *YunjingOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.AlertTopic, "yunjing.alert-topic", o.AlertTopic, ""+
		"Yujing alert log topic.")

	fs.StringSliceVar(&o.HostTopics, "yunjing.host-topics", o.HostTopics, ""+
		"Yujing host log topics.")

	fs.StringVar(&o.ConsumerGroupid, "yunjing.consumer-groupid", o.ConsumerGroupid, ""+
		"The consumer groupid for consumer log from yunjing kafka.")
}
