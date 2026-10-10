package auction_usecase

import (
	"context"
	"os"
	"time"

	"github.com/dpcamargo/distributed-auction-system-go/configuration/logger"
	"github.com/dpcamargo/distributed-auction-system-go/internal/entity/auction_entity"
	"github.com/dpcamargo/distributed-auction-system-go/internal/entity/bid_entity"
	"github.com/dpcamargo/distributed-auction-system-go/internal/internal_error"
	"github.com/dpcamargo/distributed-auction-system-go/internal/usecase/bid_usecase"
	"go.uber.org/zap"
)

type AuctionInputDTO struct {
	ProductName string           `json:"product_name" binding:"required,min=1"`
	Category    string           `json:"category" binding:"required,min=2"`
	Description string           `json:"description" binding:"required,min=10,max=200"`
	Condition   ProductCondition `json:"condition"`
}

type AuctionOutputDTO struct {
	ID          string           `json:"id"`
	ProductName string           `json:"product_name"`
	Category    string           `json:"category"`
	Description string           `json:"description"`
	Condition   ProductCondition `json:"condition"`
	Status      AuctionStatus    `json:"status"`
	Timestamp   time.Time        `json:"timestamp"`
}

type AuctionUseCaseInterface interface {
	CreateAuction(ctx context.Context, auctionInput AuctionInputDTO) *internal_error.InternalError
	FindAuctionById(ctx context.Context, id string) (*AuctionOutputDTO, *internal_error.InternalError)
	FindAuctions(ctx context.Context, status AuctionStatus, category, productName string) ([]AuctionOutputDTO, *internal_error.InternalError)
	FindWinningBidByAuctionId(ctx context.Context, auctionId string) (*WinningInfoOutputDTO, *internal_error.InternalError)
}

type WinningInfoOutputDTO struct {
	Auction AuctionOutputDTO          `json:"auction"`
	Bid     *bid_usecase.BidOutputDTO `json:"bid,omitempty"`
}

type AuctionUseCase struct {
	auctionRepositoryInterface auction_entity.AuctionRepositoryInterface
	bidRepositoryInterface     bid_entity.BidEntityRepository
}

type ProductCondition int64
type AuctionStatus int64

func NewAuctionUseCase(auctionRepository auction_entity.AuctionRepositoryInterface,
	bidRepository bid_entity.BidEntityRepository) AuctionUseCaseInterface {
	return &AuctionUseCase{
		auctionRepositoryInterface: auctionRepository,
		bidRepositoryInterface:     bidRepository,
	}
}

func (au *AuctionUseCase) CreateAuction(ctx context.Context, auctionInput AuctionInputDTO) *internal_error.InternalError {
	auction, err := auction_entity.CreateAuction(
		auctionInput.ProductName,
		auctionInput.Category,
		auctionInput.Description,
		auction_entity.ProductCondition(auctionInput.Condition),
	)
	if err != nil {
		return err
	}

	if err := au.auctionRepositoryInterface.CreateAuction(ctx, auction); err != nil {
		return err
	}

	au.scheduleAuctionExpiration(auction.ID)
	return nil
}

func (au *AuctionUseCase) scheduleAuctionExpiration(auctionID string) {
	duration := getAuctionDuration()
	go func() {
		time.Sleep(duration)
		err := au.auctionRepositoryInterface.UpdateAuctionStatus(context.Background(), auctionID, auction_entity.Completed)
		if err != nil {
			logger.Error("Error trying to expire auction", err, zap.String("auctionID", auctionID))
		}
	}()
}

func getAuctionDuration() time.Duration {
	auctionInterval := os.Getenv("AUCTION_INTERVAL")
	duration, err := time.ParseDuration(auctionInterval)
	if err != nil {
		return 5 * time.Minute
	}
	return duration
}
