import pytest
import json
import common


dataColumns = ("data", "expected")
createTestData = [
    (
        # Input data
        {
            "username": "signupUser",
            "email": "signup@test.com",
            "password": "asd123ASD",
            "group": "client"
        },
        # Expected
        {
            "error": "",
            "response": {
              "data": "OK"
            },
            "requests": {
                common.ADD_USER_PATH: {
                    'email': 'signup@test.com',
                    'password': '$2a$16$YFDXFOyicyLMt4mLmdulX\
OCIAlK9YHqVQhUPArpMwlNyM6YUuviVS',
                    'username': 'signupUser'
                },
                common.UPDATE_USER_ASSETS_PATH: {
                    'user-data': {
                        'datamap': {
                            'base_asset_path': '/user-assets/9f02/fbd5/\
15b7/465a/a941/f4fd/c11d/b23e'
                        },
                        'id': '9f02fbd5-15b7-465a-a941-f4fdc11db23e'
                    },
                    'user-id': "026eede8-0b9b-4355-ad48-8a4f6cf0b49e"
                },
                common.UPDATE_USER_SETTINGS_PATH: {
                    'user-data': {
                        'datamap': {
                            'group': 'client',
                            'privileges': {
                                'delete_user': 0,
                                'edit_page': 0,
                                'main_dashboard': 0,
                                'misuse_metrics': 0,
                                'product_stats': 0,
                                'project_stats': 1
                            }
                        },
                        'id': '8b683a4c-198a-4cfd-abb1-7a3715a51bbb'
                    },
                    'user-id': "026eede8-0b9b-4355-ad48-8a4f6cf0b49e"
                }
            }
        }),
    (
        # Input data
        {
            "username": "signupUser",
            "email": "signup@test.com",
            "password": "asd123ASD",
            "group": "client"
        },
        # Expected
        {
            "error": "User with this name already exists",
            "requests": {
                common.ADD_USER_PATH: {
                    'email': 'signup@test.com',
                    'password': '$2a$16$YFDXFOyicyLMt4mLmdulX\
OCIAlK9YHqVQhUPArpMwlNyM6YUuviVS',
                    'username': 'signupUser'
                }
            }
        })
]

ids = ['Success', 'Failure']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_AddUser(httpBackend, httpDummyUserDB, data, expected):
    # Clears all previously stored incoming requests on the dummy server
    try:
        r = httpDummyUserDB.POST("/clear-request-data", None)
    except Exception:
        pytest.fail("Failed to send POST request")
        return None

    try:
        r = httpBackend.POST(
            "/add-user",
            data)
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

    try:
        r = httpDummyUserDB.GET("/get-request-data", None)
    except Exception:
        pytest.fail("Failed to send GET request")
        return None

    try:
        response = json.loads(r.text)
        if "password" in response[common.ADD_USER_PATH]:
            expected["requests"][common.ADD_USER_PATH]["password"] = \
                response[common.ADD_USER_PATH]["password"]
    except Exception:
        pytest.fail(f"Failed to decode response text {r.text}")
        return None

    expectedRequest = expected["requests"]
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
            "error": "",
            "response": {
                "assets": {
                    "datamap": {
                        "base_asset_path": "/user-assets/9f02\
/fbd5/15b7/465a/a941/f4fd/c11d/b23e"
                    },
                    "id": "9f02fbd5-15b7-465a-a941-f4fdc11db23e"
                },
                "email": "signup@test.com",
                "id": "026eede8-0b9b-4355-ad48-8a4f6cf0b49e",
                "username": "signupUser",
                "settings": {
                    "datamap": {
                        "group": "client",
                        "privileges": {
                            "delete_user": 0,
                            'edit_page': 0,
                            "main_dashboard": 0,
                            "misuse_metrics": 0,
                            "product_stats": 0,
                            "project_stats": 1
                        }
                    },
                    "id": "8b683a4c-198a-4cfd-abb1-7a3715a51bbb"
                }
            },
            "requests": {
                common.GET_USER_BY_EMAIL_PATH: \
                "/get-user-by-email?email=signup@test.com"
            }
        }),
    (
        # Input data
        {
            "email": "invalid@test.com",
            "password": "asd123ASD",
        },
        # Expected
        {
            "error": "Incorrect email or password",
            "response": "",
            "requests": {
                common.GET_USER_BY_EMAIL_PATH: \
                "/get-user-by-email?email=invalid@test.com"
            }
        })
]

