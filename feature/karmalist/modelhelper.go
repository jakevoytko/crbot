package karmalist

import (
	"bytes"
	"math"
	"sort"
	"strconv"

	"github.com/jakevoytko/crbot/log"
	stringmap "github.com/jakevoytko/go-stringmap"
)

// Model helper add functions that abstract transformations to the Karma
// before it is sent to the Gist api. Currently the only transformation is
// sorting by magnitude. Future enhancements could be only showing hated or
// loved things. Karma would need an overhaul to list users vs other since users
// are stored by plain text Discord username with no differentiation
type ModelHelper struct {
	KarmaMap stringmap.StringMap
}

// NewModelHelper works as advertised.
func NewModelHelper(stringMap stringmap.StringMap) *ModelHelper {
	return &ModelHelper{
		KarmaMap: stringMap,
	}
}

const (
	// Error string for when the Karma map fails
	MsgKarmaMapFailed = "Error reading all the karma"
	// Error string for when karma hasn't been stored yet.
	MsgNoKarma = "Nothing has accumulated karma."
)

func (h *ModelHelper) GenerateList() string {
	all, err := h.KarmaMap.GetAll()
	if err != nil {
		log.Fatal(MsgKarmaMapFailed, err)
	}
	if len(all) <= 0 {
		return MsgNoKarma
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

	return buffer.String()
}
