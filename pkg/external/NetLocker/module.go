package NetLocker

import (
	"context"
	grpc_client "gitlab.com/devpro_studio/Paranoia/pkg/client/grpc-client"
)

type NetLocker struct {
	NamePkg    string
	grpcClient *grpc_client.GrpcClient
	client     NetLockerServiceClient
}

func (t *NetLocker) Init(cfg map[string]interface{}) error {
	t.grpcClient = grpc_client.New(t.NamePkg)
	err := t.grpcClient.Init(cfg)
	if err != nil {
		return err
	}

	t.client = NewNetLockerServiceClient(t.grpcClient.GetClient())

	return nil
}

func (t *NetLocker) Stop() error {
	return nil
}

func (t *NetLocker) Name() string {
	return t.NamePkg
}

func (t *NetLocker) Type() string {
	return "external"
}

func (t *NetLocker) Lock(ctx context.Context, key string, timeLock int64, uniqueId *string) bool {
	res, err := t.client.TryAndLock(ctx, &NetLockRequest{Key: key, TimeLock: timeLock, UniqueId: uniqueId})
	if err != nil {
		return false
	}

	return res.GetSuccess()
}

func (t *NetLocker) Unlock(ctx context.Context, key string, uniqueId *string) bool {
	res, err := t.client.Unlock(ctx, &NetUnlockRequest{Key: key, UniqueId: uniqueId})
	if err != nil {
		return false
	}

	return res.GetSuccess()
}
