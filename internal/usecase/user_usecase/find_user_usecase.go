package user_usecase

import (
	"context"

	"github.com/dpcamargo/distributed-auction-system-go/internal/entity/user_entity"
	"github.com/dpcamargo/distributed-auction-system-go/internal/internal_error"
)

type UserUseCase struct {
	UserRepository user_entity.UserRepositoryInterface
}

type UserOutputDTO struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type UserUseCaseInterface interface {
	FindUserById(ctx context.Context, id string) (*UserOutputDTO, *internal_error.InternalError)
}

func NewUserUseCase(userRepository user_entity.UserRepositoryInterface) UserUseCaseInterface {
	return &UserUseCase{
		UserRepository: userRepository,
	}
}

func (u *UserUseCase) FindUserById(ctx context.Context, id string) (*UserOutputDTO, *internal_error.InternalError) {
	userEntity, err := u.UserRepository.FindUserById(ctx, id)
	if err != nil {
		return nil, err
	}
	return &UserOutputDTO{
		ID:   userEntity.ID,
		Name: userEntity.Name,
	}, nil
}
