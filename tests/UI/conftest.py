"""
This module contains shared fixtures.
"""

import pytest
import selenium.webdriver
import requests
import time


def _pingServer():
    URL = "http://0.0.0.0:8081"
    connected = False
    timeout = 15
    while timeout > 0:
        try:
            r = requests.get(url=URL)
            if r.status_code == 200:
                connected = True
                break
        except Exception:
            timeout -= 1
            time.sleep(1)

    if connected is False:
        raise Exception("Cannot connect to test server")


@pytest.fixture(scope="session")
def browser():
    _pingServer()

    # Initialize the ChromeDriver instance
    b = selenium.webdriver.Firefox()

    # Make its calls wait up to 10 seconds for elements to appear
    b.implicitly_wait(10)

    # Return the WebDriver instance for the setup
    yield b

    # Quit the WebDriver instance for the cleanup
    b.quit()
