package dto

import (
	"github.com/RipperAcskt/broker/internal/entities/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageCreate struct {
	Body any  `bson:"body"`
	Read bool `bson:"read"`
}

type Message struct {
	ID   primitive.ObjectID `bson:"_id"`
	Body any                `bson:"body"`
	Read bool               `bson:"read"`
}

func MessageCreateModelToDTO(message models.Message) *MessageCreate {
	return &MessageCreate{
		Body: message.Body,
		Read: message.Read,
	}
}
