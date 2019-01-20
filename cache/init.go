package cache

import (
	"log"
	"os"

	"github.com/dstpierre/gosaas/queue"
	"github.com/go-redis/redis"
)

var rc *redis.Client

func init() {
	host := os.Getenv("REDIS_ADDR")
	if len(host) == 0 {
		host = "127.0.0.1:6379"
	}

	c := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: os.Getenv("REDIS_KEY"),
		DB:       0, // use default DB
	})

	if _, err := c.Ping().Result(); err != nil {
		log.Fatal("unable to connect to redis", err)
	}

	rc = c
}

// New initializes the queue service via the queue.New function.
//
// The queueProcessor flag indicates if this instance will act
// as the Pub/Sub subscriber. There must be only one subscriber.
//
// The ex parameter map[queue.TaskID]queue.Executor allow you to supply
// custom executors for your own custom task. A TaskExecutor must satisfy
// this interface.
//
// 	type TaskExecutor interface {
// 		Run(t QueueTask) error
// 	}
func New(queueProcessor, isDev bool, ex map[queue.TaskID]queue.TaskExecutor) {
	queue.New(rc, isDev, ex)

	if queueProcessor {
		queue.SetAsSubscriber()
	}
}
