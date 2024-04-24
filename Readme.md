# Paranoia framework

[![pipeline status](https://gitlab.com/devpro_studio/Paranoia/badges/master/pipeline.svg)](https://gitlab.com/devpro_studio/Paranoia/-/commits/master) 
[![coverage report](https://gitlab.com/devpro_studio/Paranoia/badges/master/coverage.svg)](https://gitlab.com/devpro_studio/Paranoia/-/commits/master) 
[![Latest Release](https://gitlab.com/devpro_studio/Paranoia/-/badges/release.svg)](https://gitlab.com/devpro_studio/Paranoia/-/releases)


## Simple start:
Import to project `go get https://gitlab.com/devpro_studio/Paranoia.git`

add to main.go

```
	s := Paranoia.
		New("base paranoia app", &config.Env{}, &logger.File{&logger.Std{}}).
		PushCache(&cache.Memory{Name: "cache"}).
		PushRepository(&myRepository{Name: "repository"}).
		PushController(&myController{Name: "controller"}).
		Init()
		
	s.Run()
	
	defer s.Stop()
```