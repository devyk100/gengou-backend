package kafka_internal

import (
	"context"
	"crypto/tls"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/scram"
	"log"
	"sync"
)

type KafkaProducer struct {
	topic     string
	addr      string
	username  string
	password  string
	mechanism *sasl.Mechanism
	staller   chan *string
	w         *kafka.Writer
	ctx       context.Context
	cancel    context.CancelFunc
	waitGroup *sync.WaitGroup
}

func (instance *KafkaProducer) CreateProducer(topic string, addr string, username string, password string) {
	instance.topic = topic
	instance.addr = addr
	instance.username = username
	instance.password = password
	mechanism, _ := scram.Mechanism(scram.SHA512, username, password)
	instance.mechanism = &mechanism
	ctx, cancel := context.WithCancel(context.Background())
	instance.ctx = ctx
	instance.cancel = cancel
	w := kafka.Writer{
		Addr:  kafka.TCP(addr),
		Topic: topic,
		Transport: &kafka.Transport{
			SASL: mechanism,
			TLS:  &tls.Config{},
		},
	}
	instance.staller = make(chan *string, 100)
	instance.w = &w
	instance.waitGroup = &sync.WaitGroup{}
	instance.waitGroup.Add(1)
	go func(instance *KafkaProducer) {
		defer instance.waitGroup.Done()
		for {
			select {
			case val, ok := <-instance.staller:
				if !ok {
					return
				}
				if val == nil {
					return
				}
				err := instance.w.WriteMessages(instance.ctx, kafka.Message{Value: []byte(*val)})
				if err != nil {
					log.Println(err.Error())
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}(instance)
	return
}

func (instance *KafkaProducer) Produce(val string) {
	instance.staller <- &val
	return
}

func (instance *KafkaProducer) CloseProducer() {
	instance.staller <- nil
	instance.waitGroup.Wait()
	close(instance.staller)
	err := instance.w.Close()
	instance.cancel()
	if err != nil {
		log.Println(err.Error())
		return
	}

}

//func Producer(topic string, addr string, username string, password string) {
//	mechanism, _ := scram.Mechanism(scram.SHA256, username, password)
//	w := kafka.Writer{
//		Addr:  kafka.TCP(addr),
//		Topic: topic,
//		Transport: &kafka.Transport{
//			SASL: mechanism,
//			TLS:  &tls.Config{},
//		},
//	}
//	err := w.WriteMessages(context.Background(), kafka.Message{Value: []byte("HAHAHAHAHAHA 09132rjvnijk nifen")})
//	if err != nil {
//		log.Println(err.Error())
//		return
//	}
//	err = w.Close()
//	if err != nil {
//		log.Println(err.Error())
//		return
//	}
//}
