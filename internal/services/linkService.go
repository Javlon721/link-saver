package services

import (
	"context"

	"github.com/Javlon721/link-saver/internal/errs"
	"github.com/Javlon721/link-saver/internal/types"
)

type LinkService struct {
	linkStore types.LinkStore
}

func NewLinkService(linkStore types.LinkStore) *LinkService {
	return &LinkService{linkStore: linkStore}
}

func (service *LinkService) RegisterLink(ctx context.Context, userID int64, params *types.LinkInfo) (*types.Link, error) {
	link, err := service.linkStore.AddLink(ctx, userID, params.Link, params.Desctibtion)

	if err != nil {
		return nil, err
	}

	return link, nil
}

func (service *LinkService) GetAll(ctx context.Context, userID int64) ([]*types.Link, error) {
	links := service.linkStore.GetAll(ctx, userID)

	if len(links) == 0 {
		return []*types.Link{}, errs.ErrLinksNotFound
	}

	return links, nil
}

func (service *LinkService) DeleteLink(ctx context.Context, userID, linkID int64) error {
	return service.linkStore.DeleteLink(ctx, userID, linkID)
}

func (service *LinkService) DeleteUserLinks(ctx context.Context, userID int64) error {
	return service.linkStore.DeleteUserLinks(ctx, userID)
}

func (service *LinkService) NewWithTx(db types.DB) *LinkService {
	linkStore := service.linkStore.NewWithTx(db)

	return &LinkService{
		linkStore: linkStore,
	}
}
