package queue

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/robfig/cron"
)

func TestQueue_Setup_Queue(t *testing.T) {
	fmt.Println("opening redis connection...")

	c := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	defer c.Close()

	if _, err := c.Ping().Result(); err != nil {
		t.Fatal("unable to connect to redis", err)
	}

	New(c, true)
	go SetAsSubscriber()

	time.AfterFunc(time.Second, func() {
		fmt.Println("enqueing something")

		err := Enqueue(TaskEmail, SendEmailParameter{
			From:    "me@testing.com",
			To:      "unit@test.com",
			Subject: "unit test",
			Body:    "<h1>unit test</h1>",
		})

		if err != nil {
			t.Fatal("unable to publish to channel", err)
		}
	})

	time.Sleep(3 * time.Second)
}

func TestQueue_ParseTask(t *testing.T) {
	cexp := "0 0 */15 0 0 0 https://google.com"
	exp, url := parseTask(cexp)
	if url != "https://google.com" {
		t.Errorf("url is %s was expecting https://google.com", url)
	} else if exp != "0 0 */15 0 0 0" {
		t.Errorf("exp is %s was expecting 0 0 */15 0 0 0", exp)
	}
}

func TestQueue_SetupCron(t *testing.T) {
	scheduler = cron.New()
	defer scheduler.Stop()

	setupCron()

	ch := make(chan bool)
	time.AfterFunc(17*time.Second, func() {
		ch <- true
	})

	<-ch
}
