package services

import (
	"context"
	"errors"

	"github.com/Javlon721/link-saver/internal/errs"
	"github.com/Javlon721/link-saver/internal/types"
)

type UserService struct {
	userStore types.UserStore
}

func NewUserService(userStore types.UserStore) *UserService {
	return &UserService{userStore: userStore}
}

func (service UserService) RegisterUser(ctx context.Context, params *types.RegisterUser) (*types.User, error) {
	_, err := service.userStore.GetUser(ctx, params.TelegramID)

	if err == nil {
		return nil, errs.ErrUserAlreadyExists
	}

	if !errors.Is(err, errs.ErrUserNotFound) {
		return nil, err
	}

	user, err := service.userStore.AddUser(context.Background(), params)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service UserService) GetUser(ctx context.Context, userID int64) (*types.User, error) {
	return service.userStore.GetUser(ctx, userID)
}

func (service UserService) DeleteUser(ctx context.Context, userID int64) error {
	user, err := service.userStore.GetUser(ctx, userID)

	if err != nil {
		return err
	}

	return service.userStore.DeleteUser(ctx, user.ID)
}

func (service UserService) NewWithTx(db types.DB) *UserService {
	userStore := service.userStore.NewWithTx(db)

	return &UserService{
		userStore: userStore,
	}
}
