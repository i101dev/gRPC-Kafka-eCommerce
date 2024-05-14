package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/IBM/sarama"
)

func StartConsumer(topic string, KAFKA_URI string, quitCh <-chan os.Signal) {

	consumer, err := connect([]string{KAFKA_URI})

	if err != nil {
		log.Fatalf("Error connecting Kafka consumer for topic %s: %v", topic, err)
		return
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			fmt.Printf("Error closing Kafka consumer for topic %s: %v\n", topic, err)
		}
	}()

	topicConsumer, err := initialize(consumer, topic, 0, 0)

	if err != nil {
		fmt.Printf("Error starting Kafka consumer for topic %s: %v\n", topic, err)
		return
	}

	doneCh := make(chan struct{})

	processMessages(topicConsumer, quitCh, doneCh)

	<-doneCh

	fmt.Printf("Stopped consuming messages for topic %s\n", topic)
}

func connect(brokerURL []string) (sarama.Consumer, error) {

	config := sarama.NewConfig()

	config.Consumer.Return.Errors = true

	conn, err := sarama.NewConsumer(brokerURL, config)

	if err != nil {
		return nil, err
	}

	return conn, nil
}

func initialize(consumer sarama.Consumer, topic string, partition int32, offset int64) (sarama.PartitionConsumer, error) {

	partitionConsumer, err := consumer.ConsumePartition(topic, partition, offset)

	if err != nil {
		return nil, err
	}

	fmt.Println("\n*** >>> Consumer started for topic:", topic)

	return partitionConsumer, nil
}

func processMessages(consumer sarama.PartitionConsumer, quitCh <-chan os.Signal, doneCh chan struct{}) {

	go func() {
		for {
			select {

			case err := <-consumer.Errors():
				fmt.Println("\n*** >>> [consumer.error] -", err)

			case msg := <-consumer.Messages():
				handleKafkaMsg(msg)

			case <-quitCh:
				fmt.Println("\nInterruption detected")
				doneCh <- struct{}{}
				return
			}
		}
	}()
}

func handleKafkaMsg(msg *sarama.ConsumerMessage) {

	var data KafkaMsg

	if err := json.Unmarshal(msg.Value, &data); err != nil {
		fmt.Println("Error decoding message:", err)
	}

	fmt.Printf("\n*** >>> [msg] - %s", data.Msg)
	fmt.Printf("\n*** >>> [val] - %d", data.Val)
}
