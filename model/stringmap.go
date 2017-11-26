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
	// Scan finds all keys that match the given pattern. It uses *-style wildcard
	// matching. It is not guaranteed to find all of the keys if keys are being
	// added and removed during the search in a separate thread.
	ScanKeys(pattern string) ([]string, error)
}
