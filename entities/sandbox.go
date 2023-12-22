package entities

const (
	SandboxMemoryMB = 1024 * 1024
	SandboxMemoryGB = 1024 * SandboxMemoryMB
)

type SandboxInstance struct {
	Code            string
	Stdin           string
	Language        string
	Stdout          string
	Stderr          string
	Timeout         int
	MemoryLimit     int
	Error           error
	ExitCode        int
	CodeFilePath    string
	StdinFilePath   string
	CompileStdout   string
	CompileStderr   string
	CompileExitCode int
	Note            string
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
	CompileCmd:  "cp /sandbox/code /sandbox/code.py",
	RunCmd:      "python3 /sandbox/code.py < /sandbox/stdin",
}

var GoInstructionBook = SandboxInstruction{
	Language:    "go",
	DockerImage: "docker.io/library/golang:1.21",
	CompileCmd:  "cd /sandbox && cp /sandbox/code /sandbox/main.go && (go mod init sandbox > /dev/null 2>&1) && go build -o main",
	RunCmd:      "/sandbox/main < /sandbox/stdin",
}

var CInstructionBook = SandboxInstruction{
	Language:    "c",
	DockerImage: "docker.io/library/gcc:12.3.0",
	CompileCmd:  "cp /sandbox/code /sandbox/main.c && cd /sandbox/ && gcc -o /sandbox/main /sandbox/main.c > /dev/null",
	RunCmd:      "/sandbox/main < /sandbox/stdin",
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
