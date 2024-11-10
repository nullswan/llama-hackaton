package llama

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	ollamaServerTimeout = 10 * time.Minute
)

func ollamaServerIsRunning() bool {
	defaultURL := "http://localhost:11434/health"
	client := &http.Client{
		Timeout: ollamaServerTimeout,
	}
	ctx, cancel := context.WithTimeout(
		context.Background(),
		ollamaServerTimeout,
	)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, defaultURL, nil)
	if err != nil {
		return false
	}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}

	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func tryStartOllama() (*exec.Cmd, error) {
	binary := llamaPath
	path, err := exec.LookPath(binary)
	if err == nil {
		cmd := exec.Command(path, "serve")
		if err := cmd.Start(); err != nil {
			return nil, fmt.Errorf("error starting ollama: %w", err)
		}

		fmt.Println("Ollama server started using binary:", path)
		return cmd, nil
	}

	localTarget := llamaPath
	if _, err := os.Stat(localTarget); os.IsNotExist(err) {
		fmt.Println("Downloading ollama...")
		err := downloadOllama(context.TODO(), localTarget)
		if err != nil {
			return nil, fmt.Errorf("error installing ollama: %w", err)
		}

		if err := os.Chmod(localTarget, 0o755); err != nil {
			return nil, fmt.Errorf(
				"error setting permissions on ollama binary: %w",
				err,
			)
		}

		fmt.Println("Ollama binary downloaded to:", localTarget)

		cmd := exec.Command(localTarget, "serve")
		if err := cmd.Start(); err != nil {
			return nil, fmt.Errorf("error starting ollama: %w", err)
		}

		fmt.Println("Ollama server started using binary:", localTarget)
		return cmd, nil
	}

	return nil, errors.New("unable to find ollama binary")
}

// this code is part of: https://github.com/redpanda-data/connect/blob/main/internal/impl/ollama/subprocess_unix.go
func downloadOllama(
	ctx context.Context,
	path string,
) error {
	var url string
	const baseURL string = "https://github.com/ollama/ollama/releases/download/v0.3.12/ollama"
	switch runtime.GOOS {
	case "darwin":
		// They ship an universal executable for darwin
		url = baseURL + "-darwin"
	case "linux":
		url = fmt.Sprintf(
			"%s-%s-%s.tgz",
			baseURL,
			runtime.GOOS,
			runtime.GOARCH,
		)
	default:
		return fmt.Errorf(
			"automatic download of ollama is not supported on %s, please download ollama manually",
			runtime.GOOS,
		)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to download ollama binary: %w", err)
	}

	httpClient := &http.Client{
		Timeout: ollamaServerTimeout,
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download ollama binary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"failed to download ollama binary: status_code=%d",
			resp.StatusCode,
		)
	}
	var binary io.Reader = resp.Body
	if strings.HasSuffix(url, ".tgz") {
		gz, err := gzip.NewReader(resp.Body)
		if err != nil {
			return fmt.Errorf(
				"unable to read tarball for ollama binary download: %w",
				err,
			)
		}
		reader := tar.NewReader(gz)
		for {
			header, err := reader.Next()
			if err == io.EOF {
				return fmt.Errorf(
					"unable to find ollama binary within tarball at %s",
					url,
				)
			} else if err != nil {
				return fmt.Errorf("unable to read tarball at %s: %w", url, err)
			}
			if !header.FileInfo().Mode().IsRegular() ||
				header.Name != "./bin/ollama" {
				continue
			}
			binary = reader
			break
		}
	}

	ollama, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o755)
	if err != nil {
		return fmt.Errorf(
			"unable to create file for ollama binary download: %w",
			err,
		)
	}
	defer ollama.Close()

	_, err = io.Copy(ollama, binary)
	if err != nil {
		return fmt.Errorf(
			"unable to download ollama binary to filesystem: %w",
			err,
		)
	}
	return fmt.Errorf("ollama binary downloaded to %s", path)
}

func getOllamaURL() string {
	if os.Getenv("OLLAMA_URL") != "" {
		return os.Getenv("OLLAMA_URL")
	}

	return "http://localhost:11434"
}
