package queue

import (
	"context"
	"fmt"
	"time"

	"example-wikipedia-scraper/internal/config"
	queueInterfaces "example-wikipedia-scraper/internal/interfaces/queue"
	"example-wikipedia-scraper/internal/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQService struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	queues   map[string]amqp.Queue
	handlers map[string]func(*queueInterfaces.Task)
	cfg      *config.RabbitMQConfig
	dlxSetup map[string]bool // Track which queues have DLX setup
}

func NewRabbitMQService(cfg *config.RabbitMQConfig) *RabbitMQService {
	svc := &RabbitMQService{
		queues:   make(map[string]amqp.Queue),
		handlers: make(map[string]func(*queueInterfaces.Task)),
		cfg:      cfg,
		dlxSetup: make(map[string]bool),
	}
	if err := svc.reconnect(); err != nil {
		panic(err)
	}
	return svc
}

func (r *RabbitMQService) reconnect() error {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", r.cfg.User, r.cfg.Password, r.cfg.Host, r.cfg.Port, r.cfg.Vhost)
	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}
	if r.conn != nil {
		r.conn.Close()
	}
	r.conn = conn
	r.channel = ch
	r.queues = make(map[string]amqp.Queue)
	return nil
}

func (r *RabbitMQService) Publish(task *queueInterfaces.Task) error {
	if r.channel == nil || r.conn == nil || r.conn.IsClosed() {
		if err := r.reconnect(); err != nil {
			return err
		}
	}
	queueName := task.Type
	r.getOrCreateQueue(queueName)
	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(5*time.Second))
	defer cancel()
	return r.channel.PublishWithContext(
		ctx,
		"",                            // exchange
		r.getFullQueueName(queueName), // routing key
		false,                         // mandatory
		false,                         // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         []byte(task.Payload),
			DeliveryMode: amqp.Persistent,
		},
	)
}

// PublishWithDelay publikuje wiadomość z opóźnieniem przez DLX
func (r *RabbitMQService) PublishWithDelay(task *queueInterfaces.Task, delay time.Duration) error {
	if r.channel == nil || r.conn == nil || r.conn.IsClosed() {
		if err := r.reconnect(); err != nil {
			return err
		}
	}

	queueName := task.Type
	r.setupDelayedQueue(queueName)

	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(5*time.Second))
	defer cancel()

	return r.channel.PublishWithContext(
		ctx,
		"",                                   // default exchange
		r.getFullQueueName(queueName)+"_dlx", // routing key = delayed queue name
		false,                                // mandatory
		false,                                // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         []byte(task.Payload),
			DeliveryMode: amqp.Persistent,
			Expiration:   fmt.Sprintf("%d", int(delay.Milliseconds())),
		},
	)
}

func (r *RabbitMQService) RegisterHandler(taskType string, handler func(*queueInterfaces.Task)) {
	r.getOrCreateQueue(taskType)
	r.handlers[taskType] = handler
}

func (r *RabbitMQService) Start() {
	if r.channel == nil || r.conn == nil || r.conn.IsClosed() {
		if err := r.reconnect(); err != nil {
			panic(err)
		}
	}
	for queueName, handler := range r.handlers {
		msgs, err := r.channel.Consume(
			r.getFullQueueName(queueName),
			"",    // consumer
			false, // autoAck
			false, // exclusive
			false, // noLocal
			false, // noWait
			nil,   // args
		)
		if err != nil {
			panic(err)
		}
		go func(qName string, h func(*queueInterfaces.Task), m <-chan amqp.Delivery) {
			for d := range m {
				payload := queueInterfaces.JSONString(d.Body)
				func() {
					defer func() {
						if rec := recover(); rec != nil {
							err := d.Reject(true)
							if err != nil {
								r.log(err, &queueInterfaces.Task{Type: qName, Payload: payload})
							}
						}
					}()
					h(&queueInterfaces.Task{Type: qName, Payload: payload})
					err := d.Ack(false)
					if err != nil {
						r.log(err, &queueInterfaces.Task{Type: qName, Payload: payload})
					}
				}()
			}
			r.log("stop consuming", nil)

		}(queueName, handler, msgs)
	}
}

func (r *RabbitMQService) Close() {
	r.channel.Close()
	r.conn.Close()
}

func (r *RabbitMQService) log(err any, task *queueInterfaces.Task) {
	logger := logger.GetLogger()
	if logger != nil {
		logger.Error("Panic recovered in RabbitMQService handler", "error", err, "task,", task)
	}
}

func (r *RabbitMQService) getOrCreateQueue(taskType string) (amqp.Queue, error) {
	if _, ok := r.queues[taskType]; !ok {
		q, err := r.channel.QueueDeclare(
			r.getFullQueueName(taskType),
			true,  // durable
			false, // autoDelete
			false, // exclusive
			false, // noWait
			nil,   // args
		)
		if err != nil {
			panic(err)
		}
		r.queues[taskType] = q
	}
	return r.queues[taskType], nil
}

func (r *RabbitMQService) getFullQueueName(taskType string) string {
	return taskType
}

func (r *RabbitMQService) setupDelayedQueue(queueName string) error {
	fullName := r.getFullQueueName(queueName)
	if r.dlxSetup[fullName] {
		return nil
	}

	// Dead Letter Exchange
	err := r.channel.ExchangeDeclare(
		fullName+"_dlx_exchange",
		"direct",
		true,  // durable
		false, // autoDelete
		false, // internal
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return err
	}

	// Delayed queue with TTL i DLX
	_, err = r.channel.QueueDeclare(
		fullName+"_dlx",
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		amqp.Table{
			"x-dead-letter-exchange":    fullName + "_dlx_exchange",
			"x-dead-letter-routing-key": fullName,
		},
	)
	if err != nil {
		return err
	}

	// Bind DLX do docelowej queue
	err = r.channel.QueueBind(
		fullName,
		fullName,
		fullName+"_dlx_exchange",
		false,
		nil,
	)
	if err != nil {
		return err
	}

	r.dlxSetup[fullName] = true
	return nil
}
