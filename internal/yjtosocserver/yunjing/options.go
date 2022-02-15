package yunjing

import (
	"github.com/spf13/pflag"
)

type KafkaOptions struct {
	Topics        []string `json:"topics"  mapstructure:"topics"`
	ConsumerGroup string   `json:"consumer-group"  mapstructure:"consumer-group"`
	OffsetOldest  bool     `json:"offset-oldest"  mapstructure:"offset-oldest"`
	Version       string   `json:"version"  mapstructure:"version"`
	Assignor      string   `json:"assignor"  mapstructure:"assignor"`
	Brokers       []string `json:"brokers"  mapstructure:"brokers"`
}

func NewYunjingOptions() *KafkaOptions {
	return &KafkaOptions{
		Topics:        []string{},
		ConsumerGroup: "cg_soc_yunjing",
	}
}

// Validate verifies flags passed to MySQLOptions.
func (o *KafkaOptions) Validate() []error {
	errs := []error{}

	return errs
}

// AddFlags adds flags related to mysql storage for a specific APIServer to the specified FlagSet.
func (o *KafkaOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringSliceVar(&o.Topics, "yunjing-kafka.topics", o.Topics, ""+
		"Yujing topics to consume.")

	fs.StringVar(&o.ConsumerGroup, "yunjing-kafka.consumer-group", o.ConsumerGroup, ""+
		"The consumer group for consume log from yunjing kafka.")

	fs.BoolVar(&o.OffsetOldest, "yunjing-kafka.offset-oldest", o.OffsetOldest, ""+
		"Is consume from oldest offset.")

	fs.StringVar(&o.Version, "yunjing-kafka.version", o.Version, ""+
		"The kafka client version.")

	fs.StringVar(&o.Assignor, "yunjing-kafka.assignor", o.Assignor, ""+
		"The kafka client assignor[range, roundrobin].")

	fs.StringSliceVar(&o.Brokers, "yunjing-kafka.brokers", o.Brokers, ""+
		"The kafka brokers addr.")
}
