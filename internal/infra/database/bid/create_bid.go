package bid

import (
	"context"
	"sync"

	"github.com/ivandersr/go-auction/config/logger"
	"github.com/ivandersr/go-auction/internal/entity/auction"
	"github.com/ivandersr/go-auction/internal/entity/bid"
	auctionMongo "github.com/ivandersr/go-auction/internal/infra/database/auction"
	"github.com/ivandersr/go-auction/internal/internal_errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type BidMongo struct {
	Id        string  `bson:"_id"`
	UserId    string  `bson:"user_id"`
	AuctionId string  `bson:"auction_id"`
	Amount    float64 `bson:"amount"`
	CreatedAt int64   `bson:"created_at"`
}

type BidRepository struct {
	Collection        *mongo.Collection
	AuctionRepository *auctionMongo.AuctionRepository
}

func NewBidRepository(database *mongo.Database, auctionRepository *auctionMongo.AuctionRepository) *BidRepository {
	return &BidRepository{
		Collection:        database.Collection("bids"),
		AuctionRepository: auctionRepository,
	}
}

func (r *BidRepository) CreateBids(ctx context.Context, bids []bid.Bid) *internal_errors.InternalError {
	var wg sync.WaitGroup
	for _, bidItem := range bids {
		wg.Add(1)
		go func(bidValue bid.Bid) {
			defer wg.Done()

			auctionFound, err := r.AuctionRepository.FindAuctionById(ctx, bidValue.AuctionId)
			if err != nil {
				logger.Error("error when trying to find auction", err)
				return
			}

			if auctionFound.Status != auction.Active {
				return
			}

			bidMongo := &BidMongo{
				Id:        bidValue.Id,
				UserId:    bidValue.UserId,
				AuctionId: bidValue.AuctionId,
				Amount:    bidValue.Amount,
				CreatedAt: bidValue.CreatedAt.Unix(),
			}

			if _, err := r.Collection.InsertOne(ctx, bidMongo); err != nil {
				logger.Error("error when trying to create bid", err)
				return
			}

		}(bidItem)
	}

	wg.Wait()
	return nil
}
