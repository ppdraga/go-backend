package main

import (
	"log"
	"os"
	"os/signal"
	"time"

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

var (
	kafkaBrokers = []string{"localhost:9092"}
	KafkaTopic   = "sarama_topic"
	enqueued     int
)

func main() {

	//producer, err := setupProducer()
	producer, err := setupSyncProducer()
	if err != nil {
		panic(err)
	} else {
		log.Println("Kafka AsyncProducer up and running!")
	}

	// Trap SIGINT to trigger a graceful shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	produceMessagesSync(producer, signals)

	log.Printf("Kafka AsyncProducer finished with %d messages produced.", enqueued)
}

// setupProducer will create a AsyncProducer and returns it
func setupProducer() (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	return sarama.NewAsyncProducer(kafkaBrokers, config)
}

// setupSyncProducer will create a AsyncProducer and returns it
func setupSyncProducer() (sarama.SyncProducer, error) {
	//config := sarama.NewConfig()
	return sarama.NewSyncProducer(kafkaBrokers, nil)
}

// produceMessages will send 'testing 123' to KafkaTopic each second, until receive a os signal to stop e.g. control + c
// by the user in terminal
func produceMessages(producer sarama.AsyncProducer, signals chan os.Signal) {
	for {
		time.Sleep(time.Second)
		message := &sarama.ProducerMessage{Topic: KafkaTopic, Value: sarama.StringEncoder("testing 123")}
		select {
		case producer.Input() <- message:
			enqueued++
			log.Println("New Message produced")
		case <-signals:
			producer.AsyncClose() // Trigger a shutdown of the producer.
			return
		}
	}
}

// produceMessages will send 'testing 123' to KafkaTopic each second, until receive a os signal to stop e.g. control + c
// by the user in terminal
func produceMessagesSync(producer sarama.SyncProducer, signals chan os.Signal) {
	for {
		time.Sleep(time.Second)
		message := &sarama.ProducerMessage{Topic: KafkaTopic, Value: sarama.StringEncoder("testing 123")}
		select {
		default:
			_, _, err := producer.SendMessage(message)
			if err != nil {
				log.Println("Error:", err.Error())
			}
			enqueued++
			log.Println("New Message produced")
		case <-signals:
			_ = producer.Close() // Trigger a shutdown of the producer.
			return
		}
	}
}
