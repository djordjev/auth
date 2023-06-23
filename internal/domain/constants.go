package domain

import "errors"

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUserNotExist = errors.New("user not exists")
var ErrInvalidToken = errors.New("invalid token")
