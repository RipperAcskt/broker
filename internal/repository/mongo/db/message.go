package db

import (
	"context"
	"fmt"
	"github.com/RipperAcskt/broker/internal/repository/mongo/db/dto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/RipperAcskt/broker/internal/entities/models"
)

const (
	messageCollection = "message"
)

type Message struct {
	mongo *mongo.Collection
}

func (m *Message) Store(ctx context.Context, message models.Message) (primitive.ObjectID, error) {
	messageDTO := dto.MessageCreateModelToDTO(message)

	result, err := m.mongo.InsertOne(ctx, messageDTO)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("insert one failed: %w", err)
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("inserted id to object id failed")
	}
	return id, nil
}

func (m *Message) FindByIDs(ctx context.Context, messageIDs []primitive.ObjectID) ([]*dto.Message, error) {
	cursor, err := m.mongo.Find(ctx, bson.M{
		"_id": bson.M{
			"$in": messageIDs,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("find failed: %w", err)
	}

	messages := make([]*dto.Message, 0, cursor.RemainingBatchLength())
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, fmt.Errorf("all failed: %w", err)
	}

	return messages, nil
}

func (m *Message) UpdateByID(ctx context.Context, message models.Message) error {
	messageDTO := dto.MessageCreateModelToDTO(message)

	id, err := primitive.ObjectIDFromHex(message.ID)
	if err != nil {
		return fmt.Errorf("object id from hex failed: %w", err)
	}

	_, err = m.mongo.UpdateByID(ctx, id, bson.M{
		"$set": messageDTO,
	})
	if err != nil {
		return fmt.Errorf("update by id failed: %w", err)
	}

	return nil
}

func NewMessage(mongo *mongo.Database) *Message {
	return &Message{
		mongo: mongo.Collection(messageCollection),
	}
}
