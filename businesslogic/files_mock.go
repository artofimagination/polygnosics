package businesslogic

import (
	"io"
	"os"

	dbModels "github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/polygnosics/rest"
)

type FileProcessorMock struct {
	CallOrder         []string
	UpdatedAsset      dbModels.Asset
	GeneratePathError error
	WriteToFileError  error
	UploadFileError   error
	RemoveFileError   error
}

func (f *FileProcessorMock) GeneratePath(asset *dbModels.Asset) error {
	f.CallOrder = append(f.CallOrder, "GeneratePath")
	asset.ID = f.UpdatedAsset.ID
	asset.DataMap = f.UpdatedAsset.DataMap
	return f.GeneratePathError
}
func (f *FileProcessorMock) WriteToFile(filename string, data string) (err error) {
	f.CallOrder = append(f.CallOrder, "WriteToFile")
	return f.WriteToFileError
}

func (f *FileProcessorMock) UploadFile(key string, fileName string, r rest.RequestInterface) (err error) {
	f.CallOrder = append(f.CallOrder, "UploadFile")
	return f.UploadFileError
}

func (f *FileProcessorMock) RemoveFile(path string) error {
	f.CallOrder = append(f.CallOrder, "RemoveFile")
	return f.RemoveFileError
}

func (f *FileProcessorMock) GenerateID() string {
	return "a4d28c75-595b-4059-801f-2a9ad127916b"
}

// IoImpl is the production implementation of io module functions.
type IoMock struct {
	CopyWritten     int64
	CopyError       error
	WriteBytesCount int
	WriteError      error
}

func (i IoMock) Copy(dst io.Writer, src io.Reader) (int64, error) {
	return i.CopyWritten, i.CopyError
}

func (i IoMock) WriteString(w io.Writer, s string) (n int, err error) {
	return i.WriteBytesCount, i.WriteError
}

// OsImpl is the production implementation of io module functions
type OsMock struct {
	CreateFile     FileImpl
	CreateError    error
	RemoveAllError error
	MkdirAllError  error
}

func (o OsMock) Create(name string) (*FileImpl, error) {
	return &o.CreateFile, o.CreateError
}

func (o OsMock) RemoveAll(path string) error {
	return o.RemoveAllError
}

func (o OsMock) MkdirAll(path string, perm os.FileMode) error {
	return o.MkdirAllError
}

type FileMock struct {
	CloseError error
	WriteError error
	SyncError  error
	Written    int
}

func (f FileMock) Close() error {
	return f.CloseError
}

func (f FileMock) Write(b []byte) (n int, err error) {
	return f.Written, f.WriteError
}

func (f FileMock) Sync() error {
	return f.SyncError
}
