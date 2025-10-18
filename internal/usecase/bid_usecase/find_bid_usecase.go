package bid_usecase

import (
	"context"

	"github.com/ivandersr/go-auction/internal/internal_errors"
)

func (bu *BidUseCase) FindBidsByAuctionId(ctx context.Context, auctionId string) ([]BidOutputDTO, *internal_errors.InternalError) {
	bids, err := bu.BidRepository.FindBidsByAuctionId(ctx, auctionId)
	if err != nil {
		return nil, err
	}
	var outputList []BidOutputDTO
	for _, bid := range bids {
		outputList = append(outputList, BidOutputDTO{
			Id:        bid.Id,
			UserId:    bid.UserId,
			AuctionId: bid.AuctionId,
			Amount:    bid.Amount,
			CreatedAt: bid.CreatedAt,
		})
	}
	return outputList, nil
}

func (bu *BidUseCase) FindWinningBidByAuctionId(ctx context.Context, id string) (*BidOutputDTO, *internal_errors.InternalError) {
	bid, err := bu.BidRepository.FindWinningBidByAuctionId(ctx, id)
	if err != nil {
		return nil, err
	}
	if bid == nil {
		return nil, nil
	}
	return &BidOutputDTO{
		Id:        bid.Id,
		UserId:    bid.UserId,
		AuctionId: bid.AuctionId,
		Amount:    bid.Amount,
		CreatedAt: bid.CreatedAt,
	}, nil
}
