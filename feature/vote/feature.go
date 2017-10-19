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
		NewBallotParser(model.Name_VoteInFavorF1, true /* inFavor */),
		NewBallotParser(model.Name_VoteInFavorYes, true /* inFavor */),
		NewBallotParser(model.Name_VoteAgainstF2, false /* inFavor */),
		NewBallotParser(model.Name_VoteAgainstNo, false /* inFavor */),
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
		NewBallotExecutor(f.modelHelper),
	}
}
