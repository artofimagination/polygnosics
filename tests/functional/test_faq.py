import pytest
import json
import common


dataColumns = ("data", "expected")
createTestData = [
    (
        # Input data
        {
            'group': "Test group 1",
            "answer": "This is answer 1",
            "question": "This is question 1",
        },
        # Expected
        {
            "requests": [
              {
                  '/add-resource': {
                      'category': 2,
                      'content': {
                          'group': 'Test group 1'
                      },
                      'id': '00000000-0000-0000-0000-000000000000'}
              }],
            "refs": {},
            "error": ""
        }
    )
]

ids = ['Success']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_AddFAQ(httpBackend, httpDummyResourceDB, data, expected):
    # Clears all previously stored incoming requests on the dummy server
    try:
        r = httpDummyResourceDB.POST("/clear-request-data", None)
    except Exception:
        pytest.fail("Failed to send POST request")
        return None

    try:
        r = httpBackend.POST(
            "/resources/create-faq-item",
            data=data)
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

    # Check the requests sent from the backend
    try:
        r = httpDummyResourceDB.GET("/get-request-data", None)
    except Exception:
        pytest.fail("Failed to send GET request")
        return None

    try:
        response = json.loads(r.text)
    except Exception:
        pytest.fail(f"Failed to decode response text {r.text}")
        return None

    # Folder is a unique UUID for FAQ questions and answers,
    # so we need to delete before compare
    # Being able to delete it shows, whether the content is there.
    (answerValue, deleted) = common.deleteElement("answer", response)
    if deleted is False:
        pytest.fail(
            "Request failed. \"answer\" is missing")
        return None

    (questionValue, deleted) = common.deleteElement("question", response)
    if deleted is False:
        pytest.fail(
            "Request failed. \"question\" is missing")
        return None

    if common.checkFileContent(answerValue, data["answer"]) is None:
        return None

    if common.checkFileContent(questionValue, data["question"]) is None:
        return None

    with open(answerValue) as file:
        lines = file.readlines()
        answer = data["answer"]
        if lines != answer:
            pytest.fail(
                f"Invalid file content\n\
                Returned: {lines}\n\
                Expected: {answer}")

    with open(questionValue) as file:
        lines = file.readlines()
        question = data["question"]
        if lines != answer:
            pytest.fail(
                f"Invalid file content\n\
                Returned: {lines}\n\
                Expected: {question}")

    expectedRequest = expected["requests"]
    if response != expectedRequest:
        pytest.fail(
            f"Request failed\n\
            Status code: {r.status_code}\n\
            Returned: {response}\n\
            Expected: {expectedRequest}")
