package tests_test

import (
	"testing"

	"github.com/wuttinanhi/code-judge-system/entities"
	"github.com/wuttinanhi/code-judge-system/services"
)

func TestSandbox(t *testing.T) {
	testServiceKit := services.CreateTestServiceKit()

	t.Run("Sandbox Go Test", func(t *testing.T) {
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
			t.Error("stdout not match")
		}
		if instance.Stderr != "" {
			t.Error("stderr not match")
		}
	})

	t.Run("Sandbox Python Test", func(t *testing.T) {
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
			t.Error("stdout not match")
		}
		if instance.Stderr != "" {
			t.Error("stderr not match")
		}
	})

	t.Run("Sandbox C Test", func(t *testing.T) {
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
			t.Error("stdout not match")
		}
		if instance.Stderr != "" {
			t.Error("stderr not match")
		}
	})

	t.Run("Sandbox OOM Python Test", func(t *testing.T) {
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
			t.Error("OOM exit code not match")
		}
	})

	t.Run("Sandbox Timeout Python Test", func(t *testing.T) {
		instance, err := testServiceKit.SandboxService.Run(&entities.SandboxInstance{
			Language:    entities.PythonInstructionBook.Language,
			Code:        entities.PythonCodeTimeoutTestCode,
			Stdin:       "1\n2\n",
			Timeout:     1000,
			MemoryLimit: entities.SandboxMemoryMB * 128,
		})
		if err != nil {
			t.Fatal(err)
		}

		// exit code must be timeout
		if instance.ExitCode != 137 {
			t.Error("timeout exit code not match expected 137 got", instance.ExitCode)
		}
		// note must be timeout
		if instance.Note != "timeout" {
			t.Error("note not match expected 'timeout' got", instance.Note)
		}
	})
}
