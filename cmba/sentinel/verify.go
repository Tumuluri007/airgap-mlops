package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Verifier hashes the mounted model file and compares it against the
// expected hash declared in the mounted ModelBinding manifest.
type Verifier struct {
	BindingFile string
	DefaultRoot string
}

// VerifyResult captures one verification outcome.
type VerifyResult struct {
	Match    bool
	Expected string
	Actual   string
	Path     string
}

// NewVerifier constructs a Verifier.
func NewVerifier(bindingFile, defaultRoot string) *Verifier {
	return &Verifier{BindingFile: bindingFile, DefaultRoot: defaultRoot}
}

// modelBinding mirrors the subset of the CRD the sentinel cares about.
type modelBinding struct {
	Spec struct {
		ModelArtifact struct {
			Path   string `yaml:"path"`
			SHA256 string `yaml:"sha256"`
		} `yaml:"modelArtifact"`
	} `yaml:"spec"`
}

// Verify reads the mounted ModelBinding, hashes the referenced file, and
// returns whether the hash matches.
func (v *Verifier) Verify() (*VerifyResult, error) {
	bb, err := os.ReadFile(v.BindingFile)
	if err != nil {
		return nil, fmt.Errorf("read binding: %w", err)
	}
	var mb modelBinding
	if err := yaml.Unmarshal(bb, &mb); err != nil {
		return nil, fmt.Errorf("parse binding: %w", err)
	}
	expected := strings.ToLower(mb.Spec.ModelArtifact.SHA256)
	if expected == "" {
		return nil, fmt.Errorf("binding does not declare modelArtifact.sha256")
	}

	path := mb.Spec.ModelArtifact.Path
	if path == "" {
		return nil, fmt.Errorf("binding does not declare modelArtifact.path")
	}

	actual, err := sha256File(path)
	if err != nil {
		return nil, err
	}

	return &VerifyResult{
		Match:    actual == expected,
		Expected: expected,
		Actual:   actual,
		Path:     path,
	}, nil
}

func sha256File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open model file %s: %w", path, err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("hash model file %s: %w", path, err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
