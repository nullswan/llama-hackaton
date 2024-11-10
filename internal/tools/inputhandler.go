package tools

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/nullswan/llama-hackaton/internal/term"
)

type InputHandler interface {
	Read(ctx context.Context, defaultValue string) (string, error)
}

type inputHandler struct {
	logger *slog.Logger
}

func NewInputHandler(
	logger *slog.Logger,
) InputHandler {
	return &inputHandler{
		logger: logger,
	}
}

func (i *inputHandler) Read(
	ctx context.Context,
	defaultValue string,
) (string, error) {
	rl, err := term.InitReadline(defaultValue)
	if err != nil {
		return "", fmt.Errorf("error initializing readline: %w", err)
	}

	inputErrCh := make(chan error)
	inputCh := make(chan string)

	go func() {
		ret, err := term.ReadInputOnce(rl)
		if err != nil {
			if rl.Closed() {
				return
			}

			select {
			case inputErrCh <- err:
			case <-ctx.Done():
			}

			return
		}

		select {
		case inputCh <- ret:
		case <-ctx.Done():
			return
		}
	}()

	defer func() {
		rl.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			return "", errors.New("context canceled")
		case line := <-inputCh:
			return line, nil
		case err := <-inputErrCh:
			return "", fmt.Errorf("error reading input: %w", err)
		}
	}
}
