package interfaces

type IService interface {
	Init() error
	Stop() error
	GetLogger() ILogger
	GetConfig() IConfig
	SetMetrics(c IMetrics)
	PushCache(c ICache) IService
	GetCache(key string) ICache
	PushDatabase(c IDatabase) IService
	GetDatabase(key string) IDatabase
	PushNoSql(c INoSql) IService
	GetNoSql(key string) INoSql
	PushController(c IController) IService
	GetController(key string) IController
	PushModule(c IModules) IService
	GetModule(key string) IModules
	PushRepository(c IRepository) IService
	GetRepository(key string) IRepository
	PushTask(c ITask) IService
	GetTask(key string) ITask
	RemoveTask(key string)
	PushServer(c IServer) IService
	GetServer(key string) IServer
	PushClient(c IClient) IService
	GetClient(key string) IClient
	PushStorage(c IStorage) IService
	GetStorage(key string) IStorage
	PushMiddleware(c IMiddleware) IService
	GetMiddleware(key string) IMiddleware
}
