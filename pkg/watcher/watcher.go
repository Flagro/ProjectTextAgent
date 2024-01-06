package watcher

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	fsWatcher    *fsnotify.Watcher
	FileModified chan string
	FileCreated  chan string
	FileDeleted  chan string
}

// NewWatcher creates and initializes a new Watcher instance.
func NewWatcher() (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		fsWatcher:    fsWatcher,
		FileModified: make(chan string),
		FileCreated:  make(chan string),
		FileDeleted:  make(chan string),
	}, nil
}

// WatchDirectory starts watching the specified directory and its subdirectories.
func (w *Watcher) WatchDirectory(rootPath string) error {
	err := w.addDirectory(rootPath)
	if err != nil {
		return err
	}

	go w.watchEvents()
	return nil
}

// addDirectory adds the specified directory and its subdirectories to the watcher.
func (w *Watcher) addDirectory(path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsDir() {
			return w.fsWatcher.Add(path)
		}
		return nil
	})
}

// watchEvents handles the file system events and sends notifications.
func (w *Watcher) watchEvents() {
	for {
		select {
		case event, ok := <-w.fsWatcher.Events:
			if !ok {
				return
			}
			w.handleEvent(event)
		case err, ok := <-w.fsWatcher.Errors:
			if !ok {
				return
			}
			log.Println("Watcher error:", err)
		}
	}
}

// handleEvent processes the received fsnotify event.
func (w *Watcher) handleEvent(event fsnotify.Event) {
	switch {
	case event.Op&fsnotify.Write == fsnotify.Write:
		w.FileModified <- event.Name
	case event.Op&fsnotify.Create == fsnotify.Create:
		fi, err := os.Stat(event.Name)
		if err != nil {
			log.Println("Error stating file:", err)
			return
		}
		if fi.Mode().IsDir() {
			w.addDirectory(event.Name)
		} else {
			w.FileCreated <- event.Name
		}
	case event.Op&fsnotify.Remove == fsnotify.Remove:
		w.FileDeleted <- event.Name
	}
}
