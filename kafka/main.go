package main

import (
	"fmt"
	"github.com/Shopify/sarama"
)

var brokers = []string{"127.0.0.1:9092"}

func newProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)

	return producer, err
}

func prepareMessage(topic, message string) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
		Value:     sarama.StringEncoder(message),
	}

	return msg
}

func main() {

	fmt.Println("kafka")
	producer, err := newProducer()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = producer.SendMessages([]*sarama.ProducerMessage{})
	if err != nil {
		fmt.Println(err)
		return
	}

}
