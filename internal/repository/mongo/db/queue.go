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
	queueCollection = "queue"
)

type Queue struct {
	mongo *mongo.Collection
}

func (q *Queue) Store(ctx context.Context, queue models.Queue) (primitive.ObjectID, error) {
	queueDTO, err := dto.QueueCreateModelToDTO(queue)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("queue model to dto failed: %w", err)
	}

	result, err := q.mongo.InsertOne(ctx, queueDTO)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("insert one failed: %w", err)
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("inserted id to object id failed")
	}
	return id, nil
}

func (q *Queue) FindByIDs(ctx context.Context, queueIDs []primitive.ObjectID) ([]*dto.Queue, error) {
	cursor, err := q.mongo.Find(ctx, bson.M{
		"_id": bson.M{
			"$in": queueIDs,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("find failed: %w", err)
	}

	queues := make([]*dto.Queue, 0, cursor.RemainingBatchLength())
	if err := cursor.All(ctx, &queues); err != nil {
		return nil, fmt.Errorf("all failed: %w", err)
	}

	return queues, nil
}

func (q *Queue) UpdateByID(ctx context.Context, queue models.Queue) error {
	queueDTO, err := dto.QueueCreateModelToDTO(queue)
	if err != nil {
		return fmt.Errorf("queue model to dto failed: %w", err)
	}

	id, err := primitive.ObjectIDFromHex(queue.ID)
	if err != nil {
		return fmt.Errorf("object id from hex failed: %w", err)
	}

	_, err = q.mongo.UpdateByID(ctx, id, bson.M{
		"$set": queueDTO,
	})
	if err != nil {
		return fmt.Errorf("update by id failed: %w", err)
	}

	return nil
}

func NewQueue(mongo *mongo.Database) *Queue {
	return &Queue{
		mongo: mongo.Collection(queueCollection),
	}
}
