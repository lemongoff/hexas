package ctx

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// UpdateGoWorkIfExist adds workDir's containing module to the active workspace
// if it is not already listed. GOWORK=off leaves workspace files untouched.
func UpdateGoWorkIfExist(workDir string) error {
	workFile, err := goWorkFile(workDir)
	if err != nil || workFile == "" {
		return err
	}

	modFile, err := goModFile(workDir)
	if err != nil {
		return err
	}
	if modFile == "" {
		return errModuleCheck
	}
	moduleDir := filepath.Dir(modFile)
	moduleInfo, err := os.Stat(moduleDir)
	if err != nil {
		return err
	}
	data, err := runGo(workDir, nil, "work", "edit", "-json")
	if err != nil {
		return err
	}
	var workspace struct {
		Use []struct{ DiskPath string }
	}
	if err := json.Unmarshal([]byte(data), &workspace); err != nil {
		return fmt.Errorf("decode go.work: %w", err)
	}
	for _, use := range workspace.Use {
		dir := use.DiskPath
		if !filepath.IsAbs(dir) {
			dir = filepath.Join(filepath.Dir(workFile), dir)
		}
		info, err := os.Stat(dir)
		if err != nil {
			return fmt.Errorf("workspace module %q: %w", dir, err)
		}
		if os.SameFile(moduleInfo, info) {
			return nil
		}
	}
	_, err = runGo(moduleDir, nil, "work", "use", ".")
	return err
}

// isGoWork detect if the workDir is in a go workspace
func isGoWork(workDir string) (bool, error) {
	file, err := goWorkFile(workDir)
	return file != "", err
}

func goWorkFile(workDir string) (string, error) {
	if len(workDir) == 0 {
		return "", errors.New("the work directory is not found")
	}
	if _, err := os.Stat(workDir); err != nil {
		return "", err
	}
	goWorkPath, err := runGo(workDir, nil, "env", "GOWORK")
	if err != nil {
		return "", err
	}
	if goWorkPath == "off" {
		return "", nil
	}
	return goWorkPath, nil
}
