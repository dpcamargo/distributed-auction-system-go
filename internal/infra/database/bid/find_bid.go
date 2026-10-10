package bid

import (
	"context"
	"time"

	"github.com/dpcamargo/distributed-auction-system-go/configuration/logger"
	"github.com/dpcamargo/distributed-auction-system-go/internal/entity/bid_entity"
	"github.com/dpcamargo/distributed-auction-system-go/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func (bd *BidRepository) FindBidByAuctionID(ctx context.Context, auctionID string) ([]bid_entity.Bid, *internal_error.InternalError) {
	filter := bson.M{"auction_id": auctionID}

	cursor, err := bd.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error("Error trying to find bids by auctionID", err, zap.String("auctionID", auctionID))
		return nil, internal_error.NewInternalServerError("Error trying to find bids by auctionID")
	}

	var bidEntitiesMongo []BidEntityMongo
	if err := cursor.All(ctx, &bidEntitiesMongo); err != nil {
		logger.Error("Error trying to find bids by auctionID", err, zap.String("auctionID", auctionID))
		return nil, internal_error.NewInternalServerError("Error trying to find bids by auctionID")
	}

	var bidEntities []bid_entity.Bid
	for _, bidEntity := range bidEntitiesMongo {
		bidEntities = append(bidEntities, bid_entity.Bid{
			ID:        bidEntity.ID,
			UserID:    bidEntity.UserID,
			AuctionID: bidEntity.AuctionID,
			Amount:    bidEntity.Amount,
			Timestamp: time.Unix(bidEntity.Timestamp, 0),
		})
	}
	return bidEntities, nil
}

func (bd *BidRepository) FindWinningBidByAuctionID(ctx context.Context, auctionID string) (*bid_entity.Bid, *internal_error.InternalError) {
	filter := bson.M{"auction_id": auctionID}

	var bidEntityMongo BidEntityMongo
	opts := options.FindOne().SetSort(bson.D{{Key: "amount", Value: -1}})
	if err := bd.Collection.FindOne(ctx, filter, opts).Decode(&bidEntityMongo); err != nil {
		logger.Error("Error trying to find winning bid", err, zap.String("auctionID", auctionID))
		return nil, internal_error.NewInternalServerError("Error trying to find find winning bid")
	}

	return &bid_entity.Bid{
		ID:        bidEntityMongo.ID,
		UserID:    bidEntityMongo.UserID,
		AuctionID: bidEntityMongo.AuctionID,
		Amount:    bidEntityMongo.Amount,
		Timestamp: time.Unix(bidEntityMongo.Timestamp, 0),
	}, nil
}
