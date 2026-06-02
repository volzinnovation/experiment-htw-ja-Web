package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go/format"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type featureIR struct {
	Name string `json:"name"`
}

type metadata struct {
	SchemaVersion      int      `json:"schema_version"`
	FeaturePath        string   `json:"feature_path"`
	IRPath             string   `json:"ir_path"`
	ImplementationHash string   `json:"implementation_hash"`
	HashScope          string   `json:"hash_scope"`
	GeneratedFiles     []string `json:"generated_files"`
}

func main() {
	os.Exit(run(os.Args, os.Stderr))
}

func run(args []string, stderr io.Writer) int {
	if len(args) != 3 {
		fmt.Fprintln(stderr, "usage: acceptance-entrypoint-generator <json-ir> <generated-test-output>")
		return 2
	}
	if err := generate(args[1], args[2]); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func generate(irPath, outputDir string) error {
	feature, err := readFeatureIR(irPath)
	if err != nil {
		return err
	}
	source, err := formattedTestSource(feature.Name, irPath)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return err
	}
	testPath := filepath.Join(outputDir, generatedFilename(irPath))
	if err := os.WriteFile(testPath, source, 0o644); err != nil {
		return err
	}
	return writeMetadata(outputDir, irPath, testPath, source)
}

func readFeatureIR(irPath string) (featureIR, error) {
	contents, err := os.ReadFile(irPath)
	if err != nil {
		return featureIR{}, err
	}
	var feature featureIR
	if err := json.Unmarshal(contents, &feature); err != nil {
		return featureIR{}, err
	}
	if feature.Name == "" {
		return featureIR{}, fmt.Errorf("feature name missing in %s", irPath)
	}
	return feature, nil
}

func formattedTestSource(featureName, irPath string) ([]byte, error) {
	source, err := format.Source([]byte(testSource(featureName, irPath)))
	if err != nil {
		return nil, err
	}
	return source, nil
}

func testSource(featureName, irPath string) string {
	return fmt.Sprintf(`package generated_test

import (
	"testing"

	"htwgo/acceptance/runtime"
	"htwgo/acceptance/steps"
)

func Test%s(t *testing.T) {
	runtime.RunFeatureFile(t, %q, steps.NewHandlers())
}
`, exportedName(featureName), filepath.ToSlash(irPath))
}

func writeMetadata(outputDir, irPath, testPath string, source []byte) error {
	sum := sha256.Sum256(source)
	generatedFile := filepath.ToSlash(testPath)
	data := metadata{
		SchemaVersion:      1,
		FeaturePath:        featurePathForIR(irPath),
		IRPath:             filepath.ToSlash(irPath),
		ImplementationHash: "sha256:" + hex.EncodeToString(sum[:]),
		HashScope:          "generated_files",
		GeneratedFiles:     []string{generatedFile},
	}
	metadataDir := filepath.Join(outputDir, "metadata")
	if err := os.MkdirAll(metadataDir, 0o755); err != nil {
		return err
	}
	encoded, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(metadataDir, metadataFilename(data.FeaturePath)), append(encoded, '\n'), 0o644)
}

func generatedFilename(irPath string) string {
	base := strings.TrimSuffix(filepath.Base(irPath), filepath.Ext(irPath))
	return normalizeIdentifier(base) + "_acceptance_test.go"
}

func exportedName(name string) string {
	parts := regexp.MustCompile(`[^A-Za-z0-9]+`).Split(name, -1)
	var builder strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		builder.WriteString(strings.ToUpper(part[:1]))
		if len(part) > 1 {
			builder.WriteString(part[1:])
		}
	}
	return builder.String()
}

func normalizeIdentifier(name string) string {
	name = strings.ToLower(name)
	return strings.Trim(regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(name, "_"), "_")
}

func featurePathForIR(irPath string) string {
	base := strings.TrimSuffix(filepath.Base(irPath), filepath.Ext(irPath))
	return filepath.ToSlash(filepath.Join("features", base+".feature"))
}

func metadataFilename(featurePath string) string {
	lower := strings.ToLower(featurePath)
	slug := regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(lower, "-")
	return strings.Trim(slug, "-") + ".json"
}
