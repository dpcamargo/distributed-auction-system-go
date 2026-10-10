package bid_usecase

import (
	"context"

	"github.com/dpcamargo/distributed-auction-system-go/internal/internal_error"
)

func (bu *BidUsecase) FindBidByAuctionID(ctx context.Context, auctionID string) ([]BidOutputDTO, *internal_error.InternalError) {
	bidList, err := bu.BidRepository.FindBidByAuctionID(ctx, auctionID)
	if err != nil {
		return nil, err
	}

	var bidOutputList []BidOutputDTO
	for _, bid := range bidList {
		bidOutputList = append(bidOutputList, BidOutputDTO{
			ID:        bid.ID,
			UserID:    bid.UserID,
			AuctionID: bid.AuctionID,
			Amount:    bid.Amount,
			Timestamp: bid.Timestamp,
		})
	}
	return bidOutputList, nil
}

func (bu *BidUsecase) FindWinningBidByAuctionID(ctx context.Context, auctionID string) (*BidOutputDTO, *internal_error.InternalError) {
	bidEntity, err := bu.BidRepository.FindWinningBidByAuctionID(ctx, auctionID)
	if err != nil {
		return nil, err
	}

	return &BidOutputDTO{
		ID:        bidEntity.ID,
		UserID:    bidEntity.UserID,
		AuctionID: bidEntity.AuctionID,
		Amount:    bidEntity.Amount,
		Timestamp: bidEntity.Timestamp,
	}, nil
}
