package service

import (
	"context"
	"github.com/RipperAcskt/broker/internal/repository/mongo/db/dto"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/RipperAcskt/broker/internal/entities/models"
)

type messagesRepository interface {
	Store(ctx context.Context, message models.Message) (primitive.ObjectID, error)
	FindByIDs(ctx context.Context, messageIDs []primitive.ObjectID) ([]*dto.Message, error)
	UpdateByID(ctx context.Context, message models.Message) error
}

type Messages struct {
	messagesRepository messagesRepository
}

func (m *Messages) Store(ctx context.Context, message models.Message) (primitive.ObjectID, error) {
	return m.messagesRepository.Store(ctx, message)
}

func (m *Messages) FindByIDs(ctx context.Context, messageIDs []primitive.ObjectID) ([]*dto.Message, error) {
	return m.messagesRepository.FindByIDs(ctx, messageIDs)
}

func (m *Messages) UpdateByID(ctx context.Context, message models.Message) error {
	return m.messagesRepository.UpdateByID(ctx, message)
}

func NewMessages(messagesRepository messagesRepository) *Messages {
	return &Messages{
		messagesRepository: messagesRepository,
	}
}
