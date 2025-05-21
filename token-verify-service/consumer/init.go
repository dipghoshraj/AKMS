package consumer

import "tvs/dbops"

type Consumer struct {
	Topic    string
	brokers  []string
	tokenOps dbops.TokenOps
}

func NewComsumer(brokers []string, topic string, ops *dbops.OpsManager) *Consumer {
	return &Consumer{
		brokers:  brokers,
		Topic:    topic,
		tokenOps: ops.TokenOps,
	}
}
