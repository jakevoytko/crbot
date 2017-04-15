package main

import (
	"errors"
	"fmt"
)

// Feature encapsulates all of the behavior necessary for a built-in
// feature.
type Feature interface {
	// Returns all parsers associated with this feature.
	Parsers() []Parser
	// FallbackParser returns the parser to execute if no other parser is
	// recognized by name. There can only be one system-wide.
	FallbackParser() Parser

	// Returns all executors associated with this feature.
	Executors() []Executor
}

// Parsers is used to multiplex on builtin ?* commands, and ensure that the
// commands are correctly formatted.
type Parser interface {
	// The user-facing name of the command. Must be unique.
	GetName() string
	// Parses the given split command line.
	Parse([]string) (*Command, error)
	// The user-facing help text for the given name.
	HelpText() string
}

// Executor actually performs the given action based on the given command.
type Executor interface {
	// The command type to execute.
	GetType() int
	// Execute the given command, for the session and channel name provided.
	Execute(DiscordSession, string, *Command)
}

// FeatureRegistry stores all of the features.
type FeatureRegistry struct {
	// FallbackFeature is the feature that should attempt to parse the command
	// line if no named feature matches.
	FallbackParser Parser

	nameToParser          map[string]Parser
	typeToExecutor        map[int]Executor
	invokableFeatureNames []string
}

// NewFeatureRegistry works as advertised.
func NewFeatureRegistry() *FeatureRegistry {
	return &FeatureRegistry{
		nameToParser:          map[string]Parser{},
		typeToExecutor:        map[int]Executor{},
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

	// Register executors.
	for _, executor := range feature.Executors() {
		if _, ok := r.typeToExecutor[executor.GetType()]; ok {
			return errors.New(fmt.Sprintf("Duplicate executor: %v", executor.GetType()))
		}
		r.typeToExecutor[executor.GetType()] = executor
	}

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

// GetExecutorByType returns the executor with the given type, or null if no
// feature handles this type.
func (r *FeatureRegistry) GetExecutorByType(commandType int) Executor {
	if f, ok := r.typeToExecutor[commandType]; ok {
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
