import pytest
from pages.index_header import IndexPageHeader

dataColumns = ("data", "expected")
createTestData = [
    ({
        'title': 'Polygnosics'
    },
        "test-image:latest"),
]

ids = ['Elements exist']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_LoadingPageHeader(browser, data, expected):
    indexPageHeader = IndexPageHeader(browser)
    indexPageHeader.load()

    assert data["title"] == indexPageHeader.title()

    try:
        indexPageHeader.waitUntilPreloaderDisappears()
    except Exception as e:
        assert f"{e}" == ""

    try:
        indexPageHeader.maximizeWindow()
        indexPageHeader.goToSignup()
    except Exception as e:
        assert f"{e}" == ""
