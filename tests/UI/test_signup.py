import pytest
from pages.signup import SignUpPage

dataColumns = ("data", "expected")
createTestData = [
    ({
        'title': 'Polygnosics - Sign up'
    },
        "test-image:latest"),
]

ids = ['Elements exist']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_LoadingSignupPage(browser, data, expected):
    signupPage = SignUpPage(browser)
    signupPage.load()

    assert data["title"] == signupPage.title()
