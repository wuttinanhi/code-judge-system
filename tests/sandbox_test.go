package tests_test

import (
	"testing"

	"github.com/wuttinanhi/code-judge-system/entities"
	"github.com/wuttinanhi/code-judge-system/services"
)

func TestSandbox(t *testing.T) {
	testServiceKit := services.CreateTestServiceKit()

	t.Run("go sandbox", func(t *testing.T) {
		instance, err := testServiceKit.SandboxService.Run(&entities.SandboxInstance{
			Language:    entities.GoInstructionBook.Language,
			Code:        entities.GoCodeExample,
			Stdin:       "1\n2\n",
			Timeout:     10000,
			MemoryLimit: entities.SandboxMemoryMB * 128,
		})
		if err != nil {
			t.Fatal(err)
		}

		if instance.Stdout != "3\n" {
			t.Fatal("stdout not match")
		}
		if instance.Stderr != "" {
			t.Fatal("stderr not match")
		}
	})

	t.Run("python sandbox", func(t *testing.T) {
		instance, err := testServiceKit.SandboxService.Run(&entities.SandboxInstance{
			Language:    entities.PythonInstructionBook.Language,
			Code:        entities.PythonCodeExample,
			Stdin:       "1\n2\n",
			Timeout:     1000,
			MemoryLimit: entities.SandboxMemoryMB * 128,
		})
		if err != nil {
			t.Fatal(err)
		}

		if instance.Stdout != "3\n" {
			t.Fatal("stdout not match")
		}
		if instance.Stderr != "" {
			t.Fatal("stderr not match")
		}
	})

	t.Run("c sandbox", func(t *testing.T) {
		instance, err := testServiceKit.SandboxService.Run(&entities.SandboxInstance{
			Language:    entities.CInstructionBook.Language,
			Code:        entities.CCodeExample,
			Stdin:       "1\n2\n",
			Timeout:     10000,
			MemoryLimit: entities.SandboxMemoryMB * 128,
		})
		if err != nil {
			t.Fatal(err)
		}

		if instance.Stdout != "3\n" {
			t.Fatal("stdout not match")
		}
		if instance.Stderr != "" {
			t.Fatal("stderr not match")
		}
	})

	t.Run("python OOM sandbox", func(t *testing.T) {
		instance, err := testServiceKit.SandboxService.Run(&entities.SandboxInstance{
			Language:    entities.PythonInstructionBook.Language,
			Code:        entities.PythonCodeOOMTestCode,
			Stdin:       "1\n2\n",
			Timeout:     1000,
			MemoryLimit: entities.SandboxMemoryMB * 128,
		})
		if err != nil {
			t.Fatal(err)
		}

		// exit code must be OOM
		if instance.ExitCode != 137 {
			t.Fatal("OOM exit code not match")
		}
	})
}
