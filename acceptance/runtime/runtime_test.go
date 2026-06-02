package runtime

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestRunFeatureFileRunsBackgroundAndScenarioForEachExample(t *testing.T) {
	dir := tempDirInPackage(t)
	path := filepath.Join(dir, "feature.json")
	writeFeature(t, path, Feature{
		Name:       "Runtime feature",
		Background: []Step{{Text: "a background"}},
		Scenarios: []Scenario{{
			Name:     "uses examples",
			Steps:    []Step{{Text: "a scenario step"}},
			Examples: []map[string]string{{"value": "first"}, {"value": "second"}},
		}},
	})

	var calls []string
	RunFeatureFile(t, path, Handlers{
		"a background": func(world *World, example map[string]string) error {
			world.State["value"] = example["value"]
			calls = append(calls, "background:"+example["value"])
			return nil
		},
		"a scenario step": func(world *World, _ map[string]string) error {
			calls = append(calls, "scenario:"+world.State["value"].(string))
			return nil
		},
	})

	want := []string{"background:first", "scenario:first", "background:second", "scenario:second"}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("calls = %v, want %v", calls, want)
	}
}

func TestRunFeatureFileRunsScenarioWithoutExamples(t *testing.T) {
	dir := tempDirInPackage(t)
	path := filepath.Join(dir, "feature.json")
	writeFeature(t, path, Feature{
		Name:      "No examples",
		Scenarios: []Scenario{{Name: "single", Steps: []Step{{Text: "a step"}}}},
	})

	calls := 0
	RunFeatureFile(t, path, Handlers{
		"a step": func(_ *World, example map[string]string) error {
			if len(example) != 0 {
				t.Fatalf("example = %v, want empty default", example)
			}
			calls++
			return nil
		},
	})

	if calls != 1 {
		t.Fatalf("calls = %d, want 1", calls)
	}
}

func TestRunFeatureFileMatchesConcreteStepToPlaceholderHandler(t *testing.T) {
	dir := tempDirInPackage(t)
	path := filepath.Join(dir, "feature.json")
	writeFeature(t, path, Feature{
		Name: "Concrete steps",
		Scenarios: []Scenario{{
			Name:  "extracts placeholder values",
			Steps: []Step{{Text: "the player moves from room 1 to room 2"}},
		}},
	})

	var gotFrom string
	var gotTo string
	RunFeatureFile(t, path, Handlers{
		"the player moves from room <from_room> to room <to_room>": func(_ *World, example map[string]string) error {
			gotFrom = example["from_room"]
			gotTo = example["to_room"]
			return nil
		},
	})

	if gotFrom != "1" || gotTo != "2" {
		t.Fatalf("extracted from=%q to=%q, want 1 and 2", gotFrom, gotTo)
	}
}

func TestTemplateMatchRejectsInvalidText(t *testing.T) {
	tests := []struct {
		name     string
		template string
		text     string
	}{
		{"trailing text", "room <room> done", "room 1 now"},
		{"malformed placeholder", "room <room", "room 1"},
		{"adjacent placeholders", "rooms <first><second>", "rooms 12"},
		{"repeated placeholder mismatch", "room <room> returns to <room>", "room 1 returns to 2"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, ok := matchTemplate(test.template, test.text); ok {
				t.Fatalf("template %q matched %q", test.template, test.text)
			}
		})
	}
}

func TestTemplateMatchAcceptsRepeatedPlaceholderMatch(t *testing.T) {
	extracted, ok := matchTemplate("room <room> returns to <room>", "room 1 returns to 1")
	if !ok {
		t.Fatal("template did not match repeated placeholder with same value")
	}
	if extracted["room"] != "1" {
		t.Fatalf("room = %q, want 1", extracted["room"])
	}
}

func TestSplitPlaceholderValueRejectsInvalidText(t *testing.T) {
	tests := []struct {
		name        string
		text        string
		nextLiteral string
	}{
		{"missing following literal", "room 1", " done"},
		{"empty value before literal", " done", " done"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, _, ok := splitPlaceholderValue(test.text, test.nextLiteral); ok {
				t.Fatalf("split %q before %q", test.text, test.nextLiteral)
			}
		})
	}
}

