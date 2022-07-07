package common

import (
	"fmt"
	"os"

	"github.com/wagoodman/go-partybus"
	syftEventParsers "github.com/zj1244/syft/syft/event/parsers"
)

// CatalogerFinishedHandler is a UI function for processing the CatalogerFinished bus event, displaying the catalog
// via the given presenter to stdout.
func CatalogerFinishedHandler(event partybus.Event) error {
	// show the report to stdout
	pres, err := syftEventParsers.ParseCatalogerFinished(event)
	if err != nil {
		return fmt.Errorf("bad CatalogerFinished event: %w", err)
	}

	if err := pres.Present(os.Stdout); err != nil {
		return fmt.Errorf("unable to show package catalog report: %w", err)
	}
	return nil
}
