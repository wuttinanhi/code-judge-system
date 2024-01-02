package tests_test

import (
	"testing"

	"github.com/wuttinanhi/code-judge-system/databases"
	"github.com/wuttinanhi/code-judge-system/entities"
	"github.com/wuttinanhi/code-judge-system/services"
)

func TestSandbox(t *testing.T) {
	db := databases.NewTempSQLiteDatabase()
	testServiceKit := services.CreateTestServiceKit(db)

	t.Run("Sandbox Go Test", func(t *testing.T) {
		sandbox, err := testServiceKit.SandboxService.CreateSandbox(
			entities.GoInstructionBook.Language,
			entities.GoCodeExample,
		)
		if err != nil {
			t.Fatal(err)
		}
		defer testServiceKit.SandboxService.CleanUp(sandbox)

		compile := testServiceKit.SandboxService.CompileSandbox(sandbox)
		if compile.Err != nil {
			t.Fatal(compile.Err)
		}

		result := testServiceKit.SandboxService.Run(sandbox, "1\n2\n", entities.SandboxMemoryMB*128, 1000)
		if result.Err != nil {
			t.Fatal(err)
		}

		if result.Stdout != "3\n" {
			t.Error("stdout not match got\n", result.Stdout)
		}
		if result.Stderr != "" {
			t.Error("stderr not match")
		}
	})

	t.Run("Sandbox Python Test", func(t *testing.T) {
		sandbox, err := testServiceKit.SandboxService.CreateSandbox(
			entities.PythonInstructionBook.Language,
			entities.PythonCodeExample,
		)
		if err != nil {
			t.Fatal(err)
		}
		defer testServiceKit.SandboxService.CleanUp(sandbox)

		compile := testServiceKit.SandboxService.CompileSandbox(sandbox)
		if compile.Err != nil {
			t.Fatal(compile.Err)
		}

		result := testServiceKit.SandboxService.Run(sandbox, "1\n2\n", entities.SandboxMemoryMB*128, 1000)
		if err != nil {
			t.Fatal(err)
		}

		if result.Stdout != "3\n" {
			t.Error("stdout not match got\n", result.Stdout)
		}
		if result.Stderr != "" {
			t.Error("stderr not match")
		}
	})

	t.Run("Sandbox C Test", func(t *testing.T) {
		sandbox, err := testServiceKit.SandboxService.CreateSandbox(
			entities.CInstructionBook.Language,
			entities.CCodeExample,
		)
		if err != nil {
			t.Fatal(err)
		}
		defer testServiceKit.SandboxService.CleanUp(sandbox)

		compile := testServiceKit.SandboxService.CompileSandbox(sandbox)
		if compile.Err != nil {
			t.Fatal(compile.Err)
		}

		result := testServiceKit.SandboxService.Run(sandbox, "1\n2\n", entities.SandboxMemoryMB*128, 1000)
		if err != nil {
			t.Fatal(err)
		}

		if result.Stdout != "3\n" {
			t.Error("stdout not match got\n", result.Stdout)
		}
		if result.Stderr != "" {
			t.Error("stderr not match")
		}
	})

	t.Run("Sandbox OOM Python Test", func(t *testing.T) {
		sandbox, err := testServiceKit.SandboxService.CreateSandbox(
			entities.PythonInstructionBook.Language,
			entities.PythonCodeOOMTestCode,
		)
		if err != nil {
			t.Fatal(err)
		}
		defer testServiceKit.SandboxService.CleanUp(sandbox)

		compile := testServiceKit.SandboxService.CompileSandbox(sandbox)
		if compile.Err != nil {
			t.Fatal(compile.Err)
		}

		result := testServiceKit.SandboxService.Run(sandbox, "1\n2\n", entities.SandboxMemoryMB*128, 1000)
		if err != nil {
			t.Fatal(err)
		}

		// exit code must be OOM
		if result.ExitCode != 137 {
			t.Error("OOM exit code not match, got", result.ExitCode)
		}
	})

	t.Run("Sandbox Timeout Python Test", func(t *testing.T) {
		sandbox, err := testServiceKit.SandboxService.CreateSandbox(
			entities.PythonInstructionBook.Language,
			entities.PythonCodeTimeoutTestCode,
		)
		if err != nil {
			t.Fatal(err)
		}
		defer testServiceKit.SandboxService.CleanUp(sandbox)

		compile := testServiceKit.SandboxService.CompileSandbox(sandbox)
		if compile.Err != nil {
			t.Fatal(compile.Err)
		}

		result := testServiceKit.SandboxService.Run(sandbox, "1\n2\n", entities.SandboxMemoryMB*128, 1000)
		if err != nil {
			t.Fatal(err)
		}

		// exit code must be timeout
		if result.ExitCode != 137 {
			t.Error("timeout exit code not match expected 137 got", result.ExitCode)
		}
		// must be timeout
		if result.Timeout != true {
			t.Error("timeout not match expected true got", result.Timeout)
		}
	})
}
