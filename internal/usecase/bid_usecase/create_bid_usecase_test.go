package bid_usecase

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/dpcamargo/distributed-auction-system-go/internal/entity/bid_entity"
	"github.com/dpcamargo/distributed-auction-system-go/internal/internal_error"
)

type mockBidRepository struct {
	mu       sync.Mutex
	inserted []bid_entity.Bid
}

func (m *mockBidRepository) CreateBid(ctx context.Context, bids []bid_entity.Bid) *internal_error.InternalError {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.inserted = append(m.inserted, bids...)
	return nil
}

func (m *mockBidRepository) FindBidByAuctionID(ctx context.Context, auctionID string) ([]bid_entity.Bid, *internal_error.InternalError) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var list []bid_entity.Bid
	for _, b := range m.inserted {
		if b.AuctionID == auctionID {
			list = append(list, b)
		}
	}
	return list, nil
}

func (m *mockBidRepository) FindWinningBidByAuctionID(ctx context.Context, auctionID string) (*bid_entity.Bid, *internal_error.InternalError) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var highest *bid_entity.Bid
	for _, b := range m.inserted {
		if b.AuctionID == auctionID {
			if highest == nil || b.Amount > highest.Amount {
				copy := b
				highest = &copy
			}
		}
	}
	if highest == nil {
		return nil, internal_error.NewNotFoundError("no bids found")
	}
	return highest, nil
}

func TestBidBatchProcessing(t *testing.T) {
	os.Setenv("BATCH_INSERT_INTERVAL", "50ms")
	os.Setenv("MAX_BATCH_SIZE", "3")
	defer os.Unsetenv("BATCH_INSERT_INTERVAL")
	defer os.Unsetenv("MAX_BATCH_SIZE")

	repo := &mockBidRepository{}
	bu := NewBidUsecase(repo)

	userID := "a0a12131-dbbd-4896-a48f-6e60e418760a"
	auctionID := "b1b23242-dbbd-4896-a48f-6e60e418760b"

	// Insert 3 bids to reach maxBatchSize immediately
	for i := 1; i <= 3; i++ {
		err := bu.CreateBid(context.Background(), BidInputDTO{
			UserID:    userID,
			AuctionID: auctionID,
			Amount:    float64(i * 10),
		})
		if err != nil {
			t.Fatalf("unexpected error creating bid: %v", err)
		}
	}

	// Short wait for goroutine to process batch
	time.Sleep(20 * time.Millisecond)

	repo.mu.Lock()
	count := len(repo.inserted)
	repo.mu.Unlock()

	if count != 3 {
		t.Errorf("expected 3 bids inserted via batch size trigger, got %d", count)
	}

	// Test FindBidByAuctionID (fixed bug: delegating to repository without recursion)
	bids, err := bu.FindBidByAuctionID(context.Background(), auctionID)
	if err != nil {
		t.Fatalf("unexpected error finding bids: %v", err)
	}
	if len(bids) != 3 {
		t.Errorf("expected 3 bids returned, got %d", len(bids))
	}

	// Test FindWinningBidByAuctionID (highest bid)
	winner, err := bu.FindWinningBidByAuctionID(context.Background(), auctionID)
	if err != nil {
		t.Fatalf("unexpected error finding winning bid: %v", err)
	}
	if winner.Amount != 30.0 {
		t.Errorf("expected winning bid amount 30.0, got %f", winner.Amount)
	}
}
