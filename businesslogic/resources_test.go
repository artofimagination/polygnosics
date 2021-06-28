package businesslogic

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/artofimagination/polygnosics/models"
	"github.com/artofimagination/polygnosics/rest"
	"github.com/artofimagination/polygnosics/rest/resourcesdb"
	"github.com/artofimagination/polygnosics/tests"
)

// Test getDataModel
type getDataModelInput struct {
	category string
}

type getDataModelExpected struct {
	reflectedType reflect.Type
	err           error
}

type getDataModelTestData struct {
	testCaseName string
	input        getDataModelInput
	expected     getDataModelExpected
}

var testDataGetDataModel = []getDataModelTestData{
	{
		testCaseName: "TutorialDataModel",
		input: getDataModelInput{
			category: "Tutorial",
		},
		expected: getDataModelExpected{
			reflectedType: reflect.TypeOf(&models.Tutorial{}),
		},
	},
	{
		testCaseName: "FAQDataModel",
		input: getDataModelInput{
			category: "FAQ",
		},
		expected: getDataModelExpected{
			reflectedType: reflect.TypeOf(&models.FAQ{}),
		},
	},
	{
		testCaseName: "FilesSectionDataModel",
		input: getDataModelInput{
			category: "FilesSection",
		},
		expected: getDataModelExpected{
			reflectedType: reflect.TypeOf(&models.FilesSection{}),
		},
	},
	{
		testCaseName: "FileContentDataModel",
		input: getDataModelInput{
			category: "FileContent",
		},
		expected: getDataModelExpected{
			reflectedType: reflect.TypeOf(&models.FileContent{}),
		},
	},
	{
		testCaseName: "Unknown",
		input: getDataModelInput{
			category: "Unknown",
		},
		expected: getDataModelExpected{
			reflectedType: nil,
			err:           errors.New("Invalid resource category: Unknown"),
		},
	},
}

func TestGetDataModel(t *testing.T) {
	// Run tests
	for _, testCase := range testDataGetDataModel {
		testCase := testCase
		t.Run(testCase.testCaseName, func(t *testing.T) {
			dataModel, err := getDataModel(testCase.input.category)
			tests.CheckResult(reflect.TypeOf(dataModel), testCase.expected.reflectedType, err, testCase.expected.err, testCase.testCaseName, t)
		})
	}
}

// Test getCategoryID
type getCategoryIDMock struct {
	db resourcesdb.Mock
}

type getCategoryIDInput struct {
	categoryString string
}

type getCategoryIDExpected struct {
	id  int
	err error
}

type getCategoryIDTestData struct {
	testCaseName string
	mock         getCategoryIDMock
	input        getCategoryIDInput
	expected     getCategoryIDExpected
}

var testDataGetCateogoryID = []getCategoryIDTestData{
	{
		testCaseName: "GetCategoriesError",
		mock: getCategoryIDMock{
			db: resourcesdb.Mock{
				CategoriesError: errors.New("GetCategoriesError"),
			},
		},
		input: getCategoryIDInput{
			categoryString: "TestCategory1",
		},
		expected: getCategoryIDExpected{
			id:  -1,
			err: errors.New("GetCategoriesError"),
		},
	},
	{
		testCaseName: "MissingCategory",
		mock: getCategoryIDMock{
			db: resourcesdb.Mock{
				Categories: []models.Category{
					{
						ID:          1,
						Name:        "TestCategory1",
						Description: "Test 1",
					},
					{
						ID:          2,
						Name:        "TestCategory2",
						Description: "Test 2",
					},
				},
			},
		},
		input: getCategoryIDInput{
			categoryString: "TestCategoryMissing",
		},
		expected: getCategoryIDExpected{
			id:  -1,
			err: fmt.Errorf("TestCategoryMissing category not found"),
		},
	},
	{
		testCaseName: "Success",
		mock: getCategoryIDMock{
			db: resourcesdb.Mock{
				Categories: []models.Category{
					{
						ID:          1,
						Name:        "TestCategory1",
						Description: "Test 1",
					},
					{
						ID:          2,
						Name:        "TestCategory2",
						Description: "Test 2",
					},
				},
			},
		},
		input: getCategoryIDInput{
			categoryString: "TestCategory1",
		},
		expected: getCategoryIDExpected{
			id: 1,
		},
	},
}

