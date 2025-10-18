package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ivandersr/go-auction/config/database/mongodb"
	"github.com/ivandersr/go-auction/internal/infra/api/web/controller/auction_controller"
	"github.com/ivandersr/go-auction/internal/infra/api/web/controller/bid_controller"
	"github.com/ivandersr/go-auction/internal/infra/api/web/controller/user_controller"
	"github.com/ivandersr/go-auction/internal/infra/database/auction"
	"github.com/ivandersr/go-auction/internal/infra/database/bid"
	"github.com/ivandersr/go-auction/internal/infra/database/user"
	"github.com/ivandersr/go-auction/internal/usecase/auction_usecase"
	"github.com/ivandersr/go-auction/internal/usecase/bid_usecase"
	"github.com/ivandersr/go-auction/internal/usecase/user_usecase"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	ctx := context.Background()
	if err := godotenv.Load("cmd/auction/.env"); err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	databaseConnection, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	userController, bidController, auctionController := initDependencies(databaseConnection)

	router := gin.Default()
	router.GET("/auctions", auctionController.FindAuctions)
	router.GET("/auctions/:auctionId", auctionController.FindAuctionById)
	router.GET("/auctions/winner/:auctionId", auctionController.FindWinningBidByAuctionId)
	router.POST("/auctions", auctionController.CreateAuction)
	router.GET("/bid/:auctionId", bidController.FindBidsByAuctionId)
	router.POST("/bid", bidController.CreateBid)
	router.GET("/user/:userId", userController.FindUserById)

	router.Run(":8080")
}

func initDependencies(database *mongo.Database) (
	userController *user_controller.UserController,
	bidController *bid_controller.BidController,
	auctionController *auction_controller.AuctionController,
) {
	auctionRepository := auction.NewAuctionRepository(database)
	bidRepository := bid.NewBidRepository(database, auctionRepository)
	userRepository := user.NewUserRepository(database)

	return user_controller.NewUserController(user_usecase.NewUserUseCase(userRepository)),
		bid_controller.NewBidController(bid_usecase.NewBidUseCase(bidRepository)),
		auction_controller.NewAuctionController(auction_usecase.NewAuctionUseCase(auctionRepository, bidRepository))
}
