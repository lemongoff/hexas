package ctx

import (
	"go/build"
	"os"
	"path/filepath"
	"testing"

	"github.com/lemongoff/hexas/core/stringx"
	"github.com/lemongoff/hexas/tools/goctl/rpc/execx"
	"github.com/lemongoff/hexas/tools/goctl/util/pathx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsGoModModes(t *testing.T) {
	root := contextTestDir(t)
	moduleDir := filepath.Join(root, "project")
	workDir := filepath.Join(moduleDir, "child")
	require.NoError(t, os.MkdirAll(workDir, 0755))
	writeContextFile(t, filepath.Join(moduleDir, "go.mod"), "module example.com/project\n\ngo 1.24.0\n")
	workFile := filepath.Join(root, "go.work")
	writeContextFile(t, workFile, "go 1.24.0\n\nuse ./project\n")

	for _, mode := range []string{"off", workFile} {
		t.Run(mode, func(t *testing.T) {
			t.Setenv("GOWORK", mode)
			t.Setenv("GO111MODULE", "on")
			ok, err := IsGoMod(workDir)
			require.NoError(t, err)
			assert.True(t, ok)
			ok, err = IsGoMod(root)
			require.NoError(t, err)
			assert.False(t, ok)
		})
	}
	t.Run("modules disabled", func(t *testing.T) {
		t.Setenv("GO111MODULE", "off")
		ok, err := IsGoMod(workDir)
		require.NoError(t, err)
		assert.False(t, ok)
	})
}

func TestIsGoModRejectsFile(t *testing.T) {
	file := filepath.Join(contextTestDir(t), "go.mod")
	writeContextFile(t, file, "module example.com/project\n")
	ok, err := IsGoMod(file)
	require.Error(t, err)
	assert.False(t, ok)
}

func TestIsGoMod(t *testing.T) {
	// create mod project
	dft := build.Default
	gp := dft.GOPATH
	if len(gp) == 0 {
		return
	}
	projectName := stringx.Rand()
	dir := filepath.Join(gp, "src", projectName)
	err := pathx.MkdirIfNotExist(dir)
	if err != nil {
		return
	}

	_, err = execx.Run("go mod init "+projectName, dir)
	assert.Nil(t, err)
	defer func() {
		_ = os.RemoveAll(dir)
	}()

	isGoMod, err := IsGoMod(dir)
	assert.Nil(t, err)
	assert.Equal(t, true, isGoMod)
}

func TestIsGoModNot(t *testing.T) {
	dft := build.Default
	gp := dft.GOPATH
	if len(gp) == 0 {
		return
	}
	projectName := stringx.Rand()
	dir := filepath.Join(gp, "src", projectName)
	err := pathx.MkdirIfNotExist(dir)
	if err != nil {
		return
	}

	defer func() {
		_ = os.RemoveAll(dir)
	}()

	isGoMod, err := IsGoMod(dir)
	assert.Nil(t, err)
	assert.Equal(t, false, isGoMod)
}

func TestIsGoModWorkDirIsNil(t *testing.T) {
	_, err := IsGoMod("")
	assert.Equal(t, err.Error(), func() string {
		return "the work directory is not found"
	}())
}
