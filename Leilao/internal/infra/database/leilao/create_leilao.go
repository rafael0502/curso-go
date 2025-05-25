package leilao

import (
	"context"
	"os"
	"time"

	"github.com/rafael0502/curso-go/Leilao/configuration/logger"
	"github.com/rafael0502/curso-go/Leilao/internal/entity/leilao_entity"
	"github.com/rafael0502/curso-go/Leilao/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LeilaoEntityMongo struct {
	Id          string                         `bson:"_id"`
	ProductName string                         `bson:"product_name"`
	Category    string                         `bson:"category"`
	Description string                         `bson:"description"`
	Condition   leilao_entity.ProductCondition `bson:"condition"`
	Status      leilao_entity.LeilaoStatus     `bson:"status"`
	Timestamp   int64                          `bson:"timestamp"`
}

type LeilaoRepository struct {
	Collection *mongo.Collection
}

func NewLeilaoRepository(database *mongo.Database) *LeilaoRepository {
	return &LeilaoRepository{
		Collection: database.Collection("leiloes"),
	}
}

func (lr *LeilaoRepository) CreateLeilao(ctx context.Context, leilaoEntity leilao_entity.Leilao) *internal_error.InternalError {
	leilaoEntityMongo := &LeilaoEntityMongo{
		Id:          leilaoEntity.Id,
		ProductName: leilaoEntity.ProductName,
		Category:    leilaoEntity.Category,
		Description: leilaoEntity.Description,
		Condition:   leilaoEntity.Condition,
		Status:      leilaoEntity.Status,
		Timestamp:   leilaoEntity.Timestamp.Unix(),
	}

	_, err := lr.Collection.InsertOne(ctx, leilaoEntityMongo)
	if err != nil {
		return internal_error.NewInternalServerError("Error creating leilao")
	}

	go func() {
		select {
		case <-time.After(getLeilaoInterval()):
			update := bson.M{"$set": bson.M{"status": leilao_entity.Completed}}
			filter := bson.M{"_id": leilaoEntityMongo.Id}

			_, err := lr.Collection.UpdateOne(ctx, filter, update)

			if err != nil {
				logger.Error("Error trying to update leilao status to completed", err)
				return
			}
		}
	}()

	return nil
}

func getLeilaoInterval() time.Duration {
	leilaoInterval := os.Getenv("LEILAO_INTERVAL")
	duration, err := time.ParseDuration(leilaoInterval)

	if err != nil {
		return time.Minute * 5
	}

	return duration
}
