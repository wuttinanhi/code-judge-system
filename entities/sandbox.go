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
