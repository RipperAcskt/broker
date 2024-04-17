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
	exchangeCollection = "exchange"
)

type Exchange struct {
	mongo *mongo.Collection
}

func (m *Exchange) Store(ctx context.Context, exchange models.Exchange) (primitive.ObjectID, error) {
	exchangeDTO, err := dto.ExchangeCreateModelToDTO(exchange)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("exchange model to dto failed: %w", err)
	}

	result, err := m.mongo.InsertOne(ctx, exchangeDTO)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("insert one failed: %w", err)
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("inserted id to object id failed")
	}
	return id, nil
}

func (m *Exchange) UpdateByID(ctx context.Context, exchange models.Exchange) error {
	exchangeDTO, err := dto.ExchangeCreateModelToDTO(exchange)
	if err != nil {
		return fmt.Errorf("exchange model to dto failed: %w", err)
	}

	id, err := primitive.ObjectIDFromHex(exchange.ID)
	if err != nil {
		return fmt.Errorf("object id from hex failed: %w", err)
	}

	_, err = m.mongo.UpdateByID(ctx, id, bson.M{
		"$set": exchangeDTO,
	})
	if err != nil {
		return fmt.Errorf("update by id failed: %w", err)
	}

	return nil
}

func (m *Exchange) All(ctx context.Context) ([]*dto.Exchange, error) {
	cursor, err := m.mongo.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("find failed: %w", err)
	}

	exchanges := make([]*dto.Exchange, 0, cursor.RemainingBatchLength())
	if err = cursor.All(ctx, &exchanges); err != nil {
		return nil, fmt.Errorf("cursor all failed: %w", err)
	}

	return exchanges, nil
}

func NewExchange(mongo *mongo.Database) *Exchange {
	return &Exchange{
		mongo: mongo.Collection(exchangeCollection),
	}
}
