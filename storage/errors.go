package storage

import "errors"

var FileNotFound = errors.New("file not found")
var TypeMismatch = errors.New("type mismatch")
