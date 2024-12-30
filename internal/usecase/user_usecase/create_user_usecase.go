package user_usecase

import (
	"context"
	"github.com/tiagoncardoso/fc-pge-auction-multithread/internal/entity/user_entity"
	"github.com/tiagoncardoso/fc-pge-auction-multithread/internal/internal_error"
)

type CreateUserUseCase struct {
	userRepositoryInterface user_entity.UserRepositoryInterface
}

func NewCreateUserUseCase(userRepositoryInterface user_entity.UserRepositoryInterface) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepositoryInterface: userRepositoryInterface,
	}
}

func (uc *CreateUserUseCase) CreateUser(ctx context.Context, name string) *internal_error.InternalError {
	err := uc.userRepositoryInterface.CreateUser(ctx, name)
	if err != nil {
		return err
	}

	return nil
}

type CreateUserUseCaseInterface interface {
	CreateUser(ctx context.Context, name string) *internal_error.InternalError
}
