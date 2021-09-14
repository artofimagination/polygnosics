import json
import pytest
import pathlib
import os

ADD_USER_PATH = "/add-user"
GET_USER_BY_EMAIL_PATH = "/get-user-by-email"
UPDATE_USER_SETTINGS_PATH = "/update-user-settings"
UPDATE_USER_ASSETS_PATH = "/update-user-assets"

ADD_RESOURCE = "/add-resource"
UPDATE_RESOURCE = "/update-resource"
DELETE_RESOURCE = "/delete-resource"
GET_RESOURCE = "/get-resource-by-id"


# getResponse unwraps the data/error from json response.
def getResponse(responseText, expected=None):
    response = json.loads(responseText)
    if "error" in response:
        error = response["error"]
        if expected is None or (expected is not None and error != expected):
            pytest.fail(f"Failed to run test.\nDetails: {error}")
        return None
    return response["data"]


# Locates and deletes item in the dictionary.
# Will return the element value and True as second value if it was deleted.
# Returns None as first value if there was no value found.
def deleteElement(key, jsonData):
    if isinstance(jsonData, list):
        for item in jsonData:
            (element, deleted) = deleteElement(key, item)
            if element is not None:
                return (element, deleted)
    else:
        if isinstance(jsonData, list) is False \
                and isinstance(jsonData, dict) is False:
            return (None, False)
        if key in jsonData:
            value = jsonData[key]
            del jsonData[key]
            return (value, True)

        for k, v in jsonData.items():
            if isinstance(jsonData, list) is False \
                    and isinstance(jsonData, dict) is False:
                return (None, False)
            (element, deleted) = deleteElement(key, v)
            if element is not None:
                return (element, deleted)
    return (None, False)


def checkFileContent(filename, expectedContent):
    p = pathlib.Path(filename)
    p.parts[2:]
    truncPath = pathlib.Path(*p.parts[2:])
    testPath = os.path.dirname(os.path.dirname(os.path.realpath(__file__)))
    print(testPath)
    print(os.path.join(
            testPath,
            os.path.normpath("dummy-resourcedb/test-data/uploads")))

    hostSidePath = \
        os.path.join(
            testPath,
            os.path.normpath("dummy-resourcedb/test-data/uploads"),
            truncPath)

    with open(hostSidePath) as file:
        lines = file.read()
        if lines != expectedContent:
            pytest.fail(
                f"Invalid file content\n\
                Returned: {lines}\n\
                Expected: {expectedContent}")
            return None
        file.close()
