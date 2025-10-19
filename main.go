package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	eveng "github.com/nawada0615/terraform-provider-eve-ng/eve-ng"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: eveng.Provider,
		ProviderAddr: "local/nawada0615/eve-ng",
	}

	if debugMode {
		opts.Debug = true
	}

	plugin.Serve(opts)
}
