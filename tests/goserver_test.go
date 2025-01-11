package tests

import (
	"gitlab.com/devpro_studio/Paranoia/cache/memory"
	"gitlab.com/devpro_studio/Paranoia/framework"
	"testing"
)

func TestEngine_PushCache(t1 *testing.T) {
	s := framework.New("test", nil, nil)

	mockCache := memory.Memory{Name: "mock"}

	t1.Run("base push test", func(t *testing.T) {
		s.PushCache(&mockCache)

		if s.GetCache("mock") == nil {
			t.Failed()
		}
	})
}
