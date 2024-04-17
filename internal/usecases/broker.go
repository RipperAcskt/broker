package usecases

import (
	"context"
	"fmt"
	"github.com/RipperAcskt/broker/internal/entities/models"
	"github.com/RipperAcskt/broker/internal/repository/mongo/db/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const bufSize = 10

type exchangeService interface {
	Store(ctx context.Context, exchange models.Exchange) (primitive.ObjectID, error)
	All(ctx context.Context) ([]*dto.Exchange, error)
	UpdateByID(ctx context.Context, exchange models.Exchange) error
}

type queueService interface {
	Store(ctx context.Context, queue models.Queue) (primitive.ObjectID, error)
	FindByIDs(ctx context.Context, queueIDs []primitive.ObjectID) ([]*dto.Queue, error)
	UpdateByID(ctx context.Context, queue models.Queue) error
}

type messagesService interface {
	Store(ctx context.Context, message models.Message) (primitive.ObjectID, error)
	FindByIDs(ctx context.Context, messageIDs []primitive.ObjectID) ([]*dto.Message, error)
	UpdateByID(ctx context.Context, message models.Message) error
}

type Broker struct {
	exchanges       []*models.Exchange
	queues          map[string][]chan models.Message
	exchangeService exchangeService
	queueService    queueService
	messagesService messagesService
}

func (b *Broker) NewExchange(ctx context.Context, key string) error {
	exchange := models.Exchange{
		Key:    key,
		Queues: make([]*models.Queue, 0, 1),
	}

	queue := models.Queue{
		Messages: make([]*models.Message, 0),
	}

	id, err := b.queueService.Store(ctx, queue)
	if err != nil {
		return fmt.Errorf("store queue failed: %w", err)
	}

	queue.ID = id.Hex()
	exchange.Queues = append(exchange.Queues, &queue)

	id, err = b.exchangeService.Store(ctx, exchange)
	if err != nil {
		return fmt.Errorf("store exchange failed: %w", err)
	}

	exchange.ID = id.Hex()
	b.exchanges = append(b.exchanges, &exchange)

	ch := make(chan models.Message, bufSize)
	b.queues[key] = append(b.queues[key], ch)

	return nil
}

func (b *Broker) NewQueue(ctx context.Context, key string, offset int) (chan models.Message, int, error) {
	foundExchange := &models.Exchange{}
	for _, exchange := range b.exchanges {
		if exchange.Key == key {
			foundExchange = exchange
			break
		}
	}
	if foundExchange.ID == "" {
		return nil, 0, fmt.Errorf("no exchange found for key %s", key)
	}

	var queue models.Queue
	var queueIndex int
	var ch = make(chan models.Message, bufSize)
	for i, queueExchange := range foundExchange.Queues {
		if !queueExchange.UserConnected {
			queue = *queueExchange
			ch = b.queues[key][i]
			queueExchange.UserConnected = true
			queueIndex = i
			break
		}
	}

	if queue.ID == "" {
		queue := models.Queue{
			Messages: make([]*models.Message, 0),
		}

		id, err := b.queueService.Store(ctx, queue)
		if err != nil {
			return nil, 0, fmt.Errorf("store queue failed: %w", err)
		}

		queue.ID = id.Hex()
		foundExchange.Queues = append(foundExchange.Queues, &queue)

		err = b.exchangeService.UpdateByID(ctx, *foundExchange)
		if err != nil {
			return nil, 0, fmt.Errorf("update by id failed: %w", err)
		}

		ch = make(chan models.Message, bufSize)

		b.queues[key] = append(b.queues[key], ch)
	}

	messagesMap := make(map[models.Message]interface{})
	messages := make([]*models.Message, 0)
	for _, queue := range foundExchange.Queues {
		if queue.Messages != nil {
			for _, message := range queue.Messages {
				if _, ok := messagesMap[*message]; !ok {
					messagesMap[*message] = nil
					messages = append(messages, message)
				}
			}
		}
	}

	go func(queue chan models.Message, messages []*models.Message) {
		if offset > 0 && offset < len(messages) {
			messages = messages[offset:]
		}

		for _, message := range messages {
			ch <- *message
		}
	}(ch, messages)

	return ch, queueIndex, nil
}

func (b *Broker) NewMessage(ctx context.Context, key string, body any) error {
	foundExchange := &models.Exchange{}
	for _, exchange := range b.exchanges {
		if exchange.Key == key {
			foundExchange = exchange
			break
		}
	}
	if foundExchange.ID == "" {
		return fmt.Errorf("no exchange found for key %s", key)
	}

	message := models.Message{
		Body: body,
	}
	id, err := b.messagesService.Store(ctx, message)
	if err != nil {
		return fmt.Errorf("store message failed: %w", err)
	}

	message = models.Message{
		ID:   id.Hex(),
		Body: body,
	}

	for i, queue := range foundExchange.Queues {
		queue.Messages = append(queue.Messages, &message)

		err = b.queueService.UpdateByID(ctx, *queue)
		if err != nil {
			return fmt.Errorf("update by id failed: %w", err)
		}

		b.queues[key][i] <- message
	}

	return nil
}

func (b *Broker) ReceiveMessage(ctx context.Context, queue chan models.Message) (models.Message, error) {
	message := <-queue

	message.Read = true
	err := b.messagesService.UpdateByID(ctx, message)
	if err != nil {
		return models.Message{}, fmt.Errorf("update message failed: %w", err)
	}

	return message, nil
}

func (b *Broker) Disconnect(key string, queueIndex int) {
	foundExchange := &models.Exchange{}
	for _, exchange := range b.exchanges {
		if exchange.Key == key {
			foundExchange = exchange
			break
		}
	}
	foundExchange.Queues[queueIndex].UserConnected = false
	b.queues[key][queueIndex] = make(chan models.Message, bufSize)
}

func NewBroker(ctx context.Context, exchangeService exchangeService, queueService queueService, messagesService messagesService) (*Broker, error) {
	exchanges, err := exchangeService.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("exchange all failed: %w", err)
	}

	ex := make([]*models.Exchange, len(exchanges))
	queuesMap := make(map[string][]chan models.Message)
	for i, exchange := range exchanges {
		queues, err := queueService.FindByIDs(ctx, exchange.Queues)
		if err != nil {
			return nil, fmt.Errorf("queue find by ids failed: %w", err)
		}

		queuesModels := make([]*models.Queue, len(queues))
		for j, queue := range queues {
			messages, err := messagesService.FindByIDs(ctx, queue.Messages)
			if err != nil {
				return nil, fmt.Errorf("queue find by ids failed: %w", err)
			}

			queuesModels[j] = &models.Queue{
				ID:       queue.ID.Hex(),
				Messages: make([]*models.Message, 0, len(messages)),
			}
			for _, message := range messages {
				queuesModels[j].Messages = append(queuesModels[j].Messages, &models.Message{
					ID:   message.ID.Hex(),
					Body: message.Body,
					Read: message.Read,
				})
			}

			queuesMap[exchange.Key] = append(queuesMap[exchange.Key], make(chan models.Message, bufSize))
		}
		ex[i] = &models.Exchange{
			ID:     exchange.ID.Hex(),
			Key:    exchange.Key,
			Queues: queuesModels,
		}
	}

	return &Broker{
		exchanges: ex,
		queues:    queuesMap,

		exchangeService: exchangeService,
		queueService:    queueService,
		messagesService: messagesService,
	}, nil
}
