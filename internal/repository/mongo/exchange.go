package mongo

import (
	"context"
	"github.com/RipperAcskt/broker/internal/repository/mongo/db"
	"github.com/RipperAcskt/broker/internal/repository/mongo/db/dto"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/RipperAcskt/broker/internal/entities/models"
)

type Exchange struct {
	db *db.Exchange
}

func (m *Exchange) Store(ctx context.Context, exchange models.Exchange) (primitive.ObjectID, error) {
	return m.db.Store(ctx, exchange)
}

func (m *Exchange) UpdateByID(ctx context.Context, exchange models.Exchange) error {
	return m.db.UpdateByID(ctx, exchange)
}

func (m *Exchange) All(ctx context.Context) ([]*dto.Exchange, error) {
	return m.db.All(ctx)
}

func NewExchange(mongo *mongo.Database) *Exchange {
	return &Exchange{
		db: db.NewExchange(mongo),
	}
}
