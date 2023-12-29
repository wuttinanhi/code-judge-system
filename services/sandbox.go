package services

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/docker/docker/api/types/mount"
	"github.com/wuttinanhi/code-judge-system/entities"
)

type SandboxService interface {
	CreateSandbox(lang, code string) (*entities.SandboxInstance, error)
	Run(instance *entities.SandboxInstance, stdin string, memoryLimit, timeLimit uint) (result *entities.SandboxRunResult)
	CleanUp(instance *entities.SandboxInstance) error
}

type sandboxService struct {
	dockerService DockerService
}

func (s *sandboxService) CopyFileToVolume(instance *entities.SandboxInstance, volumeMount []mount.Mount, fileContentMap map[string]string) error {
	// create container to store necessary files
	containerName := fmt.Sprintf("%s-copy-%d", instance.RunID, time.Now().UnixNano())
	resp, err := s.dockerService.CreateContainer(
		instance.ImageName,
		[]string{"/bin/sh", "-c", "chmod 777 -R /sandbox && sleep 9999"},
		volumeMount,
		entities.SandboxMemoryMB*256,
		containerName,
	)
	if err != nil {
		return err
	}
	defer s.dockerService.RemoveContainer(resp.ID)

	// start container
	err = s.dockerService.StartContainer(resp.ID)
	if err != nil {
		return errors.New("copy: failed to start container")
	}

	<-time.After(1 * time.Second)

	// copy file to container
	for path, content := range fileContentMap {
		err = s.dockerService.CopyToContainer(resp.ID, path, []byte(content))
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

// CreateSandbox implements SandboxService.
func (s *sandboxService) CreateSandbox(lang string, code string) (*entities.SandboxInstance, error) {
	instance := &entities.SandboxInstance{
		RunID:    strconv.Itoa(int(time.Now().UnixNano())),
		Language: lang,
	}

	instance.Instruction = entities.GetSandboxInstructionByLanguage(instance.Language)
	if instance.Instruction == nil {
		return nil, fmt.Errorf("language %s not supported", instance.Language)
	}

	instance.ImageName = instance.Instruction.DockerImage
	if instance.ImageName == "" {
		return nil, fmt.Errorf("language %s not supported", instance.Language)
	}

	// check if image exist
	exist, err := s.dockerService.ImageExist(instance.ImageName)
	if err != nil {
		return nil, err
	}

	// if image not exist, pull image
	if !exist {
		// pull image
		err := s.dockerService.PullImage(instance.ImageName)
		if err != nil {
			return nil, err
		}
	}

	// create volume
	volumeName := fmt.Sprintf("code-judge-system-%s-program", instance.RunID)
	instance.ProgramVolume, err = s.dockerService.CreateVolume(volumeName)
	if err != nil {
		return nil, errors.New("create stage: failed to create program volume")
	}

	// create volume mount
	programVolumeMount := []mount.Mount{
		{Type: mount.TypeVolume, Source: volumeName, Target: "/sandbox"},
	}

	err = s.CopyFileToVolume(instance, programVolumeMount, map[string]string{
		"/sandbox/code": code,
	})
	if err != nil {
		return nil, errors.New("create stage: failed to copy file to container")
	}

	// compile info
	compileCommand := instance.Instruction.CompileCmd
	compileTimeout := instance.Instruction.CompileTimeout

	// create container to compile
	resp, err := s.dockerService.CreateContainer(
		instance.ImageName,
		[]string{"/bin/sh", "-c", compileCommand},
		programVolumeMount,
		entities.SandboxMemoryMB*256,
		instance.RunID+"-compile",
	)
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
	waitResult := s.dockerService.WaitContainer(resp.ID, compileTimeout)
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
	if instance.CompileExitCode != 0 {
		return nil, errors.New("compile stage: failed to compile code")
	}

	return instance, nil
}

// Run implements SandboxService.
func (s *sandboxService) Run(instance *entities.SandboxInstance, stdin string, memoryLimit, timeLimit uint) (result *entities.SandboxRunResult) {
	result = &entities.SandboxRunResult{}
	runCommand := instance.Instruction.RunCmd

	// create stdin volume
	stdinVolumeName := fmt.Sprintf("code-judge-system-%s-%d-stdin", instance.RunID, time.Now().UnixNano())
	stdinVolume, err := s.dockerService.CreateVolume(stdinVolumeName)
	if err != nil {
		result.Err = errors.New("create stage: failed to create stdin volume")
		return
	}
	defer s.dockerService.DeleteVolume(stdinVolume)

	// create stdin volume mount
	stdinVolumeMount := []mount.Mount{
		{Type: mount.TypeVolume, Source: stdinVolumeName, Target: "/stdin"},
	}

	// copy stdin to container
	err = s.CopyFileToVolume(instance, stdinVolumeMount, map[string]string{
		"/stdin/stdin": stdin,
	})
	if err != nil {
		fmt.Println(err)
		result.Err = errors.New("run stage: failed to copy stdin to container")
		return
	}

	// create run volume mount
	runVolumeMount := []mount.Mount{
		{Type: mount.TypeVolume, Source: instance.ProgramVolume.Name, Target: "/sandbox"},
		{Type: mount.TypeVolume, Source: stdinVolumeName, Target: "/stdin", ReadOnly: true},
	}

	// create container to run
	containerName := fmt.Sprintf("%s-run-%d", instance.RunID, time.Now().UnixNano())
	resp, err := s.dockerService.CreateContainer(
		instance.ImageName,
		[]string{"/bin/sh", "-c", runCommand},
		runVolumeMount,
		int64(memoryLimit),
		containerName,
	)
	if err != nil {
		result.Err = errors.New("run stage: failed to create container")
		return
	}
	defer s.dockerService.RemoveContainer(resp.ID)

	// start container
	err = s.dockerService.StartContainer(resp.ID)
	if err != nil {
		result.Err = errors.New("run stage: failed to start container")
		return
	}

	// wait for container to finish
	waitResult := s.dockerService.WaitContainer(resp.ID, timeLimit)
	if waitResult == WaitResultError {
		result.Err = errors.New("run stage: failed to wait container")
		return
	}
	if waitResult == WaitResultTimeout {
		err = s.dockerService.StopContainer(resp.ID)
		if err != nil {
			result.Err = errors.New("run stage: failed to stop container")
			return
		}
	}

	// get container exit code
	exitCode, err := s.dockerService.GetContainerExitCode(resp.ID)
	if err != nil {
		result.Err = errors.New("run stage: failed to get container exit code")
		return
	}

	// get container stdout
	stdout, err := s.dockerService.GetLog(resp.ID, true, false)
	if err != nil {
		result.Err = errors.New("run stage: failed to get container stderr")
		return
	}

	// get container stderr
	stderr, err := s.dockerService.GetLog(resp.ID, false, true)
	if err != nil {
		result.Err = errors.New("run stage: failed to get container stderr")
		return
	}

	result.ExitCode = exitCode
	result.Stdout = stdout
	result.Stderr = stderr
	result.Timeout = waitResult == WaitResultTimeout

	// return instance
	return
}

// CleanUp implements SandboxService.
func (s *sandboxService) CleanUp(instance *entities.SandboxInstance) error {
	// remove volume
	err := s.dockerService.DeleteVolume(instance.ProgramVolume)
	if err != nil {
		return err
	}

	return nil
}

func NewSandboxService() SandboxService {
	dockerService := NewDockerservice()

	return &sandboxService{
		dockerService: dockerService,
	}
}
