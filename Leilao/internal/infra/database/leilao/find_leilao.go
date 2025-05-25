package leilao

import (
	"context"
	"fmt"
	"time"

	"github.com/rafael0502/curso-go/Leilao/configuration/logger"
	"github.com/rafael0502/curso-go/Leilao/internal/entity/leilao_entity"
	"github.com/rafael0502/curso-go/Leilao/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (lr *LeilaoRepository) FindLeilaoById(ctx context.Context, id string) (*leilao_entity.Leilao, *internal_error.InternalError) {
	filter := bson.M{"_id": id}
	var leilaoEntityMongo LeilaoEntityMongo

	if err := lr.Collection.FindOne(ctx, filter).Decode(&leilaoEntityMongo); err != nil {
		logger.Error(fmt.Sprintf("Error trying to find leilão by id = %s", id), err)
		return nil, internal_error.NewInternalServerError(fmt.Sprintf("Error trying to find leilão by id = %s", id))
	}

	return &leilao_entity.Leilao{
		Id:          leilaoEntityMongo.Id,
		ProductName: leilaoEntityMongo.ProductName,
		Category:    leilaoEntityMongo.Category,
		Description: leilaoEntityMongo.Description,
		Condition:   leilaoEntityMongo.Condition,
		Status:      leilaoEntityMongo.Status,
		Timestamp:   time.Unix(leilaoEntityMongo.Timestamp, 0),
	}, nil
}

func (lr *LeilaoRepository) FindLeiloes(
	ctx context.Context,
	status leilao_entity.LeilaoStatus,
	category, productName string) ([]leilao_entity.Leilao, *internal_error.InternalError) {

	filter := bson.M{}

	if status != 0 {
		filter["status"] = status
	}

	if category != "" {
		filter["category"] = category
	}

	if productName != "" {
		filter["productName"] = primitive.Regex{
			Pattern: productName,
			Options: "i",
		}
	}

	cursor, err := lr.Collection.Find(ctx, filter)

	if err != nil {
		logger.Error("Error trying to find leiloes", err)
		return nil, internal_error.NewInternalServerError("Error trying to find leiloes")
	}

	defer cursor.Close(ctx)

	var leilaoEntityMongo []LeilaoEntityMongo

	if err := cursor.All(ctx, &leilaoEntityMongo); err != nil {
		logger.Error("Error trying to find leiloes", err)
		return nil, internal_error.NewInternalServerError("Error trying to find leiloes")
	}

	var leilaoEntity []leilao_entity.Leilao
	for _, leilaoMongo := range leilaoEntityMongo {
		leilaoEntity = append(leilaoEntity, leilao_entity.Leilao{
			Id:          leilaoMongo.Id,
			ProductName: leilaoMongo.ProductName,
			Category:    leilaoMongo.Category,
			Description: leilaoMongo.Description,
			Condition:   leilaoMongo.Condition,
			Status:      leilaoMongo.Status,
			Timestamp:   time.Unix(leilaoMongo.Timestamp, 0),
		})
	}

	return leilaoEntity, nil
}
