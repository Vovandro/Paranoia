package FeatureChaos

type IFeatureChaos interface {
	Check(featureName string, seed string, attr map[string]string) bool
}
