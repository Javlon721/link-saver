package errs

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrLinkNotProvided   = errors.New("link not provided")
	ErrLinksNotFound     = errors.New("links not found")
)
