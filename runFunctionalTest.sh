pip3 install -r tests/requirements.txt
docker-compose down
docker stop $(docker ps -aq)
docker rm $(docker ps -aq)
docker system prune -f
docker-compose --file docker-compose-ui-test.yml up --build --force-recreate -d backend
python3 -m pytest -v tests/functional --soft-asserts