func TestGetCategoryID(t *testing.T) {
	// Run tests
	for _, testCase := range testDataGetCateogoryID {
		testCase := testCase
		t.Run(testCase.testCaseName, func(t *testing.T) {
			context := Context{
				ResourcesDBController: &testCase.mock.db,
			}
			input := testCase.input
			expected := testCase.expected
			categoryID, err := context.getCategoryID(input.categoryString)
			tests.CheckResult(categoryID, expected.id, err, expected.err, testCase.testCaseName, t)
		})
	}
}

// Test UpdateHandler
type updateHandlerInput struct {
	category    string
	testParams  []interface{}
	handlerFunc func(rest.RequestInterface, ...interface{}) error
}

type updateHandlerExpected struct {
	testObject models.Tutorial
}

type updateHandlerTestData struct {
	testCaseName string
	input        updateHandlerInput
	expected     updateHandlerExpected
}

var testDataUpdateHandler = []updateHandlerTestData{
	{
		testCaseName: "Success",
		input: updateHandlerInput{
			category: "Tutorial",
			testParams: []interface{}{
				"test_ID",
			},
			handlerFunc: func(r rest.RequestInterface, dataModel ...interface{}) error {
				test := dataModel[0].(*models.Tutorial)
				test.Content = "content set"
				return nil
			},
		},
		expected: updateHandlerExpected{
			testObject: models.Tutorial{
				Content: "content set",
			},
		},
	},
}

func TestUpdateHandler(t *testing.T) {
	// Run tests
	for _, testCase := range testDataUpdateHandler {
		testCase := testCase
		t.Run(testCase.testCaseName, func(t *testing.T) {
			context := Context{
				ResourcesDBController: &resourcesdb.Mock{
					UpdatedResource: models.Tutorial{},
				},
			}

			_ = context.UpdateHandler(testCase.input.category, &tests.RequestMock{}, testCase.input.handlerFunc, testCase.input.testParams...)
			tests.CheckResult(context.ResourcesDBController.(*resourcesdb.Mock).UpdatedResource, &testCase.expected.testObject, nil, nil, testCase.testCaseName, t)
		})
	}
}

// Test AddHandler
type addHandlerInput struct {
	category    string
	testParams  []interface{}
	handlerFunc func(rest.RequestInterface, ...interface{}) error
}

type addHandlerExpected struct {
	testObject models.Tutorial
}

type addHandlerTestData struct {
	testCaseName string
	input        addHandlerInput
	expected     addHandlerExpected
}

var testDataAddHandler = []addHandlerTestData{
	{
		testCaseName: "Success0Params",
		input: addHandlerInput{
			category:   "Tutorial",
			testParams: []interface{}{},
			handlerFunc: func(r rest.RequestInterface, input ...interface{}) error {
				test := input[0].(*models.Tutorial)
				test.Content = "content set"
				return nil
			},
		},
		expected: addHandlerExpected{
			testObject: models.Tutorial{
				Content: "content set",
			},
		},
	},
	{
		testCaseName: "Success1Param",
		input: addHandlerInput{
			category: "Tutorial",
			testParams: []interface{}{
				"content set1",
			},
			handlerFunc: func(r rest.RequestInterface, input ...interface{}) error {
				test := input[0].(*models.Tutorial)
				test.Content = input[1].([]interface{})[0].(string)
				return nil
			},
		},
		expected: addHandlerExpected{
			testObject: models.Tutorial{
				Content: "content set1",
			},
		},
	},
}

func TestAddHandler(t *testing.T) {
	// Run tests
	for _, testCase := range testDataAddHandler {
		testCase := testCase
		t.Run(testCase.testCaseName, func(t *testing.T) {
			context := Context{
				ResourcesDBController: &resourcesdb.Mock{
					AddedResource: models.Tutorial{},
					Categories: []models.Category{
						{
							ID:          1,
							Name:        "Tutorial",
							Description: "Test 1",
						},
					},
				},
			}

			_, _ = context.AddHandler(testCase.input.category, &tests.RequestMock{}, testCase.input.handlerFunc, testCase.input.testParams...)
			tests.CheckResult(context.ResourcesDBController.(*resourcesdb.Mock).AddedResource, &testCase.expected.testObject, nil, nil, testCase.testCaseName, t)
		})
	}
}

// Test setTutorialArticle
type setTutorialArticleMock struct {
	fileProcessor FileProcessorMock
}

