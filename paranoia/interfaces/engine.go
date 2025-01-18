package interfaces

const (
	PkgCache      = "cache"
	PkgDatabase   = "database"
	PkgClient     = "client"
	PkgServer     = "server"
	PkgStorage    = "storage"
	PkgMiddleware = "middleware"
	PkgLogger     = "logger"
)

const (
	ModuleController = "controller"
	ModuleRepository = "repository"
	ModuleService    = "service"
	ModuleMiddleware = "middleware"
	ModuleModule     = "module"
	ModuleCustom     = "custom"
)

type IEngine interface {
	Init() error
	Stop() error
	GetLogger() ILogger
	GetConfig() IConfig
	SetMetrics(c IMetrics)
	SetTrace(c ITrace)

	PushPkg(c IPkg) IEngine
	GetPkg(typePkg string, key string) IPkg

	PushModule(c IModules) IEngine
	GetModule(typePkg string, key string) IModules

	PushTask(c ITask) IEngine
	GetTask(key string) ITask
	RemoveTask(key string)
	RunTask(key string, args map[string]interface{}) error
}
