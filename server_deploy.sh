#!/bin/bash

git fetch
git checkout master
git reset --hard origin/master
sudo docker-compose build
sudo docker-compose down
sudo docker-compose up -d
