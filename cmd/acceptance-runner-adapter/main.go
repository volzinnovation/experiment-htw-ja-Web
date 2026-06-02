package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	outcomeTestSuccess         = "test_success"
	outcomeTestFailure         = "test_failure"
	outcomeInfrastructureError = "infrastructure_error"
)

type workerRequest struct {
	ID           string `json:"id"`
	FeatureJSON  string `json:"feature_json"`
	GeneratedDir string `json:"generated_dir"`
	WorkDir      string `json:"work_dir"`
	Timeout      string `json:"timeout"`
}

type workerResponse struct {
	ID       string `json:"id"`
	Outcome  string `json:"outcome"`
	Output   string `json:"output"`
	Error    string `json:"error"`
	Duration int64  `json:"duration"`
}

type testRunner func(context.Context, workerRequest) (string, error)

func main() {
	os.Exit(run(os.Stdin, os.Stdout, os.Stderr, runGeneratedAcceptance))
}

func run(stdin io.Reader, stdout, stderr io.Writer, runner testRunner) int {
	scanner := bufio.NewScanner(stdin)
	encoder := json.NewEncoder(stdout)
	for scanner.Scan() {
		request, err := decodeRequest(scanner.Bytes())
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		if err := encoder.Encode(evaluate(request, runner)); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func decodeRequest(line []byte) (workerRequest, error) {
	var request workerRequest
	if err := json.Unmarshal(line, &request); err != nil {
		return workerRequest{}, err
	}
	if request.ID == "" {
		return workerRequest{}, errors.New("request id is required")
	}
	if request.FeatureJSON == "" {
		return workerRequest{}, errors.New("feature_json is required")
	}
	if request.GeneratedDir == "" {
		request.GeneratedDir = "acceptance/generated"
	}
	return request, nil
}

func evaluate(request workerRequest, runner testRunner) workerResponse {
	started := time.Now()
	ctx, cancel := contextFor(request.Timeout)
	defer cancel()

	output, err := runner(ctx, request)
	response := workerResponse{
		ID:       request.ID,
		Output:   output,
		Duration: int64(time.Since(started)),
	}
	if err == nil {
		response.Outcome = outcomeTestSuccess
		return response
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		response.Outcome = outcomeInfrastructureError
		response.Error = ctx.Err().Error()
		return response
	}
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		response.Outcome = outcomeTestFailure
		return response
	}
	response.Outcome = outcomeInfrastructureError
	response.Error = err.Error()
	return response
}

func contextFor(timeoutText string) (context.Context, context.CancelFunc) {
	if timeoutText == "" {
		return context.WithCancel(context.Background())
	}
	timeout, err := time.ParseDuration(timeoutText)
	if err != nil || timeout <= 0 {
		return context.WithCancel(context.Background())
	}
	return context.WithTimeout(context.Background(), timeout)
}

func runGeneratedAcceptance(ctx context.Context, request workerRequest) (string, error) {
	packagePath := "./" + filepath.ToSlash(request.GeneratedDir)
	cmd := exec.CommandContext(ctx, "go", "test", packagePath)
	cmd.Env = append(os.Environ(),
		"HTW_ACCEPTANCE_FEATURE_JSON="+request.FeatureJSON,
		"GOCACHE="+localPath("tmp/go-cache"),
		"GOMODCACHE="+localPath("tmp/go-modcache"),
	)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func localPath(path string) string {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return absolute
}
