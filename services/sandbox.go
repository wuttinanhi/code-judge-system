package services

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/wuttinanhi/code-judge-system/entities"
)

type SandboxService interface {
	Run(instance *entities.SandboxInstance) (*entities.SandboxInstance, error)
}

type sandboxService struct {
	DockerClient *client.Client
}

func (s *sandboxService) pullImage(ctx context.Context, imageName string) error {
	log.Println("pulling image", imageName)
	out, err := s.DockerClient.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer out.Close()

	var logs bytes.Buffer
	io.Copy(&logs, out)

	log.Println("pulling image", imageName, "done")

	return nil
}

func (s *sandboxService) imageExist(imageName string) (bool, error) {
	ctx := context.Background()
	_, _, err := s.DockerClient.ImageInspectWithRaw(ctx, imageName)
	if err != nil {
		if client.IsErrNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func CreateFileWrapper(path, content string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func (s *sandboxService) CreateTempCodeFile(instance *entities.SandboxInstance) error {
	os.MkdirAll("/tmp/code-judge-system", os.ModePerm)
	unixTime := time.Now().UnixNano()

	codefileName := fmt.Sprintf("/tmp/code-judge-system/code-%d.py", unixTime)
	stdinfileName := fmt.Sprintf("/tmp/code-judge-system/stdin-%d.py", unixTime)

	err := CreateFileWrapper(codefileName, instance.Code)
	if err != nil {
		return err
	}

	err = CreateFileWrapper(stdinfileName, instance.Stdin)
	if err != nil {
		return err
	}

	instance.CodeFilePath = codefileName
	instance.StdinFilePath = stdinfileName

	return nil
}

func (s *sandboxService) DeleteTempCodeFile(instance *entities.SandboxInstance) error {
	err := os.Remove(instance.CodeFilePath)
	if err != nil {
		return err
	}
	err = os.Remove(instance.StdinFilePath)
	if err != nil {
		return err
	}
	return nil
}

func (s sandboxService) getLog(containerID string, showStdout, showStderr bool) (string, error) {
	ctx := context.Background()
	logs, err := s.DockerClient.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
		ShowStdout: showStdout,
		ShowStderr: showStderr,
	})
	if err != nil {
		return "", err
	}
	defer logs.Close()

	var buf bytes.Buffer
	io.Copy(&buf, logs)
	logsString := buf.String()
	logsString = strings.ReplaceAll(logsString, "\r\n", "\n")

	return logsString, nil
}

func (s sandboxService) getContainerExitCode(containerID string) (int, error) {
	ctx := context.Background()
	resp, err := s.DockerClient.ContainerInspect(ctx, containerID)
	if err != nil {
		return 0, err
	}
	return resp.State.ExitCode, nil
}

func (s sandboxService) createVolume(name string) (volume.Volume, error) {
	ctx := context.Background()
	volume, err := s.DockerClient.VolumeCreate(ctx, volume.CreateOptions{
		Name:   name,
		Driver: "local",
	})
	return volume, err
}

func (s sandboxService) deleteVolume(v volume.Volume) error {
	ctx := context.Background()
	err := s.DockerClient.VolumeRemove(ctx, v.Name, true)
	return err
}

func (s sandboxService) copyToContainer(containerID, targetPath string, content io.Reader) error {
	ctx := context.Background()

	data, err := io.ReadAll(content)
	if err != nil {
		return err
	}

	targetDir := filepath.Dir(targetPath)
	fileName := filepath.Base(targetPath)

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	tw.WriteHeader(&tar.Header{
		Name:   fileName,
		Mode:   0777,
		Size:   int64(len(data)),
		Format: tar.FormatGNU,
	})
	tw.Write(data)
	tw.Close()

	err = s.DockerClient.CopyToContainer(ctx, containerID, targetDir, &buf, types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: false,
		CopyUIDGID:                false,
	})
	return err
}

func (s sandboxService) createContainer(imageName string, command []string, volumes []mount.Mount, memoryLimit int64) (response container.CreateResponse, err error) {
	ctx := context.Background()
	response, err = s.DockerClient.ContainerCreate(ctx, &container.Config{
		Image:           imageName,
		NetworkDisabled: true,
		Tty:             true,
		AttachStdout:    true,
		AttachStderr:    true,
		AttachStdin:     true,
		OpenStdin:       true,
		Env:             []string{"PYTHONUNBUFFERED=1"},
		Entrypoint:      command,
	},
		&container.HostConfig{
			Mounts: volumes,
			Resources: container.Resources{
				Memory: memoryLimit,
			},
		},
		nil,
		nil,
		"",
	)
	return
}

func (s sandboxService) startContainer(containerID string) error {
	ctx := context.Background()
	err := s.DockerClient.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
	return err
}

func (s sandboxService) stopContainer(containerID string) error {
	ctx := context.Background()
	timeout := int(0)
	err := s.DockerClient.ContainerStop(ctx, containerID, container.StopOptions{
		Timeout: &timeout,
		Signal:  "SIGKILL",
	})
	return err
}

func (s sandboxService) removeContainer(containerID string) error {
	ctx := context.Background()
	err := s.DockerClient.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
	return err
}

