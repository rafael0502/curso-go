package bid_usecase

import (
	"context"

	"github.com/rafael0502/curso-go/Leilao/internal/internal_error"
)

func (bu *BidUseCase) FindBidByLeilaoId(ctx context.Context, leilaoId string) ([]BidOutputDTO, *internal_error.InternalError) {
	bidList, err := bu.BidRepository.FindBidByLeilaoId(ctx, leilaoId)
	if err != nil {
		return nil, err
	}

	var bidOutputList []BidOutputDTO
	for _, bid := range bidList {
		bidOutputList = append(bidOutputList, BidOutputDTO{
			Id:        bid.Id,
			UserId:    bid.UserId,
			LeilaoId:  bid.LeilaoId,
			Amount:    bid.Amount,
			Timestamp: bid.Timestamp,
		})
	}
	return bidOutputList, nil
}

func (bu *BidUseCase) FindWinningBidByLeilaoId(ctx context.Context, leilaoId string) (*BidOutputDTO, *internal_error.InternalError) {
	bidEntity, err := bu.BidRepository.FindWinningBidByLeilaoId(ctx, leilaoId)
	if err != nil {
		return nil, err
	}

	return &BidOutputDTO{
		Id:        bidEntity.Id,
		UserId:    bidEntity.UserId,
		LeilaoId:  bidEntity.LeilaoId,
		Amount:    bidEntity.Amount,
		Timestamp: bidEntity.Timestamp,
	}, nil
}
