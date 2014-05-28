package lbx

import (
	"os"
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
