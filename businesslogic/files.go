package businesslogic

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/polygnosics/rest"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// IoInterface represents an interface for any io module function.
// Required in order to allow
// implementations: IoImpl
//									IoMock
type IoInterface interface {
	Copy(dst io.Writer, src io.Reader) (int64, error)
	WriteString(w io.Writer, s string) (n int, err error)
}

// IoImpl is the production implementation of io module functions.
type IoImpl struct {
}

func (IoImpl) Copy(dst io.Writer, src io.Reader) (int64, error) {
	return io.Copy(dst, src)
}

func (IoImpl) WriteString(w io.Writer, s string) (n int, err error) {
	return io.WriteString(w, s)
}

// OsInterface represents an interface for any os module function.
// Required in order to allow
// implementations: OsImpl
//									OsMock
type OsInterface interface {
	Create(name string) (*FileImpl, error)
	RemoveAll(path string) error
	MkdirAll(path string, perm os.FileMode) error
}

// OsImpl is the production implementation of io module functions
type OsImpl struct {
}

func (OsImpl) Create(name string) (*FileImpl, error) {
	file, err := os.Create(name)
	return &FileImpl{file}, err
}

func (OsImpl) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (OsImpl) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

// File is the interface class to redefine os.File with custom, i.e. mock implementations
type File interface {
	Close() error
	Write(b []byte) (n int, err error)
	Sync() error
}

type FileImpl struct {
	File
}

type FileProcessor interface {
	GeneratePath(asset *models.Asset) error
	WriteToFile(filename string, data string) (err error)
	UploadFile(key string, fileName string, r rest.RequestInterface) (err error)
	RemoveFile(path string) error
	GenerateID() string
}

type FileProcessorImpl struct {
	FileIO IoInterface
	OsFunc OsInterface
}

// Default file paths
const (
	DefaultUserAvatarPath    = "/assets/images/avatar.jpg"
	DefaultProductAvatarPath = "/assets/images/avatar.jpg"
	DefaultProjectAvatarPath = "/assets/images/avatar.jpg"
)

var splitRegexp = regexp.MustCompile(`(\S{4})`)

func (f *FileProcessorImpl) GeneratePath(asset *models.Asset) error {
	assetIDString := strings.Replace(asset.ID.String(), "-", "", -1)
	assetStringSplit := splitRegexp.FindAllString(assetIDString, -1)
	assetPath := path.Join(assetStringSplit...)
	rootPath := os.Getenv("USER_STORE_DOCKER")
	assetPath = path.Join(rootPath, assetPath)
	if err := f.OsFunc.MkdirAll(assetPath, os.ModePerm); err != nil {
		return err
	}
	asset.DataMap[models.BaseAssetPath] = assetPath
	return nil
}

// WriteToFile dumps the data string in the file defined by filename
func (f *FileProcessorImpl) WriteToFile(filename string, data string) (err error) {
	file, err := f.OsFunc.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if errClose := file.Close(); errClose != nil {
			err = errClose
		}
	}()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return file.Sync()
}

// Remove file deletes the file defined by path.
func (f *FileProcessorImpl) RemoveFile(path string) error {
	return f.OsFunc.RemoveAll(path)
}

func (f *FileProcessorImpl) GenerateID() string {
	return uuid.New().String()
}

// UploadFile writes the multipart file in the request to the disk.
func (f *FileProcessorImpl) UploadFile(key string, fileName string, r rest.RequestInterface) (err error) {
	file, handler, err := r.FormFile(key)
	if err == http.ErrMissingFile {
		return nil
	}

	if err != nil {
		return err
	}

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("Stored Filename: %+v\n", fileName)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	defer func() {
		if errClose := file.Close(); errClose != nil {
			err = errClose
		}
	}()

	// Create file
	dst, err := f.OsFunc.Create(fileName)
	if err != nil {
		return err
	}

	// Copy the uploaded file to the created file on the file system.
	if _, err := f.FileIO.Copy(dst, file); err != nil {
		if errClose := f.OsFunc.RemoveAll(fileName); errClose != nil {
			err = errors.Wrap(errors.WithStack(err), errClose.Error())
		}
		return err
	}

	if err := dst.Close(); err != nil {
		return err
	}

	return err
}

func removeFolder(dir string) error {
	if err := removeContents(dir); err != nil {
		return err
	}
	err := os.Remove(dir)
	if err != nil {
		return err
	}
	return nil
}

func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
