package yunjingconvert

import (
	"github.com/gavinlhchen/logconvert/errors"
	"github.com/gavinlhchen/logconvert/internal/yunjingconvert/config"
	"github.com/gavinlhchen/logconvert/log"
)

type yjConvertServer struct {
	alertServer *yjToSocServer
}

type preparedYunjingConvertServer struct {
	*yjConvertServer
}

func CreateYunjingConvertServer(cfg *config.Config) (*yjConvertServer, error) {
	yjAlertConfig := buildAlertConfig(cfg)

	alertServer, err := yjAlertConfig.New()
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	server := &yjConvertServer{
		alertServer: alertServer,
	}

	return server, nil
}

func (s *yjConvertServer) PrepareRun() preparedYunjingConvertServer {
	return preparedYunjingConvertServer{s}
}

func (s preparedYunjingConvertServer) Run() error {

	log.Errorf("test log %s", "sfsf")
	return errors.New("sf")

}

func (s yjConvertServer) Run() {
	go s.alertServer.Run()
}
