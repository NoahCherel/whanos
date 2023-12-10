FROM jenkins/jenkins:lts

ENV JAVA_OPTS -Djenkins.install.runSetupWizard=false

USER root

#install docker
RUN apt-get update && \
    apt-get -qy full-upgrade && \
    apt-get install -qy curl && \
    curl -sSL https://get.docker.com/ | sh

ENV ADMIN_PASSWORD=test

ENV CASC_JENKINS_CONFIG=/var/jenkins_home/jenkins.yml

COPY ./jenkins/plugins.txt /usr/share/plugins.txt

RUN jenkins-plugin-cli --plugins -f /usr/share/plugins.txt