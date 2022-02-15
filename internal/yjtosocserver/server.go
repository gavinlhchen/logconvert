package yjtosocserver

import (
	"context"
	"github.com/gavinlhchen/logconvert/internal/yjtosocserver/config"
	"github.com/gavinlhchen/logconvert/log"
	"golang.org/x/sync/errgroup"
)

type yjToSocServer struct {
	ws []*yjToSocWorker
}

func CreateYjToSocServer(cfg *config.Config) (*yjToSocServer, error) {
	var workers []*yjToSocWorker
	for _, topic := range cfg.YunjingOptions.Topics {
		workerConfig, err := buildWorkConfig(cfg, topic)
		if err != nil {
			log.Panicf("Error creating cg group client: %s", err)
		}
		worker, err := workerConfig.New()
		if err != nil {
			log.Panicf("Error creating cg group client: %s", err)
		}
		workers = append(workers, worker)
	}

	return &yjToSocServer{ws: workers}, nil
}

func (s *yjToSocServer) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	eg, errCtx := errgroup.WithContext(ctx)
	for _, worker := range s.ws {
		w := worker
		eg.Go(func() error {
			return w.Run(errCtx, cancel)
		})
	}
	if err := eg.Wait(); err != nil {
		log.Fatal(err.Error())
	}
	return nil
}
