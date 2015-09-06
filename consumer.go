package redismq

import "time"

type Consumer struct {
	Name  string
	Queue *Queue
}

func (consumer *Consumer) startHeartbeat() {
	firstWrite := make(chan bool, 1)
	go func() {
		firstRun := true
		for {
			consumer.Queue.redisClient.Set(
				consumerHeartbeatKey(consumer.Queue.Name, consumer.Name),
				"ping",
				time.Second,
			)
			if firstRun {
				firstWrite <- true
				firstRun = false
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()
	<-firstWrite
	return
}

func (consumer *Consumer) HasUnacked() bool {
	if consumer.GetUnackedLength() != 0 {
		return true
	}
	return false
}

func (consumer *Consumer) GetUnackedLength() int64 {
	return consumer.Queue.redisClient.LLen(consumerWorkingQueueKey(consumer.Queue.Name, consumer.Name)).Val()
}

func (consumer *Consumer) ackPackage(pack *Package) error {
	return consumer.Queue.redisClient.RPop(consumerWorkingQueueKey(consumer.Queue.Name, consumer.Name)).Err()
}
