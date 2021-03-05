# python -m pip install --upgrade pip
# pip3 install -r tests/requirements.txt
# sudo apt-get install firefox
# seleniumbase install geckodriver
docker-compose down
docker stop $(docker ps -aq)
docker rm $(docker ps -aq)
docker system prune -f
docker-compose up --build --force-recreate -d polygnosics

python3 -m pytest -v tests/UI