type setTutorialArticleInput struct {
	tutorial models.Tutorial
	request  tests.RequestMock
}

type setTutorialArticleExpected struct {
	fileProcessorCallOrder []string
	tutorial               models.Tutorial
	err                    error
}

type setTutorialArticleTestData struct {
	testCaseName string
	mock         setTutorialArticleMock
	input        setTutorialArticleInput
	expected     setTutorialArticleExpected
}

var testDataSetTutorialArticle = []setTutorialArticleTestData{
	{
		testCaseName: "UpdateFromArticleToNoArticle",
		mock:         setTutorialArticleMock{},
		input: setTutorialArticleInput{
			tutorial: models.Tutorial{
				Content: "Testcontent",
			},
			request: tests.RequestMock{
				FormValues: map[string]string{
					"article":      "",
					"avatar_type":  "image",
					"avatar_video": "",
					"avatar":       "",
					"title":        "UpdatedTestTutorial",
					"short":        "Update TestTutorial",
				},
			},
		},
		expected: setTutorialArticleExpected{
			fileProcessorCallOrder: []string{
				"RemoveFile",
			},
			tutorial: models.Tutorial{
				Content: "",
			},
		},
	},
	{
		testCaseName: "UpdateFromNoArticleToArticle",
		mock:         setTutorialArticleMock{},
		input: setTutorialArticleInput{
			tutorial: models.Tutorial{
				Content: "",
			},
			request: tests.RequestMock{
				FormValues: map[string]string{
					"article":      "test.txt",
					"avatar_type":  "image",
					"avatar_video": "",
					"avatar":       "",
					"title":        "UpdatedTestTutorial",
					"short":        "Update TestTutorial",
				},
			},
		},
		expected: setTutorialArticleExpected{
			fileProcessorCallOrder: []string{
				"WriteToFile",
			},
			tutorial: models.Tutorial{
				Content: "/resources/tutorials/a4d28c75-595b-4059-801f-2a9ad127916b.txt",
			},
		},
	},
	{
		testCaseName: "UpdateArticleToArticle",
		mock:         setTutorialArticleMock{},
		input: setTutorialArticleInput{
			tutorial: models.Tutorial{
				Content: "testprevious.txt",
			},
			request: tests.RequestMock{
				FormValues: map[string]string{
					"article":      "test.txt",
					"avatar_type":  "image",
					"avatar_video": "",
					"avatar":       "",
					"title":        "UpdatedTestTutorial",
					"short":        "Update TestTutorial",
				},
			},
		},
		expected: setTutorialArticleExpected{
			fileProcessorCallOrder: []string{
				"WriteToFile",
			},
			tutorial: models.Tutorial{
				Content: "testprevious.txt",
			},
		},
	},
	{
		testCaseName: "UpdateNoArticleToNoArticle",
		mock:         setTutorialArticleMock{},
		input: setTutorialArticleInput{
			tutorial: models.Tutorial{
				Content: "",
			},
			request: tests.RequestMock{
				FormValues: map[string]string{
					"article":      "",
					"avatar_type":  "image",
					"avatar_video": "",
					"avatar":       "",
					"title":        "UpdatedTestTutorial",
					"short":        "Update TestTutorial",
				},
			},
		},
		expected: setTutorialArticleExpected{
			fileProcessorCallOrder: []string{},
			tutorial: models.Tutorial{
				Content: "",
			},
		},
	},
	{
		testCaseName: "UpdateErrorFromArticleToNoArticle",
		mock: setTutorialArticleMock{
			fileProcessor: FileProcessorMock{
				RemoveFileError: errors.New("UpdateError"),
			},
		},
		input: setTutorialArticleInput{
			tutorial: models.Tutorial{
				Content: "test.txt",
			},
			request: tests.RequestMock{
				FormValues: map[string]string{
					"article":      "",
					"avatar_type":  "image",
					"avatar_video": "",
					"avatar":       "",
					"title":        "UpdatedTestTutorial",
					"short":        "Update TestTutorial",
				},
			},
		},
		expected: setTutorialArticleExpected{
			fileProcessorCallOrder: []string{
				"RemoveFile",
			},
			tutorial: models.Tutorial{
				Content: "test.txt",
			},
			err: errors.New("UpdateError"),
		},
	},
}

