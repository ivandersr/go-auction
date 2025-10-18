package bid_controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ivandersr/go-auction/config/rest_err"
)

func (b *BidController) FindBidsByAuctionId(c *gin.Context) {
	auctionId := c.Param("auctionId")

	if err := uuid.Validate(auctionId); err != nil {
		restErr := rest_err.NewBadRequestError("invalid fields", rest_err.Cause{
			Field:   "auctionId",
			Message: "invalid UUID value",
		})

		c.JSON(restErr.Code, restErr)
		return
	}

	bidsData, err := b.bidUseCase.FindBidsByAuctionId(context.Background(), auctionId)

	if err != nil {
		restErr := rest_err.ConvertError(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	c.JSON(http.StatusOK, bidsData)
}
