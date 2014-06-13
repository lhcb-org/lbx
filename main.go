package main

import (
	"os"

	"github.com/gonuts/commander"
	"github.com/gonuts/flag"
	"github.com/lhcb-org/lbx/lbctx"
)

var g_cmd *commander.Command
var g_ctx *lbctx.Context

func init() {
	g_cmd = &commander.Command{
		UsageLine: "lbx",
		Short:     "tools for development.",
		Subcommands: []*commander.Command{
			lbx_make_cmd_init(),
			lbx_make_cmd_pkg(),
			lbx_make_cmd_version(),
			lbx_make_cmd_which(),
		},
		Flag: *flag.NewFlagSet("lbx", flag.ExitOnError),
	}

	g_ctx = lbctx.NewContext("lbx")
}

func main() {
	err := g_cmd.Flag.Parse(os.Args[1:])
	if err != nil {
		handle_err(err)
	}

	args := g_cmd.Flag.Args()
	err = g_cmd.Dispatch(args)
	handle_err(err)
}