func TestSetTutorialArticle(t *testing.T) {
	// Run tests
	for _, testCase := range testDataSetTutorialArticle {
		testCase := testCase
		t.Run(testCase.testCaseName, func(t *testing.T) {
			context := Context{
				FileProcessor: &FileProcessorMock{
					CallOrder:       make([]string, 0),
					RemoveFileError: testCase.mock.fileProcessor.RemoveFileError,
				},
			}

			tutorial := &testCase.input.tutorial
			err := context.setTutorialArticle(tutorial, &testCase.input.request)
			tests.CheckResult(*tutorial, testCase.expected.tutorial, err, testCase.expected.err, testCase.testCaseName, t)
			tests.CheckResult(context.FileProcessor.(*FileProcessorMock).CallOrder, testCase.expected.fileProcessorCallOrder, nil, nil, testCase.testCaseName, t)
		})
	}
}

// Test setTutorialAvatar
type setTutorialAvatarMock struct {
	fileProcessor FileProcessorMock
}

type setTutorialAvatarInput struct {
	tutorial models.Tutorial
	request  tests.RequestMock
}

type setTutorialAvatarExpected struct {
	fileProcessorCallOrder []string
	tutorial               models.Tutorial
	err                    error
}

type setTutorialAvatarTestData struct {
	testCaseName string
	mock         setTutorialAvatarMock
	input        setTutorialAvatarInput
	expected     setTutorialAvatarExpected
}

var testDataSetTutorialAvatar = []setTutorialAvatarTestData{
	{
		testCaseName: "UpdateFromVideoToImage",
		mock:         setTutorialAvatarMock{},
		input: setTutorialAvatarInput{
			tutorial: models.Tutorial{
				AvatarType:   "video",
				AvatarSource: "http://oldvideo",
			},
			request: tests.RequestMock{
				FormValues: map[string]string{
					"article":      "",
					"avatar_type":  "image",
					"avatar_video": "",
					"avatar":       "image.img",
					"title":        "UpdatedTestTutorial",
					"short":        "Update TestTutorial",
				},
			},
		},
		expected: setTutorialAvatarExpected{
			fileProcessorCallOrder: []string{
				"UploadFile",
			},
			tutorial: models.Tutorial{
				AvatarSource: "/resources/tutorials/a4d28c75-595b-4059-801f-2a9ad127916b.img",
				AvatarType:   "image",
			},
		},
	},
	{
		testCaseName: "UpdateFromVideoToVideo",
		mock:         setTutorialAvatarMock{},
		input: setTutorialAvatarInput{
			tutorial: models.Tutorial{
				AvatarType:   "video",
				AvatarSource: "http://oldvideo",
			},
			request: tests.RequestMock{
				FormValues: map[string]string{
					"article":      "",
					"avatar_type":  "video",
					"avatar_video": "http://video",
					"avatar":       "",
					"title":        "UpdatedTestTutorial",
					"short":        "Update TestTutorial",
				},
			},
		},
		expected: setTutorialAvatarExpected{
			fileProcessorCallOrder: []string{},
			tutorial: models.Tutorial{
				AvatarSource: "http://video",
				AvatarType:   "video",
			},
		},
	},
	{
		testCaseName: "UpdateFromImageToVideo",
		mock:         setTutorialAvatarMock{},
		input: setTutorialAvatarInput{
			tutorial: models.Tutorial{
				AvatarType:   "image",
				AvatarSource: "test.img",
			},
			request: tests.RequestMock{
				FormValues: map[string]string{
					"article":      "",
					"avatar_type":  "video",
					"avatar_video": "http://video",
					"avatar":       "",
					"title":        "UpdatedTestTutorial",
					"short":        "Update TestTutorial",
				},
			},
		},
		expected: setTutorialAvatarExpected{
			fileProcessorCallOrder: []string{
				"RemoveFile",
			},
			tutorial: models.Tutorial{
				AvatarSource: "http://video",
				AvatarType:   "video",
			},
		},
	},
	{
		testCaseName: "UpdateFromImageToImage",
		mock:         setTutorialAvatarMock{},
		input: setTutorialAvatarInput{
			tutorial: models.Tutorial{
				AvatarType:   "image",
				AvatarSource: "oldtest.img",
			},
			request: tests.RequestMock{
				FormValues: map[string]string{
					"article":      "",
					"avatar_type":  "image",
					"avatar_video": "testimage.png",
					"avatar":       "",
					"title":        "UpdatedTestTutorial",
					"short":        "Update TestTutorial",
				},
			},
		},
		expected: setTutorialAvatarExpected{
			fileProcessorCallOrder: []string{
				"UploadFile",
			},
			tutorial: models.Tutorial{
				AvatarSource: "oldtest.img",
				AvatarType:   "image",
			},
		},
	},
	{
		testCaseName: "UpdateErrorFromImageToImage",
		mock: setTutorialAvatarMock{
			fileProcessor: FileProcessorMock{
				UploadFileError: errors.New("UpdateError"),
			},
		},
		input: setTutorialAvatarInput{
			tutorial: models.Tutorial{
				AvatarType:   "image",
				AvatarSource: "oldtest.img",
			},
			request: tests.RequestMock{
				FormValues: map[string]string{
					"article":      "",
					"avatar_type":  "image",
					"avatar_video": "testimage.png",
					"avatar":       "",
					"title":        "UpdatedTestTutorial",
					"short":        "Update TestTutorial",
				},
			},
		},
		expected: setTutorialAvatarExpected{
			fileProcessorCallOrder: []string{
				"UploadFile",
			},
			tutorial: models.Tutorial{
				AvatarSource: "oldtest.img",
				AvatarType:   "image",
			},
			err: errors.New("UpdateError"),
		},
	},
}

