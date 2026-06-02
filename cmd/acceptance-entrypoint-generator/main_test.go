package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateWritesAcceptanceTestAndMetadata(t *testing.T) {
	dir := tempDirInPackage(t)
	irPath := filepath.Join(dir, "cave-topology.json")
	outputDir := filepath.Join(dir, "generated")
	writeIR(t, irPath, featureIR{Name: "Cave topology"})

	if err := generate(irPath, outputDir); err != nil {
		t.Fatal(err)
	}

	testPath := filepath.Join(outputDir, "cave_topology_acceptance_test.go")
	source, err := os.ReadFile(testPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(source), "func TestCaveTopology(t *testing.T)") {
		t.Fatalf("generated source missing exported test name:\n%s", source)
	}
	if !strings.Contains(string(source), "runtime.RunGeneratedFeatureFile") {
		t.Fatalf("generated source missing generated runtime helper:\n%s", source)
	}
	if !strings.Contains(string(source), filepath.ToSlash(irPath)) {
		t.Fatalf("generated source missing IR path:\n%s", source)
	}

	metadataPath := filepath.Join(outputDir, "metadata", "features-cave-topology-feature.json")
	metadataData, err := os.ReadFile(metadataPath)
	if err != nil {
		t.Fatal(err)
	}
	var got metadata
	if err := json.Unmarshal(metadataData, &got); err != nil {
		t.Fatal(err)
	}
	if got.FeaturePath != "features/cave-topology.feature" {
		t.Fatalf("feature path = %q", got.FeaturePath)
	}
	if len(got.GeneratedFiles) != 1 || got.GeneratedFiles[0] != filepath.ToSlash(testPath) {
		t.Fatalf("generated files = %v, want %s", got.GeneratedFiles, filepath.ToSlash(testPath))
	}
	if !strings.HasPrefix(got.ImplementationHash, "sha256:") {
		t.Fatalf("implementation hash = %q", got.ImplementationHash)
	}
}

func TestGenerateRejectsMissingFeatureName(t *testing.T) {
	dir := tempDirInPackage(t)
	irPath := filepath.Join(dir, "missing-name.json")
	writeIR(t, irPath, featureIR{})

	if err := generate(irPath, filepath.Join(dir, "generated")); err == nil {
		t.Fatal("expected missing feature name error")
	}
}

func TestRunReportsUsageAndGenerationErrors(t *testing.T) {
	var stderr bytes.Buffer
	if code := run([]string{"acceptance-entrypoint-generator"}, &stderr); code != 2 {
		t.Fatalf("usage exit code = %d, want 2", code)
	}
	if !strings.Contains(stderr.String(), "usage: acceptance-entrypoint-generator") {
		t.Fatalf("usage stderr = %q", stderr.String())
	}

	stderr.Reset()
	if code := run([]string{"acceptance-entrypoint-generator", "missing.json", "generated"}, &stderr); code != 1 {
		t.Fatalf("generation error exit code = %d, want 1", code)
	}
	if stderr.Len() == 0 {
		t.Fatal("expected generation error on stderr")
	}
}

func TestReadFeatureIRRejectsInvalidJSON(t *testing.T) {
	dir := tempDirInPackage(t)
	path := filepath.Join(dir, "invalid.json")
	if err := os.WriteFile(path, []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := readFeatureIR(path); err == nil {
		t.Fatal("expected invalid JSON error")
	}
}

func TestIdentifierHelpersNormalizeFeatureNames(t *testing.T) {
	if got, want := exportedName("same setup replay"), "SameSetupReplay"; got != want {
		t.Fatalf("exported name = %q, want %q", got, want)
	}
	if got, want := generatedFilename("build/acceptance/entity-placement.json"), "entity_placement_acceptance_test.go"; got != want {
		t.Fatalf("generated filename = %q, want %q", got, want)
	}
	if got, want := metadataFilename("features/entity-placement.feature"), "features-entity-placement-feature.json"; got != want {
		t.Fatalf("metadata filename = %q, want %q", got, want)
	}
}

func tempDirInPackage(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp(".", "generator-test-*")
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

func writeIR(t *testing.T, path string, feature featureIR) {
	t.Helper()
	data, err := json.Marshal(feature)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}
