package main

import (
	"fmt"

	"github.com/gonuts/commander"
	"github.com/gonuts/flag"
	"github.com/gonuts/logger"
)

func lbx_make_cmd_which() *commander.Command {
	cmd := &commander.Command{
		Run:       lbx_run_cmd_which,
		UsageLine: "which <PROJECT> [<PACKAGE> [<VERSION>]]",
		Short:     "return the path to a CMT project or package",
		Long: `
which returns the path to a CMT project or package.

ex:
 $ lbx which GAUDI
 /afs/cern.ch/sw/Gaudi/releases/GAUDI/GAUDI_v25r2

 $ lbx which GAUDI GaudiKernel
 /afs/cern.ch/sw/Gaudi/releases/GAUDI/GAUDI_v25r2/GaudiKernel
`,
		Flag: *flag.NewFlagSet("lbx-which", flag.ExitOnError),
	}
	add_output_level(cmd)
	cmd.Flag.Bool("i", true, "switch on/off case insensitive version")
	cmd.Flag.Bool("d", false, "print the path to the cmt/cmake directory instead of the base dir")
	cmd.Flag.Bool("user-area", true, "enable/disable the user release area when looking for projects")
	return cmd
}

func lbx_run_cmd_which(cmd *commander.Command, args []string) error {
	var err error

	g_ctx.SetLevel(logger.Level(cmd.Flag.Lookup("lvl").Value.Get().(int)))

	proj := ""
	pkg := ""
	vers := ""

	switch len(args) {
	case 1:
		proj = args[0]
	case 2:
		proj = args[0]
		pkg = args[1]
	case 3:
		proj = args[0]
		pkg = args[1]
		vers = args[2]
	default:
		g_ctx.Errorf("lbx-which: needs at least 1 argument. got=%d\n", len(args))
		return fmt.Errorf("lbx-which: invalid number of arguments")
	}

	g_ctx.Infof("which project=%q package=%q version=%q\n", proj, pkg, vers)
	return err
}
