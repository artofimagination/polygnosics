import pytest
import requests
import time


class HTTPConnectorBackend():
    def __init__(self):
        self.URL = "http://0.0.0.0:8184"
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

    def POST(self, address, json=None, files=None):
        url = self.URL + address
        return requests.post(url=url, json=json, files=files)


class HTTPConnectorDummyUserDB():
    def __init__(self):
        self.URL = "http://0.0.0.0:8183"
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
        self.URL = "http://0.0.0.0:8182"
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
