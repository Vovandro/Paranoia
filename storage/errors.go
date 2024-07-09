package storage

import "errors"

var ErrFileNotFound = errors.New("file not found")
var ErrTypeMismatch = errors.New("type mismatch")
