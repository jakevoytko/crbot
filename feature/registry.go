package feature

import (
	"errors"
	"fmt"
)

// Registry stores all of the features.
type Registry struct {
	// FallbackFeature is the feature that should attempt to parse the command
	// line if no named feature matches.
	FallbackParser Parser

	nameToParser          map[string]Parser
	interceptors          []CommandInterceptor
	typeToExecutor        map[int]Executor
	invokableFeatureNames []string
}

// NewRegistry works as advertised.
func NewRegistry() *Registry {
	return &Registry{
		nameToParser:          map[string]Parser{},
		interceptors:          []CommandInterceptor{},
		typeToExecutor:        map[int]Executor{},
		invokableFeatureNames: []string{},
	}
}

// Register attempts to register the given feature by name. If a feature with
// the given name exists, then error.
func (r *Registry) Register(feature Feature) error {
	// Register regular parsers.
	for _, parser := range feature.Parsers() {
		if _, ok := r.nameToParser[parser.GetName()]; ok {
			return fmt.Errorf("duplicate parser: %v", parser.GetName())
		}
		if len(parser.GetName()) > 0 {
			r.nameToParser[parser.GetName()] = parser
			r.invokableFeatureNames = append(r.invokableFeatureNames, parser.GetName())
		}
	}

	// Register command interceptors.
	r.interceptors = append(r.interceptors, feature.CommandInterceptors()...)

	// Register fallback parser.
	if fallback := feature.FallbackParser(); fallback != nil {
		if r.FallbackParser != nil {
			return errors.New("more than one fallback parser found")
		}
		r.FallbackParser = feature.FallbackParser()
	}

	// Register executors.
	for _, executor := range feature.Executors() {
		if _, ok := r.typeToExecutor[executor.GetType()]; ok {
			return fmt.Errorf("duplicate executor: %v", executor.GetType())
		}
		r.typeToExecutor[executor.GetType()] = executor
	}

	return nil
}

// GetParserByName returns the feature with the given name, or null if no such
// feature exists.
func (r *Registry) GetParserByName(name string) Parser {
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
func (r *Registry) GetExecutorByType(commandType int) Executor {
	if f, ok := r.typeToExecutor[commandType]; ok {
		return f
	}
	return nil
}

// IsInvokable tests that the given string is a user-invokable command. Will
// pass whether the string is prefixed by a ? or not.
func (r *Registry) IsInvokable(name string) bool {
	_, ok1 := r.nameToParser[name]
	_, ok2 := r.nameToParser["?"+name]
	return ok1 || ok2
}

// GetInvokableFeatureNames returns invokable feature names.
func (r *Registry) GetInvokableFeatureNames() []string {
	return r.invokableFeatureNames
}

// CommandInterceptors return the execution interceptors.
func (r *Registry) CommandInterceptors() []CommandInterceptor {
	return r.interceptors
}
