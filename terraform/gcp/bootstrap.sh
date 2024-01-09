./generate-ssh-key.sh

terraform apply -auto-approve -var-file="variables.tfvars"

pip install paramiko
python3 setup-docker-swarm.py
python3 worker-join-swarm.py
python3 deploy-stack.py
