package queue

import (
	queueInterfaces "example-wikipedia-scraper/internal/interfaces/queue"
)

var messageQueueServiceInstance queueInterfaces.MessageQueueServiceInterface

func InitMessageQueueService(service queueInterfaces.MessageQueueServiceInterface) {
	messageQueueServiceInstance = service
}

func GetMessageQueueService() queueInterfaces.MessageQueueServiceInterface {
	return messageQueueServiceInstance
}
