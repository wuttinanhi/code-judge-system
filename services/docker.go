package services

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

type DockerService interface {
	PullImage(imageName string) error
	ImageExist(imageName string) (bool, error)
	GetLog(containerID string, showStdout, showStderr bool) (string, error)
	GetContainerExitCode(containerID string) (int, error)
	CreateVolume(name string) (volume.Volume, error)
	DeleteVolume(v volume.Volume) error
	CopyToContainer(containerID, targetPath string, content []byte) error
	CreateContainer(imageName string, command []string, volumes []mount.Mount, memoryLimit int64, containerName string) (response container.CreateResponse, err error)
	StartContainer(containerID string) error
	StopContainer(containerID string) error
	RemoveContainer(containerID string) error
	WaitContainer(containerID string, timeout uint) string
}

type dockerService struct {
	ctx          context.Context
	DockerClient *client.Client
}

const (
	WaitResultSuccess = "success"
	WaitResultTimeout = "timeout"
	WaitResultError   = "error"
)

// WaitContainer implements DockerService.
func (s dockerService) WaitContainer(containerID string, timeout uint) string {
	resultC, errC := s.DockerClient.ContainerWait(s.ctx, containerID, container.WaitConditionNotRunning)

	select {
	case <-time.After(time.Duration(time.Duration(timeout) * time.Millisecond)):
		return WaitResultTimeout
	case <-resultC:
		return WaitResultSuccess
	case <-errC:
		return WaitResultError
	}
}

func (s dockerService) PullImage(imageName string) error {
	log.Println("pulling image", imageName)
	out, err := s.DockerClient.ImagePull(s.ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer out.Close()

	var logs bytes.Buffer
	io.Copy(&logs, out)

	log.Println("pulling image", imageName, "done")

	return nil
}

func (s dockerService) ImageExist(imageName string) (bool, error) {
	_, _, err := s.DockerClient.ImageInspectWithRaw(s.ctx, imageName)
	if err != nil {
		if client.IsErrNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s dockerService) GetLog(containerID string, showStdout, showStderr bool) (string, error) {
	logs, err := s.DockerClient.ContainerLogs(s.ctx, containerID, types.ContainerLogsOptions{
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

func (s dockerService) GetContainerExitCode(containerID string) (int, error) {
	resp, err := s.DockerClient.ContainerInspect(s.ctx, containerID)
	if err != nil {
		return 0, err
	}
	return resp.State.ExitCode, nil
}

func (s dockerService) CreateVolume(name string) (volume.Volume, error) {
	volume, err := s.DockerClient.VolumeCreate(s.ctx, volume.CreateOptions{
		Name:   name,
		Driver: "local",
	})
	return volume, err
}

func (s dockerService) DeleteVolume(v volume.Volume) error {
	err := s.DockerClient.VolumeRemove(s.ctx, v.Name, true)
	return err
}

func (s dockerService) CopyToContainer(containerID, targetPath string, data []byte) error {
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

	err := s.DockerClient.CopyToContainer(s.ctx, containerID, targetDir, &buf, types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: false,
		CopyUIDGID:                false,
	})

	return err
}

func (s dockerService) CreateContainer(imageName string, command []string, volumes []mount.Mount, memoryLimit int64, containerName string) (response container.CreateResponse, err error) {
	response, err = s.DockerClient.ContainerCreate(s.ctx, &container.Config{
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
		containerName,
	)
	return
}

func (s dockerService) StartContainer(containerID string) error {
	err := s.DockerClient.ContainerStart(s.ctx, containerID, types.ContainerStartOptions{})
	return err
}

func (s dockerService) StopContainer(containerID string) error {
	timeout := int(0)
	err := s.DockerClient.ContainerStop(s.ctx, containerID, container.StopOptions{
		Timeout: &timeout,
		Signal:  "SIGKILL",
	})
	return err
}

func (s dockerService) RemoveContainer(containerID string) error {
	err := s.DockerClient.ContainerRemove(s.ctx, containerID, types.ContainerRemoveOptions{
		RemoveVolumes: false,
		Force:         true,
	})
	return err
}

func NewDockerservice() DockerService {
	ctx := context.Background()
	dockerClient, err := client.NewClientWithOpts(client.WithHost("unix:///var/run/docker.sock"))
	if err != nil {
		panic(err)
	}

	return dockerService{
		ctx:          ctx,
		DockerClient: dockerClient,
	}
}
