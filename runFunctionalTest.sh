#!/bin/bash

pip3 install -r tests/requirements.txt

./copyTestFiles.sh
status=$?; 
if [[ $status != 0 ]]; then 
  exit $status; 
fi

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