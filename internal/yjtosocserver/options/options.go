package options

import (
	cliflag "github.com/gavinlhchen/logconvert/cli/flag"
	genericoptions "github.com/gavinlhchen/logconvert/internal/pkg/options"
	yujingoptions "github.com/gavinlhchen/logconvert/internal/yjtosocserver/yunjing"
	"github.com/gavinlhchen/logconvert/json"
	"github.com/gavinlhchen/logconvert/log"
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
