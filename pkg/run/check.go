package run

import (
	"github.com/elek/klepif/pkg/plugins"
)

func Check() error {

	ctx, ghSource, e := Init()
	if e != nil {
		return e
	}

	events, err := ghSource.GetEventsSinceLastCheck();
	if err != nil {
		panic(err)
	}
	for _, plg := range plugins.Plugins.Instances {
		for _, event := range events {
			err = plg.HandlePrEvent(&ctx, event);
			if err != nil {
				return err;
			}
		}

	}

	return nil;
}
