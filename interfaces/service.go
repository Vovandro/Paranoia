package interfaces

type IService interface {
	GetLogger() ILogger
	PushCache(c ICache) IService
	GetCache(key string) ICache
}
