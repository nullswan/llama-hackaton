package code

type Block struct {
	ID          string
	Language    string
	Code        string
	Description string
}

type ExecutionResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Block    Block
}

type Executor interface {
	Execute(code string) ExecutionResult
}
