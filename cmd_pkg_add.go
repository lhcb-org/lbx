package main

import (
	"os"
	"os/exec"

	"github.com/gonuts/commander"
	"github.com/gonuts/flag"
	"github.com/lhcb-org/lbx/lbrelease"
)

func lbx_make_cmd_pkg_add() *commander.Command {
	cmd := &commander.Command{
		Run:       lbx_run_cmd_pkg_add,
		UsageLine: "co [options] <pkg-uri> [<pkg-version>]",
		Short:     "add a package to the current workarea",
		Long: `
co adds a package to the current workarea.

ex:
 $ lbx pkg co MyPackage vXrY
`,
		Flag: *flag.NewFlagSet("lbx-pkg-co", flag.ExitOnError),
	}
	cmd.Flag.Bool("v", false, "enable verbose output")
	cmd.Flag.Bool("go", true, "use the go version")
	return cmd
}

func lbx_run_cmd_pkg_add(cmd *commander.Command, args []string) error {
	var err error

	usego := cmd.Flag.Lookup("go").Value.Get().(bool)

	if !usego {
		// FIXME(sbinet): for the moment, forward to getpack
		getpack, err := exec.LookPath("getpack")
		if err != nil {
			g_ctx.Errorf("lbx-pkg: could not locate 'getpack': %v\n", err)
			return err
		}

		bin := exec.Command(getpack, args...)
		bin.Stdout = os.Stdout
		bin.Stderr = os.Stderr
		err = bin.Run()
		return err
	}

	pkgname := ""
	pkgvers := ""

	switch len(args) {
	case 1:
		pkgname = args[0]
		pkgvers = "head"
	case 2:
		pkgname = args[0]
		pkgvers = args[1]
	default:

	}

	gp := &lbrelease.GetPack{
		ReqPkg:     pkgname,
		ReqPkgVers: pkgvers,
	}

	err = gp.Run()
	return err
}

// EOF
