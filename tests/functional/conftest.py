import pytest
import requests
import time
import os


def getAttributes():
    variables = {}
    fileName = os.path.dirname(os.path.realpath(__file__)) + \
        "/../.env.functional_test"
    with open(fileName) as envFile:
        for line in envFile:
            name, var = line.partition("=")[::2]
            variables[name.strip()] = var.strip()
        return variables


class HTTPConnectorBackend():
    def __init__(self):
        self.URL = "http://127.0.0.1:" + getAttributes()["BACKEND_SERVER_PORT"]
        connected = False
        timeout = 15
        while timeout > 0:
            try:
                r = self.GET("/", "")
                if r.status_code == 200:
                    connected = True
                break
            except Exception:
                timeout -= 1
                time.sleep(1)

        if connected is False:
            raise Exception("Cannot connect to test server")

    def GET(self, address, params):
        url = self.URL + address
        return requests.get(url=url, params=params)

    def POST(self, address, json=None, files=None, data=None):
        url = self.URL + address
        return requests.post(url=url, json=json, files=files, data=data)


class HTTPConnectorDummyUserDB():
    def __init__(self):
        self.URL = "http://127.0.0.1:" + getAttributes()["USER_DB_PORT"]
        print(self.URL)
        connected = False
        timeout = 15
        while timeout > 0:
            try:
                r = self.GET("/", "")
                if r.status_code == 200:
                    connected = True
                break
            except Exception:
                timeout -= 1
                time.sleep(1)

        if connected is False:
            raise Exception("Cannot connect to test server")

    def GET(self, address, params):
        url = self.URL + address
        return requests.get(url=url, params=params)

    def POST(self, address, json):
        url = self.URL + address
        return requests.post(url=url, json=json)


class HTTPConnectorDummyResourceDB():
    def __init__(self):
        self.URL = "http://127.0.0.1:" + getAttributes()["RESOURCE_DB_PORT"]
        connected = False
        timeout = 15
        while timeout > 0:
            try:
                r = self.GET("/", "")
                if r.status_code == 200:
                    connected = True
                break
            except Exception:
                timeout -= 1
                time.sleep(1)

        if connected is False:
            raise Exception("Cannot connect to test server")

    def GET(self, address, params):
        url = self.URL + address
        return requests.get(url=url, params=params)

    def POST(self, address, json):
        url = self.URL + address
        return requests.post(url=url, json=json)


@pytest.fixture
def httpBackend():
    return HTTPConnectorBackend()


@pytest.fixture
def httpDummyUserDB():
    return HTTPConnectorDummyUserDB()


@pytest.fixture
def httpDummyResourceDB():
    return HTTPConnectorDummyResourceDB()
