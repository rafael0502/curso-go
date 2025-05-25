package bid_entity

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rafael0502/curso-go/Leilao/internal/internal_error"
)

type Bid struct {
	Id        string
	UserId    string
	LeilaoId  string
	Amount    float64
	Timestamp time.Time
}

func CreateBid(userId, leilaoId string, amount float64) (*Bid, *internal_error.InternalError) {
	bid := &Bid{
		Id:        uuid.New().String(),
		UserId:    userId,
		LeilaoId:  leilaoId,
		Amount:    amount,
		Timestamp: time.Now(),
	}

	if err := bid.Validate(); err != nil {
		return nil, err
	}

	return bid, nil
}

func (b *Bid) Validate() *internal_error.InternalError {
	if err := uuid.Validate(b.UserId); err != nil {
		return internal_error.NewBadRequestError("Invalid UserId")
	} else if err := uuid.Validate(b.LeilaoId); err != nil {
		return internal_error.NewBadRequestError("Invalid LeilaoId")
	} else if b.Amount <= 0 {
		return internal_error.NewBadRequestError("Amount must be greater than zero")
	}

	return nil
}

type BidEntityRepository interface {
	CreateBid(ctx context.Context, bidEntities []Bid) *internal_error.InternalError
	FindBidByLeilaoId(ctx context.Context, leilaoId string) ([]Bid, *internal_error.InternalError)
	FindWinningBidByLeilaoId(ctx context.Context, leilaoId string) (*Bid, *internal_error.InternalError)
}
