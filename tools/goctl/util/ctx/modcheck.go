package ctx

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// IsGoMod reports whether workDir belongs to a module, independently of workspace membership.
func IsGoMod(workDir string) (bool, error) {
	modFile, err := goModFile(workDir)
	return modFile != "", err
}

func goModFile(workDir string) (string, error) {
	if len(workDir) == 0 {
		return "", errors.New("the work directory is not found")
	}
	info, err := os.Stat(workDir)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("work directory %q is not a directory", workDir)
	}

	// GOMOD can be the null device in workspace mode on older Go versions.
	// Disable the workspace only for this lookup; subsequent operations honor it.
	data, err := runGo(workDir, []string{"GOWORK=off"}, "env", "GOMOD")
	if err != nil {
		return "", err
	}
	if data == "" || strings.EqualFold(data, os.DevNull) {
		return "", nil
	}

	return data, nil
}
