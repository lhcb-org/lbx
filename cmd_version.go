package main

import (
	"fmt"

	"github.com/gonuts/commander"
	"github.com/gonuts/flag"
)

func lbx_make_cmd_version() *commander.Command {
	cmd := &commander.Command{
		Run:       lbx_run_cmd_version,
		UsageLine: "version",
		Short:     "print out script version",
		Long: `
version prints out the script version.

ex:
 $ lbx version
 20140428
`,
		Flag: *flag.NewFlagSet("lbx-version", flag.ExitOnError),
	}
	add_output_level(cmd)
	return cmd
}

func lbx_run_cmd_version(cmd *commander.Command, args []string) error {
	var err error
	fmt.Printf("%s\n", Version)
	return err
}