func TestSetTutorialAvatar(t *testing.T) {
	// Run tests
	for _, testCase := range testDataSetTutorialAvatar {
		testCase := testCase
		t.Run(testCase.testCaseName, func(t *testing.T) {
			context := Context{
				FileProcessor: &FileProcessorMock{
					CallOrder:       make([]string, 0),
					UploadFileError: testCase.mock.fileProcessor.UploadFileError,
				},
			}

			tutorial := &testCase.input.tutorial
			err := context.setTutorialAvatar(tutorial, &testCase.input.request)
			tests.CheckResult(*tutorial, testCase.expected.tutorial, err, testCase.expected.err, testCase.testCaseName, t)
			tests.CheckResult(context.FileProcessor.(*FileProcessorMock).CallOrder, testCase.expected.fileProcessorCallOrder, nil, nil, testCase.testCaseName, t)
		})
	}
}

// Test UpdateTutorial
type updateTutorialMock struct {
	fileProcessor FileProcessorMock
}

type updateTutorialInput struct {
	request tests.RequestMock
}

type updateTutorialExpected struct {
	err error
}

type updateTutorialTestData struct {
	testCaseName string
	mock         updateTutorialMock
	input        updateTutorialInput
	expected     updateTutorialExpected
}

var testDataUpdateTutorial = []updateTutorialTestData{
	{
		testCaseName: "SetTutorialArticleFailed",
		mock: updateTutorialMock{
			fileProcessor: FileProcessorMock{
				WriteToFileError: errors.New("SetTutorialArticleFailed"),
			},
		},
		input: updateTutorialInput{
			request: tests.RequestMock{
				FormValues: map[string]string{
					"article":      "testArticle",
					"avatar_type":  "image",
					"avatar_video": "",
					"avatar":       "",
					"title":        "UpdatedTestTutorial",
					"short":        "Update TestTutorial",
				},
			},
		},
		expected: updateTutorialExpected{
			err: errors.New("SetTutorialArticleFailed"),
		},
	},
	{
		testCaseName: "SetTutorialAvatarFailed",
		mock: updateTutorialMock{
			fileProcessor: FileProcessorMock{
				UploadFileError: errors.New("SetTutorialAvatarFailed"),
			},
		},
		input: updateTutorialInput{
			request: tests.RequestMock{
				FormValues: map[string]string{
					"article":      "",
					"avatar_type":  "image",
					"avatar_video": "",
					"avatar":       "",
					"title":        "UpdatedTestTutorial",
					"short":        "Update TestTutorial",
				},
			},
		},
		expected: updateTutorialExpected{
			err: errors.New("SetTutorialAvatarFailed"),
		},
	},
	{
		testCaseName: "Success",
		mock:         updateTutorialMock{},
		input: updateTutorialInput{
			request: tests.RequestMock{
				FormValues: map[string]string{
					"article":      "",
					"avatar_type":  "image",
					"avatar_video": "",
					"avatar":       "",
					"title":        "UpdatedTestTutorial",
					"short":        "Update TestTutorial",
				},
			},
		},
		expected: updateTutorialExpected{},
	},
}

