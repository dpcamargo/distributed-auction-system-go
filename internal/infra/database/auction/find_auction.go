package auction

import (
	"context"
	"time"

	"github.com/dpcamargo/distributed-auction-system-go/configuration/logger"
	"github.com/dpcamargo/distributed-auction-system-go/internal/entity/auction_entity"
	"github.com/dpcamargo/distributed-auction-system-go/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func (a *AuctionRepository) FindAuctionById(ctx context.Context, id string) (*auction_entity.Auction, *internal_error.InternalError) {
	filter := bson.M{"_id": id}

	var auctionEntityMongo AuctionEntityMongo
	if err := a.Collection.FindOne(ctx, filter).Decode(&auctionEntityMongo); err != nil {
		logger.Error("Error trying to find auction by id", err, zap.String("id", id))
		return nil, internal_error.NewInternalServerError("Error trying to find auction by id")
	}

	return &auction_entity.Auction{
		ID:          auctionEntityMongo.ID,
		ProductName: auctionEntityMongo.ProductName,
		Category:    auctionEntityMongo.Category,
		Description: auctionEntityMongo.Description,
		Condition:   auctionEntityMongo.Condition,
		Status:      auctionEntityMongo.Status,
		Timestamp:   time.Unix(auctionEntityMongo.Timestamp, 0),
	}, nil
}

func (a *AuctionRepository) FindAuctions(
	ctx context.Context,
	status auction_entity.AuctionStatus,
	category, productName string) ([]auction_entity.Auction, *internal_error.InternalError) {
	filter := bson.M{}
	if status != 0 {
		filter["status"] = status
	}
	if productName != "" {
		filter["productName"] = productName
	}
	if category != "" {
		filter["category"] = primitive.Regex{
			Pattern: category,
			Options: "i",
		}
	}

	cursor, err := a.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error("Error trying to find auctions", err)
		return nil, internal_error.NewInternalServerError("Error trying to find auctions")
	}
	defer cursor.Close(ctx)

	var auctionEntityMongo []AuctionEntityMongo
	if err := cursor.All(ctx, &auctionEntityMongo); err != nil {
		logger.Error("Error trying to find auctions", err)
		return nil, internal_error.NewInternalServerError("Error trying to find auctions")
	}

	var auctionEntity []auction_entity.Auction
	for _, auctionMongo := range auctionEntityMongo {
		auctionEntity = append(auctionEntity, auction_entity.Auction{
			ID:          auctionMongo.ID,
			ProductName: auctionMongo.ProductName,
			Category:    auctionMongo.Category,
			Description: auctionMongo.Description,
			Condition:   auctionMongo.Condition,
			Status:      auctionMongo.Status,
			Timestamp:   time.Unix(auctionMongo.Timestamp, 0),
		})
	}

	return auctionEntity, nil
}
