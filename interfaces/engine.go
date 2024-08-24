package interfaces

type IEngine interface {
	Init() error
	Stop() error
	GetLogger() ILogger
	GetConfig() IConfig
	SetMetrics(c IMetrics)
	PushCache(c ICache) IEngine
	GetCache(key string) ICache
	PushDatabase(c IDatabase) IEngine
	GetDatabase(key string) IDatabase
	PushNoSql(c INoSql) IEngine
	GetNoSql(key string) INoSql
	PushController(c IController) IEngine
	GetController(key string) IController
	PushModule(c IModules) IEngine
	GetModule(key string) IModules
	PushRepository(c IRepository) IEngine
	GetRepository(key string) IRepository
	PushService(c IService) IEngine
	GetService(key string) IService
	PushTask(c ITask) IEngine
	GetTask(key string) ITask
	RemoveTask(key string)
	RunTask(key string, args map[string]interface{}) error
	PushServer(c IServer) IEngine
	GetServer(key string) IServer
	PushClient(c IClient) IEngine
	GetClient(key string) IClient
	PushStorage(c IStorage) IEngine
	GetStorage(key string) IStorage
	PushMiddleware(c IMiddleware) IEngine
	GetMiddleware(key string) IMiddleware
}
