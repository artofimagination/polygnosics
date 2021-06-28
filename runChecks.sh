#!/bin/bash

cp ./tests/.env_functional_test .env
if [ ! -d ./tests/dummy-userdb/test-data/uploads ]; then
  mkdir ./tests/dummy-userdb/test-data/uploads
fi

if [ ! -d ./tests/dummy-resourcedb/test-data/uploads ]; then
  mkdir ./tests/dummy-resourcedb/test-data/uploads
  mkdir ./tests/dummy-resourcedb/test-data/uploads/files
fi

cp ./tests/dummy-userdb/test-data/testData.json ./tests/dummy-userdb/test-data/uploads/
cp ./tests/dummy-resourcedb/test-data/testData.json ./tests/dummy-resourcedb/test-data/uploads/
go clean -testcache && go test -v -count=1 -cover ./...
status=$?; 
if [[ $status != 0 ]]; then 
  exit $status; 
fi
golangci-lint run -v ./...
status=$?; 
if [[ $status != 0 ]]; then 
  exit $status; 
fi
flake8 . --count --show-source --statistics --exclude=temp
status=$?; 
if [[ $status != 0 ]]; then 
  exit $status; 
fi
./runFunctionalTest.sh
status=$?; 
if [[ $status != 0 ]]; then 
  exit $status; 
fi