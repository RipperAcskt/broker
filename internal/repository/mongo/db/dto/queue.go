package dto

import (
	"fmt"
	"github.com/RipperAcskt/broker/internal/entities/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QueueCreate struct {
	Messages []primitive.ObjectID `bson:"messages"`
}

type Queue struct {
	ID       primitive.ObjectID   `bson:"_id"`
	Messages []primitive.ObjectID `bson:"messages"`
}

func QueueCreateModelToDTO(queue models.Queue) (*QueueCreate, error) {
	messages := make([]primitive.ObjectID, 0, len(queue.Messages))
	for _, message := range queue.Messages {
		id, err := primitive.ObjectIDFromHex(message.ID)
		if err != nil {
			return nil, fmt.Errorf("objectID from hex failed: %w", err)
		}

		messages = append(messages, id)
	}

	return &QueueCreate{
		Messages: messages,
	}, nil
}
