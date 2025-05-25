package bid_usecase

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/rafael0502/curso-go/Leilao/configuration/logger"
	"github.com/rafael0502/curso-go/Leilao/internal/entity/bid_entity"
	"github.com/rafael0502/curso-go/Leilao/internal/internal_error"
)

type BidInputDTO struct {
	UserId   string  `json:"user_id"`
	LeilaoId string  `json:"leilao_id"`
	Amount   float64 `json:"amount"`
}

type BidOutputDTO struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	LeilaoId  string    `json:"leilao_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp" time_format:"2006-01-02 15:04:05"` // Format the timestamp to match the desired format
}

type BidUseCase struct {
	BidRepository bid_entity.BidEntityRepository

	timer               *time.Timer
	maxBatchSize        int
	batchInsertInterval time.Duration

	bidChannel chan bid_entity.Bid
}

func NewBidUseCase(bidRepository bid_entity.BidEntityRepository) BidUseCaseInterface {
	maxSizeInterval := getMaxBatchSizeInterval()
	maxBatchSize := getMaxBatchSize()

	bidUseCase := &BidUseCase{
		BidRepository:       bidRepository,
		maxBatchSize:        maxBatchSize,
		batchInsertInterval: maxSizeInterval,
		timer:               time.NewTimer(maxSizeInterval),
		bidChannel:          make(chan bid_entity.Bid, maxBatchSize),
	}

	bidUseCase.triggerCreateRoutine(context.Background())

	return bidUseCase
}

var bidBatch []bid_entity.Bid

type BidUseCaseInterface interface {
	CreateBid(ctx context.Context, bidInputDTO BidInputDTO) *internal_error.InternalError
	FindBidByLeilaoId(ctx context.Context, leilaoId string) ([]BidOutputDTO, *internal_error.InternalError)
	FindWinningBidByLeilaoId(ctx context.Context, leilaoId string) (*BidOutputDTO, *internal_error.InternalError)
}

func (bu *BidUseCase) triggerCreateRoutine(ctx context.Context) {
	go func() {
		defer close(bu.bidChannel)
		for {
			select {
			case bidEntity, ok := <-bu.bidChannel:
				if !ok {
					if len(bidBatch) > 0 {
						if err := bu.BidRepository.CreateBid(ctx, bidBatch); err != nil {
							logger.Error("Error trying to process bid batch list", err)
						}
					}
					return
				}

				bidBatch = append(bidBatch, bidEntity)

				if len(bidBatch) >= bu.maxBatchSize {
					if err := bu.BidRepository.CreateBid(ctx, bidBatch); err != nil {
						logger.Error("Error trying to process bid batch list", err)
					}
					bidBatch = nil
					bu.timer.Reset(bu.batchInsertInterval)
				}
			case <-bu.timer.C:
				if err := bu.BidRepository.CreateBid(ctx, bidBatch); err != nil {
					logger.Error("Error trying to process bid batch list", err)
				}
				bidBatch = nil
				bu.timer.Reset(bu.batchInsertInterval)
			}
		}

	}()
}

func (bu *BidUseCase) CreateBid(ctx context.Context, bidInputDTO BidInputDTO) *internal_error.InternalError {

	bidEntity, err := bid_entity.CreateBid(bidInputDTO.UserId, bidInputDTO.LeilaoId, bidInputDTO.Amount)
	if err != nil {
		return err
	}

	bu.bidChannel <- *bidEntity

	return nil
}

func getMaxBatchSizeInterval() time.Duration {
	batchInsertInterval := os.Getenv("BATCH_INSERT_INTERVAL")
	duration, err := time.ParseDuration(batchInsertInterval)

	if err != nil {
		return 3 * time.Minute // Default value if parsing fails
	}

	return duration
}

func getMaxBatchSize() int {
	value, err := strconv.Atoi(os.Getenv("MAX_BATCH_SIZE"))

	if err != nil {
		return 5 // Default value if parsing fails
	}

	return value
}
