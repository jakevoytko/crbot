package karma

import (
	"fmt"
	"log"
	"strconv"

	stringmap "github.com/jakevoytko/go-stringmap"
)

// ModelHelper provides helpers for working with karma storage.
type ModelHelper struct {
	karmaMap stringmap.StringMap
}

// NewModelHelper works as advertised.
func NewModelHelper(karmaMap stringmap.StringMap) *ModelHelper {
	return &ModelHelper{karmaMap: karmaMap}
}

// Increment either adds 1 to storage in the key, or increments the key in storage if it exists.
func (h *ModelHelper) Increment(target string) (int, error) {
	return h.process(target, true /* increment */)
}

// Decrement either subtracts 1 from storage in the key, or decrements the key in
// storage if it exists.
func (h *ModelHelper) Decrement(target string) (int, error) {
	return h.process(target, false /* increment */)
}

// Process handles both Increment and Decrement.
func (h *ModelHelper) process(target string, increment bool) (int, error) {
	// Get the current value of karma (if it exists) and increment/decrement it
	currentKarma := 0
	has, err := h.karmaMap.Has(target)
	if err != nil {
		log.Fatal("Couldn't check if target has karma", err)
	}
	if has {
		currentKarmaStr, err := h.karmaMap.Get(target)
		if err != nil {
			return 0, fmt.Errorf("%s: %f", "Couldn't get target's current karma", err)
		}
		currentKarma, err = strconv.Atoi(currentKarmaStr)
		if err != nil {
			return 0, fmt.Errorf("%s: %f", "Invalid karma value", err)
		}
	}

	var newKarma int
	if increment {
		newKarma = currentKarma + 1
	} else {
		newKarma = currentKarma - 1
	}

	err = h.karmaMap.Set(target, strconv.Itoa(newKarma))
	if err != nil {
		return 0, fmt.Errorf("%s: %f", "Error storing new karma value", err)
	}

	return newKarma, nil
}
