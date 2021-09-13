#!/bin/bash

pip3 install -r tests/requirements.txt

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

docker-compose down
docker stop $(docker ps -aq)
docker rm $(docker ps -aq)
docker system prune -f
docker-compose up --build --force-recreate -d backend
status=$?; 
if [[ $status != 0 ]]; then 
  exit $status; 
fi
python3 -m pytest -v tests/functional