func TestUpdateTutorial(t *testing.T) {
	// Run tests
	for _, testCase := range testDataUpdateTutorial {
		testCase := testCase
		t.Run(testCase.testCaseName, func(t *testing.T) {
			context := Context{
				FileProcessor: &FileProcessorMock{
					CallOrder:        make([]string, 0),
					WriteToFileError: testCase.mock.fileProcessor.WriteToFileError,
					UploadFileError:  testCase.mock.fileProcessor.UploadFileError,
				},
			}

			input := testCase.input
			expected := testCase.expected
			tutorial := &models.Tutorial{}
			err := context.UpdateTutorial(&input.request, tutorial)
			tests.CheckResult(nil, nil, err, expected.err, testCase.testCaseName, t)
		})
	}
}

// Test setFileItem
type setFileItemMock struct {
	fileProcessor FileProcessorMock
}

type setFileItemInput struct {
	request tests.RequestMock
}

type setFileItemExpected struct {
	fileProcessorCallOrder []string
	fileItem               models.FileContent
	err                    error
}

type setFileItemTestData struct {
	testCaseName string
	mock         setFileItemMock
	input        setFileItemInput
	expected     setFileItemExpected
}

var testDataSetFileItem = []setFileItemTestData{
	{
		testCaseName: "AddLink",
		mock:         setFileItemMock{},
		input: setFileItemInput{
			request: tests.RequestMock{
				FormValues: map[string]string{
					"type_1":      "link",
					"ref_name_1":  "TestRef",
					"repo_link_1": "http://",
				},
				FileName: "test.txt",
			},
		},
		expected: setFileItemExpected{
			fileProcessorCallOrder: []string{},
			fileItem: models.FileContent{
				Type:    "link",
				RefName: "TestRef",
				Ref:     "http://",
			},
		},
	},
	{
		testCaseName: "AddFile",
		mock:         setFileItemMock{},
		input: setFileItemInput{
			request: tests.RequestMock{
				FormValues: map[string]string{
					"type_1":        "file",
					"ref_name_1":    "TestRef",
					"upload_file_1": "test.txt",
				},
				FileName: "test.txt",
			},
		},
		expected: setFileItemExpected{
			fileProcessorCallOrder: []string{
				"UploadFile",
			},
			fileItem: models.FileContent{
				Type:             "file",
				RefName:          "TestRef",
				Ref:              "/resources/files/a4d28c75-595b-4059-801f-2a9ad127916b.txt",
				OriginalFileName: "test.txt",
			},
		},
	},
	{
		testCaseName: "AddFileError",
		mock: setFileItemMock{
			fileProcessor: FileProcessorMock{
				UploadFileError: errors.New("AddFileError"),
			},
		},
		input: setFileItemInput{
			request: tests.RequestMock{
				FormValues: map[string]string{
					"type_1":        "file",
					"ref_name_1":    "TestRef",
					"upload_file_1": "test.txt",
				},
				FileName: "test.txt",
			},
		},
		expected: setFileItemExpected{
			fileProcessorCallOrder: []string{
				"UploadFile",
			},
			fileItem: models.FileContent{
				Type:             "file",
				RefName:          "TestRef",
				Ref:              "/resources/files/a4d28c75-595b-4059-801f-2a9ad127916b.txt",
				OriginalFileName: "test.txt",
			},
			err: errors.New("AddFileError"),
		},
	},
}

func TestSetFileItem(t *testing.T) {
	// Run tests
	for _, testCase := range testDataSetFileItem {
		testCase := testCase
		t.Run(testCase.testCaseName, func(t *testing.T) {
			context := Context{
				FileProcessor: &FileProcessorMock{
					CallOrder:       make([]string, 0),
					UploadFileError: testCase.mock.fileProcessor.UploadFileError,
				},
			}

			fileItem := models.FileContent{}
			index := "1"
			err := context.setFileItem(&testCase.input.request, &fileItem, []interface{}{"", index})
			tests.CheckResult(fileItem, testCase.expected.fileItem, err, testCase.expected.err, testCase.testCaseName, t)
			tests.CheckResult(context.FileProcessor.(*FileProcessorMock).CallOrder, testCase.expected.fileProcessorCallOrder, nil, nil, testCase.testCaseName, t)
		})
	}
}

