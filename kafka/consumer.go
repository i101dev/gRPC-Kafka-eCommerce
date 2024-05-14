package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/IBM/sarama"
)

func StartKafkaConsumer(topic string, KAFKA_URI string, quitCh <-chan os.Signal) {

	consumer, err := connectConsumer([]string{KAFKA_URI})

	if err != nil {
		log.Fatalf("Error connecting Kafka consumer for topic %s: %v", topic, err)
		return
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			fmt.Printf("Error closing Kafka consumer for topic %s: %v\n", topic, err)
		}
	}()

	topicConsumer, err := startConsumer(consumer, topic, 0, 0)

	if err != nil {
		fmt.Printf("Error starting Kafka consumer for topic %s: %v\n", topic, err)
		return
	}

	doneCh := make(chan struct{})

	processMessages(topicConsumer, quitCh, doneCh)

	<-doneCh

	fmt.Printf("Stopped consuming messages for topic %s\n", topic)
}

func connectConsumer(brokerURL []string) (sarama.Consumer, error) {

	config := sarama.NewConfig()

	config.Consumer.Return.Errors = true

	conn, err := sarama.NewConsumer(brokerURL, config)

	if err != nil {
		return nil, err
	}

	return conn, nil
}

func startConsumer(consumer sarama.Consumer, topic string, partition int32, offset int64) (sarama.PartitionConsumer, error) {

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
				destructureMSG(msg)

			case <-quitCh:
				fmt.Println("\nInterruption detected")
				doneCh <- struct{}{}
				return
			}
		}
	}()
}

func destructureMSG(msg *sarama.ConsumerMessage) {

	// fmt.Printf("Msg count: %d: | Topic (%s) | Message (%s)\n", last, string(msg.Topic), msg.Value)

	var data map[string]interface{}

	if err := json.Unmarshal(msg.Value, &data); err != nil {
		fmt.Println("Error decoding message:", err)
	}

	if msgValue, ok := data["msg"].(string); ok {
		fmt.Printf("\n*** >>> Message received - %s", msgValue)
	} else {
		fmt.Println("Message does not contain 'msg' field or it's not a string")
	}
}
