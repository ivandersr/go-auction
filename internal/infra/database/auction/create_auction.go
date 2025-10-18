package auction

import (
	"context"

	"github.com/ivandersr/go-auction/config/logger"
	"github.com/ivandersr/go-auction/internal/entity/auction"
	"github.com/ivandersr/go-auction/internal/internal_errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionMongo struct {
	Id          string                   `bson:"_id"`
	ProductName string                   `bson:"product_name"`
	Category    string                   `bson:"category"`
	Description string                   `bson:"description"`
	Condition   auction.ProductCondition `bson:"condition"`
	Status      auction.AuctionStatus    `bson:"status"`
	CreatedAt   int64                    `bson:"created_at"`
}

type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (r *AuctionRepository) CreateAuction(ctx context.Context, auction *auction.Auction) *internal_errors.InternalError {
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

	return nil
}
