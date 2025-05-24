package vote

import (
	"github.com/jakevoytko/crbot/feature"
	"github.com/jakevoytko/crbot/model"
)

// Feature registers feature-specific things for moderation.
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
	return []feature.Parser{
		NewStatusParser(),
		NewDeprecatedVoteParser(),
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
		NewDeprecatedVoteExecutor(model.CommandTypeVote),
		NewDeprecatedVoteExecutor(model.CommandTypeVoteStatus),
	}
}
