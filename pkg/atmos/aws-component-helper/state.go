package aws_component_helper

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	preserveStates = flag.Bool("preserve-states", true, "preserve states")
)

type State struct {
	basepath  string
	namespace string
}

func NewState(name, rootDir string) *State {
	basePath := rootDir
	return &State{
		basepath:  basePath,
		namespace: strings.ToLower(name),
	}
}

func (s *State) Fork(name string) (*State, error) {
	namespace := fmt.Sprintf("%s.%s", s.namespace, name)
	result := NewState(namespace, s.basepath)
	if err := result.SetUp(); err != nil {
		return nil, err
	}
	if err := copyDirectoryRecursively(s.CurrentDir(), result.CurrentDir()); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *State) SetUp() error {
	if err := createDir(s.basepath, ""); err != nil {
		return err
	}

	if err := createDir(s.CurrentDir(), ""); err != nil {
		return err
	}
	return nil
}

func (s *State) Teardown() error {
	if *preserveStates {
		fmt.Printf("Preserve states %s\n", s.namespace)
		return nil
	} else {
		return os.RemoveAll(s.CurrentDir())
	}
}

func (s *State) CurrentDir() string {
	return filepath.Join(s.basepath, s.NamespaceDir())
}

func (s *State) NamespaceDir() string {
	return fmt.Sprintf(".%s", s.namespace)
}

func (s *State) BaseDir() string {
	return s.basepath
}
