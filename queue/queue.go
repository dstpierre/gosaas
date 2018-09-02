package queue

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/robfig/cron"
)

var (
	client    *redis.Client
	pubsub    *redis.PubSub
	scheduler *cron.Cron
	isDev     bool

	emailer *Email
	biller  *Billing
)

func New(rc *redis.Client, isDev bool) {
	client = rc

	emailer = &Email{}
	biller = &Billing{}
	if isDev {
		emailer.Send = emailer.sendEmailDev
	} else {
		emailer.Send = emailer.sendEmailProd
	}
}

func SetAsSubscriber() {
	scheduler = cron.New()

	pubsub = client.Subscribe("q")
	if err := pubsub.Ping("test"); err != nil {
		log.Fatal("unable to ping pubsub", err)
	}
	defer func() {
		pubsub.Close()
		scheduler.Stop()
	}()

	if _, err := pubsub.Receive(); err != nil {
		log.Fatal("unable to receive from pubsub channel", err)
	}

	// we initialize our scheduler (cron)
	go setupCron()

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

func setupCron() {
	if _, err := os.Stat("tasks.cron"); os.IsNotExist(err) {
		log.Println("no tasks.cron file found, skipping scheduler setup")
		return
	}

	b, err := ioutil.ReadFile("tasks.cron")
	if err != nil {
		log.Println("error while reading tasks.cron", err)
		return
	}

	lines := strings.Split(string(b), "\n")
	if len(lines) == 0 {
		log.Println("no tasks found in tasks.cron, skipping scheduler setup")
		return
	}

	for _, line := range lines {
		exp, url := parseTask(line)

		err := scheduler.AddFunc(exp, func() {
			req, err := http.NewRequest("POST", url, bytes.NewReader(b))
			if err != nil {
				log.Println("error while creating an HTTP request to", url)
				return
			}

			req.SetBasicAuth("todo", "here")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Println("error while executing an HTTP request to", url)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode >= 400 {
				log.Println("scheduler HTTP request to ", url, "failed with HTTP status", resp.StatusCode)
			}
		})

		if err != nil {
			log.Fatal("unable to create cron tasks", err)
		}
	}

	scheduler.Start()
}

func parseTask(s string) (exp string, url string) {
	tokens := strings.Split(s, " ")
	url = strings.Join(tokens[len(tokens)-1:], " ")
	exp = strings.Join(tokens[0:len(tokens)-1], " ")
	return
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
	case TaskCreateInvoice:
		exec = biller
	}

	if err := exec.Run(qt); err != nil {
		//TODO: better to log those critical errors
		log.Println("error while executing this task", qt.ID, err)
	}
}
