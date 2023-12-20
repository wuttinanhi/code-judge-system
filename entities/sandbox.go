package entities

const (
	SandboxRamMB = 1000
	SandboxRamGB = 1000 * SandboxRamMB
)

type SandboxInstance struct {
	Code     string
	Stdin    string
	Language string
	Stdout   string
	Stderr   string
	Timeout  int
	RamLimit int
	Error    error
	Note     string
}

type SandboxInstruction struct {
	Language    string
	DockerImage string
	CompileCmd  string
	RunCmd      string
}

var PythonInstructionBook = SandboxInstruction{
	Language:    "python",
	DockerImage: "docker.io/library/python:3.10",
	CompileCmd:  "cp /tmp/code /tmp/code.py",
	RunCmd:      "python3 /tmp/code.py",
}

var GoInstructionBook = SandboxInstruction{
	Language:    "go",
	DockerImage: "docker.io/library/golang:1.21",
	CompileCmd:  "cp /tmp/code /main.go && cd / && (go mod init sandbox > /dev/null 2>&1) && go build -o /main > /dev/null",
	RunCmd:      "/main",
}

var CInstructionBook = SandboxInstruction{
	Language:    "c",
	DockerImage: "docker.io/library/gcc:12.3.0",
	CompileCmd:  "cp /tmp/code /tmp/main.c && cd /tmp/ && gcc -o /tmp/main /tmp/main.c > /dev/null",
	RunCmd:      "/tmp/main",
}

var PythonCodeExample = `print("Hello World")`

var GoCodeExample = `package main

import "fmt"

func main() {
	fmt.Println("Hello World")
}`

var CCodeExample = `#include <stdio.h>

int main() {
	printf("Hello World");
	return 0;
}`

var LanguageInstructionMap = map[string]SandboxInstruction{
	"python": PythonInstructionBook,
	"go":     GoInstructionBook,
	"c":      CInstructionBook,
}

func GetSandboxInstructionByLanguage(language string) *SandboxInstruction {
	// check if language exist
	instruction, ok := LanguageInstructionMap[language]
	if !ok {
		return nil
	}
	return &instruction
}
