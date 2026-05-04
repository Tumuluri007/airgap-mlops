package main

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

func writeBinding(t *testing.T, dir, modelPath, sha string) string {
	t.Helper()
	bindingPath := filepath.Join(dir, "binding.yaml")
	contents := []byte("spec:\n  modelArtifact:\n    path: " + modelPath + "\n    sha256: " + sha + "\n")
	if err := os.WriteFile(bindingPath, contents, 0o644); err != nil {
		t.Fatal(err)
	}
	return bindingPath
}

func writeModel(t *testing.T, dir string, content []byte) (string, string) {
	t.Helper()
	modelPath := filepath.Join(dir, "model.pkl")
	if err := os.WriteFile(modelPath, content, 0o644); err != nil {
		t.Fatal(err)
	}
	h := sha256.Sum256(content)
	return modelPath, hex.EncodeToString(h[:])
}

func TestVerifyMatch(t *testing.T) {
	dir := t.TempDir()
	modelPath, expected := writeModel(t, dir, []byte("hello world model"))
	bindingPath := writeBinding(t, dir, modelPath, expected)

	v := NewVerifier(bindingPath, "/")
	res, err := v.Verify()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Match {
		t.Fatalf("expected match: expected=%s actual=%s", res.Expected, res.Actual)
	}
}

func TestVerifyMismatch(t *testing.T) {
	dir := t.TempDir()
	modelPath, _ := writeModel(t, dir, []byte("hello world model"))
	wrongSha := "deadbeef" + "00000000000000000000000000000000000000000000000000000000"
	bindingPath := writeBinding(t, dir, modelPath, wrongSha[:64])

	v := NewVerifier(bindingPath, "/")
	res, err := v.Verify()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Match {
		t.Fatalf("expected mismatch")
	}
}
