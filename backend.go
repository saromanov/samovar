package samovar

import (
	"gopkg.in/redis.v3"
	//"log"
	"fmt"
	"time"
)

type Message *redis.Message

const (
	TITLE = "samovar"
)

//This file should be redisbackend

type RedisBackend struct {
	client *redis.Client
	pubsub *redis.PubSub
	queues []string
}

func InitRedisBackend() *RedisBackend {
	rb := new(RedisBackend)
	rb.client = initRedis("localhost:6379")
	rb.pubsub = rb.client.PubSub()
	rb.queues = []string{"lightstore1"}
	return rb
}

func initRedis(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
}

//Add new queue name
func (backend *RedisBackend) registerQueue(queuename string) {
	backend.queues = append(backend.queues, queuename)
	backend.subscribe(queuename)
}

//Receive messages from backend
func (backend *RedisBackend) receiveMessages() (interface{}, error) {
	return backend.pubsub.ReceiveTimeout(100 * time.Millisecond)
}

//subscribe to the new channel
func (backend *RedisBackend) subscribe(queuename string) {
	err := backend.pubsub.Subscribe(fmt.Sprintf("%s_%s", TITLE, queuename))
	if err != nil {
		panic(err)
	}
}

//Publish new hob
func (backend *RedisBackend) publishJob(jobtitle string) {
	err := backend.client.Publish(fmt.Sprintf("%s_%s", TITLE, "default"), jobtitle).Err()
	if err != nil {
		panic(err)
	}
}
