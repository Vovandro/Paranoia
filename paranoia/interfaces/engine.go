package interfaces

type TypePkg string

const (
	PkgCache      TypePkg = "cache"
	PkgDatabase   TypePkg = "database"
	PkgClient     TypePkg = "client"
	PkgServer     TypePkg = "server"
	PkgStorage    TypePkg = "storage"
	PkgMiddleware TypePkg = "middleware"
	PkgLogger     TypePkg = "logger"
)

type TypeModule string

const (
	ModuleController TypeModule = "controller"
	ModuleRepository TypeModule = "repository"
	ModuleService    TypeModule = "service"
	ModuleMiddleware TypeModule = "middleware"
	ModuleModule     TypeModule = "module"
	ModuleCustom     TypeModule = "custom"
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
