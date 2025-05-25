package bid_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rafael0502/curso-go/Leilao/configuration/rest_err"
)

func (u *BidController) FindBidByLeilaoId(c *gin.Context) {
	leilaoId := c.Param("leilaoId")

	if err := uuid.Validate(leilaoId); err != nil {
		errRest := rest_err.NewBadRequestError("Invalid fields", rest_err.Causes{
			Field:   "leilaoId",
			Message: "Invalid UUID value",
		})

		c.JSON(errRest.Code, errRest)
		return
	}

	leilaoData, err := u.bidUseCase.FindBidByLeilaoId(c.Request.Context(), leilaoId)

	if err != nil {
		errRest := rest_err.ConvertError(err)

		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, leilaoData)
}
