package NetLocker

import "context"

type INetLocker interface {
	Lock(ctx context.Context, key string, timeLock int64, uniqueId *string) (bool, error)
	Unlock(ctx context.Context, key string, uniqueId *string) bool
}
