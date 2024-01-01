package entities

import (
	"github.com/docker/docker/api/types/volume"
)

const (
	SandboxMemoryMB uint = 1024 * 1024
	SandboxMemoryGB uint = 1024 * SandboxMemoryMB
)

type SandboxInstance struct {
	RunID           string
	Language        string
	ImageName       string
	ProgramVolume   volume.Volume
	Instruction     *SandboxInstruction
	CompileExitCode int
	CompileStdout   string
	CompileStderr   string
}

type SandboxRunResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Timeout  bool
	Err      error
}

type SandboxInstruction struct {
	Language       string
	DockerImage    string
	CompileCmd     string
	RunCmd         string
	CompileTimeout uint
}

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

var PythonInstructionBook = SandboxInstruction{
	Language:       "python",
	DockerImage:    "docker.io/library/python:3.10",
	CompileCmd:     "cp /sandbox/code /sandbox/code.py",
	RunCmd:         "python3 /sandbox/code.py < /stdin/stdin",
	CompileTimeout: 0,
}

var GoInstructionBook = SandboxInstruction{
	Language:       "go",
	DockerImage:    "docker.io/library/golang:1.21",
	CompileCmd:     "cd /sandbox && cp /sandbox/code /sandbox/main.go && go mod init sandbox && go build -o main",
	RunCmd:         "/sandbox/main < /stdin/stdin",
	CompileTimeout: 10000,
}

var CInstructionBook = SandboxInstruction{
	Language:       "c",
	DockerImage:    "docker.io/library/gcc:12.3.0",
	CompileCmd:     "cp /sandbox/code /sandbox/main.c && gcc -o /sandbox/main /sandbox/main.c",
	RunCmd:         "/sandbox/main < /stdin/stdin",
	CompileTimeout: 10000,
}

var PythonCodeExample = `
x = int(input())
y = int(input())
print(x + y)
`

var PythonCodeOOMTestCode = `
data = []

while True:
    data.append(' ' * 10**6)
`

var PythonCodeTimeoutTestCode = `
import time
time.sleep(2)
`

var GoCodeExample = `
package main

import (
    "fmt"
    "bufio"
    "os"
    "strconv"
    "strings"
)

func main() {
    reader := bufio.NewReader(os.Stdin)

    xStr, _ := reader.ReadString('\n')
    x, _ := strconv.Atoi(strings.TrimSpace(xStr))

    yStr, _ := reader.ReadString('\n')
    y, _ := strconv.Atoi(strings.TrimSpace(yStr))

    fmt.Println(x + y)
}`

var CCodeExample = `
#include <stdio.h>

int main() {
    int x, y;

    scanf("%d", &x);

    scanf("%d", &y);

    printf("%d\n", x + y);

    return 0;
}`
