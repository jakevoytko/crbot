package karma

import (
	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/feature"
	"github.com/jakevoytko/crbot/model"
	stringmap "github.com/jakevoytko/go-stringmap"
)

// Feature allows crbot to record changes in karma
type Feature struct {
	featureRegistry *feature.Registry
	karmaMap        stringmap.StringMap
}

// NewFeature returns a new Feature.
func NewFeature(featureRegistry *feature.Registry, karmaMap stringmap.StringMap) *Feature {
	return &Feature{
		featureRegistry: featureRegistry,
		karmaMap:        karmaMap,
	}
}

// Parsers gets the learn feature parsers.
func (f *Feature) Parsers() []feature.Parser {
	return []feature.Parser{
		NewParser(model.CommandNameKarmaIncrement, true /* increment */),
		NewParser(model.CommandNameKarmaDecrement, false /* increment */),
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
	return []feature.Executor{NewExecutor(f.karmaMap)}
}

// OnInitialLoad does nothing.
func (f *Feature) OnInitialLoad(s api.DiscordSession) error { return nil }
