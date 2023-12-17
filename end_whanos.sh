#!/bin/bash

# Stop and remove the first two containers
docker ps -aq | head -n 2 | xargs -r docker rm -f

# Remove the Docker image with the specified repository and tag
docker rmi -f whanos-jenkins

minikube stop

# Stop and remove a specific container
docker rm -f whanos-registry

# Remove Docker images with a specific repository and wildcard tag
docker images localhost:5000/whanos-project* -q | xargs -r docker rmi -f

# Kill processes using specific ports
fuser -k 8080/tcp
fuser -k 3030/tcp
