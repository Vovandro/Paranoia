package tests

import (
	"gitlab.com/devpro_studio/Paranoia"
	"gitlab.com/devpro_studio/Paranoia/cache"
	"testing"
)

func TestService_PushCache(t1 *testing.T) {
	s := Paranoia.New("test", nil, nil, nil)

	mockCache := cache.Memory{Name: "mock"}

	t1.Run("base push test", func(t *testing.T) {
		s.PushCache(&mockCache)

		if s.GetCache("mock") == nil {
			t.Failed()
		}
	})
}
