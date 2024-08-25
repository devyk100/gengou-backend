package redis_internal

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type RedisInstance struct {
	dbUrl      string
	Client     *redis.Client
	expiration time.Duration
}

func Init(dbUrl string, expiration time.Duration) (*RedisInstance, error) {
	var instance RedisInstance
	instance.dbUrl = dbUrl
	instance.expiration = expiration
	opt, err := redis.ParseURL(instance.dbUrl)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)
	instance.Client = client
	log.Println("redis connection init success")
	return &instance, nil
}

func (instance RedisInstance) Ping() error {
	instance.Client.Ping(context.Background())
	return nil
}

func (instance RedisInstance) Get(key string) (string, error) {
	val := instance.Client.Get(context.Background(), key)
	return val.Val(), nil
}

func (instance RedisInstance) Set(key string, value string, expiration time.Duration) error {
	instance.Client.Set(context.Background(), key, value, expiration)
	return nil
}

func (instance RedisInstance) Close() error {
	err := instance.Client.Close()
	if err != nil {
		return err
	}
	return nil
}

type SomeType struct {
	name string
	age  int
}

func (instance RedisInstance) HSet(key string, values ...interface{}) {
	instance.Client.HSet(context.Background(), key, values)
}

func (instance RedisInstance) HGet(key string, field string) string {
	val := instance.Client.HGet(context.Background(), key, field)
	return val.Val()
}

func (instance RedisInstance) HDel(key string, fields ...string) {
	val := instance.Client.HDel(context.Background(), key, fields...)
	fmt.Println(val)
}

func (instance RedisInstance) Subscribe(key string, callback func(string), subscribeStaller *chan struct{}) {
	pubSub := instance.Client.Subscribe(context.Background(), key)

	if subscribeStaller == nil {
		subscribeStallerVal := make(chan struct{})
		subscribeStaller = &subscribeStallerVal
	}

	go func() {
		for {
			receive, err := pubSub.ReceiveMessage(context.Background())
			if err != nil {
				log.Println(err.Error())
				return
			}
			callback(receive.Payload)
		}
	}()

	<-*subscribeStaller
}

func (instance RedisInstance) Publish(key string, value string) {
	instance.Client.Publish(context.Background(), key, value)
}
