package leilao_usecase

import (
	"context"

	"github.com/rafael0502/curso-go/Leilao/internal/entity/leilao_entity"
	"github.com/rafael0502/curso-go/Leilao/internal/internal_error"
	"github.com/rafael0502/curso-go/Leilao/internal/usecase/bid_usecase"
)

func (lu *LeilaoUseCase) FindLeilaoById(ctx context.Context, id string) (*LeilaoOutoutDTO, *internal_error.InternalError) {
	leilao, err := lu.leilaoRepositoryInterface.FindLeilaoById(ctx, id)

	if err != nil {
		return nil, err
	}

	return &LeilaoOutoutDTO{
		Id:          leilao.Id,
		ProductName: leilao.ProductName,
		Category:    leilao.Category,
		Description: leilao.Description,
		Condition:   ProductCondition(leilao.Condition),
		Status:      LeilaoStatus(leilao.Status),
		Timestamp:   leilao.Timestamp,
	}, nil
}

func (lu *LeilaoUseCase) FindLeiloes(ctx context.Context, status LeilaoStatus, category, productName string) ([]LeilaoOutoutDTO, *internal_error.InternalError) {
	leilaoEntities, err := lu.leilaoRepositoryInterface.FindLeiloes(ctx, leilao_entity.LeilaoStatus(status), category, productName)

	if err != nil {
		return nil, err
	}

	var leilaoOutouts []LeilaoOutoutDTO

	for _, value := range leilaoEntities {
		leilaoOutouts = append(leilaoOutouts, LeilaoOutoutDTO{
			Id:          value.Id,
			ProductName: value.ProductName,
			Category:    value.Category,
			Description: value.Description,
			Condition:   ProductCondition(value.Condition),
			Status:      LeilaoStatus(value.Status),
			Timestamp:   value.Timestamp,
		})
	}

	return leilaoOutouts, nil
}

func (lu *LeilaoUseCase) FindWinningBidByLeilaoId(ctx context.Context, leilaoId string) (*WinningInfoOutputDTO, *internal_error.InternalError) {
	leilao, err := lu.leilaoRepositoryInterface.FindLeilaoById(ctx, leilaoId)

	if err != nil {
		return nil, err
	}

	leilaoOutputDTO := LeilaoOutoutDTO{
		Id:          leilao.Id,
		ProductName: leilao.ProductName,
		Category:    leilao.Category,
		Description: leilao.Description,
		Condition:   ProductCondition(leilao.Condition),
		Status:      LeilaoStatus(leilao.Status),
		Timestamp:   leilao.Timestamp,
	}

	bidWinning, err := lu.bidRepositoryInterface.FindWinningBidByLeilaoId(ctx, leilaoId)

	if err != nil {
		return &WinningInfoOutputDTO{
			Leilao: leilaoOutputDTO,
			Bid:    nil,
		}, nil
	}

	bidOutputDTO := &bid_usecase.BidOutputDTO{
		Id:        bidWinning.Id,
		UserId:    bidWinning.UserId,
		LeilaoId:  bidWinning.LeilaoId,
		Amount:    bidWinning.Amount,
		Timestamp: bidWinning.Timestamp,
	}

	return &WinningInfoOutputDTO{
		Leilao: leilaoOutputDTO,
		Bid:    bidOutputDTO,
	}, nil
}
