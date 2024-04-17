package dto

import (
	"fmt"
	"github.com/RipperAcskt/broker/internal/entities/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExchangeCreate struct {
	Queues []primitive.ObjectID `bson:"queues"`
	Key    string               `bson:"key"`
}

type Exchange struct {
	ID     primitive.ObjectID   `bson:"_id"`
	Queues []primitive.ObjectID `bson:"queues"`
	Key    string               `bson:"key"`
}

func ExchangeCreateModelToDTO(exchange models.Exchange) (*ExchangeCreate, error) {
	queues := make([]primitive.ObjectID, 0, len(exchange.Queues))
	for _, queue := range exchange.Queues {
		id, err := primitive.ObjectIDFromHex(queue.ID)
		if err != nil {
			return nil, fmt.Errorf("objectID from hex failed: %w", err)
		}

		queues = append(queues, id)
	}

	return &ExchangeCreate{
		Queues: queues,
		Key:    exchange.Key,
	}, nil
}
