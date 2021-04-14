package kafka

import (
	"baselib/logger"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ConfigMap struct {
	Broker string
}

var (
	Producer *kafka.Producer
)

// 初始化生产者
func InitProducer(config *ConfigMap) {
	var err error
	Producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": config.Broker})
	if err != nil {
		panic(err)
	}
}

// 发送消息
func Send(topic string, message []byte) bool {
	delivery := make(chan kafka.Event)
	close(delivery)
	err := Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
		Headers:        []kafka.Header{{Key: "wss", Value: []byte("KafkaLogs")}},
	}, delivery)
	if err != nil {
		logger.Sugar.Errorf("Produce failed: %s", err)
	}
	event := <-delivery
	msg := event.(*kafka.Message)
	if msg.TopicPartition.Error != nil {
		logger.Sugar.Errorf("Delivery failed: %v", msg.TopicPartition.Error)
		return false
	} else {
		logger.Sugar.Infof("Delivered message [%s] to topic %s [%d] at offset %v", string(message), *msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset) // TODO 可能会有很大的消耗，所有最后是否使用视场景而定
		return true
	}
}

// 手动关闭生产者
func CloseProducer() {
	Producer.Close()
}
