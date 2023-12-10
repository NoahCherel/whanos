# Use the latest LTS version of Jenkins as the base image
FROM jenkins/jenkins:lts

# Disable Jenkins setup wizard
ENV JAVA_OPTS -Djenkins.install.runSetupWizard=false

# Switch to root user for installation
USER root

# Install necessary tools and dependencies
RUN apt-get update && \
    apt-get -qy upgrade && \
    apt-get install -qy curl

# Install Docker
RUN curl -fsSL https://get.docker.com/ | sh

# Set the initial admin password
ENV ADMIN_PASSWORD=test

# Define the location for Jenkins Configuration as an environment variable
ENV CASC_JENKINS_CONFIG=/var/jenkins_home/jenkins.yml

# Copy the list of plugins to be installed
COPY ./jenkins/plugins.txt /usr/share/plugins.txt

# Install Jenkins plugins using the plugin list
RUN jenkins-plugin-cli --plugins -f /usr/share/plugins.txt
