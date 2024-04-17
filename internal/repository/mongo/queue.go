package mongo

import (
	"context"
	"github.com/RipperAcskt/broker/internal/repository/mongo/db"
	"github.com/RipperAcskt/broker/internal/repository/mongo/db/dto"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/RipperAcskt/broker/internal/entities/models"
)

type Queue struct {
	db *db.Queue
}

func (q *Queue) Store(ctx context.Context, queue models.Queue) (primitive.ObjectID, error) {
	return q.db.Store(ctx, queue)
}

func (q *Queue) FindByIDs(ctx context.Context, queueIDs []primitive.ObjectID) ([]*dto.Queue, error) {
	return q.db.FindByIDs(ctx, queueIDs)
}

func (q *Queue) UpdateByID(ctx context.Context, queue models.Queue) error {
	return q.db.UpdateByID(ctx, queue)
}

func NewQueue(mongo *mongo.Database) *Queue {
	return &Queue{
		db: db.NewQueue(mongo),
	}
}
