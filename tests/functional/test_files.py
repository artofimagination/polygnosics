import pytest
import json
import common
import os
from pathlib import Path


# If "ref" is a filepath it contains a UUID
# that cannot be compared with test data.
# It changes every test run. So "ref" is removed from the response
# and is checked separately
# If it is a link it is compared to expected "refs"
# If it is a file the existence of file is checked.
def checkRef(result, expected):
    expectedRequest = expected["requests"]
    i = 0
    for item in result:
        if common.DELETE_RESOURCE in item or \
            common.GET_RESOURCE in item or \
                expected["refs"][i] == "-":
            continue
        (refToCheck, deleted) = common.deleteElement("ref", item)
        if deleted is False:
            pytest.fail(
                f"Request failed. \"ref\" is missing\n\
                Returned: {item}\n\
                Expected: {expectedRequest[i]}")
            return None
        if (common.ADD_RESOURCE in item and
            item[common.ADD_RESOURCE]["content"]["type"] == "link") or \
            (common.UPDATE_RESOURCE in item and
                item[common.UPDATE_RESOURCE]["content"]["type"] == "link"):
            # Check if the "ref" link is correct
            expectedRef = expected["refs"][i]
            if expectedRef != refToCheck:
                pytest.fail(
                    f"Request failed. \"ref\" is invalid\n\
                    Returned: {refToCheck}\n\
                    Expected: {expectedRef}")
                return None
        else:
            p = Path(refToCheck)
            p = p.parts[-1]
            if os.path.exists(
                os.path.join(
                    "tests/dummy-resourcedb/\
test-data/uploads/files", p)) is False:
                pytest.fail(
                    f"Request failed. \"ref\" file is missing\n\
                    Expected: {p}")
                return None
        i = i + 1
    return True


dataColumns = ("data", "expected")
createTestData = [
    (
        # Input data
        {
            'title': (None, "testFiles", 'application/json'),
            "short": (None, "testFileShort", 'application/json'),
            "type_0": (None, "link", 'application/json'),
            "ref_name_0": (None, "TestRef1", 'application/json'),
            "repo_link_0": (None, "https://", 'application/json'),
            "type_1": (None, "file", 'application/json'),
            "ref_name_1": (None, "TestRef2", 'application/json'),
        },
        # Expected
        {
            "error": "Backend -> strconv.Atoi: parsing \"\": invalid syntax",
            "requests": None,
            "refs": {}
        }
    ),
    (
        # Input data
        {
            'title': (None, "testFiles", 'application/json'),
            "short": (None, "testFileShort", 'application/json'),
            "type_0": (None, "link", 'application/json'),
            "ref_name_0": (None, "TestRef1", 'application/json'),
            "repo_link_0": (None, "http://", 'application/json'),
            "type_1": (None, "file", 'application/json'),
            "ref_name_1": (None, "TestRef2", 'application/json'),
            "count": (None, "2", 'application/json')
        },
        # Expected
        {
            "error": "",
            "response": {
              "data": "OK"
            },
            "refs": ["http://", "",  "-"],
            "requests": [
              {
                common.ADD_RESOURCE: {
                  'category': 0,
                  'content': {
                      'orig_file_name': '',
                      'ref_name': 'TestRef1',
                      'type': 'link'
                  },
                  'id': '00000000-0000-0000-0000-000000000000'
                }
              }, {
                common.ADD_RESOURCE: {
                    'category': 0,
                    'content': {
                        'orig_file_name': 'avatar-test.jpg',
                        'ref_name': 'TestRef2',
                        'type': 'file'
                    },
                    'id': '00000000-0000-0000-0000-000000000000'
                }
              }, {
                common.ADD_RESOURCE: {
                    'category': 1,
                    'content': {
                        'files': [
                            'c1e6122b-7986-417d-8bf6-ddf2dd9289f2',
                            'c1e6122b-7986-417d-8bf6-ddf2dd9289f2'
                        ],
                        'short': 'testFileShort',
                        'title': 'testFiles'
                    },
                    'id': '00000000-0000-0000-0000-000000000000'
                }
              }
            ]
        }
    ),
]

