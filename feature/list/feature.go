package list

import (
	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/feature"
	stringmap "github.com/jakevoytko/go-stringmap"
)

// Feature is a Feature that lists commands that are available.
type Feature struct {
	featureRegistry *feature.Registry
	commandMap      stringmap.StringMap
	gist            api.Gist
}

// NewFeature returns a new Feature.
func NewFeature(
	featureRegistry *feature.Registry,
	commandMap stringmap.StringMap,
	gist api.Gist) *Feature {

	return &Feature{
		featureRegistry: featureRegistry,
		commandMap:      commandMap,
		gist:            gist,
	}
}

// Parsers returns the parsers.
func (f *Feature) Parsers() []feature.Parser {
	return []feature.Parser{NewListParser()}
}

// CommandInterceptors returns nothing.
func (f *Feature) CommandInterceptors() []feature.CommandInterceptor {
	return []feature.CommandInterceptor{}
}

// FallbackParser returns nil.
func (f *Feature) FallbackParser() feature.Parser {
	return nil
}

func (f *Feature) Executors() []feature.Executor {
	return []feature.Executor{NewListExecutor(f.featureRegistry, f.commandMap, f.gist)}
}

// OnInitialLoad does nothing.
func (f *Feature) OnInitialLoad(s api.DiscordSession) error { return nil }
