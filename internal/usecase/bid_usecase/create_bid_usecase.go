package bid_usecase

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/dpcamargo/distributed-auction-system-go/configuration/logger"
	"github.com/dpcamargo/distributed-auction-system-go/internal/entity/bid_entity"
	"github.com/dpcamargo/distributed-auction-system-go/internal/internal_error"
)

type BidInputDTO struct {
	UserID    string  `json:"user_id"`
	AuctionID string  `json:"auction_id"`
	Amount    float64 `json:"amount"`
}

type BidOutputDTO struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	AuctionID string    `json:"auction_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

type BidUsecase struct {
	BidRepository bid_entity.BidEntityRepository

	timer               *time.Timer
	maxBatchSize        int
	batchInsertInterval time.Duration
	bidChannel          chan bid_entity.Bid
	bidBatch            []bid_entity.Bid
}

func NewBidUsecase(bidRepository bid_entity.BidEntityRepository) BidUseCaseInterface {
	maxSizeInterval := getMaxBatchSizeInterval()
	maxBatchSize := getMaxBatchSize()
	bidUsecase := &BidUsecase{
		BidRepository:       bidRepository,
		timer:               time.NewTimer(maxSizeInterval),
		maxBatchSize:        maxBatchSize,
		batchInsertInterval: maxSizeInterval,
		bidChannel:          make(chan bid_entity.Bid, maxBatchSize),
		bidBatch:            make([]bid_entity.Bid, 0, maxBatchSize),
	}
	bidUsecase.triggerCreateRoutine(context.Background())
	return bidUsecase
}

type BidUseCaseInterface interface {
	CreateBid(ctx context.Context, bidInputDTO BidInputDTO) *internal_error.InternalError
	FindBidByAuctionID(ctx context.Context, auctionID string) ([]BidOutputDTO, *internal_error.InternalError)
	FindWinningBidByAuctionID(ctx context.Context, auctionID string) (*BidOutputDTO, *internal_error.InternalError)
}

func (bu *BidUsecase) CreateBid(
	ctx context.Context,
	bidInputDTO BidInputDTO) *internal_error.InternalError {

	bidEntity, err := bid_entity.CreateBid(bidInputDTO.UserID, bidInputDTO.AuctionID, bidInputDTO.Amount)
	if err != nil {
		return err
	}

	bu.bidChannel <- *bidEntity

	return nil
}

func (bu *BidUsecase) triggerCreateRoutine(ctx context.Context) {
	go func() {
		defer close(bu.bidChannel)

		for {
			select {
			case bidEntity, ok := <-bu.bidChannel:
				if !ok {
					if len(bu.bidBatch) > 0 {
						if err := bu.BidRepository.CreateBid(ctx, bu.bidBatch); err != nil {
							logger.Error("error trying to process bid batch list", err)
						}
					}
					return
				}
				bu.bidBatch = append(bu.bidBatch, bidEntity)
				if len(bu.bidBatch) >= bu.maxBatchSize {
					if err := bu.BidRepository.CreateBid(ctx, bu.bidBatch); err != nil {
						logger.Error("error trying to process bid batch list", err)
					}
					bu.bidBatch = nil
					bu.timer.Reset(bu.batchInsertInterval)
				}
			case <-bu.timer.C:
				if len(bu.bidBatch) > 0 {
					if err := bu.BidRepository.CreateBid(ctx, bu.bidBatch); err != nil {
						logger.Error("error trying to process bid batch list", err)
					}
					bu.bidBatch = nil
				}
				bu.timer.Reset(bu.batchInsertInterval)
			}

		}
	}()
}

func getMaxBatchSizeInterval() time.Duration {
	batchInsertInterval := os.Getenv("BATCH_INSERT_INTERVAL")
	duration, err := time.ParseDuration(batchInsertInterval)
	if err != nil {
		return 3 * time.Minute
	}

	return duration
}

func getMaxBatchSize() int {
	batchSize, err := strconv.Atoi(os.Getenv("MAX_BATCH_SIZE"))
	if err != nil {
		return 5
	}

	return batchSize

}