ids = ['Missing Count', "Success"]


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_AddFilesSection(httpBackend, httpDummyResourceDB, data, expected):
    # Clears all previously stored incoming requests on the dummy server
    try:
        r = httpDummyResourceDB.POST("/clear-request-data", None)
    except Exception:
        pytest.fail("Failed to send POST request")
        return None

    data['upload_file_1'] = \
        (os.path.basename("./tests/dummy-resourcedb/\
test-data/avatar-test.jpg"), open("./tests/dummy-resourcedb/\
test-data/avatar-test.jpg", 'rb'), 'application/octet-stream')

    try:
        r = httpBackend.POST(
            "/resources/create-files-item",
            files=data)
    except Exception:
        pytest.fail("Failed to send POST request")
        return None

    try:
        response = common.getResponse(r.text, expected["error"])
        if response is not None:
            expectedData = expected["response"]["data"]
            if response != expectedData:
                pytest.fail(
                    f"Request failed\n\
                    Status code: {r.status_code}\n\
                    Returned: {r.text}\n\
                    Expected: {expectedData}")
    except Exception:
        pytest.fail("Failed to process request")
        return None

    # Check the requests sent from the backend
    try:
        r = httpDummyResourceDB.GET("/get-request-data", None)
    except Exception:
        pytest.fail("Failed to send GET request")
        return None

    try:
        response = json.loads(r.text)
    except Exception:
        pytest.fail(f"Failed to decode response text {r.text}")
        return None

    expectedRequest = expected["requests"]
    if len(expected["refs"]) > 0 and checkRef(response, expected) is None:
        return

    if response != expectedRequest:
        pytest.fail(
            f"Request failed\n\
            Status code: {r.status_code}\n\
            Returned: {response}\n\
            Expected: {expectedRequest}")


