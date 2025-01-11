package auction

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/tiagoncardoso/fc-pge-auction-multithread/internal/entity/auction_entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"os"
	"testing"
	"time"
)

var auctionEntity, _ = auction_entity.CreateAuction(
	"Test Product",
	"Test Category",
	"This is a test description",
	auction_entity.New,
)

func TestCreateAuction(t *testing.T) {
	dbName := "auction_test"

	err := os.Setenv("AUCTION_INTERVAL", "2s")
	defer os.Clearenv()

	assert.Nil(t, err)

	mtestDb := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mtestDb.Run("givenAuctionPayloadWhenPayloadIsValidThenCreateAndCloseItAfterTimeout", func(mt *mtest.T) {
		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "n", Value: 1},
			{Key: "acknowledged", Value: true},
		})

		dbMock := mt.Client.Database(dbName)

		auctionRepository := NewAuctionRepository(dbMock)

		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "n", Value: 1},
			{Key: "acknowledged", Value: true},
		})
		err = auctionRepository.CreateAuction(context.Background(), auctionEntity)

		time.Sleep(3 * time.Second)

		assert.Nil(t, err)

		mt.AddMockResponses(mtest.CreateCursorResponse(
			1,
			"auction_test.auctions",
			mtest.FirstBatch,
			bson.D{
				{Key: "_id", Value: auctionEntity.Id},
				{Key: "product_name", Value: "Test Product"},
				{Key: "category", Value: "Test Category"},
				{Key: "description", Value: "This is a test description"},
				{Key: "condition", Value: auction_entity.New},
				{Key: "status", Value: auction_entity.Completed},
			},
		))

		var result AuctionEntityMongo
		err = dbMock.Collection("auctions").FindOne(context.Background(), bson.M{"_id": auctionEntity.Id}).Decode(&result)

		assert.Nil(t, err)
		assert.Equal(t, auction_entity.Completed, result.Status)
		assert.Equal(t, auctionEntity.Id, result.Id)
	})
}
