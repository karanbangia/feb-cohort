package main

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

func main() {
	// Kafka configuration - equivalent to the kafkajs setup
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3

	// Enable idempotent producer to avoid duplicates
	config.Producer.Idempotent = true
	config.Net.MaxOpenRequests = 1

	// Create producer - equivalent to kafka.producer()
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	topic := "28-june"

	fmt.Println("Starting Go producer (equivalent to producer2.js)")
	fmt.Println("================================================")

	// Send 10 messages - equivalent to the for loop in JS
	for i := 0; i < 10; i++ {
		message := fmt.Sprintf("Message %d", i)

		// Create message with empty key
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(""), // Using empty key
			Value: sarama.StringEncoder(message),
		}

		// Add timestamp for tracking
		msg.Timestamp = time.Now()

		// Send the message - equivalent to producer.send() in JS
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.Printf("ERROR: Failed to send message %d: %v\n", i, err)
		} else {
			fmt.Printf("Sent: %s (partition: %d, offset: %d)\n", message, partition, offset)
		}

		// Small delay between messages for better visibility
	}

}