dataColumns = ("data", "expected")
createTestData = [
    (
        # Input data
        {
            "id":
                (None,
                    "baca221f-0b14-4122-9d30-909e9b1014de",
                    'application/json'),
            "title": (None, "testFilesUpdate", 'application/json'),
            "short": (None, "testFilesShortUpdate", 'application/json'),
            "type_8ae5ae44-7b42-46d9-be18-08ea1d698883": \
            (None, "file", 'application/json'),
            "ref_name_8ae5ae44-7b42-46d9-be18-08ea1d698883": \
            (None, "TestRef2", 'application/json'),
            "count": (None, "2", 'application/json'),
            "type_1130d177-7f68-4ef8-92ac-8bfb40aa144c": \
            (None, "link", 'application/json'),
            "ref_name_1130d177-7f68-4ef8-92ac-8bfb40aa144c": \
            (None, "TestRef1", 'application/json'),
            "repo_link_1130d177-7f68-4ef8-92ac-8bfb40aa144c": \
            (None, "http://", 'application/json'),
        },
        # Expected
        {
            "error": "",
            "response": {
              "data": "OK"
            },
            "refs": ["", "http://", "-"],
            "requests": [
              {
                common.GET_RESOURCE: \
                '/get-resource-by-id?id=baca221f-0b14-4122-9d30-909e9b1014de'
              }, {
                common.GET_RESOURCE: \
                '/get-resource-by-id?id=8ae5ae44-7b42-46d9-be18-08ea1d698883'
              }, {
                common.UPDATE_RESOURCE: {
                  'category': 0,
                  'content': {
                      'orig_file_name': 'avatar-test.jpg',
                      'ref_name': 'TestRef2',
                      'type': 'file'
                  },
                  'id': '8ae5ae44-7b42-46d9-be18-08ea1d698883'
                }
              }, {
                common.GET_RESOURCE: \
                '/get-resource-by-id?id=1130d177-7f68-4ef8-92ac-8bfb40aa144c'
              }, {
                common.UPDATE_RESOURCE: {
                    'category': 0,
                    'content': {
                        'orig_file_name': '',
                        'ref_name': 'TestRef1',
                        'type': 'link'
                    },
                    'id': '1130d177-7f68-4ef8-92ac-8bfb40aa144c'
                }
              }, {
                common.UPDATE_RESOURCE: {
                    'category': 1,
                    'content': {
                        'files': [
                            '8ae5ae44-7b42-46d9-be18-08ea1d698883',
                            '1130d177-7f68-4ef8-92ac-8bfb40aa144c'
                        ],
                        'short': 'testFilesShortUpdate',
                        'title': 'testFilesUpdate'
                    },
                    'id': 'baca221f-0b14-4122-9d30-909e9b1014de'
                }
              }
            ]
        }
    ),
    (
        # Input data
        {
            "id":
                (None,
                    "baca221f-0b14-4122-9d30-909e9b1014de",
                    'application/json'),
            "title": (None, "testFilesUpdate", 'application/json'),
            "short": (None, "testFilesShortUpdate", 'application/json'),
            "type_8ae5ae44-7b42-46d9-be18-08ea1d698883": \
            (None, "file", 'application/json'),
            "ref_name_8ae5ae44-7b42-46d9-be18-08ea1d698883": \
            (None, "TestRef2", 'application/json'),
            "count": (None, "3", 'application/json'),
            "type_1130d177-7f68-4ef8-92ac-8bfb40aa144c": \
            (None, "link", 'application/json'),
            "ref_name_1130d177-7f68-4ef8-92ac-8bfb40aa144c": \
            (None, "TestRef1", 'application/json'),
            "repo_link_1130d177-7f68-4ef8-92ac-8bfb40aa144c": \
            (None, "http://", 'application/json'),
            "type_2": (None, "link", 'application/json'),
            "ref_name_2": (None, "TestRef3", 'application/json'),
            "repo_link_2": (None, "http://test", 'application/json'),
        },
        # Expected
        {
            "error": "",
            "response": {
              "data": "OK"
            },
            "refs": ["", "http://", "http://test", "-"],
            "requests": [
              {
                common.GET_RESOURCE: \
                '/get-resource-by-id?id=baca221f-0b14-4122-9d30-909e9b1014de'
              }, {
                common.GET_RESOURCE: \
                '/get-resource-by-id?id=8ae5ae44-7b42-46d9-be18-08ea1d698883'
              }, {
                common.UPDATE_RESOURCE: {
                  'category': 0,
                  'content': {
                      'orig_file_name': 'avatar-test.jpg',
                      'ref_name': 'TestRef2',
                      'type': 'file'
                  },
                  'id': '8ae5ae44-7b42-46d9-be18-08ea1d698883'
                }
              }, {
                common.GET_RESOURCE: \
                '/get-resource-by-id?id=1130d177-7f68-4ef8-92ac-8bfb40aa144c'
              }, {
                common.UPDATE_RESOURCE: {
                    'category': 0,
                    'content': {
                        'orig_file_name': '',
                        'ref_name': 'TestRef1',
                        'type': 'link'
                    },
                    'id': '1130d177-7f68-4ef8-92ac-8bfb40aa144c'
                }
              },
              {
                common.ADD_RESOURCE: {
                    'category': 0,
                    'content': {
                        'orig_file_name': '',
                        'ref_name': 'TestRef3',
                        'type': 'link'
                    },
                    'id': '00000000-0000-0000-0000-000000000000'
                }
              },
              {
                common.UPDATE_RESOURCE: {
                    'category': 1,
                    'content': {
                        'files': [
                            '8ae5ae44-7b42-46d9-be18-08ea1d698883',
                            '1130d177-7f68-4ef8-92ac-8bfb40aa144c',
                            'c1e6122b-7986-417d-8bf6-ddf2dd9289f2',
                        ],
                        'short': 'testFilesShortUpdate',
                        'title': 'testFilesUpdate'
                    },
                    'id': 'baca221f-0b14-4122-9d30-909e9b1014de'
                }
              }
            ]
        }
    ),
    (
        # Input data
        {
            "id":
                (None,
                    "baca221f-0b14-4122-9d30-909e9b1014de",
                    'application/json'),
            "title": (None, "testFilesDelete", 'application/json'),
            "short": (None, "testFilesShortDelete", 'application/json'),
            "remove_8ae5ae44-7b42-46d9-be18-08ea1d698883": \
            (None, "checked", 'application/json'),
            "count": (None, "2", 'application/json'),
            "type_1130d177-7f68-4ef8-92ac-8bfb40aa144c": \
            (None, "link", 'application/json'),
            "ref_name_1130d177-7f68-4ef8-92ac-8bfb40aa144c": \
            (None, "TestRef1", 'application/json'),
            "repo_link_1130d177-7f68-4ef8-92ac-8bfb40aa144c": \
            (None, "http://", 'application/json'),
        },
        # Expected
        {
            "error": "",
            "response": {
              "data": "OK"
            },
            "refs": ["http://", "-"],
            "requests": [
              {
                common.GET_RESOURCE: \
                '/get-resource-by-id?id=baca221f-0b14-4122-9d30-909e9b1014de'
              }, {
                common.DELETE_RESOURCE: {
                    'id': '8ae5ae44-7b42-46d9-be18-08ea1d698883'
                }
              }, {
                common.GET_RESOURCE: \
                '/get-resource-by-id?id=1130d177-7f68-4ef8-92ac-8bfb40aa144c'
              }, {
                common.UPDATE_RESOURCE: {
                    'category': 0,
                    'content': {
                        'orig_file_name': '',
                        'ref_name': 'TestRef1',
                        'type': 'link'
                    },
                    'id': '1130d177-7f68-4ef8-92ac-8bfb40aa144c'
                }
              },
              {
                common.UPDATE_RESOURCE: {
                    'category': 1,
                    'content': {
                        'files': [
                            '1130d177-7f68-4ef8-92ac-8bfb40aa144c',
                        ],
                        'short': 'testFilesShortDelete',
                        'title': 'testFilesDelete'
                    },
                    'id': 'baca221f-0b14-4122-9d30-909e9b1014de'
                }
              }
            ]
        }
    ),
]

