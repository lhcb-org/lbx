package main

import (
	"github.com/gonuts/commander"
	"github.com/gonuts/flag"
	"github.com/lhcb-org/lbx/lbx"
)

func lbx_make_cmd_init() *commander.Command {
	cmd := &commander.Command{
		Run:       lbx_run_cmd_init,
		UsageLine: "init [options] <project-name> <project-version>",
		Short:     "initialize a local development project.",
		Long: `
init initialize a local development project.

ex:
 $ lbx init Gaudi trunk
 $ lbx init -name mydev Gaudi trunk
`,
		Flag: *flag.NewFlagSet("lbx-init", flag.ExitOnError),
	}
	add_output_level(cmd)
	add_search_path(cmd)

	cmd.Flag.String("name", "", "name of the local project (default: <project>Dev_<version>)")
	return cmd
}

func lbx_run_cmd_init(cmd *commander.Command, args []string) error {
	var err error

	proj := ""
	vers := ""

	switch len(args) {
	case 2:
		proj = args[0]
		vers = args[1]
	default:
		g_ctx.msg.Errorf("lbx-init: needs 2 args (project+version). got=%d\n", len(args))
		return err
	}

	proj = lbx.FixProjectCase(proj)

	g_ctx.msg.Infof(">>> project=%q version=%q\n", proj, vers)
	return err
}
