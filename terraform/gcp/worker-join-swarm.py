import subprocess
import paramiko
import json
import util

# get the manager server IP address
manager_server_ip = util.get_manager_server_ip()

# Get the worker IPs from Terraform output
output = subprocess.check_output(['terraform', 'output', '-json', 'worker_ips'])
worker_ips = json.loads(output.decode('utf-8').strip())

# Read the Docker Swarm token from the file
with open('./swarm-worker-token', 'r') as f:
    token = f.read().strip()

for worker_ip in worker_ips:
    # Join the Docker Swarm
    stdout, stderr = util.ssh_exec_command(worker_ip, "docker", f'docker swarm join --token {token} {manager_server_ip}:2377')
    print(stdout)
    print(stderr)

print("All workers have joined the Docker Swarm.")
