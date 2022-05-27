package options

import (
	cliflag "logconvert/cli/flag"
	genericoptions "logconvert/internal/pkg/options"
	yujingoptions "logconvert/internal/yjtosocserver/yunjing"
	"logconvert/json"
	"logconvert/log"
)

type Options struct {
	GenericServerRunOptions *genericoptions.ServerRunOptions `json:"server"   mapstructure:"server"`
	YunjingOptions          *yujingoptions.KafkaOptions      `json:"yunjing-kafka"   mapstructure:"yunjing-kafka"`
	Log                     *log.Options                     `json:"log" mapstructure:"log"`
}

func NewOptions() *Options {
	o := Options{
		GenericServerRunOptions: genericoptions.NewServerRunOptions(),
		YunjingOptions:          yujingoptions.NewYunjingOptions(),
		Log:                     log.NewOptions(),
	}

	return &o
}

// Flags returns flags for a specific APIServer by section name.
func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	o.GenericServerRunOptions.AddFlags(fss.FlagSet("generic"))
	o.YunjingOptions.AddFlags(fss.FlagSet("yunjing"))
	o.Log.AddFlags(fss.FlagSet("logs"))

	return fss
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}
