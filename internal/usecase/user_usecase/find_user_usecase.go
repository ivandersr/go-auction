package user_usecase

import (
	"context"

	"github.com/ivandersr/go-auction/internal/entity/user"
	"github.com/ivandersr/go-auction/internal/internal_errors"
)

type UserUsecase struct {
	UserRepository user.UserRepostitoryInterface
}

type UserOutputDTO struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func NewUserUseCase(userRepository user.UserRepostitoryInterface) UserUsecaseInterface {
	return &UserUsecase{UserRepository: userRepository}
}

type UserUsecaseInterface interface {
	FindUserById(ctx context.Context, id string) (*UserOutputDTO, *internal_errors.InternalError)
}

func (u *UserUsecase) FindUserById(ctx context.Context, id string) (*UserOutputDTO, *internal_errors.InternalError) {
	foundUser, err := u.UserRepository.FindUserById(ctx, id)

	if err != nil {
		return nil, err
	}

	return &UserOutputDTO{
		Id:   foundUser.Id,
		Name: foundUser.Name,
	}, nil
}
