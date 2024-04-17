package service

import (
	"context"
	"github.com/RipperAcskt/broker/internal/repository/mongo/db/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/RipperAcskt/broker/internal/entities/models"
)

type exchangeRepository interface {
	Store(ctx context.Context, exchange models.Exchange) (primitive.ObjectID, error)
	All(ctx context.Context) ([]*dto.Exchange, error)
	UpdateByID(ctx context.Context, exchange models.Exchange) error
}

type Exchange struct {
	exchangeRepository exchangeRepository
}

func (e *Exchange) All(ctx context.Context) ([]*dto.Exchange, error) {
	return e.exchangeRepository.All(ctx)
}

func (e *Exchange) Store(ctx context.Context, exchange models.Exchange) (primitive.ObjectID, error) {
	return e.exchangeRepository.Store(ctx, exchange)
}

func (e *Exchange) UpdateByID(ctx context.Context, exchange models.Exchange) error {
	return e.exchangeRepository.UpdateByID(ctx, exchange)
}

func NewExchange(exchangeRepository exchangeRepository) (*Exchange, error) {
	return &Exchange{
		exchangeRepository: exchangeRepository,
	}, nil
}
