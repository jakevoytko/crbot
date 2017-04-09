package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// RemoteGist implements Gist, and interacts with the actual Github Gist API.
type RemoteGist struct{}

// NewRemoteGist works as advertised.
func NewRemoteGist() *RemoteGist {
	return &RemoteGist{}
}

// Upload sends the contents to the Gist API.
func (g *RemoteGist) Upload(contents string) (string, error) {
	gist := simpleGist(contents)
	serializedGist, err := json.Marshal(gist)
	if err != nil {
		info("Error marshalling gist", err)
		return "", errors.New(MsgGistSerializeFail)
	}
	response, err := http.Post(
		"https://api.github.com/gists", "application/json", bytes.NewBuffer(serializedGist))
	if err != nil {
		info("Error POSTing gist", err)
		return "", errors.New(MsgGistPostFail)
	} else if response.StatusCode != 201 {
		info("Bad status code", errors.New("Code: "+strconv.Itoa(response.StatusCode)))
		return "", errors.New(MsgGistStatusCode)
	}

	responseMap := map[string]interface{}{}
	if err := json.NewDecoder(response.Body).Decode(&responseMap); err != nil {
		info("Error reading gist response", err)
		return "", errors.New(MsgGistResponseFail)
	}

	if finalUrl, ok := responseMap["html_url"]; ok {
		return finalUrl.(string), nil
	}
	return "", errors.New(MsgGistUrlFail)
}

// GistMessage has the json mapping for the gist payload.
type GistMessage struct {
	Description string           `json:"description"`
	Public      bool             `json:"public"`
	Files       map[string]*File `json:"files"`
}

// A file represents the contents of a Gist.
type File struct {
	Content string `json:"content"`
}

// simpleGist returns a Gist object with just the given contents.
func simpleGist(contents string) *GistMessage {
	return &GistMessage{
		Public:      false,
		Description: "CRBot command list",
		Files: map[string]*File{
			"commands": &File{
				Content: contents,
			},
		},
	}
}
