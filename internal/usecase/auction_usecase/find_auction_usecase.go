package auction_usecase

import (
	"context"

	"github.com/dpcamargo/distributed-auction-system-go/internal/entity/auction_entity"
	"github.com/dpcamargo/distributed-auction-system-go/internal/internal_error"
	"github.com/dpcamargo/distributed-auction-system-go/internal/usecase/bid_usecase"
)

func (au *AuctionUseCase) FindAuctions(ctx context.Context, status AuctionStatus, category, productName string) ([]AuctionOutputDTO, *internal_error.InternalError) {
	auctionEntities, err := au.auctionRepositoryInterface.FindAuctions(ctx, auction_entity.AuctionStatus(status), category, productName)
	if err != nil {
		return nil, err
	}

	var auctionOutputs []AuctionOutputDTO
	for _, auction := range auctionEntities {
		auctionOutputs = append(auctionOutputs, AuctionOutputDTO{
			ID:          auction.ID,
			ProductName: auction.ProductName,
			Category:    auction.Category,
			Description: auction.Description,
			Condition:   ProductCondition(auction.Condition),
			Status:      AuctionStatus(auction.Status),
			Timestamp:   auction.Timestamp,
		})
	}
	return auctionOutputs, nil
}

func (au *AuctionUseCase) FindAuctionById(ctx context.Context, id string) (*AuctionOutputDTO, *internal_error.InternalError) {
	auctionEntity, err := au.auctionRepositoryInterface.FindAuctionById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &AuctionOutputDTO{
		ID:          auctionEntity.ID,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   ProductCondition(auctionEntity.Condition),
		Status:      AuctionStatus(auctionEntity.Status),
		Timestamp:   auctionEntity.Timestamp,
	}, nil
}

func (au *AuctionUseCase) FindWinningBidByAuctionId(
	ctx context.Context,
	auctionID string) (*WinningInfoOutputDTO, *internal_error.InternalError) {
	auction, err := au.auctionRepositoryInterface.FindAuctionById(ctx, auctionID)
	if err != nil {
		return nil, err
	}

	auctionOutputDTO := AuctionOutputDTO{
		ID:          auction.ID,
		ProductName: auction.ProductName,
		Category:    auction.Category,
		Description: auction.Description,
		Condition:   ProductCondition(auction.Condition),
		Status:      AuctionStatus(auction.Status),
		Timestamp:   auction.Timestamp,
	}

	bidWinning, err := au.bidRepositoryInterface.FindWinningBidByAuctionID(ctx, auction.ID)
	if err != nil {
		return &WinningInfoOutputDTO{
			Auction: auctionOutputDTO,
			Bid:     nil,
		}, nil
	}

	bidOutputDTO := &bid_usecase.BidOutputDTO{
		ID:        bidWinning.ID,
		UserID:    bidWinning.UserID,
		AuctionID: bidWinning.AuctionID,
		Amount:    bidWinning.Amount,
		Timestamp: bidWinning.Timestamp,
	}
	return &WinningInfoOutputDTO{
		Auction: auctionOutputDTO,
		Bid:     bidOutputDTO,
	}, nil
}
