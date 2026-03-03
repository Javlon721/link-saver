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

func (service LinkService) RegisterLink(ctx context.Context, params *types.RegisterLink) (*types.Link, error) {
	user, err := service.userStore.GetUser(ctx, params.UserID)

	if err != nil {
		return nil, err
	}

	link, err := service.linkStore.Register(ctx, user.ID, params.Link, params.Desctibtion)

	if err != nil {
		return nil, err
	}

	return link, nil
}

func (service LinkService) GetAll(ctx context.Context, userID int64) ([]*types.Link, error) {
	user, err := service.userStore.GetUser(ctx, userID)

	if err != nil {
		return []*types.Link{}, err
	}

	links := service.linkStore.GetAll(ctx, user.ID)

	if len(links) == 0 {
		return []*types.Link{}, errs.ErrLinksNotFound
	}

	return links, nil
}
