package kafka

import (
	"fmt"
	"log"
)

type KafkaError struct {
	OriginalError error
	Message       string
}

func HandleKafkaError(err error, context string) error {
	log.Printf("*** >>> Kafka Error (%s): %v", context, err)
	return &KafkaError{
		OriginalError: err,
		Message:       context,
	}
}
func (e *KafkaError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.OriginalError)
}
