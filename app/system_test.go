package app

import (
	"testing"

	"github.com/jakevoytko/crbot/testutil"
)

func TestNewServer(t *testing.T) {
	runner := testutil.NewRunner(t)

	// Assert initial state.
	runner.AssertState()
}
