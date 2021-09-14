#!/bin/bash

cp ./tests/.env.functional_test .env
if [ ! -d ./tests/dummy-userdb/test-data/uploads ]; then
  mkdir ./tests/dummy-userdb/test-data/uploads
fi

if [ ! -d ./tests/dummy-resourcedb/test-data/uploads ]; then
  mkdir ./tests/dummy-resourcedb/test-data/uploads
  mkdir ./tests/dummy-resourcedb/test-data/uploads/files
fi

cp ./tests/dummy-userdb/test-data/testData.json ./tests/dummy-userdb/test-data/uploads/
cp ./tests/dummy-resourcedb/test-data/testData.json ./tests/dummy-resourcedb/test-data/uploads/