package fsw

import (
	"context"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/ikugo-dev/DeFence/logger"
)

func InitWatch(directoryName string) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := watchDirectory(ctx, directoryName); err != nil {
			log.Printf("Watching folder %s has stopped with error %s\n", directoryName, err)
		}
	}()

	return cancel
}

func watchDirectory(ctx context.Context, directoryName string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()
	go watch(ctx, watcher)

	if _, err := os.Stat(directoryName); os.IsNotExist(err) {
		os.Mkdir(directoryName, 0755)
	}
	if err := watcher.Add(directoryName); err != nil {
		return err
	}

	log.Printf("Watching folder %s has started\n", directoryName)
	<-ctx.Done()
	log.Printf("Watching folder %s has stopped\n", directoryName)
	return ctx.Err()
}

func watch(ctx context.Context, watcher *fsnotify.Watcher) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			handleFileEvent(event)

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("watcher error:", err)

		case <-ctx.Done():
			return
		}
	}
}

func handleFileEvent(event fsnotify.Event) {
	logger.Log("file event: %s", event.String())
	// TODO
}
