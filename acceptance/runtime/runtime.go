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
	stepExample := example
	if !ok {
		handler, stepExample, ok = matchHandler(handlers, step.Text, example)
	}
	if !ok {
		t.Fatalf("unsupported step: %s", step.Text)
	}
	if err := handler(world, stepExample); err != nil {
		t.Fatalf("%s: %v", expand(step.Text, stepExample), err)
	}
}

func matchHandler(handlers Handlers, text string, example map[string]string) (Handler, map[string]string, bool) {
	for template, handler := range handlers {
		extracted, ok := matchTemplate(template, text)
		if !ok {
			continue
		}
		merged := copyExample(example)
		for key, value := range extracted {
			merged[key] = value
		}
		return handler, merged, true
	}
	return nil, nil, false
}

func matchTemplate(template string, text string) (map[string]string, bool) {
	extracted := map[string]string{}
	remaining := text
	for strings.Contains(template, "<") {
		key, value, nextTemplate, nextRemaining, ok := consumePlaceholder(template, remaining)
		if !ok {
			return nil, false
		}
		if !recordPlaceholder(extracted, key, value) {
			return nil, false
		}
		template = nextTemplate
		remaining = nextRemaining
	}
	if template != remaining {
		return nil, false
	}
	return extracted, true
}

func consumePlaceholder(template string, remaining string) (string, string, string, string, bool) {
	start, end, ok := placeholderBounds(template)
	if !ok {
		return "", "", "", "", false
	}
	prefix := template[:start]
	if !strings.HasPrefix(remaining, prefix) {
		return "", "", "", "", false
	}
	key := template[start+1 : end]
	nextTemplate := template[end+1:]
	nextLiteral, ok := nextLiteral(nextTemplate)
	if !ok {
		return "", "", "", "", false
	}
	value, rest, ok := splitPlaceholderValue(strings.TrimPrefix(remaining, prefix), nextLiteral)
	return key, value, nextTemplate, rest, ok
}

func placeholderBounds(template string) (int, int, bool) {
	start := strings.Index(template, "<")
	end := strings.Index(template[start:], ">")
	if end == -1 {
		return 0, 0, false
	}
	return start, start + end, true
}

func nextLiteral(template string) (string, bool) {
	if strings.HasPrefix(template, "<") {
		return "", false
	}
	nextLiteral := template
	if strings.Contains(template, "<") {
		nextLiteral = template[:strings.Index(template, "<")]
	}
	return nextLiteral, true
}

func recordPlaceholder(extracted map[string]string, key string, value string) bool {
	previous, exists := extracted[key]
	if exists && previous != value {
		return false
	}
	extracted[key] = value
	return true
}

func splitPlaceholderValue(text string, nextLiteral string) (string, string, bool) {
	if nextLiteral == "" {
		return text, "", text != ""
	}
	index := strings.Index(text, nextLiteral)
	if index <= 0 {
		return "", "", false
	}
	value := text[:index]
	return value, text[index:], true
}

func copyExample(example map[string]string) map[string]string {
	copied := make(map[string]string, len(example))
	for key, value := range example {
		copied[key] = value
	}
	return copied
}

func expand(text string, example map[string]string) string {
	for key, value := range example {
		text = strings.ReplaceAll(text, "<"+key+">", value)
	}
	return text
}

// mutate4go-manifest-begin
// {"version":1,"tested_at":"2026-06-02T09:11:59-05:00","module_hash":"36cd089811f54f6aa196f26d168589e5634872c1dca1d5c9c3ed1dc408765019","functions":[{"id":"func/RunFeatureFile","name":"RunFeatureFile","line":38,"end_line":40,"hash":"8eef7028ebe1e2b51951f13392eaebef424c42e0fc143884cfa2307941f64f90"},{"id":"func/RunGeneratedFeatureFile","name":"RunGeneratedFeatureFile","line":42,"end_line":44,"hash":"380f792c822ebcbb414cbabb834a6eb9a28cedace702204100407dabb6603122"},{"id":"func/runFeatureLoader","name":"runFeatureLoader","line":46,"end_line":53,"hash":"634f04f7fcb7b4a9540c149f7565385b0c911aaa64b644b300242e583a480673"},{"id":"func/loadGeneratedFeature","name":"loadGeneratedFeature","line":55,"end_line":72,"hash":"45487fe727947297095dba8e24cf7e17027926b9da7bc7c59351f17c35b26d7d"},{"id":"func/runFeature","name":"runFeature","line":74,"end_line":90,"hash":"337a14bfcb3e7786d3dd72af5109895fe72fef48834b395190567c51ecc9da63"},{"id":"func/executionName","name":"executionName","line":92,"end_line":94,"hash":"09300b19bd1d70e9413b69fbeac780aef5a561accc5cd67d0f0371b0743fd4e2"},{"id":"func/loadFeature","name":"loadFeature","line":96,"end_line":110,"hash":"a0db71ef39c35cc38367b0cb3d24dc598dbe19d876e33cbd1545dad7c29792c3"},{"id":"func/resolvePath","name":"resolvePath","line":112,"end_line":123,"hash":"49da79cdc1bb1f483cd22bc7e9bc01282e6f2fd701b6351255782efc625a4386"},{"id":"func/runStep","name":"runStep","line":125,"end_line":138,"hash":"ade036d8998369cf9321fc3f02614d307007d886b07aed9334035013164b4476"},{"id":"func/matchHandler","name":"matchHandler","line":140,"end_line":153,"hash":"d0412965f4fa2f139215a0961686ad5897f635c70da533eae42ffc32fda18a01"},{"id":"func/matchTemplate","name":"matchTemplate","line":155,"end_line":196,"hash":"4387ef5e166ce5e80fcc613bd22910b98368623e683ef3244ff288ba34607038"},{"id":"func/splitPlaceholderValue","name":"splitPlaceholderValue","line":198,"end_line":208,"hash":"a812bbd6c48d3643b94139df1c846fa9585ccc07a36f94946571304d4427b82e"},{"id":"func/copyExample","name":"copyExample","line":210,"end_line":216,"hash":"5c2f2cae85f99e0918f59b582875dac9d37fa0e2ef49a00d5d5f6e52380d774e"},{"id":"func/expand","name":"expand","line":218,"end_line":223,"hash":"ec5283395045c329ce89208fbc2454ced54d1d1e57cb632902170174e537d84e"}]}
// mutate4go-manifest-end
