package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os/exec"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestRunProcessesJSONRequests(t *testing.T) {
	input := strings.NewReader(`{"id":"m1","feature_json":"mutated.json","generated_dir":"acceptance/generated"}` + "\n")
	var output bytes.Buffer
	runner := func(_ context.Context, request workerRequest) (string, error) {
		if request.FeatureJSON != "mutated.json" {
			t.Fatalf("feature json = %q", request.FeatureJSON)
		}
		return "ok", nil
	}

	if code := run(input, &output, &bytes.Buffer{}, runner); code != 0 {
		t.Fatalf("exit code = %d", code)
	}

	response := decodeResponse(t, output.String())
	if response.ID != "m1" || response.Outcome != outcomeTestSuccess || response.Output != "ok" {
		t.Fatalf("response = %#v", response)
	}
}

func TestEvaluateClassifiesTestFailure(t *testing.T) {
	request := workerRequest{ID: "m2", FeatureJSON: "mutated.json"}
	response := evaluate(request, func(context.Context, workerRequest) (string, error) {
		return "failed", &exec.ExitError{}
	})

	if response.Outcome != outcomeTestFailure || response.Error != "" {
		t.Fatalf("response = %#v", response)
	}
}

func TestEvaluateClassifiesInfrastructureError(t *testing.T) {
	request := workerRequest{ID: "m3", FeatureJSON: "mutated.json"}
	response := evaluate(request, func(context.Context, workerRequest) (string, error) {
		return "", errors.New("cannot start")
	})

	if response.Outcome != outcomeInfrastructureError || response.Error != "cannot start" {
		t.Fatalf("response = %#v", response)
	}
}

func TestDecodeRequestDefaultsGeneratedDir(t *testing.T) {
	request, err := decodeRequest([]byte(`{"id":"m1","feature_json":"mutated.json"}`))
	if err != nil {
		t.Fatal(err)
	}

	want := workerRequest{ID: "m1", FeatureJSON: "mutated.json", GeneratedDir: "acceptance/generated"}
	if !reflect.DeepEqual(request, want) {
		t.Fatalf("request = %#v, want %#v", request, want)
	}
}

func TestContextForFallsBackToCancelableContext(t *testing.T) {
	for _, timeout := range []string{"", "not-a-duration", "0s", "-1s"} {
		ctx, cancel := contextFor(timeout)
		cancel()
		if err := ctx.Err(); err != context.Canceled {
			t.Fatalf("contextFor(%q) err = %v, want canceled", timeout, err)
		}
	}
}

func TestContextForUsesPositiveTimeout(t *testing.T) {
	ctx, cancel := contextFor("1ms")
	defer cancel()

	select {
	case <-ctx.Done():
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timed context did not expire")
	}
	if err := ctx.Err(); err != context.DeadlineExceeded {
		t.Fatalf("context error = %v, want deadline exceeded", err)
	}
}

func TestRunRejectsMalformedInput(t *testing.T) {
	var stderr bytes.Buffer
	code := run(strings.NewReader("{\n"), &bytes.Buffer{}, &stderr, nil)

	if code != 1 || stderr.Len() == 0 {
		t.Fatalf("code = %d stderr = %q", code, stderr.String())
	}
}

func decodeResponse(t *testing.T, line string) workerResponse {
	t.Helper()
	var response workerResponse
	if err := json.Unmarshal([]byte(line), &response); err != nil {
		t.Fatal(err)
	}
	return response
}
