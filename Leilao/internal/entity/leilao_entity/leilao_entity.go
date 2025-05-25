package leilao_entity

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rafael0502/curso-go/Leilao/internal/internal_error"
)

func CreateLeilao(productName, category, description string, condition ProductCondition) (*Leilao, *internal_error.InternalError) {
	leilao := &Leilao{
		Id:          uuid.New().String(),
		ProductName: productName,
		Category:    category,
		Description: description,
		Condition:   condition,
		Status:      Active,
		Timestamp:   time.Now(),
	}

	if err := leilao.Validate(); err != nil {
		return nil, err
	}

	return leilao, nil
}

func (l *Leilao) Validate() *internal_error.InternalError {
	if len(l.ProductName) <= 1 ||
		len(l.Category) <= 2 ||
		len(l.Description) <= 10 && (l.Condition != New &&
			l.Condition != Refurbished &&
			l.Condition != Used) {
		return internal_error.NewBadRequestError("Invalid leilao object")
	}

	return nil
}

type Leilao struct {
	Id          string
	ProductName string
	Category    string
	Description string
	Condition   ProductCondition
	Status      LeilaoStatus
	Timestamp   time.Time
}

type ProductCondition int
type LeilaoStatus int

const (
	Active    LeilaoStatus = iota
	Completed              = 1
)
const (
	New         ProductCondition = iota
	Used                         = 1
	Refurbished                  = 2
)

type (
	LeilaoRepositoryInterface interface {
		CreateLeilao(ctx context.Context, leilaoEntity Leilao) *internal_error.InternalError
		FindLeilaoById(ctx context.Context, id string) (*Leilao, *internal_error.InternalError)
		FindLeiloes(ctx context.Context, status LeilaoStatus, category, productName string) ([]Leilao, *internal_error.InternalError)
	}
)
