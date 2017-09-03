package moderation

import (
	"github.com/jakevoytko/crbot/app"
	"github.com/jakevoytko/crbot/feature"
)

// Feature registers feature-specific things for moderation.
type Feature struct {
	featureRegistry *feature.Registry
	config          *app.Config
}

// NewFeature returns a new Feature.
func NewFeature(featureRegistry *feature.Registry, config *app.Config) *Feature {
	return &Feature{
		featureRegistry: featureRegistry,
		config:          config,
	}
}

// Parsers returns the parsers.
func (f *Feature) Parsers() []feature.Parser {
	return []feature.Parser{}
}

// CommandInterceptors returns command interceptors.
func (f *Feature) CommandInterceptors() []feature.CommandInterceptor {
	return []feature.CommandInterceptor{
		NewRickListCommandInterceptor(f.config),
	}
}

// FallbackParser returns nil.
func (f *Feature) FallbackParser() feature.Parser {
	return nil
}

// Executors gets the executors.
func (f *Feature) Executors() []feature.Executor {
	return []feature.Executor{NewRickListExecutor(f.featureRegistry)}
}
