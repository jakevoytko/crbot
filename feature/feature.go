package feature

// Feature encapsulates all of the behavior necessary for a built-in
// feature.
type Feature interface {
	// Returns all parsers associated with this feature.
	Parsers() []Parser
	// FallbackParser returns the parser to execute if no other parser is
	// recognized by name. There can only be one system-wide.
	FallbackParser() Parser
	// CommandInterceptors returns the command interceptors for this feature.
	CommandInterceptors() []CommandInterceptor
	// Returns all executors associated with this feature.
	Executors() []Executor
}
