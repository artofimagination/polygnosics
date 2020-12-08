package restcontrollers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"polygnosics/app/restcontrollers/auth"

	"github.com/pkg/errors"
)

const (
	DefaultAvatarPath = "/assets/images/avatar.jpg"
)

const (
	UserAvatar     = "user_avatar"
	UserBackground = "user_background"

	ProductDescription = "product_description"
)

func getContent() map[string]interface{} {
	p := make(map[string]interface{})
	p["assets"] = make(map[string]interface{})
	path := auth.UserData.Assets.GetImagePath(UserAvatar, DefaultAvatarPath)
	p["assets"].(map[string]interface{})[UserAvatar] = path
	p["texts"] = make(map[string]interface{})
	p["texts"].(map[string]interface{})["avatar-upload"] = "Upload your avatar"
	p["texts"].(map[string]interface{})["username"] = auth.UserData.Name

	return p
}

func uploadFile(destination string, r *http.Request) error {
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
