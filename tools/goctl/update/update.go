package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"path"

	hexasconfig "github.com/lemongoff/hexas-config"
	"github.com/lemongoff/hexas/core/hash"
	"github.com/lemongoff/hexas/core/logx"
	"github.com/lemongoff/hexas/tools/goctl/update/config"
	"github.com/lemongoff/hexas/tools/goctl/util/pathx"
)

const (
	contentMd5Header = "Content-Md5"
	filename         = "goctl"
)

var configFile = flag.String("f", "update/config/base.yaml", "the YAML configuration file")

func forChksumHandler(file string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !pathx.FileExists(file) {
			logx.Errorf("file %q not exist", file)
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		content, err := os.ReadFile(file)
		if err != nil {
			logx.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		chksum := hash.Md5Hex(content)
		if chksum == r.Header.Get(contentMd5Header) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Set(contentMd5Header, chksum)
		next.ServeHTTP(w, r)
	})
}

func main() {
	flag.Parse()
	manager, err := hexasconfig.NewManager(config.DefaultConfig(), hexasconfig.YAMLFile(*configFile), hexasconfig.Environment("HEXAS_"))
	if err != nil {
		log.Fatal(err)
	}
	if err := manager.Load(context.Background()); err != nil {
		log.Fatal(err)
	}
	snapshot, ok := manager.Current()
	if !ok {
		log.Fatal("configuration was not published")
	}
	c := snapshot.Value()
	fs := http.FileServer(http.Dir(c.FileDir))
	http.Handle(c.FilePath, http.StripPrefix(c.FilePath, forChksumHandler(path.Join(c.FileDir, filename), fs)))
	logx.Must(http.ListenAndServe(c.ListenOn, nil))
}
