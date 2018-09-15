package vote

import (
	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/feature"
	"github.com/jakevoytko/crbot/model"
	stringmap "github.com/jakevoytko/go-stringmap"
)

// Feature registers feature-specific things for moderation.
type Feature struct {
	featureRegistry *feature.Registry
	modelHelper     *ModelHelper
	commandChannel  chan<- *model.Command
	utcTimer        model.UTCTimer
	utcClock        model.UTCClock
}

// NewFeature returns a new Feature.
func NewFeature(featureRegistry *feature.Registry, voteMap stringmap.StringMap, clock model.UTCClock, timer model.UTCTimer, commandChannel chan<- *model.Command) *Feature {
	return &Feature{
		featureRegistry: featureRegistry,
		modelHelper:     NewModelHelper(voteMap, clock),
		utcTimer:        timer,
		utcClock:        clock,
		commandChannel:  commandChannel,
	}
}

// Parsers returns the parsers.
func (f *Feature) Parsers() []feature.Parser {
	return []feature.Parser{
		NewStatusParser(),
		NewStartVoteParser(),
		NewBallotParser(model.CommandNameVoteInFavorF1, true /* inFavor */),
		NewBallotParser(model.CommandNameVoteInFavorYes, true /* inFavor */),
		NewBallotParser(model.CommandNameVoteAgainstF2, false /* inFavor */),
		NewBallotParser(model.CommandNameVoteAgainstNo, false /* inFavor */),
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
		NewBallotExecutor(f.modelHelper),
		NewConcludeExecutor(f.modelHelper),
		NewStatusExecutor(f.modelHelper),
		NewStartVoteExecutor(f.modelHelper, f.commandChannel, f.utcTimer),
	}
}

// OnInitialLoad cleans up from any votes that were already active when crbot shut down.
func (f *Feature) OnInitialLoad(s api.DiscordSession) error {
	return HandleVotesOnInitialLoad(s, f.modelHelper, f.utcClock, f.utcTimer, f.commandChannel)
}
