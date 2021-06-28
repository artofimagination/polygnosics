package tests

import (
	"mime/multipart"
	"strings"
	"testing"

	"github.com/artofimagination/polygnosics/rest"
	"github.com/kr/pretty"
)

type OrderedTests struct {
	TestDataSet DataSet
	OrderedList OrderedTestList
}

type DataSet map[string]*Data
type OrderedTestList []string

type Data struct {
	Expected interface{}
	Input    interface{}
	Mock     interface{}
}

type MultipartFileMock struct {
	CloseError error
	ReadError  error
	ReadBytes  int
}

func (f MultipartFileMock) Close() error {
	return f.CloseError
}

func (f MultipartFileMock) Read(b []byte) (n int, err error) {
	return f.ReadBytes, f.ReadError
}

type RequestMock struct {
	FormFileError           error
	ParseFormError          error
	ParseMultiPartFormError error
	File                    MultipartFileMock
	FormValues              map[string]string
	FileName                string
}

func (r *RequestMock) FormFile(key string) (rest.MultiPartFileImpl, *multipart.FileHeader, error) {
	return rest.MultiPartFileImpl{MultiPartFile: r.File}, &multipart.FileHeader{
		Filename: r.FileName,
	}, r.FormFileError
}

func (r *RequestMock) ParseForm() error {
	return r.ParseFormError
}

func (r *RequestMock) FormValue(key string) string {
	return r.FormValues[key]
}

func (r *RequestMock) ParseMultipartForm(maxMemory int64) error {
	return r.ParseMultiPartFormError
}

var TestResultString = "\n%s test failed.\n\nReturned:\n%+v\n\nExpected:\n%+v\nDiff:\n%+v"
var TestPartialErrorMatchString = "\n%s test failed.\n\nReturned:\n%+v\n\nExpected:\n%+v"

func CheckErrPartialMatch(errA error, errB error, testCaseString string, t *testing.T) {
	if errA == nil && errB != nil {
		t.Errorf(TestPartialErrorMatchString, testCaseString, errA, errB)
	}

	if errA != nil && errB == nil {
		t.Errorf(TestPartialErrorMatchString, testCaseString, errA, errB)
	}

	if errA != nil && errB != nil && !strings.Contains(errA.Error(), errB.Error()) {
		t.Errorf(TestPartialErrorMatchString, testCaseString, errA, errB)
		return
	}
}

func CheckResult(outputA interface{}, outputB interface{}, errA interface{}, errB interface{}, testCaseString string, t *testing.T) {
	if diff := pretty.Diff(outputA, outputB); len(diff) != 0 {
		t.Errorf(TestResultString, testCaseString, outputA, outputB, diff)
		return
	}

	if diff := pretty.Diff(errA, errB); len(diff) != 0 {
		t.Errorf(TestResultString, testCaseString, errA, errB, diff)
		return
	}
}
