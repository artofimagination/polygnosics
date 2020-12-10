package page

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

func UploadFile(destination string, r *http.Request) error {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return err
	}

	file, handler, err := r.FormFile("asset")
	if err != nil {
		return err
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create file
	dst, err := os.Create(destination)
	if err != nil {
		if err2 := dst.Close(); err2 != nil {
			err = errors.Wrap(errors.WithStack(err), err2.Error())
		}
		return err
	}

	// Copy the uploaded file to the created file on the file system.
	if _, err := io.Copy(dst, file); err != nil {
		if err2 := dst.Close(); err2 != nil {
			err = errors.Wrap(errors.WithStack(err), err2.Error())
		}
		return err
	}

	return nil
}
