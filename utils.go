package main

import (
	"fmt"
	"os"

	"github.com/gonuts/commander"
	"github.com/gonuts/logger"
)

func path_exists(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func handle_err(err error) {
	if err != nil {
		if g_ctx != nil {
			g_ctx.Errorf("%v\n", err.Error())
		} else {
			fmt.Fprintf(os.Stderr, "**error** %v\n", err)
		}
		os.Exit(1)
	}
}

func add_search_path(cmd *commander.Command) {
	cmd.Flag.String("user-area", ".", "use the specified path as User_release_area instead of ${User_release_area}")
}

func add_output_level(cmd *commander.Command) {
	cmd.Flag.Bool("v", false, "enable verbose mode")
	cmd.Flag.Int("lvl", logger.INFO, "message level to print")
}

func add_platform(cmd *commander.Command) {
	var plat string
	for _, k := range []string{"BINARY_TAG", "CMTCONFIG"} {
		plat = os.Getenv(k)
		if plat != "" {
			break
		}
	}

	if plat == "" {
		// auto-detect

	}
	cmd.Flag.String("c", plat, "runtime platform")
}

// EOF
