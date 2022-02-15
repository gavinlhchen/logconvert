package yjtosocserver

import (
	"github.com/gavinlhchen/logconvert/internal/yjtosocserver/config"
)

func Run(cfg *config.Config) error {
	server, err := CreateYjToSocServer(cfg)

	if err != nil {
		return err
	}
	return server.Run()
}
