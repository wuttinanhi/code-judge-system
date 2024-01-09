import paramiko
import subprocess
import time

def get_manager_server_ip() -> str:
    # get the manager server IP address
    return subprocess.check_output(['terraform', 'output', '-raw', 'manager_first_node_public_ip']).decode('utf-8').strip()

def copy_file_to_server(target_server: str, username: str, local_file_path: str, remote_file_path: str):
    ssh = paramiko.SSHClient()
    ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    ssh.connect(hostname=target_server, username=username, key_filename="ssh/id_rsa")

    sftp = ssh.open_sftp()
    sftp.put(local_file_path, remote_file_path)

    sftp.close()
    ssh.close()

def ssh_exec_command(target_server_ip: str, username: str, command: str):
    ssh = paramiko.SSHClient()
    ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())

    ssh.connect(
        hostname=target_server_ip, 
        username=username, 
        key_filename="ssh/id_rsa",
        banner_timeout=200
    )

    _, stdout, stderr = ssh.exec_command(command)
    stdout = stdout.readline()
    stderr = stderr.readline()
    
    ssh.close()

    return stdout, stderr

def wait_for_docker(target_server_ip: str, username: str):
    while True:
        try:
            stdout, stderr = ssh_exec_command(target_server_ip, username, 'docker info')
            print(stdout)
            if stderr:
                print(stderr)
                print("Waiting for Docker to ready...")
            else:
                return
        except:
            pass
        time.sleep(1)
