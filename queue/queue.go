package queue

import (
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis"
)

var (
	client *redis.Client
	pubsub *redis.PubSub
	isDev  bool

	emailer *Email
)

func New(rc *redis.Client, isDev bool) {
	client = rc

	emailer = &Email{}
	if isDev {
		emailer.Send = emailer.sendEmailDev
	} else {
		emailer.Send = emailer.sendEmailProd
	}
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

		go process(msg)
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
	var qt QueueTask
	if err := json.Unmarshal([]byte(msg.Payload), &qt); err != nil {
		log.Fatal("unable to decode this Redis message", err)
	}

	var exec TaskExecutor

	switch qt.ID {
	case TaskEmail:
		exec = emailer
	}

	if err := exec.Run(qt); err != nil {
		//TODO: better to log those critical errors
		log.Println("error while executing this task", qt.ID, err)
	}
}
