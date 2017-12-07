package storage

import (
	"bytes"
	"net/http"
	"os"
)

// Object contain uploading file details
type Object struct {
	Body *bytes.Reader
	Type string
	Size int64
	Path string
}

func getObject(abspath, relpath string) (*Object, error) {
	file, err := os.Open(abspath)
	defer file.Close()

	if err != nil {
		return &Object{}, err
	}

	fileInfo, _ := file.Stat()
	size := fileInfo.Size()

	buffer := make([]byte, size)

	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)

	return &Object{
		Body: fileBytes,
		Size: size,
		Type: fileType,
		Path: relpath,
	}, nil
}
