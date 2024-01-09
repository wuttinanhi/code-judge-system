import paramiko
import subprocess
import time
import util

# get the manager server IP address
manager_server_ip = util.get_manager_server_ip()

# get username of the current user
username = "docker"

print(f'Connecting to manager server at {manager_server_ip}...')

# Check if host is ready and Docker is running
util.wait_for_docker(manager_server_ip, username)

# Create SSH client
ssh = paramiko.SSHClient()
ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())

# Replace with your SSH server details
ssh.connect(
    hostname=manager_server_ip, 
    username=username, 
    key_filename="ssh/id_rsa"
)

# initialize the swarm
stdin, stdout, stderr = ssh.exec_command('docker swarm init')
print(stderr.readline())
print(stdout.readline())

stdin, stdout, stderr = ssh.exec_command('docker swarm join-token -q worker')
output = stdout.readline()
print(output)
print(stderr.readline())

outlines = output
resp = ''.join(outlines)

print(f'Docker Swarm token: {resp}')

# save the token to "./swarm-worker-token"
with open('./swarm-worker-token', 'w') as f:
    f.write(resp)

ssh.close()
