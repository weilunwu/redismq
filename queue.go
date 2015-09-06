package redismq

import (
	"fmt"
	"time"

	"gopkg.in/redis.v3"
)

type Queue struct {
	redisClient *redis.Client
	Name        string
}

func CreateQueue(redisHost, redisPort, redisPassword string, redisDB int64, name string) *Queue {
	return newQueue(redisHost, redisPort, redisPassword, redisDB, name)
}

func SelectQueue(redisHost, redisPort, redisPassword string, redisDB int64, name string) (*Queue, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
		DB:       redisDB,
	})
	defer redisClient.Close()

	ok, err := redisClient.SIsMember(masterQueueKey(), name).Result()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("queue with this name doesn't exist")
	}
	return newQueue(redisHost, redisPort, redisPassword, redisDB, name), nil
}

func newQueue(redisHost, redisPort, redisPassword string, redisDB int64, name string) *Queue {
	queue := &Queue{
		Name: name,
		redisClient: redis.NewClient(&redis.Options{
			Addr:     redisHost + ":" + redisPort,
			Password: redisPassword,
			DB:       redisDB,
		}),
	}
	queue.redisClient.SAdd(masterQueueKey(), name)
	return queue
}

func (queue *Queue) Put(payload string) error {
	pack := &Package{
		CreatedAt: time.Now(),
		Payload:   payload,
		Queue:     queue,
	}
	lpush := queue.redisClient.LPush(queueInputKey(queue.Name), pack.getString())
	return lpush.Err()
}

func (queue *Queue) AddConsumer(name string) (*Consumer, error) {
	consumer := &Consumer{
		Name:  name,
		Queue: queue,
	}
	// check uniqueness
	added, err := queue.redisClient.SAdd(queueWorkersKey(name), name).Result()
	if err != nil {
		return nil, err
	}
	if added == 0 {
		if queue.isActiveConsumer(name) {
			return nil, fmt.Errorf("consumer with this name already exists")
		}
	}
	consumer.startHeartbeat()
	return consumer, nil
}

func (queue *Queue) isActiveConsumer(name string) bool {
	val := queue.redisClient.Get(consumerHeartbeatKey(queue.Name, name)).Val()
	return val == "ping"
}
