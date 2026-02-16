package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDetectModulePrefix(t *testing.T) {
	coverage := strings.Join([]string{
		"mode: set",
		"demo-source/cmd/api/main.go:8.2,12.20 2 1",
		"demo-source/internal/service/health.go:5.2,8.24 2 1",
	}, "\n")

	if got := detectModulePrefix(coverage); got != "demo-source" {
		t.Fatalf("unexpected module prefix: got %q", got)
	}
}

func TestBuildCodeFilesMapWithoutCommonPrefix(t *testing.T) {
	root := t.TempDir()
	a := filepath.Join(root, "cmd", "api", "main.go")
	b := filepath.Join(root, "internal", "service", "health.go")

	if err := os.MkdirAll(filepath.Dir(a), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(b), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(a, []byte("package main\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(b, []byte("package service\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	coverage := strings.Join([]string{
		"mode: set",
		"cmd/api/main.go:1.1,2.1 1 1",
		"internal/service/health.go:1.1,2.1 1 0",
	}, "\n")

	files := buildCodeFilesMap(coverage, "", root)
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}
}

func TestResolveCoverageDiskPathWithModulePrefix(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "internal", "service", "health.go")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("package service\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	resolved := resolveCoverageDiskPath("demo-source/internal/service/health.go", "demo-source", root)
	if resolved == "" {
		t.Fatal("expected a resolved path")
	}
	if resolved != path {
		t.Fatalf("unexpected resolved path: got %q want %q", resolved, path)
	}
}
