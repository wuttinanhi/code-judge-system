package tests_test

import (
	"fmt"
	"testing"

	"github.com/wuttinanhi/code-judge-system/entities"
	"github.com/wuttinanhi/code-judge-system/services"
)

func TestSandbox(t *testing.T) {
	testServiceKit := services.CreateTestServiceKit()

	instance, err := testServiceKit.SandboxService.Run(&entities.SandboxInstance{
		// Code:     "print(input(), input())",
		Code:     "import time\ntime.sleep(0)\nprint('hello world')",
		Stdin:    "1\n2\n",
		Timeout:  1000,
		Language: "python",
		RamLimit: 1000000000,
	})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(instance.Note)
	fmt.Println(instance.Stdout)
	fmt.Println(instance.Stderr)
}
