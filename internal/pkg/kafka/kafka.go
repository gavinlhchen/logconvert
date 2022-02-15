package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/gavinlhchen/logconvert/log"
	"time"
)

type ConsumerConfig struct {
	Brokers       []string
	Topic         string
	ConsumerGroup string
	OffsetOldest  bool
	Version       string
	Assignor      string
}

type ProducerConfig struct {
	Brokers []string
	Topic   string
}

func NewConsumer(cc *ConsumerConfig) sarama.ConsumerGroup {
	if len(cc.Brokers) == 0 {
		log.Panic("no Kafka bootstrap Brokers defined")
	}

	if len(cc.Topic) == 0 {
		log.Panic("no Topic given to be consumed")
	}

	if len(cc.ConsumerGroup) == 0 {
		log.Panic("no Kafka consumer group defined")
	}

	version, err := sarama.ParseKafkaVersion(cc.Version)
	if err != nil {
		log.Panicf("Error parsing Kafka Version: %v", err)
	}

	config := sarama.NewConfig()
	config.Version = version

	switch cc.Assignor {
	case "roundrobin":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case "range":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	default:
		log.Panicf("Unrecognized consumer group partition assignor: %s", cc.Assignor)
	}

	if cc.OffsetOldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	cg, err := sarama.NewConsumerGroup(cc.Brokers, cc.ConsumerGroup, config)
	if err != nil {
		log.Panicf("Error creating consumer group cg: %v", err)
	}
	return cg
}

func NewProducer(pc *ProducerConfig) (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	producer, err := sarama.NewAsyncProducer(pc.Brokers, config)

	if err != nil {
		return nil, err
	}

	go func() {
		for err := range producer.Errors() {
			log.Errorf("Failed to send msg:%s", err)
		}
	}()

	return producer, nil
}
