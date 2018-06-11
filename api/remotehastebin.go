package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
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
	MsgHastebinPostFail     = "Unable to connect to Hastebin service. Give it a few minutes and try again"
	MsgHastebinResponseFail = "Failure reading response from Hastebin service"
	MsgHastebinStatusCode   = "Failed to upload Hastebin :("
	MsgHastebinUrlFail      = "Failed getting url from Hastebin service"
)

func (g *RemoteHastebin) Upload(contents string) (string, error) {
	response, err := http.Post(
		"https://hastebin.com/documents", "application/x-www-form-urlencoded", bytes.NewBufferString(contents))
	if err != nil {
		log.Info("Error POSTing Hastebin", err)
		return "", errors.New(MsgHastebinPostFail)
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Info("Bad status code", errors.New("Code: "+strconv.Itoa(response.StatusCode)))
		body, _ := ioutil.ReadAll(response.Body)
		log.Info("Response body: ", errors.New(string(body)))
		return "", errors.New(MsgHastebinStatusCode)
	}

	responseMap := map[string]interface{}{}
	if err := json.NewDecoder(response.Body).Decode(&responseMap); err != nil {
		log.Info("Error reading Hastebin response", err)
		return "", errors.New(MsgHastebinResponseFail)
	}

	if finalUrl, ok := responseMap["key"]; ok {
		return "https://hastebin.com/raw/" + finalUrl.(string), nil
	}

	return "", errors.New(MsgHastebinUrlFail)
}
