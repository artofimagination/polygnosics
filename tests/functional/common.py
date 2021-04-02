import json
import pytest

ADD_USER_PATH = "/add-user"
GET_USER_BY_EMAIL_PATH = "/get-user-by-email"
UPDATE_USER_SETTINGS_PATH = "/update-user-settings"
UPDATE_USER_ASSETS_PATH = "/update-user-assets"


# getResponse unwraps the data/error from json response.
def getResponse(responseText, expected=None):
    response = json.loads(responseText)
    if "error" in response:
        error = response["error"]
        if expected is None or (expected is not None and error != expected):
            pytest.fail(f"Failed to run test.\nDetails: {error}")
        return None
    return response["data"]
