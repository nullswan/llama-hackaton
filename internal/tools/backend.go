package tools

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/nullswan/llama-hackaton/internal/chat"
	"github.com/nullswan/llama-hackaton/internal/completion"
	"github.com/nullswan/llama-hackaton/internal/llama"
)

type TextToJSONBackend struct {
	backend *llama.TextToJSONProvider
	logger  *slog.Logger
}

func NewTextToJSONBackend(
	backend *llama.TextToJSONProvider,
	logger *slog.Logger,
) TextToJSONBackend {
	return TextToJSONBackend{
		backend: backend,
		logger:  logger,
	}
}

func (t TextToJSONBackend) Do(
	ctx context.Context,
	conversation *chat.Conversation,
) (string, error) {
	messages := conversation.GetMessages()

	outCh := make(chan completion.Completion)
	go func() {
		if err := t.backend.GenerateCompletion(ctx, messages, outCh); err != nil {
			if strings.Contains(err.Error(), "context canceled") {
				return
			}
			t.logger.With("error", err).
				Error("Error generating completion")
		}
	}()

	defer close(outCh)

	for cmpl := range outCh {
		if !completion.IsTombStone(cmpl) {
			continue
		}

		content := strings.ReplaceAll(cmpl.Content(), "```json", "")
		return strings.ReplaceAll(content, "```", ""), nil
	}

	return "", errors.New("completion channel closed")
}
