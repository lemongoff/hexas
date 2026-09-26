package ctx

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var errModuleCheck = errors.New("the work directory must be found in the go mod or the $GOPATH")

// ProjectContext is a structure for the project,
// which contains WorkDir, Name, Path and Dir
type ProjectContext struct {
	WorkDir string
	// Name is the root name of the project
	// eg: go-zero、greet
	Name string
	// Path identifies which module a project belongs to, which is module value if it's a go mod project,
	// or else it is the root name of the project, eg: github.com/lemongoff/hexas、greet
	Path string
	// Dir is the path of the project, eg: /Users/keson/goland/go/go-zero、/Users/keson/go/src/greet
	Dir string
}

// Prepare checks the project which module belongs to,and returns the path and module.
// workDir parameter is the directory of the source of generating code,
// where can be found the project path and the project module,
func Prepare(workDir string) (*ProjectContext, error) {
	return PrepareWithModule(workDir, "")
}

// PrepareWithModule checks the project which module belongs to,and returns the path and module.
// workDir parameter is the directory of the source of generating code,
// where can be found the project path and the project module,
// moduleName parameter is the custom module name to use if creating a new go.mod
// Existing module and workspace errors are returned without initializing a module.
func PrepareWithModule(workDir string, moduleName string) (*ProjectContext, error) {
	if workDir == "" {
		return nil, errors.New("the work directory is not found")
	}
	workDir, err := filepath.Abs(workDir)
	if err != nil {
		return nil, err
	}
	ctx, err := background(workDir)
	if err == nil {
		return ctx, nil
	}
	if !errors.Is(err, errModuleCheck) {
		return nil, err
	}

	var name string
	if len(moduleName) > 0 {
		name = moduleName
	} else {
		name = filepath.Base(workDir)
	}

	_, err = runGo(workDir, nil, "mod", "init", name)
	if err != nil {
		return nil, err
	}

	return background(workDir)
}

// runGo preserves command failures and passes arguments without a shell.
func runGo(workDir string, env []string, args ...string) (string, error) {
	cmd := exec.Command("go", args...)
	cmd.Dir = workDir
	if len(env) > 0 {
		cmd.Env = append(os.Environ(), env...)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("go %s: %w: %s", strings.Join(args, " "), err,
			strings.TrimSpace(stderr.String()))
	}
	return strings.TrimSpace(string(output)), nil
}

func background(workDir string) (*ProjectContext, error) {
	isGoMod, err := IsGoMod(workDir)
	if err != nil {
		return nil, err
	}

	if isGoMod {
		return projectFromGoMod(workDir)
	}
	return projectFromGoPath(workDir)
}
