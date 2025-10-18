package bid

import (
	"context"
	"fmt"
	"time"

	"github.com/ivandersr/go-auction/config/logger"
	"github.com/ivandersr/go-auction/internal/entity/bid"
	"github.com/ivandersr/go-auction/internal/internal_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *BidRepository) FindBidsByAuctionId(ctx context.Context, auctionId string) ([]bid.Bid, *internal_errors.InternalError) {
	filter := bson.M{
		"auction_id": auctionId,
	}
	var errorMessage string

	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		errorMessage = fmt.Sprintf("error when trying to find bids by auction id %s", auctionId)
		logger.Error(errorMessage, err)
		return nil, internal_errors.NewInternalServerError(errorMessage)
	}
	defer cursor.Close(ctx)

	var bidMongoList []BidMongo
	if err := cursor.All(ctx, &bidMongoList); err != nil {
		errorMessage = fmt.Sprintf("error when trying to find bids by auction id %s", auctionId)
		logger.Error(errorMessage, err)
		return nil, internal_errors.NewInternalServerError(errorMessage)
	}

	var bidList []bid.Bid

	for _, bidMongo := range bidMongoList {
		bidList = append(bidList, bid.Bid{
			Id:        bidMongo.Id,
			UserId:    bidMongo.UserId,
			AuctionId: bidMongo.AuctionId,
			Amount:    bidMongo.Amount,
			CreatedAt: time.Unix(bidMongo.CreatedAt, 0),
		})
	}

	return bidList, nil
}

func (r *BidRepository) FindWinningBidByAuctionId(ctx context.Context, auctionId string) (*bid.Bid, *internal_errors.InternalError) {
	filter := bson.M{
		"auction_id": auctionId,
	}
	opts := options.FindOne().SetSort(bson.D{{Key: "amount", Value: -1}})
	var errorMessage string
	var winningBidMongo BidMongo

	err := r.Collection.FindOne(ctx, filter, opts).Decode(&winningBidMongo)

	if err != nil {
		errorMessage = fmt.Sprintf("error when trying to find winning bid by auction id %s", auctionId)
		logger.Error(errorMessage, err)
		return nil, internal_errors.NewInternalServerError(errorMessage)
	}

	return &bid.Bid{
		Id:        winningBidMongo.Id,
		UserId:    winningBidMongo.UserId,
		AuctionId: winningBidMongo.AuctionId,
		Amount:    winningBidMongo.Amount,
		CreatedAt: time.Unix(winningBidMongo.CreatedAt, 0),
	}, nil
}
