package queue

import "time"

type MessageQueueServiceInterface interface {
	Publish(task *Task) error
	PublishWithDelay(task *Task, delay time.Duration) error
	RegisterHandler(taskType string, handler func(*Task))
	Start()
	Close()
}

type JSONString string

type Task struct {
	Type    string
	Payload JSONString
}
