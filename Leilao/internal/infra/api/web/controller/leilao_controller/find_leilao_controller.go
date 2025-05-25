package leilao_controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rafael0502/curso-go/Leilao/configuration/rest_err"
	"github.com/rafael0502/curso-go/Leilao/internal/usecase/leilao_usecase"
	"golang.org/x/net/context"
)

func (l *LeilaoController) FindLeilaoById(c *gin.Context) {
	leilaoId := c.Param("leilaoId")

	if err := uuid.Validate(leilaoId); err != nil {
		errRest := rest_err.NewBadRequestError("Invalid fields", rest_err.Causes{
			Field:   "leilaoId",
			Message: "Invalid UUID value",
		})

		c.JSON(errRest.Code, errRest)
		return
	}

	leilaoData, err := l.leilaoUseCase.FindLeilaoById(context.Background(), leilaoId)

	if err != nil {
		restErr := rest_err.ConvertError(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	c.JSON(http.StatusOK, leilaoData)
}

func (l *LeilaoController) FindLeiloes(c *gin.Context) {
	status := c.Query("status")
	category := c.Query("category")
	productName := c.Query("productName")

	statusNumber, errConv := strconv.Atoi(status)

	if errConv != nil {
		errRest := rest_err.NewBadRequestError("Error trying to validate status leilao")
		c.JSON(errRest.Code, errRest)
		return
	}

	leiloes, err := l.leilaoUseCase.FindLeiloes(context.Background(), leilao_usecase.LeilaoStatus(statusNumber), category, productName)

	if err != nil {
		restErr := rest_err.ConvertError(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	c.JSON(http.StatusOK, leiloes)
}

func (l *LeilaoController) FindWinningBidByLeilaoid(c *gin.Context) {
	leilaoId := c.Param("leilaoId")

	if err := uuid.Validate(leilaoId); err != nil {
		errRest := rest_err.NewBadRequestError("Invalid fields", rest_err.Causes{
			Field:   "leilaoId",
			Message: "Invalid UUID value",
		})

		c.JSON(errRest.Code, errRest)
		return
	}

	leilaoData, err := l.leilaoUseCase.FindWinningBidByLeilaoId(context.Background(), leilaoId)

	if err != nil {
		restErr := rest_err.ConvertError(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	c.JSON(http.StatusOK, leilaoData)
}
