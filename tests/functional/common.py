import json
import pytest


# getResponse unwraps the data/error from json response.
def getResponse(responseText, expected=None):
    response = json.loads(responseText)
    if "error" in response:
        error = response["error"]
        if expected is None or (expected is not None and error != expected):
            pytest.fail(f"Failed to run test.\nDetails: {error}")
        return None
    return response["data"]
