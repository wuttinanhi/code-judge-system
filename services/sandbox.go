package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types/mount"
	"github.com/wuttinanhi/code-judge-system/entities"
)

type SandboxService interface {
	Run(instance *entities.SandboxInstance) (*entities.SandboxInstance, error)
}

type sandboxService struct {
	dockerService DockerService
}

// Run implements SandboxService.
func (s *sandboxService) Run(instance *entities.SandboxInstance) (*entities.SandboxInstance, error) {
	instruction := entities.GetSandboxInstructionByLanguage(instance.Language)
	if instruction == nil {
		return nil, fmt.Errorf("language %s not supported", instance.Language)
	}

	imageName := instruction.DockerImage
	if imageName == "" {
		return nil, fmt.Errorf("language %s not supported", instance.Language)
	}

	// check if image exist
	exist, err := s.dockerService.ImageExist(imageName)
	if err != nil {
		return nil, err
	}

	// if image not exist, pull image
	if !exist {
		// pull image
		err := s.dockerService.PullImage(imageName)
		if err != nil {
			return nil, err
		}
	}

	// create volume mount
	runID := strconv.Itoa(int(time.Now().UnixNano()))
	volumeID := "code-judge-system-" + runID
	volumeMount := []mount.Mount{
		{Type: mount.TypeVolume, Source: volumeID, Target: "/sandbox"},
	}

	// create volume
	volume, err := s.dockerService.CreateVolume(volumeID)
	if err != nil {
		return nil, errors.New("failed to create volume")
	}
	defer s.dockerService.DeleteVolume(volume)

	// create container to create necessary files
	resp, err := s.dockerService.CreateContainer(
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
	defer s.dockerService.RemoveContainer(resp.ID)

	// start container
	err = s.dockerService.StartContainer(resp.ID)
	if err != nil {
		return nil, errors.New("create stage: failed to start container")
	}

	time.Sleep(1 * time.Second)

	// copy code file to container
	err = s.dockerService.CopyToContainer(resp.ID, "/sandbox/code", strings.NewReader(instance.Code))
	if err != nil {
		return nil, err
	}
	// copy stdin file to container
	err = s.dockerService.CopyToContainer(resp.ID, "/sandbox/stdin", strings.NewReader(instance.Stdin))
	if err != nil {
		return nil, err
	}

	// stop container
	err = s.dockerService.StopContainer(resp.ID)
	if err != nil {
		return nil, errors.New("create stage: failed to stop container")
	}

	// remove container
	err = s.dockerService.StopContainer(resp.ID)
	if err != nil {
		return nil, errors.New("create stage: failed to remove container")
	}

	// compile stage
	compileCommand := instruction.CompileCmd

	// create container to compile
	resp, err = s.dockerService.CreateContainer(imageName, []string{"/bin/bash", "-c", compileCommand}, volumeMount, entities.SandboxMemoryMB*256)
	if err != nil {
		return nil, errors.New("compile stage: failed to create container")
	}
	defer s.dockerService.RemoveContainer(resp.ID)

	// start container
	err = s.dockerService.StartContainer(resp.ID)
	if err != nil {
		return nil, errors.New("compile stage: failed to start container")
	}

	// wait container to finish
	waitResult := s.dockerService.WaitContainer(resp.ID, instance.Timeout)
	if waitResult == WaitResultError {
		return nil, errors.New("compile stage: failed to compile code")
	}

	// get container exit code
	exitCode, err := s.dockerService.GetContainerExitCode(resp.ID)
	if err != nil {
		return nil, errors.New("compile stage: failed to get container exit code")
	}

	compileStdOut, err := s.dockerService.GetLog(resp.ID, true, false)
	if err != nil {
		return nil, errors.New("compile stage: failed to get container log")
	}

	compileStdErr, err := s.dockerService.GetLog(resp.ID, false, true)
	if err != nil {
		return nil, errors.New("compile stage: failed to get container log")
	}

	instance.CompileExitCode = exitCode
	instance.CompileStdout = compileStdOut
	instance.CompileStderr = compileStdErr

	// if exit code is not 0, return error
	if exitCode != 0 {
		return nil, errors.New("compile stage: failed to compile code")
	}

	// stop and remove container
	err = s.dockerService.RemoveContainer(resp.ID)
	if err != nil {
		return nil, errors.New("compile stage: failed to stop and remove container")
	}

	// run stage
	runCommand := instruction.RunCmd

	// create container to run
	resp, err = s.dockerService.CreateContainer(imageName, []string{"/bin/bash", "-c", runCommand}, volumeMount, int64(instance.MemoryLimit))
	if err != nil {
		return nil, errors.New("run stage: failed to create container")
	}
	defer s.dockerService.RemoveContainer(resp.ID)

	// start container
	err = s.dockerService.StartContainer(resp.ID)
	if err != nil {
		return nil, errors.New("run stage: failed to start container")
	}

	// waitResult for container to finish
	waitResult = s.dockerService.WaitContainer(resp.ID, instance.Timeout)
	if waitResult == WaitResultError {
		return nil, errors.New("run stage: failed to wait container")
	}
	if waitResult == WaitResultTimeout {
		err = s.dockerService.StopContainer(resp.ID)
		if err != nil {
			return nil, errors.New("run stage: failed to stop container")
		}
	}

	// get container exit code
	exitCode, err = s.dockerService.GetContainerExitCode(resp.ID)
	if err != nil {
		return nil, errors.New("run stage: failed to get container exit code")
	}

	// get container stdout
	stdout, err := s.dockerService.GetLog(resp.ID, true, false)
	if err != nil {
		return nil, errors.New("run stage: failed to get container stderr")
	}

	// get container stderr
	stderr, err := s.dockerService.GetLog(resp.ID, false, true)
	if err != nil {
		return nil, errors.New("run stage: failed to get container stderr")
	}

	instance.ExitCode = exitCode
	instance.Stdout = stdout
	instance.Stderr = stderr
	instance.Note = waitResult

	// return instance
	return instance, nil
}

func NewSandboxService() SandboxService {
	dockerService := NewDockerservice()

	return &sandboxService{
		dockerService: dockerService,
	}
}
