package main

import (
	"fmt"
	"path/filepath"

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
	case 1:
		proj = args[0]
		vers = "trunk"
	case 2:
		proj = args[0]
		vers = args[1]
	default:
		g_ctx.Errorf("lbx-init: needs 2 args (project+version). got=%d\n", len(args))
		return fmt.Errorf("lbx-init: invalid number of arguments")
	}

	proj = lbx.FixProjectCase(proj)

	dirname := cmd.Flag.Lookup("name").Value.Get().(string)
	local_proj, local_vers := dirname, "HEAD"
	if dirname == "" {
		dirname = proj + "Dev_" + vers
		local_proj = proj + "Dev"
		local_vers = vers
	}

	usr_area := cmd.Flag.Lookup("user-area").Value.Get().(string)
	if usr_area == "" {
		g_ctx.Errorf("lbx-init: user area not defined (env.var. User_release_area or option -user-area)\n")
		return fmt.Errorf("lbx-init: user-area not defined")
	}
	projdir := filepath.Join(usr_area, dirname)
	if path_exists(projdir) {
		g_ctx.Errorf("lbx-init: directory %q already exists\n", projdir)
		return fmt.Errorf("lbx-init: invalid project dir")
	}

	g_ctx.Infof(">>> project=%q version=%q\n", proj, vers)
	g_ctx.Infof("local-proj=%q\n", local_proj)
	g_ctx.Infof("local-vers=%q\n", local_vers)
	g_ctx.Infof("proj-dir=%q\n", projdir)
	return err
}
