#!/bin/bash

# create directory for SSH key pair if it doesn't exist
mkdir -p ssh

# Generate a new SSH key pair with a comment
ssh-keygen -m PEM -t rsa -b 4096 -f ./ssh/id_rsa -q -N "" -C "docker"

echo "SSH key pair generated."
