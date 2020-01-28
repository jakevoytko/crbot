package karmalist

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/jakevoytko/crbot/feature/karma"
	"github.com/jakevoytko/crbot/feature/karmalist"
	"github.com/jakevoytko/crbot/testutil"
)

func TestKarmaList_PositiveSort(t *testing.T) {
	runner := testutil.NewRunner(t)
	var buffer bytes.Buffer
	buffer.WriteString(karmalist.MsgListKarma)
	buffer.WriteString("\n")
	buffer.WriteString("Peas: 2")
	buffer.WriteString("\n")
	buffer.WriteString("Carrots: 1")
	buffer.WriteString("\n")
	expected := buffer.String()

	karmaModelHelper := karma.NewModelHelper(runner.KarmaMap)
	karmaModelHelper.Increment("Carrots")
	karmaModelHelper.Increment("Peas")
	karmaModelHelper.Increment("Peas")
	karmalistModelHelper := karmalist.NewModelHelper(runner.KarmaMap)
	generated := karmalistModelHelper.GenerateList()

	if generated != expected {
		t.Fatalf(fmt.Sprintf("Gist failure, got `%v` expected `%v`", generated, expected))
	}
}

func TestKarmaList_MagnitudeSort(t *testing.T) {
	runner := testutil.NewRunner(t)
	var buffer bytes.Buffer
	buffer.WriteString(karmalist.MsgListKarma)
	buffer.WriteString("\n")
	buffer.WriteString("Errors: -3")
	buffer.WriteString("\n")
	buffer.WriteString("Peas: 2")
	buffer.WriteString("\n")
	buffer.WriteString("Carrots: 1")
	buffer.WriteString("\n")
	expected := buffer.String()

	karmaModelHelper := karma.NewModelHelper(runner.KarmaMap)
	karmaModelHelper.Increment("Carrots")
	karmaModelHelper.Increment("Peas")
	karmaModelHelper.Increment("Peas")
	karmaModelHelper.Decrement("Errors")
	karmaModelHelper.Decrement("Errors")
	karmaModelHelper.Decrement("Errors")
	karmalistModelHelper := karmalist.NewModelHelper(runner.KarmaMap)
	generated := karmalistModelHelper.GenerateList()

	if generated != expected {
		t.Fatalf(fmt.Sprintf("Gist failure, got `%v` expected `%v`", generated, expected))
	}
}
func TestKarmaList_NoData(t *testing.T) {
	runner := testutil.NewRunner(t)
	var buffer bytes.Buffer
	buffer.WriteString(karmalist.MsgNoKarma)
	expected := buffer.String()

	karmalistModelHelper := karmalist.NewModelHelper(runner.KarmaMap)
	generated := karmalistModelHelper.GenerateList()

	if generated != expected {
		t.Fatalf(fmt.Sprintf("Gist failure, got `%v` expected `%v`", generated, expected))
	}
}
