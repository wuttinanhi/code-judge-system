package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/wuttinanhi/code-judge-system/entities"
)

type SandboxService interface {
	Run(instance *entities.SandboxInstance) (*entities.SandboxInstance, error)
}

type sandboxService struct {
	DockerClient *client.Client
}

func (s *sandboxService) pullImage(imageName string) error {
	ctx := context.Background()
	out, err := s.DockerClient.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()
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

func (s *sandboxService) createTempCodeFile(code string) (string, error) {
	os.MkdirAll("/tmp/code-judge-system", os.ModePerm)
	unixTime := time.Now().UnixNano()
	fileName := fmt.Sprintf("/tmp/code-judge-system/%d.py", unixTime)

	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.WriteString(code)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

// Run implements SandboxService.
func (s *sandboxService) Run(instance *entities.SandboxInstance) (*entities.SandboxInstance, error) {
	ctx := context.Background()

	imageName := "docker.io/library/python:3.10"

	// check if image exist
	exist, err := s.imageExist(imageName)
	if err != nil {
		return nil, err
	}

	// if image not exist, pull image
	if !exist {
		// pull image
		err := s.pullImage(imageName)
		if err != nil {
			return nil, err
		}
	}

	// create temp code file
	// codeFile, err := s.createTempCodeFile(instance.Code)
	// if err != nil {
	// 	return nil, err
	// }

	// create mount host code file to container at /tmp/code.py (read only)
	hostConfig := &container.HostConfig{
		// Mounts: []mount.Mount{
		// 	{
		// 		Type:     mount.TypeBind,
		// 		ReadOnly: true,
		// 		Source:   codeFile,
		// 		Target:   "/tmp/code.py",
		// 	},
		// },
	}

	// create container
	resp, err := s.DockerClient.ContainerCreate(ctx, &container.Config{
		Image:           imageName,
		NetworkDisabled: true,
		Tty:             true,
		AttachStdout:    true,
		AttachStderr:    true,
		AttachStdin:     true,
		OpenStdin:       true,
		Env:             []string{"PYTHONUNBUFFERED=1"},
		// Entrypoint:      []string{"python3", "/tmp/code.py"},
		Entrypoint: []string{"python3", "-c", instance.Code},
	},
		hostConfig,
		nil,
		nil,
		"",
	)
	if err != nil {
		return nil, err
	}

	// start container
	if err := s.DockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	// create channel to wait for container to finish or timeout
	waitChannel := make(chan string, 1)

	// wait for timeout
	go func() {
		time.Sleep(time.Millisecond * time.Duration(instance.Timeout))
		waitChannel <- "timeout"
	}()

	// wait for container to finish
	go func() {
		resultC, errC := s.DockerClient.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

		select {
		case err := <-errC:
			waitChannel <- err.Error()
		case <-resultC:
			waitChannel <- "success"
		}
	}()

	// wait for waitChannel
	exitReason := <-waitChannel
	instance.Note = exitReason

	// stop container
	timeout := int(0)
	err = s.DockerClient.ContainerStop(ctx, resp.ID, container.StopOptions{
		Timeout: &timeout,
	})
	if err != nil {
		return nil, err
	}

	// get stdout container log
	stdoutLog, err := s.DockerClient.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
	})
	if err != nil {
		return nil, err
	}
	defer stdoutLog.Close()

	// get stderr container log
	stdErrLog, err := s.DockerClient.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
		ShowStderr: true,
	})
	if err != nil {
		return nil, err
	}
	defer stdErrLog.Close()

	// create buffer to store logs
	var stdoutBuffer bytes.Buffer
	var stderrBuffer bytes.Buffer

	// copy logs to buffer
	io.Copy(&stdoutBuffer, stdoutLog)
	io.Copy(&stderrBuffer, stdErrLog)

	// assign logs to instance
	instance.Stdout = stdoutBuffer.String()
	instance.Stderr = stderrBuffer.String()

	// remove container
	// err = s.DockerClient.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{
	// 	RemoveVolumes: true,
	// 	Force:         true,
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// return instance
	return instance, nil
}

func NewSandboxService() SandboxService {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	return &sandboxService{
		DockerClient: dockerClient,
	}
}
