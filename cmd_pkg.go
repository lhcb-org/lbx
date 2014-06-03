package main

import (
	"github.com/gonuts/commander"
	"github.com/gonuts/flag"
)

func lbx_make_cmd_pkg() *commander.Command {
	cmd := &commander.Command{
		UsageLine: "pkg [options]",
		Short:     "add, remove or inspect sub-packages",
		Subcommands: []*commander.Command{
			lbx_make_cmd_pkg_add(),
			// lbx_make_cmd_pkg_create(),
			// lbx_make_cmd_pkg_ls(),
			// lbx_make_cmd_pkg_rm(),
		},
		Flag: *flag.NewFlagSet("lbx-pkg", flag.ExitOnError),
	}
	return cmd
}

// EOF
