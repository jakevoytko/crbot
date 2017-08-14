package main

// ModerationFeature is a Feature that executes moderation events.
type ModerationFeature struct {
	featureRegistry *FeatureRegistry
}

// NewModerationFeature returns a new ModerationFeature.
func NewModerationFeature(featureRegistry *FeatureRegistry) *ModerationFeature {
	return &ModerationFeature{
		featureRegistry: featureRegistry,
	}
}

// Parsers returns the parsers.
func (f *ModerationFeature) Parsers() []Parser {
	return []Parser{}
}

// FallbackParser returns nil.
func (f *ModerationFeature) FallbackParser() Parser {
	return nil
}

// Executors gets the executors.
func (f *ModerationFeature) Executors() []Executor {
	return []Executor{NewRickListExecutor(f.featureRegistry)}
}

const (
	MsgRickList = "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
)

// RickListExecutor prints a rick roll.
type RickListExecutor struct {
	featureRegistry *FeatureRegistry
}

// NewRickListExecutor works as advertised.
func NewRickListExecutor(featureRegistry *FeatureRegistry) *RickListExecutor {
	return &RickListExecutor{
		featureRegistry: featureRegistry,
	}
}

// GetType returns the type.
func (e *RickListExecutor) GetType() int {
	return Type_RickList
}

// Execute replies over the given channel with a rick roll.
func (e *RickListExecutor) Execute(s DiscordSession, channel string, command *Command) {
	if _, err := s.ChannelMessageSend(channel, MsgRickList); err != nil {
		info("Failed to send ricklist message", err)
	}
}
