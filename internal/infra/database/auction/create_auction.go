package auction

import (
	"context"
	"os"
	"time"

	"github.com/ivandersr/go-auction/config/logger"
	auctionEntity "github.com/ivandersr/go-auction/internal/entity/auction"
	"github.com/ivandersr/go-auction/internal/internal_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionMongo struct {
	Id          string                         `bson:"_id"`
	ProductName string                         `bson:"product_name"`
	Category    string                         `bson:"category"`
	Description string                         `bson:"description"`
	Condition   auctionEntity.ProductCondition `bson:"condition"`
	Status      auctionEntity.AuctionStatus    `bson:"status"`
	CreatedAt   int64                          `bson:"created_at"`
}

type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (r *AuctionRepository) CreateAuction(ctx context.Context, auction *auctionEntity.Auction) *internal_errors.InternalError {
	auctionMongo := &AuctionMongo{
		Id:          auction.Id,
		ProductName: auction.ProductName,
		Category:    auction.Category,
		Description: auction.Description,
		Condition:   auction.Condition,
		Status:      auction.Status,
		CreatedAt:   auction.CreatedAt.Unix(),
	}

	_, err := r.Collection.InsertOne(ctx, auctionMongo)
	if err != nil {
		logger.Error("error when trying to create auction", err)
		return internal_errors.NewInternalServerError("error when trying to create auction")
	}

	go func() {
		time.Sleep(getAuctionInterval())
		update := bson.M{"$set": bson.M{"status": auctionEntity.Completed}}
		filter := bson.M{"_id": auctionMongo.Id}

		_, err := r.Collection.UpdateOne(ctx, filter, update)
		if err != nil {
			logger.Error("error when trying to update auction", err)
			return
		}
	}()

	return nil
}

func getAuctionInterval() time.Duration {
	auctionInterval := os.Getenv("AUCTION_INTERVAL")
	duration, err := time.ParseDuration(auctionInterval)
	if err != nil {
		return 10 * time.Minute
	}
	return duration
}
