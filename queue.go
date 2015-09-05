package redismq

import "gopkg.in/redis.v3"

type Queue struct {
	redisClient *redis.Client
	Name        string
}

func CreateQueue(redisHost, redisPort, redisPassword string, redisDB int64, name string) *Queue {
	return newQueue(redisHost, redisPort, redisPassword, redisDB, name)
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
	return queue
}
