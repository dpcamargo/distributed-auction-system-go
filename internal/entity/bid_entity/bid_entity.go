package bid_entity

import (
	"context"
	"time"

	"github.com/dpcamargo/distributed-auction-system-go/internal/internal_error"
	"github.com/google/uuid"
)

type Bid struct {
	ID        string
	UserID    string
	AuctionID string
	Amount    float64
	Timestamp time.Time
}

func CreateBid(userId, auctionId string, amount float64) (*Bid, *internal_error.InternalError) {
	bid := &Bid{
		ID:        uuid.New().String(),
		UserID:    userId,
		AuctionID: auctionId,
		Amount:    amount,
		Timestamp: time.Now(),
	}

	if err := bid.Validate(); err != nil {
		return nil, err
	}

	return bid, nil
}

func (b *Bid) Validate() *internal_error.InternalError {
	if err := uuid.Validate(b.UserID); err != nil {
		return internal_error.NewBadRequestError("userID is not a valid id")
	} else if err := uuid.Validate(b.AuctionID); err != nil {
		return internal_error.NewBadRequestError("AuctionID is not a valid id")
	} else if b.Amount <= 0 {
		return internal_error.NewBadRequestError("Amount is not a valid value")
	}
	return nil
}

type BidEntityRepository interface {
	CreateBid(ctx context.Context, bidEntities []Bid) *internal_error.InternalError
	FindBidByAuctionID(ctx context.Context, auctionID string) ([]Bid, *internal_error.InternalError)
	FindWinningBidByAuctionID(ctx context.Context, auctionID string) (*Bid, *internal_error.InternalError)
}
