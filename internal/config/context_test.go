package config

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestDir(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	SetConfigDir(dir)
	t.Cleanup(func() { SetConfigDir("") })
}

func TestLoadContexts_Empty(t *testing.T) {
	setupTestDir(t)
	contexts, err := LoadContexts()
	if err != nil {
		t.Fatal(err)
	}
	if len(contexts) != 0 {
		t.Fatalf("expected 0 contexts, got %d", len(contexts))
	}
}

func TestSaveAndLoad(t *testing.T) {
	setupTestDir(t)
	want := []Context{
		{Name: "proj", Path: "/tmp/proj/.beads/dolt", Last: true},
	}
	if err := SaveContexts(want); err != nil {
		t.Fatal(err)
	}
	got, err := LoadContexts()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 context, got %d", len(got))
	}
	if got[0].Path != want[0].Path {
		t.Errorf("Path = %q, want %q", got[0].Path, want[0].Path)
	}
	if !got[0].Last {
		t.Error("Last should be true")
	}
}

func TestAddContext_Deduplication(t *testing.T) {
	setupTestDir(t)
	AddContext("/a/.beads/dolt")
	AddContext("/a/.beads/dolt")
	contexts, _ := LoadContexts()
	if len(contexts) != 1 {
		t.Fatalf("expected 1 context after dedup, got %d", len(contexts))
	}
}

func TestAddContext_SetsLast(t *testing.T) {
	setupTestDir(t)
	AddContext("/a/.beads/dolt")
	AddContext("/b/.beads/dolt")
	contexts, _ := LoadContexts()
	if len(contexts) != 2 {
		t.Fatalf("expected 2 contexts, got %d", len(contexts))
	}
	for _, c := range contexts {
		if c.Path == "/b/.beads/dolt" && !c.Last {
			t.Error("/b should be last")
		}
		if c.Path == "/a/.beads/dolt" && c.Last {
			t.Error("/a should not be last")
		}
	}
}

func TestLastContext(t *testing.T) {
	setupTestDir(t)
	AddContext("/a/.beads/dolt")
	AddContext("/b/.beads/dolt")
	c, ok := LastContext()
	if !ok {
		t.Fatal("expected to find last context")
	}
	if c.Path != "/b/.beads/dolt" {
		t.Errorf("Path = %q, want /b/.beads/dolt", c.Path)
	}
}

func TestLastContext_None(t *testing.T) {
	setupTestDir(t)
	_, ok := LastContext()
	if ok {
		t.Error("expected no last context")
	}
}

func TestRemoveContext(t *testing.T) {
	setupTestDir(t)
	AddContext("/a/.beads/dolt")
	AddContext("/b/.beads/dolt")
	RemoveContext("/a/.beads/dolt")
	contexts, _ := LoadContexts()
	if len(contexts) != 1 {
		t.Fatalf("expected 1 context, got %d", len(contexts))
	}
	if contexts[0].Path != "/b/.beads/dolt" {
		t.Errorf("remaining = %q", contexts[0].Path)
	}
}

func TestAddContext_Name(t *testing.T) {
	setupTestDir(t)
	AddContext("/Users/me/projects/myproj/.beads/dolt")
	contexts, _ := LoadContexts()
	if contexts[0].Name != "myproj" {
		t.Errorf("Name = %q, want 'myproj'", contexts[0].Name)
	}
}

func TestConfigDir_Creates(t *testing.T) {
	tmp := t.TempDir()
	sub := filepath.Join(tmp, "nested", "kbd")
	SetConfigDir(sub)
	defer SetConfigDir("")

	dir, err := ConfigDir()
	if err != nil {
		t.Fatal(err)
	}
	if dir != sub {
		t.Errorf("dir = %q", dir)
	}
	info, err := os.Stat(sub)
	if err != nil {
		t.Fatal("dir not created:", err)
	}
	if !info.IsDir() {
		t.Error("not a directory")
	}
}
