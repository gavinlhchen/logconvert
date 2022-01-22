package yunjingconvert

import (
	"github.com/gavinlhchen/logconvert/errors"
	"github.com/gavinlhchen/logconvert/internal/yunjingconvert/config"
)

type yunjingConvertServer struct{}

type preparedYunjingConvertServer struct {
	*yunjingConvertServer
}

func CreateYunjingConvertServer(cfg *config.Config) (*yunjingConvertServer, error) {
	server := &yunjingConvertServer{}
	return server, nil
}

func (s *yunjingConvertServer) PrepareRun() preparedYunjingConvertServer {
	return preparedYunjingConvertServer{s}
}

func (s preparedYunjingConvertServer) Run() error {

	return errors.New("sf")

}
