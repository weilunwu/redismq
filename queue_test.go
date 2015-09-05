package redismq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	redisHost     = "localhost"
	redisPort     = "6379"
	redisPassword = ""
	redisDB       = int64(4)

	name = "testRedismq"
)

func TestSetup(t *testing.T) {
	queue := newQueue(redisHost, redisPort, redisPassword, redisDB, name)
	assert.Nil(t, queue.redisClient.FlushDb().Err())
}
