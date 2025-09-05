package FeatureChaos

import "sync"

// Mock implements IFeatureChaos with hook and in-memory flags map
type Mock struct {
	CheckFunc func(featureName string, seed string, attr map[string]string) bool

	// Seedable states: per feature -> bool or optional per-seed map
	Flags map[string]bool

	NamePkg string

	mu    sync.Mutex
	Calls []struct{ Feature, Seed string }
}

func (m *Mock) record(feature, seed string) {
	m.mu.Lock()
	m.Calls = append(m.Calls, struct{ Feature, Seed string }{Feature: feature, Seed: seed})
	m.mu.Unlock()
}

func (t *Mock) Init(cfg map[string]interface{}) error {
	return nil
}

func (t *Mock) Stop() error {
	return nil
}

func (t *Mock) Name() string {
	return t.NamePkg
}

func (t *Mock) Type() string {
	return "external"
}

func (m *Mock) Check(featureName string, seed string, attr map[string]string) bool {
	m.record(featureName, seed)
	if m.CheckFunc != nil {
		return m.CheckFunc(featureName, seed, attr)
	}
	if m.Flags == nil {
		return false
	}
	if v, ok := m.Flags[featureName]; ok {
		return v
	}
	return false
}

var _ IFeatureChaos = (*Mock)(nil)