// Test AddFileSection
type addFileSectionMock struct {
	fileProcessor FileProcessorMock
}

type addFileSectionInput struct {
	request tests.RequestMock
}

type addFileSectionExpected struct {
	fileProcessorCallOrder []string
	filesSection           models.FilesSection
	err                    error
}

type addFileSectionTestData struct {
	testCaseName string
	mock         addFileSectionMock
	input        addFileSectionInput
	expected     addFileSectionExpected
}

var testDataAddFileSection = []addFileSectionTestData{
	{
		testCaseName: "AddSectionWithMultipleItems",
		mock:         addFileSectionMock{},
		input: addFileSectionInput{
			request: tests.RequestMock{
				FormValues: map[string]string{
					"title":       "TestTitle",
					"short":       "Test description",
					"count":       "2",
					"type_1":      "link",
					"ref_name_1":  "TestRef1",
					"repo_link_1": "http://",
					"type_2":      "file",
					"ref_name_2":  "TestRef2",
					"repo_link_2": "test.com",
				},
			},
		},
		expected: addFileSectionExpected{
			fileProcessorCallOrder: []string{},
			filesSection: models.FilesSection{
				Title:            "TestTitle",
				ShortDescription: "Test description",
				ContentIDList:    []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000000"},
			},
		},
	},
	{
		testCaseName: "AddSectionWithZeroItems",
		mock:         addFileSectionMock{},
		input: addFileSectionInput{
			request: tests.RequestMock{
				FormValues: map[string]string{
					"title": "TestTitle",
					"short": "Test description",
					"count": "0",
				},
			},
		},
		expected: addFileSectionExpected{
			fileProcessorCallOrder: []string{},
			filesSection: models.FilesSection{
				Title:            "TestTitle",
				ShortDescription: "Test description",
				ContentIDList:    []string{},
			},
		},
	},
	{
		testCaseName: "AddSectionWithNoCountFormvalue",
		mock:         addFileSectionMock{},
		input: addFileSectionInput{
			request: tests.RequestMock{
				FormValues: map[string]string{
					"title": "TestTitle",
					"short": "Test description",
				},
			},
		},
		expected: addFileSectionExpected{
			fileProcessorCallOrder: []string{},
			filesSection:           models.FilesSection{},
			err: &strconv.NumError{
				Func: "Atoi",
				Num:  "",
				Err:  errors.New("invalid syntax"),
			},
		},
	},
	{
		testCaseName: "AddSectionWithItemsOverLimit",
		mock:         addFileSectionMock{},
		input: addFileSectionInput{
			request: tests.RequestMock{
				FormValues: map[string]string{
					"title":       "TestTitle",
					"short":       "Test description",
					"count":       "3",
					"type_1":      "link",
					"ref_name_1":  "TestRef1",
					"repo_link_1": "http://",
					"type_2":      "file",
					"ref_name_2":  "TestRef2",
					"repo_link_2": "test.com",
					"type_3":      "file",
					"ref_name_3":  "TestRef3",
					"repo_link_3": "test.com",
				},
			},
		},
		expected: addFileSectionExpected{
			fileProcessorCallOrder: []string{},
			filesSection: models.FilesSection{
				Title:            "TestTitle",
				ShortDescription: "Test description",
				ContentIDList:    []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000000"},
			},
		},
	},
}

func TestAddFileSection(t *testing.T) {
	// Run tests
	for _, testCase := range testDataAddFileSection {
		testCase := testCase
		t.Run(testCase.testCaseName, func(t *testing.T) {
			context := Context{
				FileProcessor: &FileProcessorMock{
					CallOrder:       make([]string, 0),
					UploadFileError: testCase.mock.fileProcessor.UploadFileError,
				},
				ResourcesDBController: &resourcesdb.Mock{
					Categories: []models.Category{
						{
							ID:          1,
							Name:        "FileContent",
							Description: "Test 1",
						},
					},
				},
			}

			filesSection := models.FilesSection{}
			MaxFileAttachements = 2
			err := context.AddFileSection(&testCase.input.request, &filesSection)
			tests.CheckResult(filesSection, testCase.expected.filesSection, err, testCase.expected.err, testCase.testCaseName, t)
			tests.CheckResult(context.FileProcessor.(*FileProcessorMock).CallOrder, testCase.expected.fileProcessorCallOrder, nil, nil, testCase.testCaseName, t)
		})
	}
}

