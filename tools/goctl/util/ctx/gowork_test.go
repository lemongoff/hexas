package ctx

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsGoWorkDisabled(t *testing.T) {
	t.Setenv("GOWORK", "off")
	dir := contextTestDir(t)
	writeContextFile(t, filepath.Join(dir, "go.work"), "go 1.24.0\n")
	ok, err := isGoWork(dir)
	require.NoError(t, err)
	assert.False(t, ok)
	require.NoError(t, UpdateGoWorkIfExist(dir))
}

func TestPrepareRegistersNewWorkspaceModule(t *testing.T) {
	root := contextTestDir(t)
	workFile := filepath.Join(root, "go.work")
	writeContextFile(t, workFile, "go 1.24.0\n")
	t.Setenv("GOWORK", workFile)
	moduleDir := filepath.Join(root, "newproject")
	require.NoError(t, os.Mkdir(moduleDir, 0755))

	ctx, err := PrepareWithModule(moduleDir, "example.com/newproject")
	require.NoError(t, err)
	require.NotNil(t, ctx)
	assert.Equal(t, "example.com/newproject", ctx.Path)
	content, err := os.ReadFile(workFile)
	require.NoError(t, err)
	assert.Contains(t, string(content), "./newproject")
}

func TestPrepareRegistersModuleFromSubdirectory(t *testing.T) {
	root := contextTestDir(t)
	moduleDir := filepath.Join(root, "project")
	workDir := filepath.Join(moduleDir, "subdir")
	require.NoError(t, os.MkdirAll(workDir, 0755))
	writeContextFile(t, filepath.Join(moduleDir, "go.mod"), "module example.com/project\n\ngo 1.24.0\n")
	workFile := filepath.Join(root, "go.work")
	writeContextFile(t, workFile, "go 1.24.0\n")
	t.Setenv("GOWORK", workFile)

	ctx, err := Prepare(workDir)
	require.NoError(t, err)
	require.NotNil(t, ctx)
	assert.Equal(t, moduleDir, ctx.Dir)
	assert.Equal(t, "example.com/project", ctx.Path)
	assert.NoFileExists(t, filepath.Join(workDir, "go.mod"))
}

func Test_isGoWork(t *testing.T) {
	t.Setenv("GOWORK", "")
	dir := contextTestDir(t)

	gowork, err := isGoWork(dir)
	require.NoError(t, err)
	assert.False(t, gowork)

	writeContextFile(t, filepath.Join(dir, "go.work"), "go 1.24.0\n")

	gowork, err = isGoWork(dir)
	require.NoError(t, err)
	assert.True(t, gowork)

	subDir := filepath.Join(dir, "child")
	require.NoError(t, os.Mkdir(subDir, 0755))

	gowork, err = isGoWork(subDir)
	require.NoError(t, err)
	assert.True(t, gowork)
}
