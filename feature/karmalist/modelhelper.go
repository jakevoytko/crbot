package karmalist

import (
	"bytes"
	"math"
	"sort"
	"strconv"

	"github.com/jakevoytko/crbot/log"
	stringmap "github.com/jakevoytko/go-stringmap"
)

// ModelHelper adds functions that abstract transformations to the Karma
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
	// MsgKarmaMapFailed is an error string for when the Karma map fails
	MsgKarmaMapFailed = "error reading karma map"
	// MsgNoKarma is an error string for when karma hasn't been stored yet.
	MsgNoKarma = "nothing has accumulated karma"
)

// GenerateList returns a string of all the user:karma pairs in the map sorted by
// magnitude. If there is no karma, it returns an error string.
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

	// Sort by karma and then alphabetically so that the same ordering is returned
	// regardless of the ordering returned from redis
	sort.SliceStable(karmaStore, func(i, j int) bool {
		if karmaStore[i].absKarma > karmaStore[j].absKarma {
			return true
		} else if karmaStore[i].absKarma < karmaStore[j].absKarma {
			return false
		}

		return karmaStore[i].displayKarma > karmaStore[j].displayKarma
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
