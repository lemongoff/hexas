package gogen

import (
	_ "embed"

	"github.com/JellyGoFF/FF-Hexas/tools/goctl/api/spec"
	"github.com/JellyGoFF/FF-Hexas/tools/goctl/config"
	"github.com/JellyGoFF/FF-Hexas/tools/goctl/internal/version"
	"github.com/JellyGoFF/FF-Hexas/tools/goctl/util/format"
)

//go:embed svc_test.tpl
var svcTestTemplate string

func genServiceContextTest(dir, rootPkg, projectPkg string, cfg *config.Config, api *spec.ApiSpec) error {
	filename, err := format.FileNamingFormat(cfg.NamingFormat, contextFilename)
	if err != nil {
		return err
	}

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          contextDir,
		filename:        filename + "_test.go",
		templateName:    "svcTestTemplate",
		category:        category,
		templateFile:    svcTestTemplateFile,
		builtinTemplate: svcTestTemplate,
		data: map[string]any{
			"projectPkg": projectPkg,
			"version":    version.BuildVersion,
		},
	})
}
