package NetLocker

import (
	"context"
	"errors"

	"gitlab.com/devpro_studio/go_utils/decode"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NetLocker struct {
	NamePkg string
	conn    *grpc.ClientConn
	client  NetLockerServiceClient
}

type config struct {
	Url string `yaml:"url"`
}

func New(name string) *NetLocker {
	return &NetLocker{
		NamePkg: name,
	}
}

func (t *NetLocker) Init(cfg map[string]interface{}) error {
	var c config
	if err := decode.Decode(cfg, &c, "yaml", decode.DecoderStrongFoundDst); err != nil {
		return err
	}
	if c.Url == "" {
		return errors.New("url is required")
	}

	conn, err := grpc.Dial(
		c.Url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return err
	}

	t.conn = conn
	t.client = NewNetLockerServiceClient(conn)

	return nil
}

func (t *NetLocker) Stop() error {
	if t.conn != nil {
		return t.conn.Close()
	}
	return nil
}

func (t *NetLocker) Name() string {
	return t.NamePkg
}

func (t *NetLocker) Type() string {
	return "external"
}

func (t *NetLocker) Lock(ctx context.Context, key string, timeLock int64, uniqueId *string) (bool, error) {
	res, err := t.client.TryAndLock(ctx, &NetLockRequest{Key: key, TimeLock: timeLock, UniqueId: uniqueId})
	if err != nil {
		return false, err
	}

	return res.GetSuccess(), nil
}

func (t *NetLocker) Unlock(ctx context.Context, key string, uniqueId *string) bool {
	res, err := t.client.Unlock(ctx, &NetUnlockRequest{Key: key, UniqueId: uniqueId})
	if err != nil {
		return false
	}

	return res.GetSuccess()
}
