package ctx

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/lemongoff/hexas/tools/goctl/rpc/execx"
	"github.com/lemongoff/hexas/tools/goctl/util/pathx"
)

const goModuleWithoutGoFiles = "command-line-arguments"

var errInvalidGoMod = errors.New("invalid go module")

// Module contains the relative data of go module,
// which is the result of the command go list
type Module struct {
	Path      string
	Main      bool
	Dir       string
	GoMod     string
	GoVersion string
}

func (m *Module) validate() error {
	if m.Path == "" || m.Path == goModuleWithoutGoFiles || m.Dir == "" {
		return errInvalidGoMod
	}
	return nil
}

// projectFromGoMod is used to find the go module and project file path
// the workDir flag specifies which folder we need to detect based on
// only valid for go mod project
func projectFromGoMod(workDir string) (*ProjectContext, error) {
	if len(workDir) == 0 {
		return nil, errors.New("the work directory is not found")
	}
	if _, err := os.Stat(workDir); err != nil {
		return nil, err
	}

	workDir, err := pathx.ReadLink(workDir)
	if err != nil {
		return nil, err
	}

	if err := UpdateGoWorkIfExist(workDir); err != nil {
		return nil, err
	}

	m, err := getRealModule(workDir, execx.Run)
	if err != nil {
		return nil, err
	}
	if err := m.validate(); err != nil {
		return nil, err
	}

	var ret ProjectContext
	ret.WorkDir = workDir
	ret.Name = filepath.Base(m.Dir)
	dir, err := pathx.ReadLink(m.Dir)
	if err != nil {
		return nil, err
	}

	ret.Dir = dir
	ret.Path = m.Path
	return &ret, nil
}

func getRealModule(workDir string, execRun execx.RunFunc) (*Module, error) {
	if workDir == "" {
		return nil, errors.New("the work directory is not found")
	}
	workDir, err := filepath.Abs(workDir)
	if err != nil {
		return nil, err
	}
	workDir, err = pathx.ReadLink(workDir)
	if err != nil {
		return nil, err
	}
	data, err := execRun("go list -json -m", workDir)
	if err != nil {
		return nil, err
	}

	modules, err := decodePackages(strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	var matched *Module
	var matchedDir string
	for i := range modules {
		m := &modules[i]
		if err := m.validate(); err != nil {
			return nil, err
		}
		realDir, err := filepath.Abs(m.Dir)
		if err != nil {
			return nil, err
		}
		realDir, err = pathx.ReadLink(realDir)
		if err != nil {
			return nil, fmt.Errorf("failed to read go.mod, dir: %s, error: %w", m.Dir, err)
		}
		rel, err := filepath.Rel(realDir, workDir)
		if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
			continue
		}
		if matched == nil || len(realDir) > len(matchedDir) {
			matched = m
			matchedDir = realDir
		}
	}
	if matched != nil {
		return matched, nil
	}

	return nil, errors.New("no matched module")
}

func decodePackages(reader io.Reader) ([]Module, error) {
	br := bufio.NewReader(reader)
	if _, err := br.ReadSlice('{'); err != nil {
		return nil, err
	}

	if err := br.UnreadByte(); err != nil {
		return nil, err
	}

	var modules []Module
	decoder := json.NewDecoder(br)
	for decoder.More() {
		var m Module
		if err := decoder.Decode(&m); err != nil {
			return nil, fmt.Errorf("invalid module: %w", err)
		}

		modules = append(modules, m)
	}

	return modules, nil
}
