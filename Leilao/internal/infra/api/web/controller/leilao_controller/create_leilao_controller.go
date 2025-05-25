package leilao_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rafael0502/curso-go/Leilao/configuration/rest_err"
	"github.com/rafael0502/curso-go/Leilao/internal/infra/api/web/validation"
	"github.com/rafael0502/curso-go/Leilao/internal/usecase/leilao_usecase"
	"golang.org/x/net/context"
)

type LeilaoController struct {
	leilaoUseCase leilao_usecase.LeilaoUseCaseInterface
}

func NewLeilaoController(leilaoUseCase leilao_usecase.LeilaoUseCaseInterface) *LeilaoController {
	return &LeilaoController{
		leilaoUseCase: leilaoUseCase,
	}
}

func (l *LeilaoController) CreateLeilao(c *gin.Context) {
	var leilaoInputDTO leilao_usecase.LeilaoInputDTO

	if err := c.ShouldBindJSON(&leilaoInputDTO); err != nil {
		restErr := validation.ValidationErr(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	err := l.leilaoUseCase.CreateLeilao(context.Background(), leilaoInputDTO)

	if err != nil {
		restErr := rest_err.ConvertError(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	c.Status(http.StatusCreated)
}
