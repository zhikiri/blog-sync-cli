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
		Type: getMIMEType(abspath),
		Path: relpath,
	}, nil
}

func getMIMEType(abspath string) string {
	if strings.HasSuffix(abspath, ".css") {
		return "text/css"
	} else if strings.HasSuffix(abspath, ".html") {
		return "text/html"
	} else if strings.HasSuffix(abspath, ".js") {
		return "application/javascript"
	} else if strings.HasSuffix(abspath, ".png") {
		return "image/png"
	} else if strings.HasSuffix(abspath, ".svg") {
		return "image/svg+xml"
	} else {
		return "text/plain"
	}
}
