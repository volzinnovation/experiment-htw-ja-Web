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
	runFeatureLoader(t, path, handlers, loadFeature)
}

func RunGeneratedFeatureFile(t *testing.T, path string, handlers Handlers) {
	runFeatureLoader(t, path, handlers, loadGeneratedFeature)
}

func runFeatureLoader(t *testing.T, path string, handlers Handlers, loader func(string) (Feature, error)) {
	t.Helper()
	feature, err := loader(path)
	if err != nil {
		t.Fatal(err)
	}
	runFeature(t, feature, handlers)
}

func loadGeneratedFeature(path string) (Feature, error) {
	feature, err := loadFeature(path)
	if err != nil {
		return Feature{}, err
	}
	overridePath := os.Getenv("HTW_ACCEPTANCE_FEATURE_JSON")
	if overridePath == "" {
		return feature, nil
	}
	override, err := loadFeature(overridePath)
	if err != nil {
		return Feature{}, err
	}
	if override.Name != feature.Name {
		return feature, nil
	}
	return override, nil
}

func runFeature(t *testing.T, feature Feature, handlers Handlers) {
	t.Helper()
	for _, scenario := range feature.Scenarios {
		examples := scenario.Examples
		if len(examples) == 0 {
			examples = []map[string]string{{}}
		}
		for i, example := range examples {
			t.Run(executionName(scenario.Name, i), func(t *testing.T) {
				world := &World{State: map[string]any{}}
				for _, step := range append(feature.Background, scenario.Steps...) {
					runStep(t, world, handlers, step, example)
				}
			})
		}
	}
}

func executionName(scenarioName string, index int) string {
	return fmt.Sprintf("%s/example_%d", scenarioName, index+1)
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

// mutate4go-manifest-begin
// {"version":1,"tested_at":"2026-06-02T08:11:06-05:00","module_hash":"3cef193460973c7af9f05b69b27c9bd1788213e49c888fee10bedea3617bea9c","functions":[{"id":"func/RunFeatureFile","name":"RunFeatureFile","line":38,"end_line":40,"hash":"8eef7028ebe1e2b51951f13392eaebef424c42e0fc143884cfa2307941f64f90"},{"id":"func/RunGeneratedFeatureFile","name":"RunGeneratedFeatureFile","line":42,"end_line":44,"hash":"380f792c822ebcbb414cbabb834a6eb9a28cedace702204100407dabb6603122"},{"id":"func/runFeatureLoader","name":"runFeatureLoader","line":46,"end_line":53,"hash":"634f04f7fcb7b4a9540c149f7565385b0c911aaa64b644b300242e583a480673"},{"id":"func/loadGeneratedFeature","name":"loadGeneratedFeature","line":55,"end_line":72,"hash":"45487fe727947297095dba8e24cf7e17027926b9da7bc7c59351f17c35b26d7d"},{"id":"func/runFeature","name":"runFeature","line":74,"end_line":90,"hash":"337a14bfcb3e7786d3dd72af5109895fe72fef48834b395190567c51ecc9da63"},{"id":"func/executionName","name":"executionName","line":92,"end_line":94,"hash":"09300b19bd1d70e9413b69fbeac780aef5a561accc5cd67d0f0371b0743fd4e2"},{"id":"func/loadFeature","name":"loadFeature","line":96,"end_line":110,"hash":"a0db71ef39c35cc38367b0cb3d24dc598dbe19d876e33cbd1545dad7c29792c3"},{"id":"func/resolvePath","name":"resolvePath","line":112,"end_line":123,"hash":"49da79cdc1bb1f483cd22bc7e9bc01282e6f2fd701b6351255782efc625a4386"},{"id":"func/runStep","name":"runStep","line":125,"end_line":134,"hash":"f37433a52fba0b648511106966bd1caf9197f873cae0e693061471913e6481fc"},{"id":"func/expand","name":"expand","line":136,"end_line":141,"hash":"ec5283395045c329ce89208fbc2454ced54d1d1e57cb632902170174e537d84e"}]}
// mutate4go-manifest-end
