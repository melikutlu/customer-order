package consumer

import (
	"context"
	"github.com/labstack/gommon/log"
	"github.com/segmentio/kafka-go"
	"time"
)

type Consumer struct {
	reader *kafka.Reader
	dialer *kafka.Dialer
	Topic  string
}

func (c *Consumer) CreateConnection(brokers []string) {
	c.dialer = &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	c.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:   brokers,
		Topic:     c.Topic,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
		MaxWait:   time.Second * 1,
		Dialer:    c.dialer,
	})

	c.reader.SetOffset(kafka.LastOffset)
}

//config := &kafka.ConfigMap{ "bootstrap.servers": "localhost:9092", "group.id":  "my-group", "auto.offset.reset": "earliest", "enable.auto.commit": true, // Enable auto commit }

func (c *Consumer) Read(callback func(string, error)) {
	for {

		message, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			log.Errorf("Error reading message: %v", err)
			continue
		}

		uuid := string(message.Value)
		callback(uuid, nil)
	}
}
func (c *Consumer) Close() error {
	return c.reader.Close()
}
