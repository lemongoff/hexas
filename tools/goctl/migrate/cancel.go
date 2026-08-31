//go:build linux || darwin || freebsd

package migrate

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/JellyGoFF/FF-Hexas/core/syncx"
	"github.com/JellyGoFF/FF-Hexas/tools/goctl/util/console"
)

func cancelOnSignals() {
	doneChan := syncx.NewDoneChan()
	defer doneChan.Close()

	go func(dc *syncx.DoneChan) {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTSTP, syscall.SIGQUIT)
		defer signal.Stop(c)
		select {
		case <-c:
			console.Error(`
migrate failed, reason: "User Canceled"`)
			os.Exit(0)
		case <-dc.Done():
			return
		}
	}(doneChan)
}
