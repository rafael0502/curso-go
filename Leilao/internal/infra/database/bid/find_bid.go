package bid

import (
	"context"
	"fmt"
	"time"

	"github.com/rafael0502/curso-go/Leilao/configuration/logger"
	"github.com/rafael0502/curso-go/Leilao/internal/entity/bid_entity"
	"github.com/rafael0502/curso-go/Leilao/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (bd *BidRepository) FindBidByLeilaoId(ctx context.Context, leilaoId string) ([]bid_entity.Bid, *internal_error.InternalError) {
	filter := bson.M{"leilao_id": leilaoId}
	cursor, err := bd.Collection.Find(ctx, filter)

	if err != nil {
		logger.Error(fmt.Sprintf("Error finding bids by leilao id %s", leilaoId), err)
		return nil, internal_error.NewInternalServerError(fmt.Sprintf("Error finding bids by leilao id %s", leilaoId))
	}

	var bidEntitiesMongo []BidEntityMongo

	if err := cursor.All(ctx, &bidEntitiesMongo); err != nil {
		logger.Error(fmt.Sprintf("Error finding bids for leilao id %s", leilaoId), err)
		return nil, internal_error.NewInternalServerError(fmt.Sprintf("Error finding bids for leilao id %s", leilaoId))
	}

	var bidEntities []bid_entity.Bid
	for _, bidEntityMongo := range bidEntitiesMongo {
		bidEntities = append(bidEntities, bid_entity.Bid{
			Id:        bidEntityMongo.Id,
			UserId:    bidEntityMongo.UserId,
			LeilaoId:  bidEntityMongo.LeilaoId,
			Amount:    bidEntityMongo.Amount,
			Timestamp: time.Unix(bidEntityMongo.Timestamp, 0),
		})
	}

	return bidEntities, nil
}

func (bd *BidRepository) FindWinningBidByLeilaoId(ctx context.Context, leilaoId string) (*bid_entity.Bid, *internal_error.InternalError) {
	filter := bson.M{"leilaoId": leilaoId}
	opts := options.FindOne().SetSort(bson.D{{"amount", -1}})

	var bidEntityMongo BidEntityMongo
	err := bd.Collection.FindOne(ctx, filter, opts).Decode(&bidEntityMongo)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, internal_error.NewNotFoundError(fmt.Sprintf("No winning bid found for leilao id %s", leilaoId))
		}
		logger.Error(fmt.Sprintf("Error finding winning bid for leilao id %s", leilaoId), err)
		return nil, internal_error.NewInternalServerError(fmt.Sprintf("Error finding winning bid for leilao id %s", leilaoId))
	}

	return &bid_entity.Bid{
		Id:        bidEntityMongo.Id,
		UserId:    bidEntityMongo.UserId,
		LeilaoId:  bidEntityMongo.LeilaoId,
		Amount:    bidEntityMongo.Amount,
		Timestamp: time.Unix(bidEntityMongo.Timestamp, 0),
	}, nil
}
