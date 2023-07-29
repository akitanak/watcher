package app

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func logEvent(event fsnotify.Event) {
	log.Printf("%s is %s", event.Name, event.Op)
}
