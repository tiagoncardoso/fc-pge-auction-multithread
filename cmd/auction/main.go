package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tiagoncardoso/fc-pge-auction-multithread/configuration/database/mongodb"
	"github.com/tiagoncardoso/fc-pge-auction-multithread/internal/entity/auction_entity"
	"github.com/tiagoncardoso/fc-pge-auction-multithread/internal/infra/api/web/controller/auction_controller"
	"github.com/tiagoncardoso/fc-pge-auction-multithread/internal/infra/api/web/controller/bid_controller"
	"github.com/tiagoncardoso/fc-pge-auction-multithread/internal/infra/api/web/controller/user_controller"
	"github.com/tiagoncardoso/fc-pge-auction-multithread/internal/infra/database/auction"
	"github.com/tiagoncardoso/fc-pge-auction-multithread/internal/infra/database/bid"
	"github.com/tiagoncardoso/fc-pge-auction-multithread/internal/infra/database/user"
	"github.com/tiagoncardoso/fc-pge-auction-multithread/internal/usecase/auction_usecase"
	"github.com/tiagoncardoso/fc-pge-auction-multithread/internal/usecase/bid_usecase"
	"github.com/tiagoncardoso/fc-pge-auction-multithread/internal/usecase/user_usecase"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load("cmd/auction/.env"); err != nil {
		log.Fatal("Error trying to load env variables")
		return
	}

	databaseConnection, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	router := gin.Default()

	userController, bidController, auctionsController, createUserController := initDependencies(databaseConnection)

	router.GET("/auction", auctionsController.FindAuctions)
	router.GET("/auction/:auctionId", auctionsController.FindAuctionById)
	router.POST("/auction", auctionsController.CreateAuction)
	router.GET("/auction/winner/:auctionId", auctionsController.FindWinningBidByAuctionId)
	router.POST("/bid", bidController.CreateBid)
	router.GET("/bid/:auctionId", bidController.FindBidByAuctionId)
	router.GET("/user/:userId", userController.FindUserById)
	router.POST("/user", createUserController.CreateUser)

	router.Run(":8080")
}

func initDependencies(database *mongo.Database) (
	userController *user_controller.UserController,
	bidController *bid_controller.BidController,
	auctionController *auction_controller.AuctionController,
	createUserController *user_controller.CreateUserController,
) {

	auctionTimeControl := make(chan auction_entity.AuctionTimeoutControl)

	auctionRepository := auction.NewAuctionRepository(database, auctionTimeControl)
	bidRepository := bid.NewBidRepository(database, auctionRepository, auctionTimeControl)
	userRepository := user.NewUserRepository(database)

	userController = user_controller.NewUserController(
		user_usecase.NewUserUseCase(userRepository))
	createUserController = user_controller.NewCreateUserController(user_usecase.NewCreateUserUseCase(userRepository))
	auctionController = auction_controller.NewAuctionController(
		auction_usecase.NewAuctionUseCase(auctionRepository, bidRepository))
	bidController = bid_controller.NewBidController(bid_usecase.NewBidUseCase(bidRepository))

	return
}
