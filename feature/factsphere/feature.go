package factsphere

import (
	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/feature"
)

// Feature is a Feature that prints a random fact to the user.
type Feature struct {
	featureRegistry *feature.Registry
}

// NewFeature returns a new Feature.
func NewFeature(featureRegistry *feature.Registry) *Feature {
	return &Feature{
		featureRegistry: featureRegistry,
	}
}

// Parsers returns the parsers.
func (f *Feature) Parsers() []feature.Parser {
	return []feature.Parser{NewFactSphereParser()}
}

// CommandInterceptors returns nothing.
func (f *Feature) CommandInterceptors() []feature.CommandInterceptor {
	return []feature.CommandInterceptor{}
}

// FallbackParser returns nil.
func (f *Feature) FallbackParser() feature.Parser {
	return nil
}

// Executors gets the executors.
func (f *Feature) Executors() []feature.Executor {
	return []feature.Executor{NewFactSphereExecutor()}
}

// OnInitialLoad does nothing.
func (f *Feature) OnInitialLoad(s api.DiscordSession) error { return nil }
