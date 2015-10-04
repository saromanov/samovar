package backend

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
	rb.client = InitRedis("localhost:6379")
	rb.pubsub = rb.client.PubSub()
	rb.queues = []string{"lightstore1"}
	return rb
}

func InitRedis(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
}

//Add new queue name
func (backend *RedisBackend) RegisterQueue(queuename string) {
	backend.queues = append(backend.queues, queuename)
	backend.Subscribe(queuename)
}

//Receive messages from backend
func (backend *RedisBackend) ReceiveMessages() (interface{}, error) {
	return backend.pubsub.ReceiveTimeout(100 * time.Millisecond)
}

//subscribe to the new channel
func (backend *RedisBackend) Subscribe(queuename string) {
	err := backend.pubsub.Subscribe(fmt.Sprintf("%s_%s", TITLE, queuename))
	if err != nil {
		panic(err)
	}
}

//Publish new hob
func (backend *RedisBackend) PublishJob(jobtitle string, queuename string) {
	fmt.Println(fmt.Sprintf("%s_%s", TITLE, queuename))
	err := backend.client.Publish(fmt.Sprintf("%s_%s", TITLE, queuename), jobtitle).Err()
	if err != nil {
		panic(err)
	}
}
