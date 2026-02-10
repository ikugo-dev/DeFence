package fsw

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/ikugo-dev/DeFence/internal/algorithms"
	"github.com/ikugo-dev/DeFence/internal/logger"
)

var encryptionKey []byte
var encryptionAlgorithm string

func InitWatch(dir string, key []byte, algorithm string) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	encryptionKey = key
	encryptionAlgorithm = algorithm

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Log("Failed to create watcher: %s", err)
		cancel()
		return func() {}
	}

	if err := watcher.Add(dir); err != nil {
		logger.Log("Failed to watch directory: %s", err)
		cancel()
		return func() {}
	}

	go func() {
		defer watcher.Close()
		watch(ctx, watcher)
	}()

	return cancel
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
	if !event.Has(fsnotify.Create) {
		return
	}
	logger.Log("File event: %s", event.String())

	data, err := algorithms.EncryptFile(event.Name, encryptionKey, encryptionAlgorithm)
	if err != nil {
		logger.Log("Failed to encrypt file %s: %s", event.Name, err)
		return
	}

	outPath := filepath.Join("./X", filepath.Base(event.Name))
	if err := os.WriteFile(outPath, data, 0644); err != nil {
		logger.Log("Failed to write encrypted file %s: %s", outPath, err)
	}
	logger.Log("Encrypted %s -> %s (Algorithm: %s)", filepath.Base(event.Name), filepath.Base(outPath), encryptionAlgorithm)
}
