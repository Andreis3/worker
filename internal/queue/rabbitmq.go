package queue

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConfig struct {
	URL       string
	TopicName string
	Timeout   time.Duration
}

type RabbitMQConnection struct {
	cfg  RabbitMQConfig
	conn *amqp.Connection
}

func newRabbitMQConnection(cfg RabbitMQConfig) (*RabbitMQConnection, error) {
	var err error
	r := RabbitMQConnection{
		cfg: cfg,
	}

	r.conn, err = amqp.Dial(r.cfg.URL)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (r *RabbitMQConnection) Publish(msg []byte) error {
	c, err := r.conn.Channel()
	if err != nil {
		return err
	}

	mp := amqp.Publishing{
		ContentType:  "text/plain",
		Body:         msg,
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = c.PublishWithContext(ctx, r.cfg.TopicName, "", false, false, mp)
	if err != nil {
		return err
	}
	return nil
}

func (r *RabbitMQConnection) Consume(data chan<- QueueDto) error {
	c, err := r.conn.Channel()
	if err != nil {
		return err
	}

	q, err := c.QueueDeclare(r.cfg.TopicName, false, false, false, false, nil)
	if err != nil {
		return err
	}

	msgs, err := c.Consume(q.Name, "", true, false, false, false, nil)

	if err != nil {
		return err
	}

	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
		qd := new(QueueDto)
		err := qd.Unmarshal(d.Body)
		if err != nil {
			log.Printf("Error unmarshalling message: %s", err)
			continue
		}
		data <- *qd
	}

	return nil
}
