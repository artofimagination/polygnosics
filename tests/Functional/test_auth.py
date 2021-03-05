import pytest

dataColumns = ("data", "expected")
createTestData = [
    (
        # Input data
        {
            "product": {
                "name": "testProductAddProject",
            },
            "user": {
                "name": "testUserOwnerAddProject",
                "email": "testEmailOwnerAddProject",
                "password": "testPassword"
            },
            "project": {
                "name": "testProjectAddProject",
                "visibility": "Public"
            }
        },
        # Expected
        {
            "name": "testProjectAddProject",
            "visibility": "Public"
        }),
    # Input data
    ({
      "project": {
        "name": "testProjectMissingUser",
        "visibility": "Public"
      },
      "user_id": "c34a7368-344a-11eb-adc1-0242ac120002",
      "product_id": "c34a7368-344a-11eb-adc1-0242ac120002"
    },
      # Expected
      "Failed to create product: Error 1452: " \
      "Cannot add or update a child row: a foreign key constraint fails" \
      " (`user_database`.`projects`, CONSTRAINT " \
      "`projects_ibfk_1` FOREIGN KEY (`products_id`) " \
      "REFERENCES `products` (`id`))")

]

ids = ['No existing project', 'Missing product']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_Signup(httpConnection, data, expected):
    pass
