package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Context represents a saved beads database location.
type Context struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Last bool   `json:"last"`
}

// configDir can be overridden in tests.
var configDir string

// SetConfigDir overrides the config directory (for testing).
func SetConfigDir(dir string) { configDir = dir }

// ConfigDir returns ~/.kbd, creating it if needed.
func ConfigDir() (string, error) {
	if configDir != "" {
		return configDir, os.MkdirAll(configDir, 0o755)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".kbd")
	return dir, os.MkdirAll(dir, 0o755)
}

func contextsPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "contexts.json"), nil
}

// LoadContexts reads saved contexts from ~/.kbd/contexts.json.
func LoadContexts() ([]Context, error) {
	path, err := contextsPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var contexts []Context
	if err := json.Unmarshal(data, &contexts); err != nil {
		return nil, err
	}
	return contexts, nil
}

// SaveContexts writes contexts to ~/.kbd/contexts.json.
func SaveContexts(contexts []Context) error {
	path, err := contextsPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(contexts, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// AddContext adds or updates a context, marks it as last-used.
func AddContext(dbPath string) error {
	contexts, _ := LoadContexts()

	// Clear all last flags
	for i := range contexts {
		contexts[i].Last = false
	}

	// Check for existing entry
	found := false
	for i, c := range contexts {
		if c.Path == dbPath {
			contexts[i].Last = true
			found = true
			break
		}
	}

	if !found {
		name := filepath.Base(filepath.Dir(filepath.Dir(dbPath))) // parent of .beads dir
		contexts = append(contexts, Context{
			Name: name,
			Path: dbPath,
			Last: true,
		})
	}

	return SaveContexts(contexts)
}

// RemoveContext removes a context by path.
func RemoveContext(dbPath string) error {
	contexts, _ := LoadContexts()
	var filtered []Context
	for _, c := range contexts {
		if c.Path != dbPath {
			filtered = append(filtered, c)
		}
	}
	return SaveContexts(filtered)
}

// LastContext returns the most recently used context.
func LastContext() (Context, bool) {
	contexts, _ := LoadContexts()
	for _, c := range contexts {
		if c.Last {
			return c, true
		}
	}
	return Context{}, false
}
