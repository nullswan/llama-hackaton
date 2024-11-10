package llama

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/cenkalti/backoff/v4"
)

const llamaPath = "ollama"

func LoadTextToJSONProvider(
	model string,
) (*TextToJSONProvider, error) {
	var cmd *exec.Cmd
	if !ollamaServerIsRunning() {
		var err error
		cmd, err = tryStartOllama()
		if err != nil {
			ollamaOutput := llamaPath
			const maxDownloadRetries = 3
			err = backoff.Retry(func() error {
				fmt.Printf(
					"Download ollama to %s\n",
					ollamaOutput,
				)
				return downloadOllama(
					context.TODO(),
					ollamaOutput,
				)
			}, backoff.WithMaxRetries(backoff.NewConstantBackOff(time.Second), maxDownloadRetries))
			if err != nil {
				return nil, fmt.Errorf("error installing ollama: %w", err)
			}
		}
	}
	url := getOllamaURL()

	ollamaConfig := NewOlamaProviderConfig(
		url,
		model,
	)
	p, err := NewTextToJSONProvider(
		ollamaConfig,
		cmd,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating ollama provider: %w", err)
	}

	return p, nil
}
