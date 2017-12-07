package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	aws "github.com/aws/aws-sdk-go/service/s3"
	"github.com/zhikiri/blog-sync-cli/config"
	s3 "github.com/zhikiri/blog-sync-cli/storage"
)

func isFileExist(abspath string) bool {
	if _, err := os.Stat(abspath); err == nil {
		return true
	}
	return false
}

func isFileChanged(abspath string, che string) (bool, error) {
	calc, err := getChecksum(abspath)
	if err != nil {
		return false, err
	}
	// If file is empty
	if len(calc) == 0 {
		return false, nil
	}

	checksum := make([]byte, 16)
	if _, err := hex.Decode(checksum, bytes.TrimSpace([]byte(che))); err != nil {
		return false, err
	}

	return !bytes.Equal(checksum, calc), nil
}

func syncChangedFiles(s *config.Setup, con *aws.S3, che map[string]string) ([]string, []error) {
	var errs []error

	var ver []string
	var err error
	var changed bool

	for relpath, sum := range che {
		abspath := filepath.Join(s.SourcePath, relpath)
		// If file was deleted
		if !isFileExist(abspath) {
			s3.DeleteFile(relpath, con, s.AWS.Bucket)
			continue
		}
		// If file was changed
		if changed, err = isFileChanged(abspath, sum); err == nil && changed == true {
			s3.UpdateFile(abspath, relpath, con, s.AWS.Bucket)
		}
		if err != nil {
			errs = append(errs, err)
		}
		ver = append(ver, abspath)
	}
	return ver, errs
}

func syncAddedFiles(s *config.Setup, con *aws.S3, ver []string) []error {
	var errs []error
	filepath.Walk(s.SourcePath, func(abspath string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}
		for _, p := range ver {
			// If file was already versioned
			if p == abspath {
				return nil
			}
		}
		// If file was added
		relpath := path.Clean(strings.Replace(abspath, s.SourcePath, "", -1))
		s3.UpdateFile(abspath, relpath, con, s.AWS.Bucket)
		return nil
	})
	return errs
}

func getChecksum(abspath string) ([]byte, error) {
	if !isFileExist(abspath) {
		return make([]byte, 0), errors.New("File not found")
	}

	f, err := os.Open(abspath)
	if err != nil {
		return make([]byte, 0), errors.New("Can't open file")
	}
	defer f.Close()

	h := md5.New()
	_, err = io.Copy(h, f)

	return h.Sum(nil), err
}
