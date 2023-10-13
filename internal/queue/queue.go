package queue

import (
	"fmt"
	"log"
	"reflect"
)

const (
	RabbitMQ QueueType = iota
)

type QueueType int

type QueueConnection interface {
	Publish([]byte) error
	Consume() error
}
type Queue struct {
	cfg any
	qc  QueueConnection
}

func New(qt QueueType, cfg any) (*Queue, error) {
	q := new(Queue)
	rt := reflect.TypeOf(cfg)
	switch qt {
	case RabbitMQ:
		if rt != reflect.TypeOf(RabbitMQConfig{}) {
			return nil, fmt.Errorf("invalid config type")
		}
		fmt.Println("NOT IMPLEMENTED YET")
	default:
		log.Fatal("Queue type not implemented")

	}

	return q, nil
}

func (q *Queue) Publish(msg []byte) error {
	return q.qc.Publish(msg)
}

func (q *Queue) Consume() error {
	return q.qc.Consume()
}