ids = ["Success", "UpdateWithNewFile", "UpdateRemoveFile"]


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_UpdateFilesSection(httpBackend, httpDummyResourceDB, data, expected):
    testFilePath = "tests/dummy-resourcedb/test-data/uploads/files/test.txt"
    if os.path.exists(testFilePath) is False:
        with open(testFilePath, 'w') as f:
            f.write('Create a new text file!')

    # Clears all previously stored incoming requests on the dummy server
    try:
        r = httpDummyResourceDB.POST("/clear-request-data", None)
    except Exception:
        pytest.fail("Failed to send POST request")
        return None

    data['upload_file_8ae5ae44-7b42-46d9-be18-08ea1d698883'] = \
        (os.path.basename("./tests/dummy-resourcedb/\
test-data/avatar-test.jpg"), open("./tests/dummy-resourcedb/\
test-data/avatar-test.jpg", 'rb'), 'application/octet-stream')

    try:
        r = httpBackend.POST(
            "/resources/edit-files-item",
            files=data)
    except Exception:
        pytest.fail("Failed to send POST request")
        return None

    try:
        response = common.getResponse(r.text, expected["error"])
        if response is not None:
            expectedData = expected["response"]["data"]
            if response != expectedData:
                pytest.fail(
                    f"Request failed\n\
                    Status code: {r.status_code}\n\
                    Returned: {r.text}\n\
                    Expected: {expectedData}")
    except Exception:
        pytest.fail("Failed to process request")
        return None

    # Check the requests sent from the backend
    try:
        r = httpDummyResourceDB.GET("/get-request-data", None)
    except Exception:
        pytest.fail("Failed to send GET request")
        return None

    try:
        response = json.loads(r.text)
    except Exception:
        pytest.fail(f"Failed to decode response text {r.text}")
        return None

    expectedRequest = expected["requests"]

    if len(expected["refs"]) > 0 and checkRef(response, expected) is None:
        return

    if response != expectedRequest:
        pytest.fail(
            f"Request failed\n\
            Status code: {r.status_code}\n\
            Returned: {response}\n\
            Expected: {expectedRequest}")
