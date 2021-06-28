package businesslogic

import (
	"errors"
	"net/http"
	"os"
	"testing"

	dbModels "github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/polygnosics/tests"
	"github.com/google/uuid"
)

// Test UploadFile
type uploadFileMock struct {
	ioMock IoMock
	osMock OsMock
}

type uploadFileInput struct {
	request tests.RequestMock
}

type uploadFileExpected struct {
	err error
}

type uploadFileTestData struct {
	testCaseName string
	mock         uploadFileMock
	input        uploadFileInput
	expected     uploadFileExpected
}

var testDataUploadFile = []uploadFileTestData{
	{
		testCaseName: "MissingInputFile",
		input: uploadFileInput{
			request: tests.RequestMock{
				FormFileError: http.ErrMissingFile},
		},
		mock: uploadFileMock{},
		expected: uploadFileExpected{
			err: nil,
		},
	},
	{
		testCaseName: "FormFileError",
		input: uploadFileInput{
			request: tests.RequestMock{
				FormFileError: errors.New("FormFileError"),
			},
		},
		mock: uploadFileMock{},
		expected: uploadFileExpected{
			err: errors.New("FormFileError"),
		},
	},
	{
		testCaseName: "CreateOutputFileError",
		input: uploadFileInput{
			request: tests.RequestMock{},
		},
		mock: uploadFileMock{
			osMock: OsMock{
				CreateError: errors.New("CreateOutputFileError"),
			},
		},
		expected: uploadFileExpected{
			err: errors.New("CreateOutputFileError"),
		},
	},
	{
		testCaseName: "CopyInputToOutputFileError",
		input: uploadFileInput{
			request: tests.RequestMock{},
		},
		mock: uploadFileMock{
			osMock: OsMock{
				CreateFile: FileImpl{FileMock{}},
			},
			ioMock: IoMock{
				CopyError: errors.New("CopyInputToOutputFileError"),
			},
		},
		expected: uploadFileExpected{
			err: errors.New("CopyInputToOutputFileError"),
		},
	},
	{
		testCaseName: "CopyInputToOutputFileAndRemoveError",
		input: uploadFileInput{
			request: tests.RequestMock{},
		},
		mock: uploadFileMock{
			osMock: OsMock{
				CreateFile:     FileImpl{FileMock{}},
				RemoveAllError: errors.New("CopyInputToOutputFileAndRemoveError"),
			},
			ioMock: IoMock{
				CopyError: errors.New("CopyInputToOutputFileError"),
			},
		},
		expected: uploadFileExpected{
			err: errors.New("CopyInputToOutputFileAndRemoveError"),
		},
	},
	{
		testCaseName: "CloseOutputFileError",
		input: uploadFileInput{
			request: tests.RequestMock{},
		},
		mock: uploadFileMock{
			osMock: OsMock{
				CreateFile: FileImpl{FileMock{
					CloseError: errors.New("CloseOutputFileError"),
				}},
			},
		},
		expected: uploadFileExpected{
			err: errors.New("CloseOutputFileError"),
		},
	},
	{
		testCaseName: "CloseInputFileError",
		input: uploadFileInput{
			request: tests.RequestMock{
				File: tests.MultipartFileMock{
					CloseError: errors.New("CloseInputFileError"),
				},
			},
		},
		mock: uploadFileMock{
			osMock: OsMock{
				CreateFile: FileImpl{FileMock{}},
			},
		},
		expected: uploadFileExpected{
			err: errors.New("CloseInputFileError"),
		},
	},
	{
		testCaseName: "Success",
		input: uploadFileInput{
			request: tests.RequestMock{
				File: tests.MultipartFileMock{},
			},
		},
		mock: uploadFileMock{
			osMock: OsMock{
				CreateFile: FileImpl{FileMock{}},
			},
		},
		expected: uploadFileExpected{},
	},
}

func TestUploadFile(t *testing.T) {
	// Run tests
	for _, testCase := range testDataUploadFile {
		testCase := testCase
		t.Run(testCase.testCaseName, func(t *testing.T) {
			fileProcessor := FileProcessorImpl{
				FileIO: testCase.mock.ioMock,
				OsFunc: testCase.mock.osMock,
			}
			expected := testCase.expected.err
			request := testCase.input.request
			err := fileProcessor.UploadFile("testKey", "test.txt", &request)
			tests.CheckErrPartialMatch(err, expected, testCase.testCaseName, t)
		})
	}
}

// Test WriteToFile
type writeToFileMock struct {
	ioMock IoMock
	osMock OsMock
}

type writeToFileExpected struct {
	err error
}

