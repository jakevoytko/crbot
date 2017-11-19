package testutil

import "errors"

// InMemoryGist is a fake for the Gist interface.
type InMemoryGist struct {
	Messages []string
	FailNext bool
}

// NewInMemoryGist works as advertised.
func NewInMemoryGist() *InMemoryGist {
	return &InMemoryGist{
		Messages: []string{},
		FailNext: false,
	}
}

const (
	GistSuccessUrl = "https://www.example.com/success"
)

// Upload stores the message, or returns an error if FailNext is set. Resets FailNext.
func (g *InMemoryGist) Upload(content string) (string, error) {
	if g.FailNext {
		g.FailNext = false
		return "", errors.New("Gist upload failed")
	}

	g.Messages = append(g.Messages, content)
	return GistSuccessUrl, nil
}
