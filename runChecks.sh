#!/bin/bash

cp ./tests/.env_functional_test .env
if [ ! -d ./tests/dummy-userdb/uploads ]; then
  mkdir ./tests/dummy-userdb/uploads
fi

cp ./tests/dummy-userdb/test-data/testData.json ./tests/dummy-userdb/test-data/uploads/
go clean -testcache && go test -v -count=1 -cover ./...
status=$?; 
if [[ $status != 0 ]]; then 
  exit $status; 
fi
golangci-lint run -v .
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