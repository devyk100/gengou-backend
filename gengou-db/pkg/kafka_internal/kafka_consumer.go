package kafka_internal

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/scram"
	"log"
	"time"
)

type KafkaConsumer struct {
	topic       string
	groupId     string
	brokers     []string
	username    string
	password    string
	mechanism   *sasl.Mechanism
	r           *kafka.Reader
	isRunning   bool
	ctx         context.Context
	cancel      context.CancelFunc
	messageChan chan kafka.Message
}

func (instance *KafkaConsumer) CreateConsumer(topic string, groupId string, brokers []string, username string, password string) {
	mechanism, err := scram.Mechanism(scram.SHA512, username, password)
	if err != nil {
		log.Println(err.Error())
	}
	instance.ctx, instance.cancel = context.WithCancel(context.Background())
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupId,
		Topic:   topic,
		Dialer: &kafka.Dialer{
			SASLMechanism: mechanism,
			TLS:           &tls.Config{},
		},
		CommitInterval:        1 * time.Millisecond,
		WatchPartitionChanges: true,
	})
	instance.password = password
	instance.brokers = brokers
	instance.topic = topic
	instance.groupId = groupId
	instance.r = r
	instance.mechanism = &mechanism
	instance.isRunning = false
	instance.messageChan = make(chan kafka.Message, 100)
	fmt.Println("Connected as a consumer")
	//ctx, cancel := context.WithTimeout(context.Background(), time.Hour*1000)
	//defer cancel()
}

func (instance *KafkaConsumer) StartConsumer() {
	instance.isRunning = true
	go func(instance *KafkaConsumer) {
		for instance.isRunning {
			select {
			case <-instance.ctx.Done():
				return
			default:
				message, err := instance.r.ReadMessage(context.Background())
				if err != nil {
					log.Println(err.Error())
				}
				instance.messageChan <- message
			}
		}
	}(instance)
}

func (instance *KafkaConsumer) StopConsumer() {
	instance.isRunning = false
}

func (instance *KafkaConsumer) GetReader() *kafka.Reader {
	return instance.r
}

func (instance *KafkaConsumer) GetMessageChan() *chan kafka.Message {
	return &instance.messageChan
}

func (instance *KafkaConsumer) CloseConsumer() {
	instance.isRunning = false
	err := instance.r.Close()
	instance.cancel()
	if err != nil {
		log.Println(err.Error())
		return
	}
}

//func Consumer(topic string, groupId string, brokers []string, username string, password string) {
//	mechanism, _ := scram.Mechanism(scram.SHA512, username, password)
//	r := kafka.NewReader(kafka.ReaderConfig{
//		Brokers: brokers,
//		GroupID: groupId,
//		Topic:   topic,
//		Dialer: &kafka.Dialer{
//			SASLMechanism: mechanism,
//			TLS:           &tls.Config{},
//		},
//		CommitInterval:        1 * time.Millisecond,
//		WatchPartitionChanges: true,
//	})
//	ctx, cancel := context.WithTimeout(context.Background(), time.Hour*1000) // Increase the timeout
//	defer cancel()
//	i := 0
//	for {
//		i++
//		if i > 15 {
//			r.Close()
//		}
//		message, _ := r.ReadMessage(ctx)
//		//if message.Offset == 0 {
//		//	continue
//		//
//		r.Close()
//		fmt.Println(message.Partition, message.Offset, string(message.Value))
//	}
//	r.Close()
//}
