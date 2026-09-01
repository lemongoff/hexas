package gogen

import (
	_ "embed"
	"strconv"

	"github.com/lemongoff/hexas/tools/goctl/api/spec"
	"github.com/lemongoff/hexas/tools/goctl/config"
)

const (
	defaultPort        = 8888
	bootstrapConfigDir = "config"
)

//go:embed etc.tpl
var etcTemplate string

func genEtc(dir string, _ *config.Config, api *spec.ApiSpec) error {
	service := api.Service
	host := "0.0.0.0"
	port := strconv.Itoa(defaultPort)

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          bootstrapConfigDir,
		filename:        "base.yaml",
		templateName:    "etcTemplate",
		category:        category,
		templateFile:    etcTemplateFile,
		builtinTemplate: etcTemplate,
		data: map[string]string{
			"serviceName": service.Name,
			"host":        host,
			"port":        port,
		},
	})
}
