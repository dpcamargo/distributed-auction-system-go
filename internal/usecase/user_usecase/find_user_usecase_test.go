package user_usecase

import (
	"context"
	"testing"

	"github.com/dpcamargo/distributed-auction-system-go/internal/entity/user_entity"
	"github.com/dpcamargo/distributed-auction-system-go/internal/internal_error"
)

type mockUserRepository struct {
	users map[string]*user_entity.User
}

func (m *mockUserRepository) FindUserById(ctx context.Context, userID string) (*user_entity.User, *internal_error.InternalError) {
	u, ok := m.users[userID]
	if !ok {
		return nil, internal_error.NewNotFoundError("User not found")
	}
	return u, nil
}

func TestFindUserById(t *testing.T) {
	repo := &mockUserRepository{
		users: map[string]*user_entity.User{
			"user-1": {ID: "user-1", Name: "Alice"},
		},
	}

	uc := NewUserUseCase(repo)

	u, err := uc.FindUserById(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("expected user to be found, got %v", err)
	}
	if u.Name != "Alice" {
		t.Errorf("expected name Alice, got %s", u.Name)
	}

	_, err = uc.FindUserById(context.Background(), "non-existent")
	if err == nil {
		t.Fatalf("expected error for non-existent user")
	}
}
