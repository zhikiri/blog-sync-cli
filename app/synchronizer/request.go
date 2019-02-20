package synchronizer

import (
	"bytes"
	"crypto/md5"
	"github.com/zhikiri/blog-sync-cli/app/config"
	"github.com/zhikiri/blog-sync-cli/app/storage"
	"io"
	"log"
	"os"
	"path"
)

type request struct {
	file        storage.File
	absFilePath string
	store       storage.Storage
}

func newRequest(file storage.File, store storage.Storage, settings config.Settings) request {
	return request{
		file:        file,
		store:       store,
		absFilePath: path.Join(settings.Source, file.Path),
	}
}

func (r *request) isFileExists() (bool, error) {
	_, err := os.Stat(r.absFilePath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (r *request) getChecksum() ([]byte, error) {
	file, err := os.Open(r.absFilePath)
	if err != nil {
		return make([]byte, 0), err
	}
	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)

	return hash.Sum(nil), err
}

func (r *request) isFileChanged() (bool, error) {
	checksum, err := r.getChecksum()
	if err != nil || len(checksum) == 0 {
		return false, err
	}
	return !bytes.Equal(checksum, r.file.Checksum), nil
}

func (r *request) syncDelete() bool {
	exists, err := r.isFileExists()
	if err != nil {
		log.Printf("[ERROR] Cannot check file existence, %+v", err)
	}
	if !exists {
		// If file is not exists locally it means that it should be deleted from the storage
		r.store.DelFile(r.file)
		return true
	}
	return false
}

func (r *request) syncUpdate() bool {
	changed, err := r.isFileChanged()
	if err != nil {
		log.Printf("[ERROR] Cannot detect changes of the file, %+v", err)
	}
	if changed {
		r.store.PutFile(r.file, r.absFilePath)
		return true
	}
	return false
}

func (r *request) syncCreate() bool {
	if len(r.file.Checksum) > 0 {
		// If checksum is not calculated it means that file is new
		return false
	}
	if err := r.store.PutFile(r.file, r.absFilePath); err != nil {
		log.Println(err)
	}
	return true
}

func (r *request) synchronize() {
	if r.syncDelete() {
		return
	}
	if r.syncCreate() {
		return
	}
	if r.syncUpdate() {
		return
	}
}
