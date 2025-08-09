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
	daemonInternalMode := flag.Bool("daemon-internal", false, "Run in daemon mode with autosave enabled")
	bootMode := flag.Bool("b", false, "Run in boot mode")
	switcherMode := flag.Bool("s", false, "Run in switcher mode")

	flag.Parse()

	cfg, err := config.LoadConfig()
	if err != nil {
		handlers.HandleError(err)
	}

	if *daemonMode {
		daemon.StartDaemon(&cfg)
		return
	}
	if *daemonInternalMode {
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

	if *switcherMode {
		err = switcher.RunSwitcher(&cfg)
		if err != nil {
			handlers.HandleError(err)
		}
	}
}
