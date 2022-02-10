package yunjingconvert

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/gavinlhchen/logconvert/internal/yunjingconvert/config"
	"github.com/gavinlhchen/logconvert/log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type consumerConfig struct {
	brokers       []string
	topics        []string
	consumerGroup string
	offsetOldest  bool
	version       string
	assignor      string
}

type producerConfig struct {
	brokers []string
	topics  []string
}

type yjToSocConfig struct {
	consumer consumerConfig
	producer producerConfig
}

func buildAlertConfig(cfg *config.Config) *yjToSocConfig {
	alertKafkaConsumer := consumerConfig{
		brokers:       cfg.YunjingOptions.Brokers,
		topics:        []string{cfg.YunjingOptions.AlertTopic},
		consumerGroup: cfg.YunjingOptions.ConsumerGroup,
		offsetOldest:  cfg.YunjingOptions.OffsetOldest,
		version:       cfg.YunjingOptions.Version,
		assignor:      cfg.YunjingOptions.Assignor,
	}

	alertProducver := producerConfig{
		brokers: cfg.YunjingOptions.Brokers,
		topics:  []string{cfg.YunjingOptions.AlertTopic},
	}

	yjTsc := yjToSocConfig{
		consumer: alertKafkaConsumer,
		producer: alertProducver,
	}

	return &yjTsc
}

type yjToSocServer struct {
	consumerTopic []string
	producerTopic []string
	consumer      sarama.ConsumerGroup
	producer      sarama.AsyncProducer
}

func (yjTsc *yjToSocConfig) New() (*yjToSocServer, error) {
	client := newConsumer(yjTsc)

	producer, err := newProducer(yjTsc.producer.brokers)
	if err != nil {
		return nil, err
	}
	yjServer := yjToSocServer{
		consumerTopic: yjTsc.consumer.topics,
		producerTopic: yjTsc.producer.topics,
		consumer:      client,
		producer:      producer,
	}

	return &yjServer, nil
}

func newConsumer(yjTsc *yjToSocConfig) sarama.ConsumerGroup {
	if len(yjTsc.consumer.brokers) == 0 {
		log.Panic("no Kafka bootstrap brokers defined")
	}

	if len(yjTsc.consumer.topics) == 0 {
		log.Panic("nno topics given to be consumed")
	}

	if len(yjTsc.consumer.consumerGroup) == 0 {
		log.Panic("no Kafka consumer group defined")
	}

	version, err := sarama.ParseKafkaVersion(yjTsc.consumer.version)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	config := sarama.NewConfig()
	config.Version = version

	switch yjTsc.consumer.assignor {
	case "roundrobin":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case "range":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	default:
		log.Panicf("Unrecognized consumer group partition assignor: %s", yjTsc.consumer.assignor)
	}

	//if yjTsc.consumer.offsetOldest {
	//	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	//}

	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	client, err := sarama.NewConsumerGroup(yjTsc.consumer.brokers, yjTsc.consumer.consumerGroup, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}
	return client
}

func newProducer(brokerList []string) (sarama.AsyncProducer, error) {
	// For the access log, we are looking for AP semantics, with high throughput.
	// By creating batches of compressed messages, we reduce network I/O at a cost of more latency.
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	return nil, err

	// We will just log to STDOUT if we're not able to produce messages.
	// Note: messages will only be returned here after all retry attempts are exhausted.
	go func() {
		for err := range producer.Errors() {
			log.Infof("Failed to write access log entry:%s", err)
		}
	}()

	return producer, nil
}

func (yjTss yjToSocServer) Run() {
	ctx, cancel := context.WithCancel(context.Background())

	consumer := Consumer{
		ready: make(chan bool),
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := yjTss.consumer.Consume(ctx, yjTss.consumerTopic, &consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	log.Info("Sarama consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	keepRunning := true
	for keepRunning {
		select {
		case <-ctx.Done():
			log.Info("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Info("terminating: via signal")
			keepRunning = false
		}
	}
	cancel()
	wg.Wait()
	if err := yjTss.consumer.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}

}

type Consumer struct {
	ready chan bool
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		log.Infof("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		session.MarkMessage(message, "")
	}

	return nil
}
