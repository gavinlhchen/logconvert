package yunjingconvert

import "github.com/gavinlhchen/logconvert/internal/yunjingconvert/config"

func Run(cfg *config.Config) error {
	server, error := CreateYunjingConvertServer(cfg)

	if error != nil {
		return error
	}

	return server.PrepareRun().Run()
}
