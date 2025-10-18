package bid

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ivandersr/go-auction/internal/internal_errors"
)

type Bid struct {
	Id        string
	UserId    string
	AuctionId string
	Amount    float64
	CreatedAt time.Time
}

type BidRepositoryInterface interface {
	CreateBids(ctx context.Context, bids []Bid) *internal_errors.InternalError
	FindBidsByAuctionId(ctx context.Context, auctionId string) ([]Bid, *internal_errors.InternalError)
	FindWinningBidByAuctionId(ctx context.Context, auctionId string) (*Bid, *internal_errors.InternalError)
}

func CreateBid(userId string, auctionId string, amount float64) (*Bid, *internal_errors.InternalError) {
	newBid := &Bid{
		Id:        uuid.New().String(),
		UserId:    userId,
		AuctionId: auctionId,
		Amount:    amount,
		CreatedAt: time.Now(),
	}

	if err := newBid.Validate(); err != nil {
		return nil, err
	}

	return newBid, nil
}

func (b Bid) Validate() *internal_errors.InternalError {
	if err := uuid.Validate(b.UserId); err != nil {
		return internal_errors.NewBadRequestError("invalid user id")
	}
	if err := uuid.Validate(b.AuctionId); err != nil {
		return internal_errors.NewBadRequestError("invalid auction id")
	}
	if b.Amount <= 0 {
		return internal_errors.NewBadRequestError("invalid amount")
	}
	return nil
}