func TestRunFeatureFileUsesOneBasedExampleNames(t *testing.T) {
	if got, want := executionName("uses examples", 0), "uses examples/example_1"; got != want {
		t.Fatalf("execution name = %q, want %q", got, want)
	}
}

func TestLoadGeneratedFeatureUsesMatchingOverride(t *testing.T) {
	dir := tempDirInPackage(t)
	basePath := filepath.Join(dir, "base.json")
	overridePath := filepath.Join(dir, "override.json")
	writeFeature(t, basePath, Feature{Name: "Generated feature", Scenarios: []Scenario{{Name: "base"}}})
	writeFeature(t, overridePath, Feature{Name: "Generated feature", Scenarios: []Scenario{{Name: "override"}}})
	t.Setenv("HTW_ACCEPTANCE_FEATURE_JSON", overridePath)

	feature, err := loadGeneratedFeature(basePath)
	if err != nil {
		t.Fatal(err)
	}

	if feature.Scenarios[0].Name != "override" {
		t.Fatalf("scenario = %q, want override", feature.Scenarios[0].Name)
	}
}

func TestLoadGeneratedFeatureIgnoresMismatchedOverride(t *testing.T) {
	dir := tempDirInPackage(t)
	basePath := filepath.Join(dir, "base.json")
	overridePath := filepath.Join(dir, "override.json")
	writeFeature(t, basePath, Feature{Name: "Base feature", Scenarios: []Scenario{{Name: "base"}}})
	writeFeature(t, overridePath, Feature{Name: "Other feature", Scenarios: []Scenario{{Name: "override"}}})
	t.Setenv("HTW_ACCEPTANCE_FEATURE_JSON", overridePath)

	feature, err := loadGeneratedFeature(basePath)
	if err != nil {
		t.Fatal(err)
	}

	if feature.Scenarios[0].Name != "base" {
		t.Fatalf("scenario = %q, want base", feature.Scenarios[0].Name)
	}
}

func TestLoadFeatureReportsInvalidJSON(t *testing.T) {
	dir := tempDirInPackage(t)
	path := filepath.Join(dir, "invalid.json")
	if err := os.WriteFile(path, []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := loadFeature(path); err == nil {
		t.Fatal("expected invalid JSON error")
	}
}

func TestResolvePathReportsMissingFile(t *testing.T) {
	if _, err := resolvePath("missing-feature.json"); err == nil {
		t.Fatal("expected missing feature error")
	}
}

func TestResolvePathFindsFileFromGeneratedPackageDepth(t *testing.T) {
	dir := tempDirInPackage(t)
	path := filepath.Join(dir, "feature.json")
	writeFeature(t, path, Feature{Name: "Deep feature"})
	deepDir := filepath.Join("runtime-deep-test", "one", "two", "three", "four")
	if err := os.MkdirAll(deepDir, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll("runtime-deep-test"); err != nil {
			t.Fatal(err)
		}
	})
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatal(err)
		}
	})
	if err := os.Chdir(deepDir); err != nil {
		t.Fatal(err)
	}

	resolved, err := resolvePath(path)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasSuffix(filepath.ToSlash(resolved), filepath.ToSlash(path)) {
		t.Fatalf("resolved path = %q, want suffix %q", resolved, path)
	}
}

func TestExpandReplacesExamplePlaceholders(t *testing.T) {
	got := expand("room <room> has <hazard>", map[string]string{"room": "1", "hazard": "Bats"})
	want := "room 1 has Bats"
	if got != want {
		t.Fatalf("expanded text = %q, want %q", got, want)
	}
}

func tempDirInPackage(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp(".", "runtime-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Fatal(err)
		}
	})
	return dir
}

func writeFeature(t *testing.T, path string, feature Feature) {
	t.Helper()
	data, err := json.Marshal(feature)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}
