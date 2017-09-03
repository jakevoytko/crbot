package model

// StringMap stores key/value string pairs. It is always synchronous, but may be
// stored outside the memory space of the program. For instance, in Redis.
type StringMap interface {
	// Has returns whether or not key is present.
	Has(key string) (bool, error)
	// Get returns the given key. Error if key is not present.
	Get(key string) (string, error)
	// Set sets the given key. Allowed to overwrite.
	Set(key, value string) error
	// Delete deletes the given key. Error if key is not present.
	Delete(key string) error
	// GetAll returns every entry as a map.
	GetAll() (map[string]string, error)
}
