package auction

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ivandersr/go-auction/config/logger"
	"github.com/ivandersr/go-auction/internal/entity/auction"
	"github.com/ivandersr/go-auction/internal/internal_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *AuctionRepository) FindAuctionById(ctx context.Context, id string) (*auction.Auction, *internal_errors.InternalError) {
	fiter := bson.M{"_id": id}

	var auctionMongo AuctionMongo
	if err := r.Collection.FindOne(ctx, fiter).Decode(&auctionMongo); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			errorMessage := fmt.Sprintf("auction not found with id %s", id)
			logger.Error(errorMessage, err)
			return nil, internal_errors.NewNotFoundError(errorMessage)
		}
		errorMessage := fmt.Sprintf("error when trying to find auction with id %s", id)
		logger.Error(errorMessage, err)
		return nil, internal_errors.NewInternalServerError(errorMessage)
	}

	return &auction.Auction{
		Id:          auctionMongo.Id,
		ProductName: auctionMongo.ProductName,
		Description: auctionMongo.Description,
		Category:    auctionMongo.Category,
		Condition:   auctionMongo.Condition,
		Status:      auctionMongo.Status,
		CreatedAt:   time.Unix(auctionMongo.CreatedAt, 0),
	}, nil
}

func (r *AuctionRepository) FindAuctions(
	ctx context.Context,
	status auction.AuctionStatus,
	category, productName string) ([]auction.Auction, *internal_errors.InternalError) {
	filter := bson.M{}
	var errorMessage string

	if status != 0 {
		filter["status"] = status
	}

	if category != "" {
		filter["category"] = category
	}

	if productName != "" {
		filter["product_name"] = primitive.Regex{
			Pattern: productName,
			Options: "i",
		}
	}

	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		errorMessage = "error when trying to find auctions"
		logger.Error(errorMessage, err)
		return nil, internal_errors.NewInternalServerError(errorMessage)
	}
	defer cursor.Close(ctx)

	var auctionMongoList []AuctionMongo
	if err := cursor.All(ctx, &auctionMongoList); err != nil {
		errorMessage = "error when trying to find auctions"
		logger.Error(errorMessage, err)
		return nil, internal_errors.NewInternalServerError(errorMessage)
	}

	var auctionList []auction.Auction

	for _, auctionMongo := range auctionMongoList {
		auctionList = append(auctionList, auction.Auction{
			Id:          auctionMongo.Id,
			ProductName: auctionMongo.ProductName,
			Description: auctionMongo.Description,
			Category:    auctionMongo.Category,
			Condition:   auctionMongo.Condition,
			Status:      auctionMongo.Status,
			CreatedAt:   time.Unix(auctionMongo.CreatedAt, 0),
		})
	}

	return auctionList, nil
}
