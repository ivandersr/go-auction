package auction_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ivandersr/go-auction/config/database/mongodb"
	auctionEntity "github.com/ivandersr/go-auction/internal/entity/auction"
	auctionMongo "github.com/ivandersr/go-auction/internal/infra/database/auction"
)

func TestCreateAuctionCompletesExpiredAuctions(t *testing.T) {
	ctx := context.Background()
	t.Setenv("AUCTION_INTERVAL", "5s")
	t.Setenv("MONGODB_URL", "mongodb://admin:admin@mongodb:27017/auctions_test?authSource=admin")
	t.Setenv("MONGODB_DB", "auctions_test")
	testDB, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		panic(err)
	}
	auctionRepository := auctionMongo.NewAuctionRepository(testDB)

	newAuction := auctionEntity.Auction{
		Id:          uuid.New().String(),
		ProductName: "Test",
		Category:    "Test",
		Description: "Test Description",
		Condition:   auctionEntity.New,
		Status:      auctionEntity.Active,
		CreatedAt:   time.Now(),
	}

	err = auctionRepository.CreateAuction(ctx, &newAuction)
	fmt.Println(err)
	require.Nil(t, err)

	auction, err := auctionRepository.FindAuctionById(ctx, newAuction.Id)
	require.Nil(t, err)
	require.Equal(t, auction.Status, auctionEntity.Active)

	time.Sleep(6 * time.Second)

	auction, err = auctionRepository.FindAuctionById(ctx, newAuction.Id)
	require.Nil(t, err)
	require.Equal(t, auction.Status, auctionEntity.Completed)
}
