package ui

import (
	"github.com/zj1244/syft/internal/log"
	"github.com/zj1244/syft/internal/ui/common"
	syftEvent "github.com/zj1244/syft/syft/event"
	"github.com/wagoodman/go-partybus"
)

// LoggerUI is a UI function that leverages the displays all application logs to the screen.
func LoggerUI(workerErrs <-chan error, subscription *partybus.Subscription) error {
	events := subscription.Events()
eventLoop:
	for {
		select {
		case err := <-workerErrs:
			if err != nil {
				return err
			}
		case e, ok := <-events:
			if !ok {
				// event bus closed...
				break eventLoop
			}

			// ignore all events except for the final event
			if e.Type == syftEvent.CatalogerFinished {
				err := common.CatalogerFinishedHandler(e)
				if err != nil {
					log.Errorf("unable to show catalog image finished event: %+v", err)
				}

				// this is the last expected event
				break eventLoop
			}
		}
	}

	return nil
}
