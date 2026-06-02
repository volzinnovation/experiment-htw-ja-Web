package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type Feature struct {
	Name       string     `json:"name"`
	Background []Step     `json:"background"`
	Scenarios  []Scenario `json:"scenarios"`
}

type Scenario struct {
	Name     string              `json:"name"`
	Steps    []Step              `json:"steps"`
	Examples []map[string]string `json:"examples"`
}

type Step struct {
	Keyword    string   `json:"keyword"`
	Text       string   `json:"text"`
	Parameters []string `json:"parameters"`
}

type World struct {
	State map[string]any
}

type Handler func(*World, map[string]string) error

type Handlers map[string]Handler

func RunFeatureFile(t *testing.T, path string, handlers Handlers) {
	t.Helper()
	feature, err := loadFeature(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, scenario := range feature.Scenarios {
		examples := scenario.Examples
		if len(examples) == 0 {
			examples = []map[string]string{{}}
		}
		for i, example := range examples {
			name := fmt.Sprintf("%s/example_%d", scenario.Name, i+1)
			t.Run(name, func(t *testing.T) {
				world := &World{State: map[string]any{}}
				for _, step := range append(feature.Background, scenario.Steps...) {
					runStep(t, world, handlers, step, example)
				}
			})
		}
	}
}

func loadFeature(path string) (Feature, error) {
	resolved, err := resolvePath(path)
	if err != nil {
		return Feature{}, err
	}
	data, err := os.ReadFile(resolved)
	if err != nil {
		return Feature{}, err
	}
	var feature Feature
	if err := json.Unmarshal(data, &feature); err != nil {
		return Feature{}, err
	}
	return feature, nil
}

func resolvePath(path string) (string, error) {
	for i := 0; i < 6; i++ {
		candidate := path
		for j := 0; j < i; j++ {
			candidate = filepath.Join("..", candidate)
		}
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("could not resolve %s from test working directory", path)
}

func runStep(t *testing.T, world *World, handlers Handlers, step Step, example map[string]string) {
	t.Helper()
	handler, ok := handlers[step.Text]
	if !ok {
		t.Fatalf("unsupported step: %s", step.Text)
	}
	if err := handler(world, example); err != nil {
		t.Fatalf("%s: %v", expand(step.Text, example), err)
	}
}

func expand(text string, example map[string]string) string {
	for key, value := range example {
		text = strings.ReplaceAll(text, "<"+key+">", value)
	}
	return text
}
