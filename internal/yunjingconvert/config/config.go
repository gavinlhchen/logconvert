package config

import "github.com/gavinlhchen/logconvert/internal/yunjingconvert/options"

type Config struct {
	*options.Options
}

func CreateConfigFromOptions(opts *options.Options) (*Config, error) {
	return &Config{opts}, nil
}
