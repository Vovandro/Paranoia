package FeatureChaos

import (
	"context"
	"errors"
	"time"

	fc_sdk_go "gitlab.com/devpro_studio/FeatureChaos/sdk/fc_sdk_go"
	"gitlab.com/devpro_studio/go_utils/decode"
)

type FeatureChaos struct {
	NamePkg  string
	cfg      FeatureChaosConfig
	fcClient *fc_sdk_go.Client
}

func New(name string) *FeatureChaos {
	return &FeatureChaos{
		NamePkg: name,
	}
}

type FeatureChaosConfig struct {
	Host        string `yaml:"host"`
	ServiceName string `yaml:"service_name"`
}

func (t *FeatureChaos) Init(cfg map[string]interface{}) error {
	var err error
	err = decode.Decode(cfg, &t.cfg, "yaml", decode.DecoderStrongFoundDst)

	if err != nil {
		return err
	}

	if t.cfg.Host == "" {
		return errors.New("host is required")
	}

	if t.cfg.ServiceName == "" {
		return errors.New("service_name is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.fcClient, err = fc_sdk_go.New(ctx, t.cfg.Host, t.cfg.ServiceName, fc_sdk_go.Options{
		AutoSendStats: true,
	})

	return err
}

func (t *FeatureChaos) Stop() error {
	return t.fcClient.Close()
}

func (t *FeatureChaos) Name() string {
	return t.NamePkg
}

func (t *FeatureChaos) Type() string {
	return "external"
}

func (t *FeatureChaos) Check(featureName string, seed string, attr map[string]string) bool {
	return t.fcClient.IsEnabled(featureName, seed, attr)
}
