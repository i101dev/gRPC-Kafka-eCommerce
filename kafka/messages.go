package kafka

import (
	"github.com/IBM/sarama"
)

type OrderMsg struct {
	Msg string `form:"msg" json:"msg"`
	Val int64  `form:"val" json:"val"`
}
type ProductMsg struct {
	Msg string `form:"msg" json:"msg"`
	Val int64  `form:"val" json:"val"`
}
type UserMsg struct {
	Msg string `form:"msg" json:"msg"`
	Val int64  `form:"val" json:"val"`
}

func PushMsgToQueue(kafkaURI string, kafkaTopic string, msgInBytes []byte) error {

	brokerURLs := []string{kafkaURI}
	producer, err := ConnectProducer(brokerURLs)

	if err != nil {
		return err
	}

	defer producer.Close()

	err = SendMessage(kafkaTopic, msgInBytes, producer)

	return err
}

func ConnectProducer(brokerURLs []string) (sarama.SyncProducer, error) {

	config := sarama.NewConfig()

	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokerURLs, config)

	if err != nil {
		return nil, HandleKafkaError(err, "Failed to connect to Kafka broker")
	}

	return producer, nil
}

func SendMessage(kafkaTopic string, messageInByts []byte, producer sarama.SyncProducer) error {

	producerMsg := &sarama.ProducerMessage{
		Topic:    kafkaTopic,
		Value:    sarama.StringEncoder(messageInByts),
		Metadata: []string{1: "metadata 1", 2: "metadata 2", 3: "metadata 3"},
	}

	// fmt.Printf("\n*** >>> TOPIC: (%s) - BYTES: (%+v)", kafkaTopic, messageInByts)

	if _, _, err := producer.SendMessage(producerMsg); err != nil {
		return HandleKafkaError(err, "Failed to send message to Kafka")
	}

	return nil
}
