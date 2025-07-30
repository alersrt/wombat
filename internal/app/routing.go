package app

import (
	"fmt"
	"log/slog"
)

func (receiver *Application) route() {
	for update := range receiver.sourceChan {
		if !receiver.tagsRegex.MatchString(update.Content) {
			slog.Info(fmt.Sprintf("Tag not found"))
			return
		}
		receiver.targetChan <- update
	}
}
