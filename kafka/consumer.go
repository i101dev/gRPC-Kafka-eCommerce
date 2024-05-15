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

	processMessages(topic, topicConsumer, quitCh, doneCh)

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

func processMessages(topic string, consumer sarama.PartitionConsumer, quitCh <-chan os.Signal, doneCh chan struct{}) {

	go func() {
		for {
			select {

			case err := <-consumer.Errors():
				fmt.Println("\n*** >>> [consumer.error] -", err)

			case msg := <-consumer.Messages():
				handleKafkaMsg(topic, msg)

			case <-quitCh:
				fmt.Println("\nInterruption detected")
				doneCh <- struct{}{}
				return
			}
		}
	}()
}

func handleKafkaMsg(topic string, msg *sarama.ConsumerMessage) {

	switch topic {
	case "orders":
		var data = new(OrderMsg)
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			fmt.Printf("Error decoding (%s) message: %+v", topic, err)
		}
		fmt.Printf("\n*** >>> (ORDERS) - msg: (%s) - val: (%d)\n", data.Msg, data.Val)
	//
	//
	case "products":
		var data = new(ProductMsg)
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			fmt.Printf("Error decoding (%s) message: %+v", topic, err)
		}
		fmt.Printf("\n*** >>> (PRODUCTS) - msg: (%s) - val: (%d)\n", data.Msg, data.Val)
		//
		//
	case "users":
		var data = new(UserMsg)
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			fmt.Printf("Error decoding (%s) message: %+v", topic, err)
		}
		fmt.Printf("\n*** >>> (USERS) - msg: (%s) - val: (%d)\n", data.Msg, data.Val)
		//
		//
	default:
		break
	}

}
