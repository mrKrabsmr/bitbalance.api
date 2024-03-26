package services

import "errors"

var (
    ErrNotFound = errors.New("resource was not found")
    ErrDuplicate = errors.New("resource already exists")
)
