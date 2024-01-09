import util

manager_server_ip = util.get_manager_server_ip() 

print(f"Starting to deploy stack (manager ip: {manager_server_ip})")

# Copy docker-compose.yml and .env to manager server
util.copy_file_to_server(
    manager_server_ip,
    "docker",
    '../../docker-compose.yml',
    '/home/docker/docker-compose.yml'
)

util.copy_file_to_server(
    manager_server_ip,
    "docker",
    '../../.env.prod',
    '/home/docker/.env'
)

stdout, stderr = util.ssh_exec_command(
    manager_server_ip, 
    "docker", 
    "export $(grep -v '^#' /home/docker/.env | xargs) > /dev/null 2>&1 && docker stack deploy --compose-file /home/docker/docker-compose.yml cjs",
)

print(stdout)
print(stderr)

print("Stack deployed successfully")
