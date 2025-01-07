package aws_component_helper

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Global flag to determine whether to preserve state files during teardown
var (
	preserveStates = flag.Bool("preserve-states", true, "preserve states")
)

// State represents the file-based state for a component, including its base directory and namespace
type State struct {
	basepath  string // Base directory for the state
	namespace string // Namespace identifier for the state
}

// NewState initializes a new State instance with the given namespace and root directory
func NewState(name, rootDir string) *State {
	basePath := rootDir
	return &State{
		basepath:  basePath,
		namespace: strings.ToLower(name), // Namespace is converted to lowercase for consistency
	}
}

// Fork creates a new State as a "child" of the current state, copying the contents of the current state directory
func (s *State) Fork(name string) (*State, error) {
	// Create a new namespace by appending the child name to the current namespace
	namespace := fmt.Sprintf("%s.%s", s.namespace, name)
	result := NewState(namespace, s.basepath)

	// Set up the new state directory
	if err := result.SetUp(); err != nil {
		return nil, err
	}

	// Copy the current state directory to the new state directory
	if err := copyDirectoryRecursively(s.CurrentDir(), result.CurrentDir()); err != nil {
		// Clean up the partially created state
		_ = os.RemoveAll(result.CurrentDir())
		return nil, err
	}

	return result, nil
}

// SetUp creates the necessary directories for the state
func (s *State) SetUp() error {
	// Create the base directory if it doesn't exist
	if err := createDir(s.basepath, ""); err != nil {
		return err
	}

	// Create the namespace-specific directory within the base directory
	if err := createDir(s.CurrentDir(), ""); err != nil {
		return err
	}
	return nil
}

// Teardown removes the state directory unless the preserveStates flag is set
func (s *State) Teardown() error {
	if *preserveStates {
		// If preserving states, log the action and skip directory removal
		fmt.Printf("Preserve states %s\n", s.namespace)
		return nil
	} else {
		// Otherwise, remove the state directory
		return os.RemoveAll(s.CurrentDir())
	}
}

// CurrentDir returns the full path to the namespace-specific state directory
func (s *State) CurrentDir() string {
	return filepath.Join(s.basepath, s.NamespaceDir())
}

// NamespaceDir returns the directory name for the namespace (prefixed with a dot)
func (s *State) NamespaceDir() string {
	return fmt.Sprintf(".%s", s.namespace)
}

// BaseDir returns the base directory for the state
func (s *State) BaseDir() string {
	return s.basepath
}
