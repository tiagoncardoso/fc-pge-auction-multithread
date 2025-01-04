package auction

import (
	"context"
	"fmt"
	"github.com/tiagoncardoso/fc-pge-auction-multithread/configuration/logger"
	"github.com/tiagoncardoso/fc-pge-auction-multithread/internal/entity/auction_entity"
	"github.com/tiagoncardoso/fc-pge-auction-multithread/internal/internal_error"
)

func (ar *AuctionRepository) CloseAuction(ctx context.Context, id string) (string, *internal_error.InternalError) {
	resp, err := ar.Collection.UpdateOne(
		ctx,
		map[string]interface{}{"_id": id},
		map[string]interface{}{"$set": map[string]interface{}{"status": auction_entity.Completed}})
	if err != nil {
		logger.Error("Error trying to update auction status", err)
		return "", internal_error.NewInternalServerError("Error trying to update auction status")
	}

	return fmt.Sprintf("%v", resp.UpsertedID), nil
}
