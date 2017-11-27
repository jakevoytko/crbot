package karma

import (
	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/feature"
	"github.com/jakevoytko/crbot/model"
)

// Feature allows crbot to record changes in karma
type Feature struct {
	featureRegistry *feature.Registry
	karmaMap        model.StringMap
}

// NewFeature returns a new Feature.
func NewFeature(featureRegistry *feature.Registry, karmaMap model.StringMap) *Feature {
	return &Feature{
		featureRegistry: featureRegistry,
		karmaMap:        karmaMap,
	}
}

// Parsers gets the learn feature parsers.
func (f *Feature) Parsers() []feature.Parser {
	return []feature.Parser{
		NewKarmaParser(model.Name_KarmaIncrement, true),
		NewKarmaParser(model.Name_KarmaDecrement, false),
	}
}

// CommandInterceptors returns nothing.
func (f *Feature) CommandInterceptors() []feature.CommandInterceptor {
	return []feature.CommandInterceptor{}
}

// FallbackParser returns nil
func (f *Feature) FallbackParser() feature.Parser {
	return nil
}

// Executors gets the executors.
func (f *Feature) Executors() []feature.Executor {
	return []feature.Executor{NewKarmaExecutor(f.karmaMap)}
}

// OnInitialLoad does nothing.
func (f *Feature) OnInitialLoad(s api.DiscordSession) error { return nil }