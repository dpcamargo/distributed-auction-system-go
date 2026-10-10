package auction_usecase

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/dpcamargo/distributed-auction-system-go/internal/entity/auction_entity"
	"github.com/dpcamargo/distributed-auction-system-go/internal/entity/bid_entity"
	"github.com/dpcamargo/distributed-auction-system-go/internal/internal_error"
)

type mockAuctionRepository struct {
	mu       sync.Mutex
	auctions map[string]*auction_entity.Auction
}

func newMockAuctionRepository() *mockAuctionRepository {
	return &mockAuctionRepository{
		auctions: make(map[string]*auction_entity.Auction),
	}
}

func (m *mockAuctionRepository) CreateAuction(ctx context.Context, a *auction_entity.Auction) *internal_error.InternalError {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.auctions[a.ID] = a
	return nil
}

func (m *mockAuctionRepository) FindAuctionById(ctx context.Context, id string) (*auction_entity.Auction, *internal_error.InternalError) {
	m.mu.Lock()
	defer m.mu.Unlock()
	a, ok := m.auctions[id]
	if !ok {
		return nil, internal_error.NewNotFoundError("auction not found")
	}
	return a, nil
}

func (m *mockAuctionRepository) FindAuctions(ctx context.Context, status auction_entity.AuctionStatus, category, productName string) ([]auction_entity.Auction, *internal_error.InternalError) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var list []auction_entity.Auction
	for _, a := range m.auctions {
		if status != 0 && a.Status != status {
			continue
		}
		list = append(list, *a)
	}
	return list, nil
}

func (m *mockAuctionRepository) UpdateAuctionStatus(ctx context.Context, id string, status auction_entity.AuctionStatus) *internal_error.InternalError {
	m.mu.Lock()
	defer m.mu.Unlock()
	a, ok := m.auctions[id]
	if !ok {
		return internal_error.NewNotFoundError("auction not found")
	}
	a.Status = status
	return nil
}

type mockBidRepository struct{}

func (m *mockBidRepository) CreateBid(ctx context.Context, bids []bid_entity.Bid) *internal_error.InternalError {
	return nil
}

func (m *mockBidRepository) FindBidByAuctionID(ctx context.Context, auctionID string) ([]bid_entity.Bid, *internal_error.InternalError) {
	return nil, nil
}

func (m *mockBidRepository) FindWinningBidByAuctionID(ctx context.Context, auctionID string) (*bid_entity.Bid, *internal_error.InternalError) {
	return nil, nil
}

func TestAuctionAutomaticExpiration(t *testing.T) {
	os.Setenv("AUCTION_INTERVAL", "50ms")
	defer os.Unsetenv("AUCTION_INTERVAL")

	repo := newMockAuctionRepository()
	bidRepo := &mockBidRepository{}
	useCase := NewAuctionUseCase(repo, bidRepo)

	input := AuctionInputDTO{
		ProductName: "MacBook Pro M1",
		Category:    "Electronics",
		Description: "Used in great working condition",
		Condition:   ProductCondition(auction_entity.Used),
	}

	err := useCase.CreateAuction(context.Background(), input)
	if err != nil {
		t.Fatalf("expected no error creating auction, got %v", err)
	}

	auctions, err := useCase.FindAuctions(context.Background(), AuctionStatus(auction_entity.Active), "", "")
	if err != nil || len(auctions) != 1 {
		t.Fatalf("expected 1 active auction immediately after creation, got %d", len(auctions))
	}
	auctionID := auctions[0].ID

	// Wait for automatic expiration interval to fire
	time.Sleep(100 * time.Millisecond)

	updatedAuction, err := useCase.FindAuctionById(context.Background(), auctionID)
	if err != nil {
		t.Fatalf("expected to find auction, got error: %v", err)
	}
	if updatedAuction.Status != AuctionStatus(auction_entity.Completed) {
		t.Errorf("expected auction status to be Completed (%d), got %d", auction_entity.Completed, updatedAuction.Status)
	}
}
