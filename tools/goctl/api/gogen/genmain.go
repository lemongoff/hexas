package gogen

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/lemongoff/hexas/tools/goctl/api/spec"
	"github.com/lemongoff/hexas/tools/goctl/config"
	"github.com/lemongoff/hexas/tools/goctl/internal/version"
	"github.com/lemongoff/hexas/tools/goctl/util/format"
	"github.com/lemongoff/hexas/tools/goctl/util/pathx"
	"github.com/lemongoff/hexas/tools/goctl/vars"
)

//go:embed main.tpl
var mainTemplate string

func genMain(dir, rootPkg, projectPkg string, cfg *config.Config, api *spec.ApiSpec) error {
	name := strings.ToLower(api.Service.Name)
	filename, err := format.FileNamingFormat(cfg.NamingFormat, name)
	if err != nil {
		return err
	}

	configName := filename
	if strings.HasSuffix(filename, "-api") {
		filename = strings.ReplaceAll(filename, "-api", "")
	}

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          "",
		filename:        filename + ".go",
		templateName:    "mainTemplate",
		category:        category,
		templateFile:    mainTemplateFile,
		builtinTemplate: mainTemplate,
		data: map[string]string{
			"importPackages": genMainImports(rootPkg),
			"serviceName":    configName,
			"projectPkg":     projectPkg,
			"version":        version.BuildVersion,
		},
	})
}

func genMainImports(parentPkg string) string {
	var imports []string
	imports = append(imports, fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, configDir)))
	imports = append(imports, fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, handlerDir)))
	imports = append(imports, fmt.Sprintf("\"%s\"\n", pathx.JoinPackages(parentPkg, contextDir)))
	imports = append(imports, `hexasconfig "github.com/lemongoff/hexas-config"`)
	imports = append(imports, fmt.Sprintf("\"%s/rest\"", vars.ProjectOpenSourceURL))
	return strings.Join(imports, "\n\t")
}