type writeToFileTestData struct {
	testCaseName string
	mock         writeToFileMock
	expected     writeToFileExpected
}

var testDataWriteToFile = []writeToFileTestData{
	{
		testCaseName: "CreateFileError",
		mock: writeToFileMock{
			osMock: OsMock{
				CreateError: errors.New("CreateFileError"),
			},
		},
		expected: writeToFileExpected{
			err: errors.New("CreateFileError"),
		},
	},
	{
		testCaseName: "WriteError",
		mock: writeToFileMock{
			osMock: OsMock{
				CreateFile: FileImpl{FileMock{
					WriteError: errors.New("WriteError"),
				}},
			},
			ioMock: IoMock{},
		},
		expected: writeToFileExpected{
			err: errors.New("WriteError"),
		},
	},
	{
		testCaseName: "SyncError",
		mock: writeToFileMock{
			osMock: OsMock{
				CreateFile: FileImpl{FileMock{
					SyncError: errors.New("SyncError"),
				}},
			},
		},
		expected: writeToFileExpected{
			err: errors.New("SyncError"),
		},
	},
	{
		testCaseName: "CloseError",

		mock: writeToFileMock{
			osMock: OsMock{
				CreateFile: FileImpl{FileMock{
					CloseError: errors.New("CloseError"),
				}},
			},
		},
		expected: writeToFileExpected{
			err: errors.New("CloseError"),
		},
	},

	{
		testCaseName: "Success",

		mock: writeToFileMock{
			osMock: OsMock{
				CreateFile: FileImpl{FileMock{}},
			},
		},
		expected: writeToFileExpected{},
	},
}

func TestWriteToFile(t *testing.T) {
	// Run tests
	for _, testCase := range testDataWriteToFile {
		testCase := testCase
		t.Run(testCase.testCaseName, func(t *testing.T) {
			fileProcessor := FileProcessorImpl{
				FileIO: testCase.mock.ioMock,
				OsFunc: testCase.mock.osMock,
			}
			expected := testCase.expected.err
			err := fileProcessor.WriteToFile("testKey.txt", "testText")
			tests.CheckErrPartialMatch(err, expected, testCase.testCaseName, t)
		})
	}
}

// Test GeneratePath
type generatePathMock struct {
	envVarString string
	osMock       OsMock
}

type generatePathInput struct {
	asset dbModels.Asset
}

type generatePathExpected struct {
	err   error
	asset dbModels.Asset
}

type generatePathTestData struct {
	testCaseName string
	mock         generatePathMock
	input        generatePathInput
	expected     generatePathExpected
}

var testDataGeneratePath = []generatePathTestData{
	{
		testCaseName: "MakeDirError",
		mock: generatePathMock{
			envVarString: "/test",
			osMock: OsMock{
				MkdirAllError: errors.New("MakeDirError"),
			},
		},
		input: generatePathInput{
			asset: dbModels.Asset{
				ID:      uuid.MustParse("8c976c5e-cc29-4198-a386-93562c6a75cb"),
				DataMap: make(dbModels.DataMap),
			},
		},
		expected: generatePathExpected{
			asset: dbModels.Asset{
				ID:      uuid.MustParse("8c976c5e-cc29-4198-a386-93562c6a75cb"),
				DataMap: make(dbModels.DataMap),
			},
			err: errors.New("MakeDirError"),
		},
	},
	{
		testCaseName: "Success",
		mock: generatePathMock{
			envVarString: "/test",
			osMock: OsMock{
				MkdirAllError: errors.New("Success"),
			},
		},
		input: generatePathInput{
			asset: dbModels.Asset{
				ID:      uuid.MustParse("8c976c5e-cc29-4198-a386-93562c6a75cb"),
				DataMap: make(dbModels.DataMap),
			},
		},
		expected: generatePathExpected{
			asset: dbModels.Asset{
				ID:      uuid.MustParse("8c976c5e-cc29-4198-a386-93562c6a75cb"),
				DataMap: make(dbModels.DataMap),
			},
			err: errors.New("Success"),
		},
	},
}

func TestGeneratePath(t *testing.T) {
	// Run tests
	for _, testCase := range testDataGeneratePath {
		testCase := testCase
		t.Run(testCase.testCaseName, func(t *testing.T) {
			fileProcessor := FileProcessorImpl{
				OsFunc: testCase.mock.osMock,
			}
			os.Setenv("USER_STORE_DOCKER", testCase.mock.envVarString)
			input := testCase.input.asset
			expected := testCase.expected
			err := fileProcessor.GeneratePath(&input)
			tests.CheckResult(input, expected.asset, err, expected.err, testCase.testCaseName, t)
		})
	}
}
