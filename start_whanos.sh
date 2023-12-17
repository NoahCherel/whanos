#!/bin/bash

#minikube start

# Start Kubernetes API server
./kubernetes/kube_binary_api_server &

# Set insecure registries configuration
echo '{ "insecure-registries":    ["localhost:5000"] }' | sudo tee /etc/docker/daemon.json
sudo systemctl restart docker

# Remove and checkout Jenkins directory
sudo rm -rf jenkins/
git checkout jenkins/

# Build Jenkins Docker image
docker build . -f jenkins/Dockerfile -t whanos-jenkins

# Run Jenkins container
cd jenkins
docker run -d \
     -v $(pwd):/var/jenkins_home \
     -v $(pwd)/../images:/images \
     --net=host \
     -v /var/run/docker.sock:/var/run/docker.sock \
     `docker images -aq | head -n 1`

# Register Jenkins Docker image in local registry
cd ..
docker exec -it `docker ps -aq | head -n 1` \
     docker run -d \
      -p 5000:5000 \
       --restart=always \
       --name whanos-registry \
       registry:2

# Access Jenkins container shell
docker exec -it `docker ps -aq | head -n 2 | tail -n +2` /bin/bash
