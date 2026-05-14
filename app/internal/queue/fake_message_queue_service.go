package queue

import (
	queueInterfaces "example-wikipedia-scraper/internal/interfaces/queue"
	"sync"
	"time"
)

type FakeMessageQueueService struct {
	queue    chan *queueInterfaces.Task
	handlers map[string]func(*queueInterfaces.Task)
	mu       sync.RWMutex
}

func NewFakeMessageQueueService(buffer int) *FakeMessageQueueService {
	return &FakeMessageQueueService{
		queue:    make(chan *queueInterfaces.Task, buffer),
		handlers: make(map[string]func(*queueInterfaces.Task)),
	}
}

func (f *FakeMessageQueueService) Publish(task *queueInterfaces.Task) error {
	f.queue <- task
	return nil
}

// PublishWithDelay publikuje z opóźnieniem
func (f *FakeMessageQueueService) PublishWithDelay(task *queueInterfaces.Task, delay time.Duration) error {
	go func() {
		time.Sleep(delay)
		f.queue <- task
	}()
	return nil
}

func (f *FakeMessageQueueService) RegisterHandler(taskType string, handler func(*queueInterfaces.Task)) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.handlers[taskType] = handler
}

func (f *FakeMessageQueueService) Start() {
	go func() {
		for task := range f.queue {
			defer func() {
				if r := recover(); r != nil {
					f.queue <- task
				}
			}()
			f.mu.RLock()
			handler, ok := f.handlers[task.Type]
			f.mu.RUnlock()
			if ok && handler != nil {
				handler(task)
			}
		}
	}()
}

func (f *FakeMessageQueueService) Close() {
	close(f.queue)
}
