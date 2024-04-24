# Paranoia framework

[![pipeline status](https://gitlab.com/devpro_studio/Paranoia/badges/master/pipeline.svg)](https://gitlab.com/devpro_studio/Paranoia/-/commits/master) 
[![coverage report](https://gitlab.com/devpro_studio/Paranoia/badges/master/coverage.svg)](https://gitlab.com/devpro_studio/Paranoia/-/commits/master) 
[![Latest Release](https://gitlab.com/devpro_studio/Paranoia/-/badges/release.svg)](https://gitlab.com/devpro_studio/Paranoia/-/releases)


## Simple start:
Import to project `go get https://gitlab.com/devpro_studio/Paranoia.git`

add to main.go

```
	s := Paranoia.
		New("test", &config.Env{}, &logger.File{&logger.Std{}}).
		PushCache(&cache.Memory{}).
		PushRepository(&myRepository{}).
		PushController(&myController{}).
		Init()
		
	s.Run()
	
	defer s.Stop()
```