package auction_usecase

import (
	"context"

	"github.com/ivandersr/go-auction/internal/entity/auction"
	"github.com/ivandersr/go-auction/internal/internal_errors"
	"github.com/ivandersr/go-auction/internal/usecase/bid_usecase"
)

func (a *AuctionUseCase) FindAuctionById(ctx context.Context, id string) (*AuctionOutputDTO, *internal_errors.InternalError) {
	foundAuction, err := a.auctionRepositoryInterface.FindAuctionById(ctx, id)
	if err != nil {
		return nil, err
	}
	if foundAuction != nil {
		return &AuctionOutputDTO{
			Id:          foundAuction.Id,
			ProductName: foundAuction.ProductName,
			Category:    foundAuction.Category,
			Description: foundAuction.Description,
			Condition:   ProductCondition(foundAuction.Condition),
			Status:      AuctionStatus(foundAuction.Status),
			CreatedAt:   foundAuction.CreatedAt,
		}, nil
	}
	return nil, nil
}

func (a *AuctionUseCase) FindAuctions(ctx context.Context,
	status AuctionStatus,
	category, productName string) ([]AuctionOutputDTO, *internal_errors.InternalError) {
	foundAuctions, err := a.auctionRepositoryInterface.FindAuctions(ctx, auction.AuctionStatus(status), category, productName)
	if err != nil {
		return nil, err
	}

	var outputList []AuctionOutputDTO
	for _, found := range foundAuctions {
		outputList = append(outputList, AuctionOutputDTO{
			Id:          found.Id,
			ProductName: found.ProductName,
			Category:    found.Category,
			Description: found.Description,
			Condition:   ProductCondition(found.Condition),
			Status:      AuctionStatus(found.Status),
			CreatedAt:   found.CreatedAt,
		})
	}

	return outputList, nil
}

func (a *AuctionUseCase) FindWinningBidByAuctionId(ctx context.Context, id string) (*WinningInfoOutputDTO, *internal_errors.InternalError) {
	foundAuction, err := a.auctionRepositoryInterface.FindAuctionById(ctx, id)
	if err != nil {
		return nil, err
	}

	auctionOutput := AuctionOutputDTO{
		Id:          foundAuction.Id,
		ProductName: foundAuction.ProductName,
		Category:    foundAuction.Category,
		Description: foundAuction.Description,
		Condition:   ProductCondition(foundAuction.Condition),
		Status:      AuctionStatus(foundAuction.Status),
		CreatedAt:   foundAuction.CreatedAt,
	}

	winningBid, err := a.bidRepositoryInterface.FindWinningBidByAuctionId(ctx, foundAuction.Id)
	if err != nil {
		return &WinningInfoOutputDTO{
			Auction: auctionOutput,
			Bid:     nil,
		}, err
	}

	return &WinningInfoOutputDTO{
		Auction: auctionOutput,
		Bid: &bid_usecase.BidOutputDTO{
			Id:        winningBid.Id,
			UserId:    winningBid.UserId,
			AuctionId: winningBid.AuctionId,
			Amount:    winningBid.Amount,
			CreatedAt: winningBid.CreatedAt,
		},
	}, nil
}
