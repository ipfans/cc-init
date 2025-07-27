package main

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
	"strings"
)

// TemplateManager manages the embedded template files
type TemplateManager struct {
	fs     embed.FS
	prefix string
}

// NewTemplateManager creates a new TemplateManager
func NewTemplateManager(embedFS embed.FS, prefix string) *TemplateManager {
	return &TemplateManager{
		fs:     embedFS,
		prefix: prefix,
	}
}

// Walk walks through all files in the embedded filesystem
func (tm *TemplateManager) Walk(fn func(path string, entry fs.DirEntry, err error) error) error {
	return fs.WalkDir(tm.fs, tm.prefix, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fn(path, d, err)
		}

		// Skip the root directory itself
		if path == tm.prefix {
			return nil
		}

		// Calculate relative path from prefix
		relPath := strings.TrimPrefix(path, tm.prefix+"/")
		if relPath == path {
			// This shouldn't happen, but handle it gracefully
			relPath = path
		}

		return fn(relPath, d, nil)
	})
}

// ReadFile reads a file from the embedded filesystem
func (tm *TemplateManager) ReadFile(relPath string) ([]byte, error) {
	fullPath := path.Join(tm.prefix, relPath)
	data, err := tm.fs.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded file %s: %w", fullPath, err)
	}
	return data, nil
}

// GetFileInfo gets file info from the embedded filesystem
func (tm *TemplateManager) GetFileInfo(relPath string) (fs.FileInfo, error) {
	fullPath := path.Join(tm.prefix, relPath)
	file, err := tm.fs.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open embedded file %s: %w", fullPath, err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat embedded file %s: %w", fullPath, err)
	}

	return info, nil
}

// HasTemplates checks if the template directory exists and has content
func (tm *TemplateManager) HasTemplates() bool {
	entries, err := tm.fs.ReadDir(tm.prefix)
	if err != nil {
		return false
	}
	return len(entries) > 0
}

// ListTemplates returns a list of all template files
func (tm *TemplateManager) ListTemplates() ([]string, error) {
	var templates []string
	
	err := tm.Walk(func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		if !entry.IsDir() {
			templates = append(templates, path)
		}
		
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to list templates: %w", err)
	}
	
	return templates, nil
}

// GetDefaultFileMode returns the default file mode for a template file
func (tm *TemplateManager) GetDefaultFileMode(relPath string) fs.FileMode {
	// Check if it's a known executable type
	if strings.HasSuffix(relPath, ".sh") || strings.HasSuffix(relPath, ".bash") {
		return 0755
	}
	
	// Default mode for regular files
	return 0644
}

// GetDefaultDirMode returns the default directory mode
func (tm *TemplateManager) GetDefaultDirMode() fs.FileMode {
	return 0755
}