// Test removeFileItem
type removeFileItemInput struct {
	request tests.RequestMock
	id      string
}

type removeFileItemExpected struct {
	fileProcessorCallOrder []string
	fileContentIDList      []string
	err                    error
}

type removeFileItemTestData struct {
	testCaseName string
	input        removeFileItemInput
	expected     removeFileItemExpected
}

var testDataRemoveFileItem = []removeFileItemTestData{
	{
		testCaseName: "Success",
		input: removeFileItemInput{
			id: "af09d4f9-9365-43fa-b06b-56965c54a6bf",
			request: tests.RequestMock{
				FormValues: map[string]string{
					"upload_file_af09d4f9-9365-43fa-b06b-56965c54a6bf": "test.com",
				},
			},
		},
		expected: removeFileItemExpected{
			fileProcessorCallOrder: []string{
				"RemoveFile",
			},
			fileContentIDList: []string{
				"e12b2ac5-7bad-49fb-b9f8-c2d1f5328b3e",
				"a4d28c75-595b-4059-801f-2a9ad127916b",
			},
		},
	},
}

func TestRemoveFileItem(t *testing.T) {
	// Run tests
	for _, testCase := range testDataRemoveFileItem {
		testCase := testCase
		t.Run(testCase.testCaseName, func(t *testing.T) {
			context := Context{
				FileProcessor: &FileProcessorMock{
					CallOrder: make([]string, 0),
				},
				ResourcesDBController: &resourcesdb.Mock{},
			}

			contentIDList := []string{
				"e12b2ac5-7bad-49fb-b9f8-c2d1f5328b3e",
				"af09d4f9-9365-43fa-b06b-56965c54a6bf",
				"a4d28c75-595b-4059-801f-2a9ad127916b",
			}

			contentIDList, err := context.removeFileItem(testCase.input.id, &testCase.input.request, contentIDList)
			tests.CheckResult(contentIDList, testCase.expected.fileContentIDList, err, testCase.expected.err, testCase.testCaseName, t)
			tests.CheckResult(context.FileProcessor.(*FileProcessorMock).CallOrder, testCase.expected.fileProcessorCallOrder, nil, nil, testCase.testCaseName, t)
		})
	}
}

// Test UpdateFAQ
type updateFAQInput struct {
	request tests.RequestMock
}

type updateFAQExpected struct {
	fileProcessorCallOrder []string
	faq                    models.FAQ
	err                    error
}

type updateFAQTestData struct {
	testCaseName string
	input        updateFAQInput
	expected     updateFAQExpected
}

var testDataUpdateFAQ = []updateFAQTestData{
	{
		testCaseName: "Success",
		input: updateFAQInput{
			request: tests.RequestMock{
				FormValues: map[string]string{
					"group":    "TestGroup",
					"answer":   "TestAnswer",
					"question": "TestQuestion",
				},
			},
		},
		expected: updateFAQExpected{
			fileProcessorCallOrder: []string{
				"WriteToFile", "WriteToFile",
			},
			faq: models.FAQ{
				Group:    "TestGroup",
				Question: "/resources/faq/a4d28c75-595b-4059-801f-2a9ad127916b/question.txt",
				Answer:   "/resources/faq/a4d28c75-595b-4059-801f-2a9ad127916b/answer.txt",
			},
		},
	},
}

func TestUpdateFAQ(t *testing.T) {
	// Run tests
	for _, testCase := range testDataUpdateFAQ {
		testCase := testCase
		t.Run(testCase.testCaseName, func(t *testing.T) {
			context := Context{
				FileProcessor: &FileProcessorMock{
					CallOrder: make([]string, 0),
				},
				ResourcesDBController: &resourcesdb.Mock{},
			}

			faq := &models.FAQ{}
			err := context.UpdateFAQ(&testCase.input.request, faq)
			tests.CheckResult(*faq, testCase.expected.faq, err, testCase.expected.err, testCase.testCaseName, t)
			tests.CheckResult(context.FileProcessor.(*FileProcessorMock).CallOrder, testCase.expected.fileProcessorCallOrder, nil, nil, testCase.testCaseName, t)
		})
	}
}
