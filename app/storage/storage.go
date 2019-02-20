package storage

import (
	"bytes"
	"os"
	"strings"
)

type Storage interface {
	GetFiles() ([]File, error)
	PutFile(file File, path string) error
	DelFile(file File) error
}

type File struct {
	Path     string
	Checksum []byte
	Body     *bytes.Reader
	Size     int64
	Type     string
}

func (f *File) loadFileDataFrom(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	info, _ := file.Stat()
	f.Size = info.Size()

	buffer := make([]byte, f.Size)
	file.Read(buffer)

	f.Body = bytes.NewReader(buffer)
	f.Type = getMIMEType(path)
	return nil
}

func getMIMEType(path string) string {
	if strings.HasSuffix(path, ".css") {
		return "text/css"
	} else if strings.HasSuffix(path, ".html") {
		return "text/html"
	} else if strings.HasSuffix(path, ".js") {
		return "application/javascript"
	} else if strings.HasSuffix(path, ".png") {
		return "image/png"
	} else if strings.HasSuffix(path, ".svg") {
		return "image/svg+xml"
	} else {
		return "text/plain"
	}
}
