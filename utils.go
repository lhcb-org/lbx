package main

import (
	"fmt"
	"os"

	"github.com/gonuts/commander"
	"github.com/gonuts/logger"
)

// Getenv returns the environment variable associated with key k.
// if it doesn't exist, it returns val.
func Getenv(k, val string) string {
	v := os.Getenv(k)
	if v != "" {
		return v
	}
	return val
}

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
	cmd.Flag.String("dev-dirs", "", "path-list to prepend to the projects-search path")

	cmd.Flag.String("nightly", "", "specify a nightly to use. e.g. slotname,Tue")
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
		plat = "x86_64-linux-gcc-opt"
	}

	cmd.Flag.String("c", plat, "runtime platform")
}

// EOF
