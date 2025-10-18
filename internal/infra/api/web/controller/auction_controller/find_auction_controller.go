package auction_controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ivandersr/go-auction/config/rest_err"
	"github.com/ivandersr/go-auction/internal/usecase/auction_usecase"
)

func (a *AuctionController) FindAuctionById(c *gin.Context) {
	auctionId := c.Param("auctionId")

	if err := uuid.Validate(auctionId); err != nil {
		restErr := rest_err.NewBadRequestError("invalid fields", rest_err.Cause{
			Field:   "auctionId",
			Message: "invalid UUID value",
		})

		c.JSON(restErr.Code, restErr)
		return
	}

	auctionData, err := a.auctionUseCase.FindAuctionById(context.Background(), auctionId)

	if err != nil {
		restErr := rest_err.ConvertError(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	c.JSON(http.StatusOK, auctionData)
}

func (a *AuctionController) FindAuctions(c *gin.Context) {
	status := c.Query("status")
	category := c.Query("category")
	productName := c.Query("product_name")

	statusNumber, errConv := strconv.Atoi(status)
	if errConv != nil {
		errRest := rest_err.NewBadRequestError("invalid status value")
		c.JSON(errRest.Code, errRest)
		return
	}

	auctions, err := a.auctionUseCase.FindAuctions(context.Background(), auction_usecase.AuctionStatus(statusNumber), category, productName)

	if err != nil {
		restErr := rest_err.ConvertError(err)
		c.JSON(restErr.Code, restErr)
		return
	}

	c.JSON(http.StatusOK, auctions)
}

func (a *AuctionController) FindWinningBidByAuctionId(c *gin.Context) {
	auctionId := c.Param("auctionId")

	if err := uuid.Validate(auctionId); err != nil {
		restErr := rest_err.NewBadRequestError("invalid fields", rest_err.Cause{
			Field:   "auctionId",
			Message: "invalid UUID value",
		})

		c.JSON(restErr.Code, restErr)
		return
	}

	auctionData, err := a.auctionUseCase.FindWinningBidByAuctionId(context.Background(), auctionId)

	if err != nil {
		restErr := rest_err.ConvertError(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	c.JSON(http.StatusOK, auctionData)
}
