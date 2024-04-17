package service

import (
	"context"
	"github.com/RipperAcskt/broker/internal/repository/mongo/db/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/RipperAcskt/broker/internal/entities/models"
)

type queueRepository interface {
	Store(ctx context.Context, queue models.Queue) (primitive.ObjectID, error)
	FindByIDs(ctx context.Context, queueIDs []primitive.ObjectID) ([]*dto.Queue, error)
	UpdateByID(ctx context.Context, queue models.Queue) error
}

type Queue struct {
	queueRepository queueRepository
}

func (q *Queue) Store(ctx context.Context, queue models.Queue) (primitive.ObjectID, error) {
	return q.queueRepository.Store(ctx, queue)
}

func (q *Queue) FindByIDs(ctx context.Context, queueIDs []primitive.ObjectID) ([]*dto.Queue, error) {
	return q.queueRepository.FindByIDs(ctx, queueIDs)
}

func (q *Queue) UpdateByID(ctx context.Context, queue models.Queue) error {
	return q.queueRepository.UpdateByID(ctx, queue)
}

func NewQueue(queueRepository queueRepository) (*Queue, error) {
	return &Queue{
		queueRepository: queueRepository,
	}, nil
}
