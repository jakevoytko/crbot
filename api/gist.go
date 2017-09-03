package api

// Gist is a wrapper around a simple Gist uploader. Returns the URL on success.
type Gist interface {
	Upload(contents string) (string, error)
}
