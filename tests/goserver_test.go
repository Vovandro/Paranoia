package tests

import (
	"goServer"
	"goServer/cache"
	"testing"
)

func TestService_PushCache(t1 *testing.T) {
	s := goServer.Service{}
	s.New("test", nil, nil)

	mockCache := cache.Mock{Name: "mock"}

	t1.Run("base push test", func(t *testing.T) {
		s.PushCache(&mockCache)

		if s.GetCache("mock") == nil {
			t.Failed()
		}
	})
}
