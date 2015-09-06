package redismq

func masterQueueKey() string {
	return "redismq::queues"
}

func queueInputKey(queue string) string {
	return "redismq::" + queue
}

func queueWorkersKey(queue string) string {
	return queue + "::workers"
}

func queueWorkingPrefix(queue string) string {
	return "redismq::" + queue + "::working"
}

func consumerWorkingQueueKey(queue, consumer string) string {
	return queueWorkingPrefix(queue) + "::" + consumer
}

func consumerHeartbeatKey(queue, consumer string) string {
	return consumerWorkingQueueKey(queue, consumer) + "::heartbeat"
}
