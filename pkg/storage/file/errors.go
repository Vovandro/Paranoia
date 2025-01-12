package file

import "errors"

var ErrFileNotFound = errors.New("file not found")
var ErrTypeMismatch = errors.New("type mismatch")
var ErrNotSupported = errors.New("operation not supported")
