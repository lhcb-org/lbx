package lbx

import (
	"github.com/gonuts/logger"
)

type Context struct {
	msg *logger.Logger
}

func NewContext(name string) *Context {
	return &Context{
		msg: logger.New(name),
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
