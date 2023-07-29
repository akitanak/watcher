package app

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

// logEvent logs event.
func logEvent(event fsnotify.Event) {
	log.Printf("%s was %s", event.Name, convertToPastPerticiple(event.Op))
}

// convertToPastPerticiple converts fsnotify.Op to past perticiple.
func convertToPastPerticiple(op fsnotify.Op) string {
	switch op {
	case fsnotify.Create:
		return "created"
	case fsnotify.Write:
		return "written"
	case fsnotify.Remove:
		return "removed"
	default:
		return ""
	}
}
