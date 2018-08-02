package queue

import "time"

type TaskID int

const (
	TaskEmail TaskID = iota
)

type QueueTask struct {
	ID      TaskID      `json:"id"`
	Data    interface{} `json:"data"`
	Created time.Time   `json:"created"`
}
