package services

import (
	"context"

	"github.com/Javlon721/link-saver/internal/errs"
	"github.com/Javlon721/link-saver/internal/types"
)

type LinkService struct {
	linkStore types.LinkStore
	userStore types.UserStore
}

func NewLinkService(linkStore types.LinkStore, userStore types.UserStore) *LinkService {
	return &LinkService{linkStore: linkStore, userStore: userStore}
}

func (service LinkService) RegisterLink(ctx context.Context, telegramID int64, params *types.LinkInfo) (*types.Link, error) {
	user, err := service.userStore.GetUser(ctx, telegramID)

	if err != nil {
		return nil, err
	}

	link, err := service.linkStore.AddLink(ctx, user.ID, params.Link, params.Desctibtion)

	if err != nil {
		return nil, err
	}

	return link, nil
}

func (service LinkService) GetAll(ctx context.Context, telegramID int64) ([]*types.Link, error) {
	user, err := service.userStore.GetUser(ctx, telegramID)

	if err != nil {
		return []*types.Link{}, err
	}

	links := service.linkStore.GetAll(ctx, user.ID)

	if len(links) == 0 {
		return []*types.Link{}, errs.ErrLinksNotFound
	}

	return links, nil
}

func (service LinkService) DeleteLink(ctx context.Context, telegramID, linkID int64) error {
	user, err := service.userStore.GetUser(ctx, telegramID)

	if err != nil {
		return err
	}

	return service.linkStore.DeleteLink(ctx, user.ID, linkID)
}

func (service LinkService) DeleteUserLinks(ctx context.Context, telegramID int64) error {
	user, err := service.userStore.GetUser(ctx, telegramID)

	if err != nil {
		return err
	}

	return service.linkStore.DeleteUserLinks(ctx, user.ID)
}

func (service LinkService) NewWithTx(db types.DB) *LinkService {
	userStore := service.userStore.NewWithTx(db)
	linkStore := service.linkStore.NewWithTx(db)

	return &LinkService{
		userStore: userStore,
		linkStore: linkStore,
	}
}