ids = ['Success', 'Failure']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_Signin(httpBackend, httpDummyUserDB, data, expected):
    # Clears all previously stored incoming requests on the dummy server
    try:
        httpDummyUserDB.POST("/clear-request-data", None)
    except Exception:
        pytest.fail("Failed to send POST request")
        return None

    try:
        r = httpBackend.GET("/auth_login", data)
    except Exception:
        pytest.fail("Failed to send POST request")
        return None

    try:
        expectedData = expected["response"]
        response = common.getResponse(r.text, expected["error"])
        if response is not None and response != expectedData:
            pytest.fail(
                f"Request failed\n\
                Status code: {r.status_code}\n\
                Returned: {response}\n\
                Expected: {expectedData}")
    except Exception:
        pytest.fail("Failed to process request")
        return None

    try:
        r = httpDummyUserDB.GET("/get-request-data", None)
    except Exception:
        pytest.fail("Failed to send GET request")
        return None

    try:
        response = json.loads(r.text)
    except Exception:
        pytest.fail(f"Failed to decode response text {r.text}")
        return None

    expectedRequest = expected["requests"]
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
            "username": "root",
            "email": "root@test.com",
            "password": "asd123ASD",
            "group": "root"},
        # Expected
        {
            "error": "",
            "response": True,
            "requests": {
                common.GET_USER_BY_EMAIL_PATH: \
                "/get-user-by-email?email=root@test.com"
            }
        }),
    (
        # Input data
        {
            "id": "f9ebc23d-81cc-4bf2-b908-7e88c58ebe91"
        },
        # Expected
        {
            "error": "",
            "response": False,
            "requests": {
                common.GET_USER_BY_EMAIL_PATH: \
                "/get-user-by-email?email=root@test.com"
            }
        })
]

ids = ['Success', 'Failure']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_DetectRootUser(httpBackend, httpDummyUserDB, data, expected):
    if "username" in data:
        try:
            r = httpDummyUserDB.POST(
                "/add-user",
                data)
        except Exception:
            pytest.fail("Failed to send POST request")
            return None
    else:
        try:
            r = httpDummyUserDB.POST(
                "/delete-user",
                data)
        except Exception:
            pytest.fail("Failed to send POST request")
            return None

    # Clears all previously stored incoming requests on the dummy server
    try:
        httpDummyUserDB.POST("/clear-request-data", None)
    except Exception:
        pytest.fail("Failed to send POST request")
        return None

    try:
        r = httpBackend.GET("/detect-root-user", None)
    except Exception:
        pytest.fail("Failed to send POST request")
        return None

    try:
        expectedData = expected["response"]
        response = common.getResponse(r.text, expected["error"])
        if response is not None and response != expectedData:
            pytest.fail(
                f"Request failed\n\
                Status code: {r.status_code}\n\
                Returned: {response}\n\
                Expected: {expectedData}")
    except Exception:
        pytest.fail("Failed to process request")
        return None

    try:
        r = httpDummyUserDB.GET("/get-request-data", None)
    except Exception:
        pytest.fail("Failed to send GET request")
        return None

    try:
        response = json.loads(r.text)
    except Exception:
        pytest.fail(f"Failed to decode response text {r.text}")
        return None

    expectedRequest = expected["requests"]
    if response != expectedRequest:
        pytest.fail(
            f"Request failed\n\
            Status code: {r.status_code}\n\
            Returned: {response}\n\
            Expected: {expectedRequest}")
