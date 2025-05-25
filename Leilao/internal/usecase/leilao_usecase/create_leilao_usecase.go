package leilao_usecase

import (
	"context"
	"time"

	"github.com/rafael0502/curso-go/Leilao/internal/entity/bid_entity"
	"github.com/rafael0502/curso-go/Leilao/internal/entity/leilao_entity"
	"github.com/rafael0502/curso-go/Leilao/internal/internal_error"
	"github.com/rafael0502/curso-go/Leilao/internal/usecase/bid_usecase"
)

type LeilaoInputDTO struct {
	ProductName string           `json:"product_name" binding: "required, min=1"`
	Category    string           `json:"category" binding: "required, min=2"`
	Description string           `json:"description" binding: "required, min=10, max=200"`
	Condition   ProductCondition `json:"condition"`
}

type LeilaoOutoutDTO struct {
	Id          string           `json:"id"`
	ProductName string           `json:"product_name"`
	Category    string           `json:"category"`
	Description string           `json:"description"`
	Condition   ProductCondition `json:"condition"`
	Status      LeilaoStatus     `json:"status"`
	Timestamp   time.Time        `json:"timestamp" time_format:"2006-01-02 15:04:05"`
}

type WinningInfoOutputDTO struct {
	Leilao LeilaoOutoutDTO           `json:"leilao"`
	Bid    *bid_usecase.BidOutputDTO `json:"bid,omitempty"`
}

func NewLeilaoUseCase(leilaoRepositoryInterface leilao_entity.LeilaoRepositoryInterface,
	bidRepositoryInterface bid_entity.BidEntityRepository) LeilaoUseCaseInterface {
	return &LeilaoUseCase{
		leilaoRepositoryInterface: leilaoRepositoryInterface,
		bidRepositoryInterface:    bidRepositoryInterface,
	}
}

type LeilaoUseCaseInterface interface {
	CreateLeilao(ctx context.Context, leilaoInput LeilaoInputDTO) *internal_error.InternalError
	FindLeilaoById(ctx context.Context, id string) (*LeilaoOutoutDTO, *internal_error.InternalError)
	FindLeiloes(ctx context.Context, status LeilaoStatus, category, productName string) ([]LeilaoOutoutDTO, *internal_error.InternalError)
	FindWinningBidByLeilaoId(ctx context.Context, leilaoId string) (*WinningInfoOutputDTO, *internal_error.InternalError)
}

type ProductCondition int64
type LeilaoStatus int64

type LeilaoUseCase struct {
	leilaoRepositoryInterface leilao_entity.LeilaoRepositoryInterface
	bidRepositoryInterface    bid_entity.BidEntityRepository
}

func (lu *LeilaoUseCase) CreateLeilao(ctx context.Context, leilaoInput LeilaoInputDTO) *internal_error.InternalError {
	leilao, err := leilao_entity.CreateLeilao(leilaoInput.ProductName, leilaoInput.Category, leilaoInput.Description, leilao_entity.ProductCondition(leilaoInput.Condition))

	if err != nil {
		return err
	}

	if err := lu.leilaoRepositoryInterface.CreateLeilao(ctx, *leilao); err != nil {
		return err
	}

	return nil
}