// Run implements SandboxService.
func (s *sandboxService) Run(instance *entities.SandboxInstance) (*entities.SandboxInstance, error) {
	ctx := context.Background()

	instruction := entities.GetSandboxInstructionByLanguage(instance.Language)
	if instruction == nil {
		return nil, fmt.Errorf("language %s not supported", instance.Language)
	}

	imageName := instruction.DockerImage
	if imageName == "" {
		return nil, fmt.Errorf("language %s not supported", instance.Language)
	}

	// check if image exist
	exist, err := s.imageExist(imageName)
	if err != nil {
		return nil, err
	}

	// if image not exist, pull image
	if !exist {
		// pull image
		err := s.pullImage(ctx, imageName)
		if err != nil {
			return nil, err
		}
	}

	runID := strconv.Itoa(int(time.Now().UnixNano()))
	volumeID := "code-judge-system-" + runID
	volumeMount := []mount.Mount{
		{Type: mount.TypeVolume, Source: volumeID, Target: "/sandbox"},
	}

	// create volume
	volume, err := s.createVolume(volumeID)
	if err != nil {
		return nil, errors.New("failed to create volume")
	}
	defer s.deleteVolume(volume)

	// create container to create necessary files
	resp, err := s.createContainer(
		imageName,
		[]string{"/bin/bash", "-c", "chmod 777 -R /sandbox && sleep 9999"},
		volumeMount,
		entities.SandboxMemoryMB*128,
	)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	defer s.removeContainer(resp.ID)

	// start container
	err = s.startContainer(resp.ID)
	if err != nil {
		return nil, errors.New("create stage: failed to start container")
	}

	time.Sleep(1 * time.Second)

	// copy code file to container
	err = s.copyToContainer(resp.ID, "/sandbox/code", strings.NewReader(instance.Code))
	if err != nil {
		return nil, err
	}
	// copy stdin file to container
	err = s.copyToContainer(resp.ID, "/sandbox/stdin", strings.NewReader(instance.Stdin))
	if err != nil {
		return nil, err
	}

	// stop container
	err = s.stopContainer(resp.ID)
	if err != nil {
		return nil, errors.New("create stage: failed to stop container")
	}

	// remove container
	err = s.stopContainer(resp.ID)
	if err != nil {
		return nil, errors.New("create stage: failed to remove container")
	}

	// compile stage
	compileCommand := instruction.CompileCmd

	// create container to compile
	resp, err = s.createContainer(imageName, []string{"/bin/bash", "-c", compileCommand}, volumeMount, entities.SandboxMemoryMB*256)
	if err != nil {
		return nil, errors.New("compile stage: failed to create container")
	}
	defer s.removeContainer(resp.ID)

	// start container
	err = s.startContainer(resp.ID)
	if err != nil {
		return nil, errors.New("compile stage: failed to start container")
	}

	// grab container wait channel
	resultC, errC := s.DockerClient.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	if err != nil {
		return nil, errors.New("compile stage: failed to wait container")
	}

	// wait for container to finish or timeout
	select {
	case <-time.After(time.Duration(instance.Timeout) * time.Millisecond):
		// timeout
		err = s.stopContainer(resp.ID)
		if err != nil {
			return nil, errors.New("compile stage: failed to stop container")
		}
		return nil, errors.New("compile stage: timeout")
	case <-resultC:
		// container finish
		// do nothing
	case <-errC:
		// container error
		return nil, errors.New("compile stage: failed to compile code")
	}

	// get container exit code
	exitCode, err := s.getContainerExitCode(resp.ID)
	if err != nil {
		return nil, errors.New("compile stage: failed to get container exit code")
	}

	// if exit code is not 0, return error
	if exitCode != 0 {
		logs, err := s.getLog(resp.ID, true, true)
		if err != nil {
			return nil, errors.New("compile stage: failed to get container log")
		}
		instance.CompileExitCode = exitCode
		instance.CompileStderr = logs
		return nil, errors.New("compile stage: failed to compile code")
	}

	// stop and remove container
	err = s.removeContainer(resp.ID)
	if err != nil {
		return nil, errors.New("compile stage: failed to stop and remove container")
	}

	// run stage
	runCommand := instruction.RunCmd

	// create container to run
	resp, err = s.createContainer(imageName, []string{"/bin/bash", "-c", runCommand}, volumeMount, int64(instance.MemoryLimit))
	if err != nil {
		return nil, errors.New("run stage: failed to create container")
	}
	defer s.removeContainer(resp.ID)

	// start container
	err = s.startContainer(resp.ID)
	if err != nil {
		return nil, errors.New("run stage: failed to start container")
	}

	// grab container wait channel
	resultC, errC = s.DockerClient.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	if err != nil {
		return nil, errors.New("run stage: failed to wait container")
	}

	result := ""

	// wait for container to finish or timeout
	select {
	case <-time.After(time.Duration(instance.Timeout) * time.Millisecond):
		// timeout
		err = s.stopContainer(resp.ID)
		if err != nil {
			return nil, errors.New("run stage: failed to stop container")
		}
		result = "timeout"
	case <-resultC:
		// container finish
		result = "finish"
	case <-errC:
		// container error
		return nil, errors.New("run stage: failed to run code")
	}

	// get container exit code
	exitCode, err = s.getContainerExitCode(resp.ID)
	if err != nil {
		return nil, errors.New("run stage: failed to get container exit code")
	}

	// get container stdout
	stdout, err := s.getLog(resp.ID, true, false)
	if err != nil {
		return nil, errors.New("run stage: failed to get container stdout")
	}

	// get container stderr
	stderr, err := s.getLog(resp.ID, false, true)
	if err != nil {
		return nil, errors.New("run stage: failed to get container stderr")
	}

	instance.ExitCode = exitCode
	instance.Stdout = stdout
	instance.Stderr = stderr
	instance.Note = result

	// return instance
	return instance, nil
}

func NewSandboxService() SandboxService {
	// dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	dockerClient, err := client.NewClientWithOpts(client.WithHost("unix:///var/run/docker.sock"))
	if err != nil {
		panic(err)
	}

	return &sandboxService{
		DockerClient: dockerClient,
	}
}
