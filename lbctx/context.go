package lbctx

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gonuts/logger"
)

type Context struct {
	msg          *logger.Logger
	ProjectsPath []string // default (project) search path
}

func NewContext(name string) *Context {

	return &Context{
		msg:          logger.New(name),
		ProjectsPath: defaultProjectsPath(),
	}
}

func (ctx *Context) SetLevel(lvl logger.Level) {
	ctx.msg.SetLevel(lvl)
}

func (ctx *Context) Errorf(format string, args ...interface{}) (int, error) {
	return ctx.msg.Errorf(format, args...)
}

func (ctx *Context) Warnf(format string, args ...interface{}) (int, error) {
	return ctx.msg.Warnf(format, args...)
}

func (ctx *Context) Infof(format string, args ...interface{}) (int, error) {
	return ctx.msg.Infof(format, args...)
}

func (ctx *Context) Debugf(format string, args ...interface{}) (int, error) {
	return ctx.msg.Debugf(format, args...)
}

func (ctx *Context) Verbosef(format string, args ...interface{}) (int, error) {
	return ctx.msg.Verbosef(format, args...)
}

// defaultProjectsPath returns the default search path for projects
func defaultProjectsPath() []string {
	projpath := make([]string, 0)
	for _, k := range []string{"CMAKE_PREFIX_PATH", "CMTPROJECTPATH", "LHCBPROJECTPATH"} {
		v := os.Getenv(k)
		if v == "" {
			continue
		}
		vv := strings.Split(v, string(os.PathListSeparator))
		for _, p := range vv {
			projpath = append(projpath, p)
		}
	}

	return projpath
}

// FindProject finds a Gaudi-based project among the Context.ProjectsPath.
func (ctx *Context) FindProject(name, version, platform string) (string, error) {

	// standard project suffixes
	suffixes := []string{
		fmt.Sprintf("%s_%s", name, version),
		filepath.Join(
			strings.ToUpper(name),
			fmt.Sprintf("%s_%s", strings.ToUpper(name), version),
		),
	}

	// special case: with the default 'latest' version, we allow the plain name
	if version == "latest" {
		suffixes = append([]string{name}, suffixes...)
	}

	bindir := filepath.Join("InstallArea", platform)
	for _, path := range ctx.ProjectsPath {
		for _, suffix := range suffixes {
			dir := filepath.Join(path, suffix, bindir)
			ctx.Debugf("checking [%s]...\n", dir)
			_, err := os.Stat(dir)
			if err == nil {
				ctx.Debugf("checking [%s]... [OK]\n", dir)
				return dir, nil
			}
			ctx.Debugf("checking [%s]... [ERR]\n", dir)
		}
	}

	return "", fmt.Errorf(
		"lbx: no such project(name=%q, version=%q, platform=%q, path=%s)",
		name, version, platform, ctx.ProjectsPath,
	)
}
