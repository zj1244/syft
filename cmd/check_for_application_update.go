package cmd

import (
	"github.com/wagoodman/go-partybus"
	"github.com/zj1244/syft/internal"
	"github.com/zj1244/syft/internal/bus"
	"github.com/zj1244/syft/internal/log"
	"github.com/zj1244/syft/internal/version"
	"github.com/zj1244/syft/syft/event"
)

func checkForApplicationUpdate() {
	if appConfig.CheckForAppUpdate {
		isAvailable, newVersion, err := version.IsUpdateAvailable()
		if err != nil {
			// this should never stop the application
			log.Errorf(err.Error())
		}
		if isAvailable {
			log.Infof("new version of %s is available: %s (current version is %s)", internal.ApplicationName, newVersion, version.FromBuild().Version)

			bus.Publish(partybus.Event{
				Type:  event.AppUpdateAvailable,
				Value: newVersion,
			})
		} else {
			log.Debugf("no new %s update available", internal.ApplicationName)
		}
	}
}
