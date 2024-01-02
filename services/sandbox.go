package services

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/docker/docker/api/types/mount"
	"github.com/wuttinanhi/code-judge-system/entities"
)

type SandboxService interface {
	CreateSandbox(lang, code string) (*entities.SandboxInstance, error)
	CompileSandbox(instance *entities.SandboxInstance) (result *entities.SandboxRunResult)
	Run(instance *entities.SandboxInstance, stdin string, memoryLimit, timeLimit uint) (result *entities.SandboxRunResult)
	CleanUp(instance *entities.SandboxInstance) error
	ValidateMemoryLimit(memoryLimit uint) (err error)
	ValidateTimeLimit(timeLimit uint) (err error)
}

type sandboxService struct {
	dockerService DockerService
	memoryLimit   uint
	timeLimit     uint
}

func generateID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func (s *sandboxService) CopyFileToVolume(instance *entities.SandboxInstance, volumeMount []mount.Mount, fileContentMap map[string]string) error {
	// create container to store necessary files
	containerName := fmt.Sprintf("%s-copy-%s", instance.RunID, generateID())
	resp, err := s.dockerService.CreateContainer(
		instance.ImageName,
		// chmod 777 -R /sandbox &&
		[]string{"/bin/sh", "-c", "sleep 9999"},
		volumeMount,
		int64(entities.SandboxMemoryMB*512),
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

	// wait for container to ready
	// time.Sleep(1 * time.Second)

	// copy file to container
	for path, content := range fileContentMap {
		err = s.dockerService.CopyToContainer(resp.ID, path, []byte(content))
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

// CreateSandbox implements SandboxService.
func (s *sandboxService) CreateSandbox(lang string, code string) (*entities.SandboxInstance, error) {
	instance := &entities.SandboxInstance{
		RunID:    generateID(),
		Language: lang,
		Code:     code,
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

	return instance, nil
}

func (s *sandboxService) CompileSandbox(instance *entities.SandboxInstance) (result *entities.SandboxRunResult) {
	log.Println("start compiling sandbox", instance.RunID)

	result = &entities.SandboxRunResult{}

	volumeName := fmt.Sprintf("code-judge-system-%s-program", instance.RunID)
	volume, err := s.dockerService.CreateVolume(volumeName)
	if err != nil {
		result.Err = errors.New("compile stage: failed to create program volume")
		return
	}
	instance.ProgramVolume = volume

	programVolumeMount := []mount.Mount{
		{Type: mount.TypeVolume, Source: volumeName, Target: "/sandbox"},
	}

	err = s.CopyFileToVolume(instance, programVolumeMount, map[string]string{
		"/sandbox/code": instance.Code,
	})
	if err != nil {
		result.Err = errors.New("compile stage: failed to copy code to container")
		return
	}

	compileCommand := instance.Instruction.CompileCmd
	compileTimeout := instance.Instruction.CompileTimeout

	resp, err := s.dockerService.CreateContainer(
		instance.ImageName,
		[]string{"/bin/sh", "-c", compileCommand},
		programVolumeMount,
		int64(entities.SandboxMemoryGB*1),
		instance.RunID+"-compile",
	)
	if err != nil {
		result.Err = errors.New("compile stage: failed to create container")
		return
	}
	defer s.dockerService.RemoveContainer(resp.ID)

	err = s.dockerService.StartContainer(resp.ID)
	if err != nil {
		result.Err = errors.New("compile stage: failed to start container")
		return
	}

	waitResult := s.dockerService.WaitContainer(resp.ID, compileTimeout)
	if waitResult == WaitResultError {
		result.Err = errors.New("compile stage: failed to compile code")
		return
	}

	exitCode, err := s.dockerService.GetContainerExitCode(resp.ID)
	if err != nil {
		result.Err = errors.New("compile stage: failed to get container exit code")
		return
	}

	compileStdOut, err := s.dockerService.GetLog(resp.ID, true, false)
	if err != nil {
		result.Err = errors.New("compile stage: failed to get stdout")
		return
	}

	compileStdErr, err := s.dockerService.GetLog(resp.ID, false, true)
	if err != nil {
		result.Err = errors.New("compile stage: failed to get stderr")
		return
	}

	instance.CompileExitCode = exitCode
	instance.CompileStdout = compileStdOut
	instance.CompileStderr = compileStdErr

	if instance.CompileExitCode != 0 {
		result.Err = errors.New("compile stage: failed to compile code")
		return
	}

	log.Println("compiling done", instance.RunID)
	return
}

// Run implements SandboxService.
func (s *sandboxService) Run(instance *entities.SandboxInstance, stdin string, memoryLimit, timeLimit uint) (result *entities.SandboxRunResult) {
	result = &entities.SandboxRunResult{}

	maxMemoryErr := s.ValidateMemoryLimit(memoryLimit)
	if maxMemoryErr != nil {
		result.Err = errors.New("run stage: max memory exceeded sandbox limit")
		return
	}

	maxTimeLimitErr := s.ValidateTimeLimit(timeLimit)
	if maxTimeLimitErr != nil {
		result.Err = errors.New("run stage: max run time exceeded sandbox limit")
		return
	}

	runCommand := instance.Instruction.RunCmd

	// create stdin volume
	stdinVolumeName := fmt.Sprintf("code-judge-system-%s-%s-stdin", instance.RunID, generateID())
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
	containerName := fmt.Sprintf("%s-run-%s", instance.RunID, generateID())
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

func (s *sandboxService) ValidateMemoryLimit(memoryLimit uint) (err error) {
	if memoryLimit > entities.SandboxMemoryMB*s.memoryLimit {
		err = errors.New("run stage: too large memory limit")
	}
	return
}

func (s *sandboxService) ValidateTimeLimit(timeLimit uint) (err error) {

	if timeLimit > s.timeLimit {
		err = errors.New("run stage: too large time limit")
	}
	return
}

func NewSandboxService(memoryLimit uint, timeLimit uint) SandboxService {
	dockerService := NewDockerservice()

	return &sandboxService{
		dockerService: dockerService,
		memoryLimit:   memoryLimit,
		timeLimit:     timeLimit,
	}
}
