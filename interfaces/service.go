package interfaces

type IService interface {
	GetLogger() ILogger
	PushCache(c ICache)
	GetCache(key string) ICache
}
