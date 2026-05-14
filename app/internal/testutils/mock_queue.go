package testutils

import (
	queueInterface "example-wikipedia-scraper/internal/interfaces/queue"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockQueueService struct {
	mock.Mock
}

func (m *MockQueueService) Publish(task *queueInterface.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockQueueService) RegisterHandler(queueName string, handler func(*queueInterface.Task)) {
	m.Called(queueName, handler)
}

func (m *MockQueueService) Close() {
	m.Called()
}

func (m *MockQueueService) Start() {
	m.Called()
}

func (m *MockQueueService) PublishWithDelay(task *queueInterface.Task, delay time.Duration) error {
	args := m.Called(task, delay)
	return args.Error(0)
}
