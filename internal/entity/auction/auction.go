package auction

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ivandersr/go-auction/internal/internal_errors"
)

type Auction struct {
	Id          string
	ProductName string
	Category    string
	Description string
	Condition   ProductCondition
	Status      AuctionStatus
	CreatedAt   time.Time
}

type ProductCondition int
type AuctionStatus int

const (
	New ProductCondition = iota
	Used
	Refurbished
)

const (
	Active AuctionStatus = iota
	Completed
)

type AuctionRepositoryInterface interface {
	CreateAuction(ctx context.Context, auctionEntity *Auction) *internal_errors.InternalError
	FindAuctions(ctx context.Context, status AuctionStatus, category, productName string) ([]Auction, *internal_errors.InternalError)
	FindAuctionById(ctx context.Context, id string) (*Auction, *internal_errors.InternalError)
}

func CreateAuction(productName, category, description string, condition ProductCondition) (*Auction, *internal_errors.InternalError) {
	auction := &Auction{
		Id:          uuid.New().String(),
		ProductName: productName,
		Category:    category,
		Description: description,
		Condition:   condition,
		Status:      Active,
		CreatedAt:   time.Now(),
	}

	if err := auction.Validate(); err != nil {
		return nil, err
	}

	return auction, nil
}

func (a *Auction) Validate() *internal_errors.InternalError {
	if len(a.ProductName) <= 1 ||
		len(a.Category) <= 1 ||
		len(a.Description) <= 10 || (a.Condition != New && a.Condition != Used && a.Condition != Refurbished) {
		return internal_errors.NewBadRequestError("invalid auction data")
	}

	return nil
}
