package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/gonuts/commander"
	"github.com/gonuts/flag"
	"github.com/gonuts/gas"
	"github.com/gonuts/logger"
	"github.com/gonuts/toml"
	"github.com/lhcb-org/lbx/lbctx"
)

func lbx_make_cmd_init() *commander.Command {
	cmd := &commander.Command{
		Run:       lbx_run_cmd_init,
		UsageLine: "init [options] <project-name> [<project-version>]",
		Short:     "initialize a local development project.",
		Long: `
init initialize a local development project.

ex:
 $ lbx init Gaudi
 $ lbx init Gaudi HEAD
`,
		Flag: *flag.NewFlagSet("lbx-init", flag.ExitOnError),
	}
	add_output_level(cmd)
	add_search_path(cmd)
	add_platform(cmd)

	return cmd
}

func lbx_run_cmd_init(cmd *commander.Command, args []string) error {
	var err error

	g_ctx.SetLevel(logger.Level(cmd.Flag.Lookup("lvl").Value.Get().(int)))

	proj := ""
	vers := ""

	switch len(args) {
	case 1:
		proj = args[0]
		vers = "HEAD"
	case 2:
		proj = args[0]
		vers = args[1]
	default:
		g_ctx.Errorf("lbx-init: needs 2 args (project+version). got=%d\n", len(args))
		return fmt.Errorf("lbx-init: invalid number of arguments")
	}

	proj = lbctx.FixProjectCase(proj)

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	dirname := filepath.Base(pwd)
	local_proj, local_vers := dirname, "HEAD"
	if dirname == "" {
		dirname = proj + "Dev_" + vers
		local_proj = proj + "Dev"
		local_vers = vers
	}

	local_projdir := pwd
	platform := cmd.Flag.Lookup("c").Value.Get().(string)

	if nightly := cmd.Flag.Lookup("nightly").Value.Get().(string); nightly != "" {
		slice := strings.Split(nightly, ",")
		slot := slice[0]
		day := time.Now().Format("Mon")
		if len(slice) > 1 {
			day = slice[1]
		}
		nightly_bases := []string{
			Getenv("LHCBNIGHTLIES", "/afs/cern.ch/lhcb/software/nightlies"),
			filepath.Clean(filepath.Join(Getenv("LCG_release_area", "/afs/cern.ch/sw/lcg/app/releases"), "..", "nightlies")),
		}
		slot_dir := ""
		for i := range nightly_bases {
			dir := filepath.Join(nightly_bases[i], slot)
			if path_exists(dir) {
				slot_dir = dir
				break
			}
		}
		if slot_dir == "" {
			err = fmt.Errorf("lbx-init: could not find slot %q in [%s].", slot, nightly_bases)
			g_ctx.Errorf("%v\n", err)
			return err
		}
		g_ctx.ProjectsPath = append([]string{filepath.Join(slot_dir, day)}, g_ctx.ProjectsPath...)
	}

	// prepend dev-dirs to the search path
	devdirs := cmd.Flag.Lookup("dev-dirs").Value.Get().(string)
	if devdirs != "" {
		g_ctx.ProjectsPath = append(strings.Split(devdirs, string(os.PathListSeparator)), g_ctx.ProjectsPath...)
	}

	g_ctx.Infof(">>> project=%q version=%q\n", proj, vers)
	g_ctx.Infof("local-proj=%q\n", local_proj)
	g_ctx.Infof("local-vers=%q\n", local_vers)
	g_ctx.Infof("local-dir=%q\n", local_projdir)
	g_ctx.Infof("platform=%q\n", platform)

	projdir, err := g_ctx.FindProject(proj, vers, platform)
	if err != nil {
		g_ctx.Errorf("lbx-init: problem finding project: %v\n", err)
		return err
	}
	g_ctx.Infof("using [%s] [%s] from [%s]\n", proj, vers, projdir)

	use_cmake := path_exists(filepath.Join(projdir, proj+"Config.cmake"))
	if !use_cmake {
		g_ctx.Warnf("%s %s does NOT seem to be a CMake-based project\n", proj, vers)
	}

	// create the local dev project
	templates_dir, err := gas.Abs("github.com/lhcb-org/lbx/templates")
	if err != nil {
		g_ctx.Errorf("lbx-init: problem locating templates: %v\n", err)
		return err
	}
	templates := []string{
		"CMakeLists.txt", "toolchain.cmake", "Makefile",
		"searchPath.cmake",
	}

	data := map[string]interface{}{
		"Project":       proj,
		"Version":       vers,
		"SearchPath":    strings.Join(g_ctx.ProjectsPath, " "),
		"SearchPathEnv": strings.Join(g_ctx.ProjectsPath, string(os.PathListSeparator)),
		"UseCMake": func() string {
			if use_cmake {
				return "yes"
			}
			return ""
		}(),
		"PROJECT":      strings.ToUpper(proj),
		"LocalProject": local_proj,
		"LocalVersion": local_vers,
		"CMTProject":   dirname,
		"Slot":         "",
		"Day":          "",
	}

	if nightly := cmd.Flag.Lookup("nightly").Value.Get().(string); nightly != "" {
		templates = append(templates, "nightly.cmake")
		slice := strings.Split(nightly, ",")
		data["Slot"] = slice[0]
		if len(slice) > 1 {
			data["Day"] = slice[1]
		} else {
			data["Day"] = time.Now().Format("Mon")
		}
	}

	for _, tmpl := range templates {
		fname := filepath.Join(templates_dir, tmpl)
		t := template.Must(template.New(tmpl).ParseFiles(fname))
		oname := filepath.Join(local_projdir, tmpl)
		dest, err := os.Create(oname)
		if err != nil {
			g_ctx.Errorf("error creating file [%s]: %v\n", oname, err)
			return err
		}
		defer dest.Close()
		err = t.Execute(dest, data)
		if err != nil {
			g_ctx.Errorf("error running template: %v\n", err)
			return err
		}
	}

	// report success.
	fmt.Printf(`
Successfully create the local project [%[1]s] in [%[2]s]

To start working:

  > cd %[3]s
  > lbx pkg co MyPackage vXrY

then

  > make
  > make test
  > make install

You can customize the configuration by editing the file "CMakeLists.txt"
`, dirname, pwd, local_projdir)

	err = os.MkdirAll(".lbx", 0755)
	if err != nil {
		return err
	}

	conf, err := os.Create(".lbx/config.toml")
	if err != nil {
		return err
	}
	defer conf.Close()

	g_ctx.Project = proj
	g_ctx.Version = vers
	g_ctx.Platform = platform

	err = toml.NewEncoder(conf).Encode(g_ctx)
	return err
}
