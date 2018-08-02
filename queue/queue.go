package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
)

var (
	client *redis.Client
	pubsub *redis.PubSub
)

func New(rc *redis.Client) {
	client = rc
}

func SetAsSubscriber() {
	pubsub = client.Subscribe("q")
	if err := pubsub.Ping("test"); err != nil {
		log.Fatal("unable to ping pubsub", err)
	}
	defer pubsub.Close()

	if _, err := pubsub.Receive(); err != nil {
		log.Fatal("unable to receive from pubsub channel", err)
	}

	ch := pubsub.Channel()

	for {
		msg, ok := <-ch
		if !ok {
			log.Fatal("redis pub/sub is down")
			break
		}

		process(msg)
	}
}

func Enqueue(id TaskID, data interface{}) error {
	qt := QueueTask{
		ID:      id,
		Data:    data,
		Created: time.Now(),
	}

	b, err := json.Marshal(qt)
	if err != nil {
		return err
	}
	return client.Publish("q", string(b)).Err()
}

func process(msg *redis.Message) {
	fmt.Println(msg.Channel, msg.Payload)
}
