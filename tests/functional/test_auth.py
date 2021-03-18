import pytest
import json


dataColumns = ("data", "expected")
createTestData = [
    (
        # Input data
        {
            "username": "signupUser",
            "email": "signup@test.com",
            "password": "asd123ASD"
        },
        # Expected
        {
            "data": "Registration successful",
            "request": {
                "username": "signupUser",
                "email": "signup@test.com",
                "password": "asd123ASD",
                "group": "client"
            }
        }),
    (
        # Input data
        {
            "username": "root",
            "email": "signup@test.com",
            "password": "asd123ASD",
        },
        # Expected
        {
            "data": "User already exists",
            "request": {
                "username": "root",
                "email": "signup@test.com",
                "password": "asd123ASD",
                "group": "client"
            }
        })
]

ids = ['Success', 'Failure']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_Signup(httpFrontend, httpDummyBackend, data, expected):
    try:
        r = httpFrontend.POST(
            "/auth_signup",
            data)
    except Exception:
        pytest.fail("Failed to send POST request")
        return None

    expectedData = expected["data"]
    if r.text != expectedData:
        pytest.fail(
            f"Request failed\n\
    Status code: {r.status_code}\n\
    Returned: {r.text}\n\
    Expected: {expectedData}")

    try:
        r = httpDummyBackend.GET("/get-request-data", None)
    except Exception:
        pytest.fail("Failed to send GET request")
        return None

    try:
        response = json.loads(r.text)
    except Exception:
        pytest.fail(f"Failed to decode response text {r.text}")
        return None

    expectedRequest = expected["request"]
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
            "email": "signup@test.com",
            "password": "asd123ASD",
        },
        # Expected
        {
            "data": "Login successful",
            "request": {
                'uri': '/login?email=signup@test.com&password=asd123ASD'
            }
        }),
    (
        # Input data
        {
            "email": "incorrect@test.com",
            "password": "asd123ASD",
        },
        # Expected
        {
            "data": "Failed to login. Incorrect email or password",
            "request": {
                'uri': '/login?email=incorrect@test.com&password=asd123ASD'
            }
        }),
]

ids = ['Success', 'Failure']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_Signin(httpFrontend, httpDummyBackend, data, expected):
    try:
        r = httpFrontend.POST("/auth_login", data)
    except Exception:
        pytest.fail("Failed to send POST request")
        return None

    expectedData = expected["data"]
    if r.text != expectedData:
        pytest.fail(
            f"Request failed\n\
            Status code: {r.status_code}\n\
            Returned: {r.text}\n\
            Expected: {expectedData}")

    try:
        r = httpDummyBackend.GET("/get-request-data", None)
    except Exception:
        pytest.fail("Failed to send GET request")
        return None

    try:
        response = json.loads(r.text)
    except Exception:
        pytest.fail(f"Failed to decode response text {r.text}")
        return None

    expectedRequest = expected["request"]
    if response != expectedRequest:
        pytest.fail(
            f"Request failed\n\
            Status code: {r.status_code}\n\
            Returned: {response}\n\
            Expected: {expectedRequest}")
