package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
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
	out, err := s.DockerClient.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer out.Close()

	var logs bytes.Buffer
	io.Copy(&logs, out)

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

func (s *sandboxService) deleteTempCodeFile(filePath string) error {
	return os.Remove(filePath)
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

	// create temp code file
	codeFilePath, err := s.createTempCodeFile(instance.Code)
	if err != nil {
		return nil, err
	}
	defer s.deleteTempCodeFile(codeFilePath)

	// create mount host code file to container at /tmp/code.py (read only)
	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{Type: mount.TypeBind, ReadOnly: true, Source: codeFilePath, Target: "/tmp/code"},
		},
	}

	// compile command
	compileCommand := instruction.CompileCmd
	runCommand := instruction.RunCmd

	// merge two command together
	mergedCommand := fmt.Sprintf("%s && %s", compileCommand, runCommand)

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
		Entrypoint:      []string{"/bin/sh", "-c", mergedCommand},
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

	instance.Stdout, err = s.getLog(resp.ID, true, false)
	if err != nil {
		return nil, err
	}
	instance.Stderr, err = s.getLog(resp.ID, false, true)
	if err != nil {
		return nil, err
	}

	// remove container
	err = s.DockerClient.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
	if err != nil {
		return nil, err
	}

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
