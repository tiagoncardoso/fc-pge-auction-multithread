package auction

import (
	"context"
	"fmt"
	"github.com/tiagoncardoso/fc-pge-auction-multithread/configuration/logger"
	"github.com/tiagoncardoso/fc-pge-auction-multithread/internal/entity/auction_entity"
	"github.com/tiagoncardoso/fc-pge-auction-multithread/internal/internal_error"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	auctionData, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	var wg sync.WaitGroup
	auctionTimeout := calcAuctionTimeout(auctionEntity.Timestamp.Unix())

	go ar.AuctionStatusWatch(ctx, fmt.Sprintf("%v", auctionData.InsertedID), auctionTimeout, &wg)

	return nil
}

func (ar *AuctionRepository) AuctionStatusWatch(ctx context.Context, auctionId string, auctionTimeout time.Time, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			if ar.AuctionTimeoutChecker(ctx, auctionId, auctionTimeout) {
				return
			}
		}
	}
}

func (ar *AuctionRepository) AuctionTimeoutChecker(ctx context.Context, auctionId string, endTime time.Time) bool {
	logger.Info(fmt.Sprintf("Checking auction %s status", auctionId))
	if time.Now().After(endTime) {
		auId, err := ar.CloseAuction(ctx, auctionId)
		if err != nil {
			logger.Error(fmt.Sprintf("Error closing auction %s", auId), err)
			return true
		}

		logger.Info(fmt.Sprintf("Auction %s closed", auctionId))
		return true
	}

	return false
}

func calcAuctionTimeout(auctionTimeStamp int64) time.Time {
	timeout := time.Unix(auctionTimeStamp, 0)
	auctionInterval := os.Getenv("AUCTION_INTERVAL")

	interval, err := time.ParseDuration(auctionInterval)
	if err != nil {
		interval = 60 * time.Second
	}

	return timeout.Add(interval)
}
