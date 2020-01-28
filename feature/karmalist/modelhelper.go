package karmalist

import (
	"bytes"
	"math"
	"sort"
	"strconv"

	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/log"
	stringmap "github.com/jakevoytko/go-stringmap"
)

// Model helper add functions that abstract transformations to the Karma
// before it is sent to the Gist api. Currently the only transformation is
// sorting by magnitude. Future enhancements could be only showing hated or
// loved things. Karma would need an overhaul to list users vs other since users
// are stored by Discord username with no designation
type ModelHelper struct {
	StringMap stringmap.StringMap
	gist      api.Gist
}

// NewModelHelper works as advertised.
// I have no idea why I'm bothering to take arguments or return anything here.
func NewModelHelper(stringMap stringmap.StringMap, gist api.Gist) *ModelHelper {
	return &ModelHelper{
		StringMap: stringMap,
		gist:      gist,
	}
}

func (h *ModelHelper) GetGistUrl(karmaMap stringmap.StringMap) (string, error) {
	all, err := karmaMap.GetAll()
	if err != nil {
		log.Fatal("Error reading all the karma", err)
	}

	type sortableKarma struct {
		displayKarma string
		absKarma     int
	}
	var karmaStore []sortableKarma

	// Sort karma by absolute value so that stronger feelings are at the top of the list
	for k, v := range all {
		displayKarma := k + ": " + v
		floatKarma, _ := strconv.ParseFloat(v, 32)
		absKarma := int(math.Abs(floatKarma))
		karmaStore = append(karmaStore, sortableKarma{displayKarma, absKarma})
	}

	sort.Slice(karmaStore, func(i, j int) bool {
		return karmaStore[i].absKarma > karmaStore[j].absKarma
	})

	var buffer bytes.Buffer
	buffer.WriteString(MsgListKarma)
	buffer.WriteString("\n")
	for _, kv := range karmaStore {
		buffer.WriteString(kv.displayKarma)
		buffer.WriteString("\n")
	}
	return h.gist.Upload(buffer.String())
}
