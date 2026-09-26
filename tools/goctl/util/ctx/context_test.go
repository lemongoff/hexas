package ctx

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBackground(t *testing.T) {
	root := contextTestDir(t)
	moduleDir := filepath.Join(root, "tools", "generator")
	workDir := filepath.Join(moduleDir, "util", "ctx")
	require.NoError(t, os.MkdirAll(workDir, 0755))
	writeContextFile(t, filepath.Join(root, "go.mod"), "module example.com/framework\n\ngo 1.24.0\n")
	writeContextFile(t, filepath.Join(moduleDir, "go.mod"), "module example.com/framework/tools/generator\n\ngo 1.24.0\n")
	workFile := filepath.Join(root, "go.work")
	workContent := "go 1.24.0\n\nuse (\n\t.\n\t./tools/generator\n)\n"
	writeContextFile(t, workFile, workContent)
	t.Setenv("GOWORK", workFile)
	t.Chdir(workDir)

	ctx, err := Prepare(".")
	require.NoError(t, err)
	require.NotNil(t, ctx)
	assert.Equal(t, workDir, ctx.WorkDir)
	assert.Equal(t, moduleDir, ctx.Dir)
	assert.Equal(t, "example.com/framework/tools/generator", ctx.Path)
	content, err := os.ReadFile(workFile)
	require.NoError(t, err)
	assert.Equal(t, workContent, string(content))
}

func TestPrepareRelativeDirectoryWithoutModule(t *testing.T) {
	root := contextTestDir(t)
	projectDir := filepath.Join(root, "newproject")
	require.NoError(t, os.Mkdir(projectDir, 0755))
	t.Setenv("GOWORK", "off")
	t.Chdir(projectDir)

	ctx, err := Prepare(".")
	require.NoError(t, err)
	require.NotNil(t, ctx)
	assert.Equal(t, "newproject", ctx.Path)
	assert.Equal(t, projectDir, ctx.Dir)
}

func TestPreparePreservesModuleError(t *testing.T) {
	root := contextTestDir(t)
	workDir := filepath.Join(root, "child")
	require.NoError(t, os.Mkdir(workDir, 0755))
	writeContextFile(t, filepath.Join(root, "go.mod"), "module example.com/broken\ninvalid-directive value\n")
	t.Setenv("GOWORK", "off")

	ctx, err := Prepare(workDir)
	require.Error(t, err)
	assert.Nil(t, ctx)
	assert.Contains(t, err.Error(), "invalid-directive")
	assert.NoFileExists(t, filepath.Join(workDir, "go.mod"))
}

func TestPreparePreservesWorkspaceError(t *testing.T) {
	root := contextTestDir(t)
	workDir := filepath.Join(root, "child")
	require.NoError(t, os.Mkdir(workDir, 0755))
	writeContextFile(t, filepath.Join(root, "go.mod"), "module example.com/project\n\ngo 1.24.0\n")
	workFile := filepath.Join(root, "go.work")
	writeContextFile(t, workFile, "go 1.24.0\ninvalid-directive value\n")
	t.Setenv("GOWORK", workFile)

	ctx, err := Prepare(workDir)
	require.Error(t, err)
	assert.Nil(t, ctx)
	assert.Contains(t, err.Error(), "go.work")
	assert.Contains(t, err.Error(), "invalid-directive")
	assert.NoFileExists(t, filepath.Join(workDir, "go.mod"))
}

func TestPrepareMissingGo(t *testing.T) {
	root := contextTestDir(t)
	t.Setenv("PATH", t.TempDir())
	ctx, err := Prepare(root)
	require.Error(t, err)
	assert.Nil(t, ctx)
	assert.True(t, errors.Is(err, exec.ErrNotFound))
	assert.NoFileExists(t, filepath.Join(root, "go.mod"))
}

func TestPrepareExistingModuleWithWorkspaceDisabled(t *testing.T) {
	root := contextTestDir(t)
	moduleFile := filepath.Join(root, "go.mod")
	moduleContent := "module example.com/project\n\ngo 1.24.0\n"
	writeContextFile(t, moduleFile, moduleContent)
	workFile := filepath.Join(root, "go.work")
	writeContextFile(t, workFile, "invalid workspace\n")
	t.Setenv("GOWORK", "off")

	ctx, err := PrepareWithModule(root, "example.com/unused")
	require.NoError(t, err)
	require.NotNil(t, ctx)
	assert.Equal(t, "example.com/project", ctx.Path)
	content, err := os.ReadFile(moduleFile)
	require.NoError(t, err)
	assert.Equal(t, moduleContent, string(content))
	content, err = os.ReadFile(workFile)
	require.NoError(t, err)
	assert.Equal(t, "invalid workspace\n", string(content))
}

