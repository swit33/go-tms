package main

import (
	"flag"
	"go-tms/pkg/boot"
	"go-tms/pkg/config"
	"go-tms/pkg/daemon"
	"go-tms/pkg/handlers"
	"go-tms/pkg/switcher"
)

func main() {

	daemonMode := flag.Bool("d", false, "Run in daemon mode with autosave enabled")
	bootMode := flag.Bool("b", false, "Run in boot mode")

	flag.Parse()

	cfg, err := config.LoadConfig()
	if err != nil {
		handlers.HandleError(err)
	}

	if *daemonMode {
		daemon.RunDaemon(&cfg)
		return
	}

	if *bootMode {
		err := boot.RunBoot()
		if err != nil {
			handlers.HandleError(err)
		}
		return
	}

	err = switcher.RunSwitcher(&cfg)
	if err != nil {
		handlers.HandleError(err)
	}
}
