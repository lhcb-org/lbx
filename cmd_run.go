package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gonuts/commander"
	"github.com/gonuts/flag"
	"github.com/lhcb-org/lbx/lbenv"
)

func lbx_make_cmd_run() *commander.Command {
	cmd := &commander.Command{
		Run:       lbx_run_cmd_run,
		UsageLine: "run",
		Short:     "run a command with the proper runtime environment",
		Long: `
run runs a command with the proper runtime environment.

ex:
 $ lbx run echo "hello"
 hello

 $ lbx run gaudi.py some/jobo.py
`,
		Flag: *flag.NewFlagSet("lbx-run", flag.ExitOnError),
	}
	add_output_level(cmd)
	cmd.Flag.Bool("use-grid", false, "enable auto selection of LHCbGrid project")
	cmd.Flag.String("runtime-projects", "", "comma-separated list of runtime projects to add to the environment (e.g.: \"Foo:v1r2,Bar,Baz:v42\"")
	cmd.Flag.String("overriding-projects", "", "comma-separated list of projects to override packages (e.g: \"Foo:v1r2,Bar,Baz:v42\")")
	return cmd
}

func lbx_run_cmd_run(cmd *commander.Command, args []string) error {
	var err error

	switch len(args) {
	case 0:
		g_ctx.Errorf("lbx-run: needs at least one arg (prog-name). got=%d\n", len(args))
		return fmt.Errorf("lbx-run: invalid number of arguments")
	default:
	}

	type pair struct {
		Project string
		Version string
	}
	projects := make([]pair, 0, 2)
	if cmd.Flag.Lookup("use-grid").Value.Get().(bool) {
		projects = append(projects, pair{
			Project: "LHCbGrid",
			Version: "latest",
		})
	}
	for _, p := range strings.Split(cmd.Flag.Lookup("overriding-projects").Value.Get().(string), ",") {
		if p == "" {
			continue
		}
		str := strings.Split(p, ":")
		name := str[0]
		vers := "latest"
		if len(str) > 1 {
			vers = str[1]
		}
		projects = append(projects, pair{
			Project: name,
			Version: vers,
		})
	}

	projects = append(projects, pair{
		Project: g_ctx.Project,
		Version: g_ctx.Version,
	})

	for _, p := range strings.Split(cmd.Flag.Lookup("runtime-projects").Value.Get().(string), ",") {
		if p == "" {
			continue
		}
		str := strings.Split(p, ":")
		name := str[0]
		vers := "latest"
		if len(str) > 1 {
			vers = str[1]
		}
		projects = append(projects, pair{
			Project: name,
			Version: vers,
		})
	}

	// FIXME: add special search-path

	// set the environment XML search path
	xmlenvpath := make([]string, 0, len(projects))
	for _, p := range projects {
		// FIXME: use ExpandVersionAlias
		proj := p.Project
		vers := p.Version
		paths, err := g_ctx.EnvXMLPath(proj, vers, g_ctx.Platform)
		if err != nil {
			g_ctx.Errorf("lbx-run: error looking up ENVXMLPATH: %v\n", err)
			return err
		}
		xmlenvpath = append(xmlenvpath, paths...)
	}
	//fmt.Printf("xml: %v\n", xmlenvpath)

	env := lbenv.New()
	env.SearchPath = xmlenvpath
	env.LoadFromSystem = true

	// load from environment
	for _, val := range os.Environ() {
		kv := strings.Split(val, "=")
		k := kv[0]
		v := os.Getenv(k)
		err = env.Set(k, v)
		if err != nil {
			g_ctx.Errorf("lbx-run: problem initializing env.var %q: %v\n", k, err)
			return err
		}
	}

	// FIXME: handle the extra data packages

	// load the xml files
	for _, p := range projects {
		name := p.Project + "Environment.xml"
		err = env.LoadXMLByName(name)
		if err != nil {
			g_ctx.Errorf("lbx-run: problem loading [%s]: %v\n", name, err)
			return err
		}
	}

	// set the library search path correctly for the non-Linux platforms
	if env.Has("LD_LIBRARY_PATH") {
		k := ""
		// replace LD_LIBRARY_PATH with the corresponding one
		switch runtime.GOOS {
		case "windows":
			k = "PATH"
		case "darwin":
			k = "DYLD_LIBRARY_PATH"
		}
		if k != "" {
			vv := env.Get("LD_LIBRARY_PATH")
			if env.Has(k) {
				v := env.Get(k)
				err = env.Set(k, v.Value+string(os.PathListSeparator)+vv.Value)
			} else {
				err = env.Set(k, vv.Value)
			}

			if err != nil {
				return err
			}
			err = env.Unset("LD_LIBRARY_PATH")
			if err != nil {
				return err
			}
		}
	}
	// extend the prompt variable
	ps1 := os.Getenv("PS1")
	err = env.Set("PS1", fmt.Sprintf("[lbx] %s", ps1))
	if err != nil {
		return err
	}

	bin := exec.Command(args[0], args[1:]...)
	bin.Env = env.Env()
	bin.Stdin = os.Stdin
	bin.Stdout = os.Stdout
	bin.Stderr = os.Stderr

	//fmt.Printf("sub-command: %v\n", bin.Args)
	err = bin.Run()
	return err
}
