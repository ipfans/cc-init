package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// FileSystem defines the interface for file system operations
type FileSystem interface {
	Exists(path string) bool
	CreateDir(path string, perm os.FileMode) error
	CreateFile(path string, content []byte, perm os.FileMode) error
	Walk(root string, fn WalkFunc) error
	Stat(path string) (fs.FileInfo, error)
}

// WalkFunc is the type of function called for each file or directory visited by Walk
type WalkFunc func(path string, info fs.FileInfo, err error) error

// OSFileSystem implements FileSystem interface for the real OS file system
type OSFileSystem struct{}

// NewOSFileSystem creates a new OSFileSystem instance
func NewOSFileSystem() *OSFileSystem {
	return &OSFileSystem{}
}

// Exists checks if a file or directory exists
func (fs *OSFileSystem) Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// CreateDir creates a directory if it doesn't exist
func (fs *OSFileSystem) CreateDir(path string, perm os.FileMode) error {
	if fs.Exists(path) {
		// Check if it's actually a directory
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("failed to stat existing path: %w", err)
		}
		if !info.IsDir() {
			return fmt.Errorf("path exists but is not a directory: %s", path)
		}
		return nil
	}
	return os.MkdirAll(path, perm)
}

// CreateFile creates a file if it doesn't exist
func (fs *OSFileSystem) CreateFile(path string, content []byte, perm os.FileMode) error {
	if fs.Exists(path) {
		// Check if it's actually a file
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("failed to stat existing path: %w", err)
		}
		if info.IsDir() {
			return fmt.Errorf("path exists but is a directory: %s", path)
		}
		return nil
	}

	// Ensure parent directory exists
	dir := filepath.Dir(path)
	if err := fs.CreateDir(dir, 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	// Write the file
	return os.WriteFile(path, content, perm)
}

// Walk walks the file tree rooted at root
func (fs *OSFileSystem) Walk(root string, fn WalkFunc) error {
	return filepath.Walk(root, filepath.WalkFunc(fn))
}

// Stat returns file info for the given path
func (fs *OSFileSystem) Stat(path string) (fs.FileInfo, error) {
	return os.Stat(path)
}

// DryRunFileSystem wraps another FileSystem and simulates operations without making changes
type DryRunFileSystem struct {
	wrapped FileSystem
	logger  *Logger
}

// NewDryRunFileSystem creates a new DryRunFileSystem
func NewDryRunFileSystem(wrapped FileSystem, logger *Logger) *DryRunFileSystem {
	return &DryRunFileSystem{
		wrapped: wrapped,
		logger:  logger,
	}
}

// Exists delegates to the wrapped filesystem
func (fs *DryRunFileSystem) Exists(path string) bool {
	return fs.wrapped.Exists(path)
}

// CreateDir simulates directory creation
func (fs *DryRunFileSystem) CreateDir(path string, perm os.FileMode) error {
	if fs.Exists(path) {
		info, err := fs.wrapped.Stat(path)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return fmt.Errorf("path exists but is not a directory: %s", path)
		}
		fs.logger.Info("Would skip existing directory: %s", path)
		return nil
	}
	fs.logger.Info("Would create directory: %s (mode: %v)", path, perm)
	return nil
}

// CreateFile simulates file creation
func (fs *DryRunFileSystem) CreateFile(path string, content []byte, perm os.FileMode) error {
	if fs.Exists(path) {
		info, err := fs.wrapped.Stat(path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return fmt.Errorf("path exists but is a directory: %s", path)
		}
		fs.logger.Info("Would skip existing file: %s", path)
		return nil
	}
	fs.logger.Info("Would create file: %s (mode: %v, size: %d bytes)", path, perm, len(content))
	return nil
}

// Walk delegates to the wrapped filesystem
func (fs *DryRunFileSystem) Walk(root string, fn WalkFunc) error {
	return fs.wrapped.Walk(root, fn)
}

// Stat delegates to the wrapped filesystem
func (fs *DryRunFileSystem) Stat(path string) (fs.FileInfo, error) {
	return fs.wrapped.Stat(path)
}