func TestPrepareInvalidModuleName(t *testing.T) {
	root := contextTestDir(t)
	t.Setenv("GOWORK", "off")
	ctx, err := PrepareWithModule(root, "invalid module name")
	require.Error(t, err)
	assert.Nil(t, ctx)
	var exitErr *exec.ExitError
	assert.ErrorAs(t, err, &exitErr)
	assert.NoFileExists(t, filepath.Join(root, "go.mod"))
}

func contextTestDir(t *testing.T) string {
	t.Helper()
	dir, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)
	return dir
}

func writeContextFile(t *testing.T, name, content string) {
	t.Helper()
	require.NoError(t, os.WriteFile(name, []byte(content), 0644))
}

func TestBackgroundNilWorkDir(t *testing.T) {
	workDir := ""
	_, err := Prepare(workDir)
	assert.NotNil(t, err)
}

func TestPrepareWithModule(t *testing.T) {
	tests := []struct {
		name       string
		moduleName string
		expectMod  string
	}{
		{
			name:       "custom module name",
			moduleName: "github.com/example/testmodule",
			expectMod:  "github.com/example/testmodule",
		},
		{
			name:       "simple module name",
			moduleName: "simplemodule",
			expectMod:  "simplemodule",
		},
		{
			name:       "empty module name uses directory",
			moduleName: "",
			expectMod:  "", // Will be set to directory name
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory for testing
			tempDir, err := os.MkdirTemp("", "goctl-ctx-test-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			testDir := filepath.Join(tempDir, "testproject")
			err = os.MkdirAll(testDir, 0755)
			require.NoError(t, err)

			ctx, err := PrepareWithModule(testDir, tt.moduleName)
			require.NoError(t, err)
			require.NotNil(t, ctx)

			// Check that the context has expected values
			assert.NotEmpty(t, ctx.WorkDir)
			assert.NotEmpty(t, ctx.Name)
			assert.NotEmpty(t, ctx.Path)
			assert.NotEmpty(t, ctx.Dir)

			// Check that go.mod was created
			goModPath := filepath.Join(testDir, "go.mod")
			assert.FileExists(t, goModPath)

			// Verify module name in go.mod
			content, err := os.ReadFile(goModPath)
			require.NoError(t, err)

			expectedModule := tt.expectMod
			if expectedModule == "" {
				expectedModule = "testproject" // directory name fallback
			}

			assert.Contains(t, string(content), "module "+expectedModule)
			assert.Equal(t, expectedModule, ctx.Path)
		})
	}
}

func TestPrepareWithModule_ExistingGoMod(t *testing.T) {
	// Create a temporary directory with existing go.mod
	tempDir, err := os.MkdirTemp("", "goctl-ctx-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	testDir := filepath.Join(tempDir, "existingproject")
	err = os.MkdirAll(testDir, 0755)
	require.NoError(t, err)

	// Create existing go.mod file
	existingGoMod := `module existing.com/project

go 1.21
`
	goModPath := filepath.Join(testDir, "go.mod")
	err = os.WriteFile(goModPath, []byte(existingGoMod), 0644)
	require.NoError(t, err)

	// PrepareWithModule should use existing go.mod, not create new one
	ctx, err := PrepareWithModule(testDir, "github.com/new/module")
	require.NoError(t, err)
	require.NotNil(t, ctx)

	// Should use existing module name, not the provided one
	assert.Equal(t, "existing.com/project", ctx.Path)

	// Verify go.mod still contains original content
	content, err := os.ReadFile(goModPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "module existing.com/project")
	assert.NotContains(t, string(content), "module github.com/new/module")
}

func TestPrepareWithModule_InvalidWorkDir(t *testing.T) {
	_, err := PrepareWithModule("/non/existent/path", "github.com/example/test")
	assert.Error(t, err)
}

func TestPrepare_CallsPrepareWithModule(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "goctl-ctx-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	testDir := filepath.Join(tempDir, "testproject")
	err = os.MkdirAll(testDir, 0755)
	require.NoError(t, err)

	// Test that Prepare calls PrepareWithModule with empty string
	ctx1, err1 := Prepare(testDir)
	require.NoError(t, err1)

	// Clean up go.mod to test again
	os.Remove(filepath.Join(testDir, "go.mod"))

	ctx2, err2 := PrepareWithModule(testDir, "")
	require.NoError(t, err2)

	// Should produce identical results
	assert.Equal(t, ctx1.Path, ctx2.Path)
	assert.Equal(t, ctx1.Name, ctx2.Name)
}
