package vote

import (
	"github.com/jakevoytko/crbot/feature"
	"github.com/jakevoytko/crbot/model"
)

// Feature registers feature-specific things for moderation.
type Feature struct {
	featureRegistry *feature.Registry
	modelHelper     *ModelHelper
}

// NewFeature returns a new Feature.
func NewFeature(featureRegistry *feature.Registry, voteMap model.StringMap, clock model.UTCClock) *Feature {
	return &Feature{
		featureRegistry: featureRegistry,
		modelHelper:     NewModelHelper(voteMap, clock),
	}
}

// Parsers returns the parsers.
func (f *Feature) Parsers() []feature.Parser {
	return []feature.Parser{
		NewStatusParser(),
		NewVoteParser(),
	}
}

// CommandInterceptors returns command interceptors.
func (f *Feature) CommandInterceptors() []feature.CommandInterceptor {
	return []feature.CommandInterceptor{}
}

// FallbackParser returns nil.
func (f *Feature) FallbackParser() feature.Parser {
	return nil
}

// Executors gets the executors.
func (f *Feature) Executors() []feature.Executor {
	return []feature.Executor{
		NewStatusExecutor(f.modelHelper),
		NewVoteExecutor(f.modelHelper),
	}
}
