package auction_usecase

import (
	"context"
	"time"

	"github.com/ivandersr/go-auction/internal/entity/auction"
	"github.com/ivandersr/go-auction/internal/entity/bid"
	"github.com/ivandersr/go-auction/internal/internal_errors"
	"github.com/ivandersr/go-auction/internal/usecase/bid_usecase"
)

type AuctionInputDTO struct {
	ProductName string           `json:"product_name" binding:"required,min=1"`
	Category    string           `json:"category" binding:"required,min=2"`
	Description string           `json:"description" binding:"required,min=10,max=200"`
	Condition   ProductCondition `json:"condition"`
}

type AuctionOutputDTO struct {
	Id          string           `json:"id"`
	ProductName string           `json:"product_name"`
	Category    string           `json:"category"`
	Description string           `json:"description"`
	Condition   ProductCondition `json:"condition"`
	Status      AuctionStatus    `json:"status"`
	CreatedAt   time.Time        `json:"created_at" time_format:"2006-01-02 15:04:05"`
}

type WinningInfoOutputDTO struct {
	Auction AuctionOutputDTO          `json:"auction"`
	Bid     *bid_usecase.BidOutputDTO `json:"bid,omitempty"`
}

type ProductCondition int64
type AuctionStatus int64

type AuctionUseCase struct {
	auctionRepositoryInterface auction.AuctionRepositoryInterface
	bidRepositoryInterface     bid.BidRepositoryInterface
}

type AuctionUseCaseInterface interface {
	CreateAuction(ctx context.Context, input *AuctionInputDTO) *internal_errors.InternalError
	FindWinningBidByAuctionId(ctx context.Context, id string) (*WinningInfoOutputDTO, *internal_errors.InternalError)
	FindAuctions(ctx context.Context, status AuctionStatus, category, productName string) ([]AuctionOutputDTO, *internal_errors.InternalError)
	FindAuctionById(ctx context.Context, id string) (*AuctionOutputDTO, *internal_errors.InternalError)
}

func NewAuctionUseCase(auctionRepository auction.AuctionRepositoryInterface, bidRepository bid.BidRepositoryInterface) *AuctionUseCase {
	return &AuctionUseCase{auctionRepositoryInterface: auctionRepository, bidRepositoryInterface: bidRepository}
}

func (a *AuctionUseCase) CreateAuction(ctx context.Context, input *AuctionInputDTO) *internal_errors.InternalError {
	newAuction, err := auction.CreateAuction(
		input.ProductName,
		input.Category,
		input.Description,
		auction.ProductCondition(input.Condition))

	if err != nil {
		return err
	}

	err = a.auctionRepositoryInterface.CreateAuction(ctx, newAuction)
	if err != nil {
		return err
	}

	return nil
}
