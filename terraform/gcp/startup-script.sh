#!/bin/bash

# install docker
curl -fsSL https://get.docker.com | bash

# Create the user 'docker'
useradd -m docker

# Create the group 'docker' if it doesn't exist
getent group docker || groupadd docker

# Add the user 'docker' to the group 'docker'
usermod -aG docker docker
