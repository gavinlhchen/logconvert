package config

import (
	"github.com/gavinlhchen/logconvert/internal/pkg/soc"
	"github.com/gavinlhchen/logconvert/internal/yjtosocserver/options"
)

type Config struct {
	*options.Options
	*soc.IsaGlobal
}

func CreateConfigFromOptions(opts *options.Options) (*Config, error) {

	return &Config{opts, soc.NewConfig(opts.GenericServerRunOptions.IsaGlobalConfigPath)}, nil
}
