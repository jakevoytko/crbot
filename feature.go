package main

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Feature encapsulates all of the behavior necessary for a built-in
// feature.
// TODO(jvoytko): Refactor to provide parsers and command executors.
type Feature interface {
	GetName() string
	GetType() int
	// Returns whether a user can execute this command by typing the name.
	Invokable() bool
	// Parses the given split command line.
	Parse([]string) (*Command, error)
	// Execute the given command, for the session and channel name provided.
	Execute(*discordgo.Session, string, *Command)
}

// FeatureRegistry stores all of the features.
type FeatureRegistry struct {
	// FallbackFeature is the feature that should attempt to parse the command
	// line if no named feature matches.
	FallbackFeature Feature

	nameToFeature         map[string]Feature
	typeToFeature         map[int]Feature
	nameToType            map[string]int
	invokableFeatureNames []string
}

// NewFeatureRegistry works as advertised.
func NewFeatureRegistry() *FeatureRegistry {
	return &FeatureRegistry{
		nameToFeature:         map[string]Feature{},
		typeToFeature:         map[int]Feature{},
		nameToType:            map[string]int{},
		invokableFeatureNames: []string{},
	}
}

// Register attempts to register the given feature by name. If a feature with
// the given name exists, then error.
func (r *FeatureRegistry) Register(feature Feature) error {
	if feature.Invokable() {
		if _, ok := r.nameToFeature[feature.GetName()]; ok {
			return errors.New(fmt.Sprintf("Duplicate feature: %v", feature.GetName()))
		}
		r.nameToFeature[feature.GetName()] = feature
		r.invokableFeatureNames = append(r.invokableFeatureNames, feature.GetName())
	}

	if _, ok := r.typeToFeature[feature.GetType()]; ok {
		return errors.New(fmt.Sprintf("Duplicate type: %v", feature.GetType()))
	}
	r.typeToFeature[feature.GetType()] = feature
	if feature.Invokable() {
		r.nameToType[feature.GetName()] = feature.GetType()
	}
	return nil
}

// GetFeatureByName returns the feature with the given name, or null if no such
// feature exists.
func (r *FeatureRegistry) GetFeatureByName(name string) Feature {
	if f, ok := r.nameToFeature[name]; ok {
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

// For the given feature name, returns the corresponding handled command
// type. Handles missing and present leading ?. If it's unrecognized, returns
// Type_None.
func (r *FeatureRegistry) GetTypeFromName(name string) int {
	if t, ok := r.nameToType[name]; ok {
		return t
	}
	if t, ok := r.nameToType["?"+name]; ok {
		return t
	}
	return Type_None
}

// Gets the invokable feature names.
func (r *FeatureRegistry) GetInvokableFeatureNames() []string {
	return r.invokableFeatureNames
}
