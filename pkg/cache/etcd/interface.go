package etcd

import "errors"

type IConfig interface {
	GetMapInterface(key string, def map[string]interface{}) map[string]interface{}
}

var ErrKeyNotFound = errors.New("key not found")
var ErrTypeMismatch = errors.New("type mismatch")
