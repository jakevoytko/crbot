package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/jakevoytko/crbot/log"
)

// RemoteHastebin implements Gist, and interacts with the Hastebin API.
type RemoteHastebin struct{}

// NewRemoteHastebin works as advertised.
func NewRemoteHastebin() *RemoteHastebin {
	return &RemoteHastebin{}
}

const (
	msgHastebinPostFail     = "unable to connect to Hastebin service. Give it a few minutes and try again"
	msgHastebinResponseFail = "failure reading response from Hastebin service"
	msgHastebinStatusCode   = "failed to upload Hastebin :("
	msgHastebinURLFail      = "failed getting url from Hastebin service"
)

// Upload uploads the given string to hastebin and returns the URL of the hastebin on success.
func (g *RemoteHastebin) Upload(contents string) (string, error) {
	response, err := http.Post(
		"https://hastebin.com/documents", "application/x-www-form-urlencoded", bytes.NewBufferString(contents))
	if err != nil {
		log.Info("Error POSTing Hastebin", err)
		return "", errors.New(msgHastebinPostFail)
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Info("Bad status code", errors.New("Code: "+strconv.Itoa(response.StatusCode)))
		body, _ := io.ReadAll(response.Body)
		log.Info("Response body: ", errors.New(string(body)))
		return "", errors.New(msgHastebinStatusCode)
	}

	responseMap := map[string]interface{}{}
	if err := json.NewDecoder(response.Body).Decode(&responseMap); err != nil {
		log.Info("Error reading Hastebin response", err)
		return "", errors.New(msgHastebinResponseFail)
	}

	if finalURL, ok := responseMap["key"]; ok {
		return "https://hastebin.com/raw/" + finalURL.(string), nil
	}

	return "", errors.New(msgHastebinURLFail)
}
