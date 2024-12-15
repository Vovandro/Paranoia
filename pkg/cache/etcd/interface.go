package etcd

import "errors"

type IConfigItem interface {
	GetConfigItem(typeName string, name string, dst interface{}) error
}

var ErrKeyNotFound = errors.New("key not found")
var ErrTypeMismatch = errors.New("type mismatch")
