package synchronizer

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/zhikiri/bsync/app/config"
	"github.com/zhikiri/bsync/app/storage"
)

func SyncWith(settings config.Settings, store storage.Storage) error {
	files := make(map[string]storage.File)
	if err := getFilesFromStorage(files, store); err != nil {
		return err
	}
	log.Printf("[INFO] There are %d files in storage", len(files))
	if err := getLocalFiles(files, settings.Source); err != nil {
		return err
	}

	var wg sync.WaitGroup

	pool := make(chan request)
	makePoolOf(5, pool, &wg)

	for path, file := range files {
		if !isFileShouldBeIgnored(path, settings.Ignore) {
			pool <- newRequest(file, store, settings)
		}
	}
	close(pool)
	wg.Wait()

	return nil
}

func getFilesFromStorage(cache map[string]storage.File, store storage.Storage) error {
	files, err := store.GetFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		cache[file.Path] = file
	}
	return nil
}

func getLocalFiles(cache map[string]storage.File, source string) error {
	var relPath string
	return filepath.Walk(source, func(absPath string, osFile os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if osFile.IsDir() {
			return nil
		}
		relPath = getRelativePath(absPath, source)
		if _, exist := cache[relPath]; exist == false {
			cache[relPath] = storage.File{Path: relPath, Checksum: make([]byte, 0)}
		}
		return nil
	})
}

func getRelativePath(absPath, source string) string {
	return path.Clean(strings.Replace(absPath, source, "", -1))
}

func isFileShouldBeIgnored(path string, ignoreList []string) bool {
	ext := filepath.Ext(path)
	for _, ignore := range ignoreList {
		if ext == ignore {
			return true
		}
	}
	return false
}

func makePoolOf(size int, files chan request, wg *sync.WaitGroup) {
	for w := 1; w <= size; w++ {
		go worker(w, files, wg)
		wg.Add(1)
	}
}

func worker(id int, pool chan request, wg *sync.WaitGroup) {
	for request := range pool {
		if request.file.Path != "" {
			request.synchronize()
		}
	}
	defer wg.Done()
}
