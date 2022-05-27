package yjtosocserver

import (
	"context"
	"github.com/Shopify/sarama"
	"logconvert/errors"
	"logconvert/internal/pkg/kafka"
	"logconvert/internal/yjtosocserver/config"
	"logconvert/internal/yjtosocserver/yunjing"
	"logconvert/log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type workerConfig struct {
	cConfig  *kafka.ConsumerConfig
	pConfig  *kafka.ProducerConfig
	cHandler yunjing.YjToSocMsgHandler
}

type yjToSocWorker struct {
	cTopic    string
	cg        sarama.ConsumerGroup
	cgHandler *yunjing.ConsumerGroupHandler
}

func buildWorkConfig(cfg *config.Config, topic string) (*workerConfig, error) {
	cc := &kafka.ConsumerConfig{
		Brokers:       cfg.YunjingOptions.Brokers,
		Topic:         topic,
		ConsumerGroup: cfg.YunjingOptions.ConsumerGroup,
		OffsetOldest:  cfg.YunjingOptions.OffsetOldest,
		Version:       cfg.YunjingOptions.Version,
		Assignor:      cfg.YunjingOptions.Assignor,
	}

	pc := &kafka.ProducerConfig{
		Brokers: cfg.IsaGlobal.Component.KafkaNodes,
		Topic:   cfg.GenericServerRunOptions.RawEventTopic,
	}

	ch, err := yunjing.New(topic)
	if err != nil {
		return nil, err
	}

	c := &workerConfig{
		cConfig:  cc,
		pConfig:  pc,
		cHandler: ch,
	}

	return c, nil
}

func (yjTsc *workerConfig) New() (*yjToSocWorker, error) {
	cGroup := kafka.NewConsumer(yjTsc.cConfig)
	producer, err := kafka.NewProducer(yjTsc.pConfig)

	if err != nil {
		return nil, err
	}

	worker := &yjToSocWorker{
		cTopic: yjTsc.cConfig.Topic,
		cg:     cGroup,
		cgHandler: &yunjing.ConsumerGroupHandler{
			Ready:        make(chan bool),
			Producer:     producer,
			ProduceTopic: yjTsc.pConfig.Topic,
			MsgHandler:   yjTsc.cHandler,
		},
	}

	return worker, nil
}

func (yjTss *yjToSocWorker) Run(ctx context.Context, cancel context.CancelFunc) error {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := yjTss.cg.Consume(ctx, []string{yjTss.cTopic}, yjTss.cgHandler); err != nil {
				log.Panicf("Error from cg: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
			yjTss.cgHandler.Ready = make(chan bool)
		}
	}()

	<-yjTss.cgHandler.Ready
	log.Infof("start consume topic %s !...", yjTss.cTopic)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	keepRunning := true
	for keepRunning {
		select {
		case <-ctx.Done():
			log.Infof("terminating consume topic %s: context cancelled", yjTss.cTopic)
			keepRunning = false
		case <-sigterm:
			log.Infof("terminating consume topic %s: via signal", yjTss.cTopic)
			keepRunning = false
		}
	}
	cancel()
	wg.Wait()
	if err := yjTss.cg.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
	return errors.Errorf("stop consume topic %s", yjTss.cTopic)
}
