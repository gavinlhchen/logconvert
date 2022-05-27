package config

import (
	"logconvert/internal/pkg/soc"
	"logconvert/internal/yjtosocserver/options"
)

type Config struct {
	*options.Options
	*soc.IsaGlobal
}

func CreateConfigFromOptions(opts *options.Options) (*Config, error) {
	return &Config{opts, soc.NewConfig(opts.GenericServerRunOptions.IsaGlobalConfigPath)}, nil
}
