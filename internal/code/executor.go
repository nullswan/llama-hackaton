package code

import (
	"runtime"
	"sync"
)

var (
	executors                = make(map[string]Executor)
	onceExecutorRegistration sync.Once
)

func registerExecutor(language string, executor Executor) {
	executors[language] = executor
}

func ExecuteCodeBlock(block Block) ExecutionResult {
	onceExecutorRegistration.Do(
		func() {
			initBashExecutor()
			initPythonExecutor()
			initOsascriptExecutor()
		},
	)

	executor, ok := executors[block.Language]
	if !ok {
		return ExecutionResult{
			Stderr:   "Unsupported language: " + block.Language,
			ExitCode: 1,
		}
	}

	if block.Language == "osascript" && runtime.GOOS != "darwin" {
		return ExecutionResult{
			Stderr:   "Osascript is only supported on macOS",
			ExitCode: 1,
		}
	}

	if block.Language == "powershell" && runtime.GOOS != "windows" {
		return ExecutionResult{
			Stderr:   "Powershell is only supported on Windows",
			ExitCode: 1,
		}
	}

	r := executor.Execute(block.Code)
	r.Block = block

	return r
}
