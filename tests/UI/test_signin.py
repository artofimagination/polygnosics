import pytest
from pages.signin import SigninPage

dataColumns = ("data", "expected")
createTestData = [
    ({
        'title': 'Polygnosics - Sign in'
    },
        "test-image:latest"),
]

ids = ['Elements exist']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_LoadingSigninPage(browser, data, expected):
    signinPage = SigninPage(browser)
    signinPage.load()

    assert data["title"] == signinPage.title()
