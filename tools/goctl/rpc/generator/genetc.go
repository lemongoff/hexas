package generator

import (
	_ "embed"
	"path/filepath"
	"strings"

	conf "github.com/lemongoff/hexas/tools/goctl/config"
	"github.com/lemongoff/hexas/tools/goctl/rpc/parser"
	"github.com/lemongoff/hexas/tools/goctl/util"
	"github.com/lemongoff/hexas/tools/goctl/util/pathx"
	"github.com/lemongoff/hexas/tools/goctl/util/stringx"
)

//go:embed etc.tpl
var etcTemplate string

// GenEtc generates the yaml configuration file of the rpc service,
// including host, port monitoring configuration items and etcd configuration
func (g *Generator) GenEtc(ctx DirContext, _ parser.Proto, cfg *conf.Config) error {
	dir := ctx.GetBootstrapConfig()
	fileName := filepath.Join(dir.Filename, "base.yaml")

	text, err := pathx.LoadTemplate(category, etcTemplateFileFile, etcTemplate)
	if err != nil {
		return err
	}

	return util.With("bootstrapConfig").Parse(text).SaveTo(map[string]any{
		"serviceName": strings.ToLower(stringx.From(ctx.GetServiceName().Source()).ToCamel()),
	}, fileName, false)
}
