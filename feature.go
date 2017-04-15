package main

import (
	"errors"
	"fmt"
)

// Feature encapsulates all of the behavior necessary for a built-in
// feature.
// TODO(jvoytko): Refactor to provide parsers and command executors.
type Feature interface {
	GetType() int
	// Returns all parsers associated with this feature.
	Parsers() []Parser
	// FallbackParser returns the parser to execute if no other parser is
	// recognized by name. There can only be one system-wide.
	FallbackParser() Parser
	// Execute the given command, for the session and channel name provided.
	Execute(DiscordSession, string, *Command)
}

// Parsers is used to multiplex on builtin ?* commands, and ensure that the
// commands are correctly formatted.
type Parser interface {
	// The user-facing name of the command. Must be unique.
	GetName() string
	// Parses the given split command line.
	Parse([]string) (*Command, error)
	// HelpText returns user-facing help text for the message.
	HelpText() string
}

// FeatureRegistry stores all of the features.
type FeatureRegistry struct {
	// FallbackFeature is the feature that should attempt to parse the command
	// line if no named feature matches.
	FallbackParser Parser

	nameToParser          map[string]Parser
	typeToFeature         map[int]Feature
	invokableFeatureNames []string
}

// NewFeatureRegistry works as advertised.
func NewFeatureRegistry() *FeatureRegistry {
	return &FeatureRegistry{
		nameToParser:          map[string]Parser{},
		typeToFeature:         map[int]Feature{},
		invokableFeatureNames: []string{},
	}
}

// Register attempts to register the given feature by name. If a feature with
// the given name exists, then error.
func (r *FeatureRegistry) Register(feature Feature) error {
	// Register regular parsers.
	for _, parser := range feature.Parsers() {
		if _, ok := r.nameToParser[parser.GetName()]; ok {
			return errors.New(fmt.Sprintf("Duplicate parser: %v", parser.GetName()))
		}
		if len(parser.GetName()) > 0 {
			r.nameToParser[parser.GetName()] = parser
			r.invokableFeatureNames = append(r.invokableFeatureNames, parser.GetName())
		}
	}

	// Register fallback parser.
	if fallback := feature.FallbackParser(); fallback != nil {
		if r.FallbackParser != nil {
			return errors.New("More than one fallback parser found")
		}
		r.FallbackParser = feature.FallbackParser()
	}

	// Set up internals.
	if _, ok := r.typeToFeature[feature.GetType()]; ok {
		return errors.New(fmt.Sprintf("Duplicate type: %v", feature.GetType()))
	}
	r.typeToFeature[feature.GetType()] = feature
	return nil
}

// GetParserByName returns the feature with the given name, or null if no such
// feature exists.
func (r *FeatureRegistry) GetParserByName(name string) Parser {
	if f, ok := r.nameToParser[name]; ok {
		return f
	}
	if f, ok := r.nameToParser["?"+name]; ok {
		return f
	}
	return nil
}

// GetFeatureByType returns the feature with the given type, or null if no
// feature handles this type.
func (r *FeatureRegistry) GetFeatureByType(featureType int) Feature {
	if f, ok := r.typeToFeature[featureType]; ok {
		return f
	}
	return nil
}

// IsInvokable tests that the given string is a user-invokable command. Will
// pass whether the string is prefixed by a ? or not.
func (r *FeatureRegistry) IsInvokable(name string) bool {
	_, ok1 := r.nameToParser[name]
	_, ok2 := r.nameToParser["?"+name]
	return ok1 || ok2
}

// Gets the invokable feature names.
func (r *FeatureRegistry) GetInvokableFeatureNames() []string {
	return r.invokableFeatureNames
}
