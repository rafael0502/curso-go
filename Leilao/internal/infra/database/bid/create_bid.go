package bid

import (
	"context"
	"sync"

	"github.com/rafael0502/curso-go/Leilao/configuration/logger"
	"github.com/rafael0502/curso-go/Leilao/internal/entity/bid_entity"
	"github.com/rafael0502/curso-go/Leilao/internal/entity/leilao_entity"
	"github.com/rafael0502/curso-go/Leilao/internal/infra/database/leilao"
	"github.com/rafael0502/curso-go/Leilao/internal/internal_error"
	"go.mongodb.org/mongo-driver/mongo"
)

type BidEntityMongo struct {
	Id        string  `bson:"_id"`
	UserId    string  `bson:"user_id"`
	LeilaoId  string  `bson:"leilao_id"`
	Amount    float64 `bson:"amount"`
	Timestamp int64   `bson:"timestamp"`
}

type BidRepository struct {
	Collection       *mongo.Collection
	LeilaoRepository *leilao.LeilaoRepository
}

func NewBidRepository(database *mongo.Database, leilaoRepository *leilao.LeilaoRepository) *BidRepository {
	return &BidRepository{
		Collection:       database.Collection("bids"),
		LeilaoRepository: leilaoRepository,
	}
}

func (bd *BidRepository) CreateBid(ctx context.Context, bidEntities []bid_entity.Bid) *internal_error.InternalError {
	var wg sync.WaitGroup

	for _, bid := range bidEntities {
		wg.Add(1)
		go func(bidValue bid_entity.Bid) {
			defer wg.Done()

			leilaoEntity, err := bd.LeilaoRepository.FindLeilaoById(ctx, bidValue.LeilaoId)

			if err != nil {
				logger.Error("Error finding leilao by id", err)
				return
			}

			if leilaoEntity.Status != leilao_entity.Active {
				return
			}

			bidEntityMongo := BidEntityMongo{
				Id:        bidValue.Id,
				UserId:    bidValue.UserId,
				LeilaoId:  bidValue.LeilaoId,
				Amount:    bidValue.Amount,
				Timestamp: bidValue.Timestamp.Unix(),
			}

			if _, err := bd.Collection.InsertOne(ctx, bidEntityMongo); err != nil {
				logger.Error("Error inserting bid into MongoDB", err)
				return
			}
		}(bid)
	}

	wg.Wait()
	return nil
}
