#!/bin/bash

echo "Starting Whanos"

# Function to get the latest commit hash
get_latest_commit_hash() {
    git log -n 1 --pretty=format:"%H"
}

# Make a copy of the Dockerfile
# Then
# Change the first line of Dockerfile
# Exemple FROM node:latest to FROM localhost:5000/whanos-node:latest
change_dockerfile() {
    local language=$1
    cp ./Dockerfile ./Dockerfile.bak
    sed -i "1s/.*/FROM localhost:5000\/whanos-$language:latest/" ./Dockerfile
}

# Function to build Docker image based on language
build_docker_image() {
    local language=$1
    if [ -f "./Dockerfile" ]; then
        echo "Using base image"
        change_dockerfile "$language"
        docker build . -t whanos-project-$1
    else
        echo "Using standalone image"
        docker build . -t whanos-project-$1 -f /images/$language/Dockerfile.standalone
    fi
}

# Function to deploy on Kubernetes
deploy_on_kubernetes() {
    if [ -f "./whanos.yml" ]; then
        echo "Deploying on Kubernetes"
        file_content=$(cat ./whanos.yml | base64 -w 0)
        echo $file_content
        curl -H "Content-Type: application/json" -X POST -d "{\"image\":\"localhost:5000/whanos-project-$1\",\"config\":\"$file_content\",\"name\":\"$1\"}" http://localhost:3030/deployments
    fi
}

echo "Starting Whanos"

# Check if there is a stored commit hash for the project
if [ -f "/usr/share/jenkins_hash/JENKINS_HASH_$1" ]; then
    last_commit=$(cat "/usr/share/jenkins_hash/JENKINS_HASH_$1")
fi

# Compare the last stored commit hash with the latest commit in the repository
if [ "$last_commit" != "$(get_latest_commit_hash)" ]; then
    echo "Changes occurred, containerization needed"
    language=$(/var/jenkins_home/getLanguage.sh .)

    # Check for an error while detecting the programming language
    if [ $? -eq 1 ]; then
        echo "Error occurred getting language"
        exit 1
    fi

    echo "Detected language: $language"

    # Build the Docker image based on the detected language
    build_docker_image "$language"

    # Tag, push, pull, and clean up the Docker image
    echo "Tagging $1"
    docker tag whanos-project-$1 localhost:5000/whanos-project-$1
    echo "Pushing"
    docker push localhost:5000/whanos-project-$1
    echo "Pulling"
    docker pull localhost:5000/whanos-project-$1
    #docker rmi whanos-project-$1

    # Deploy on Kubernetes if a configuration file exists
    deploy_on_kubernetes "$1"

    # Update the stored commit hash
    mkdir -p /usr/share/jenkins_hash
    get_latest_commit_hash > "/usr/share/jenkins_hash/JENKINS_HASH_$1"
else
    echo "No changes occurred"
fi