package mongo

import (
	"context"
	"github.com/RipperAcskt/broker/internal/repository/mongo/db"
	"github.com/RipperAcskt/broker/internal/repository/mongo/db/dto"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/RipperAcskt/broker/internal/entities/models"
)

type Message struct {
	db *db.Message
}

func (m *Message) Store(ctx context.Context, message models.Message) (primitive.ObjectID, error) {
	return m.db.Store(ctx, message)
}

func (m *Message) FindByIDs(ctx context.Context, messageIDs []primitive.ObjectID) ([]*dto.Message, error) {
	return m.db.FindByIDs(ctx, messageIDs)
}

func (m *Message) UpdateByID(ctx context.Context, message models.Message) error {
	return m.db.UpdateByID(ctx, message)
}

func NewMessage(mongo *mongo.Database) *Message {
	return &Message{
		db: db.NewMessage(mongo),
	}
}
