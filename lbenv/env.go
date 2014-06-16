package lbenv

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

// default search path for the environment XML files
var g_envxmlpath = []string{"."}

func init() {
	if envxml := os.Getenv("ENVXMLPATH"); envxml != "" {
		for _, path := range splitpath(envxml) {
			g_envxmlpath = append(g_envxmlpath, path)
		}
	}
}

// ShellType is the type of the shell (sh,csh,bat)
type ShellType int

const (
	ShType  ShellType = iota // Sh-shell (bash, sh, dash, zsh, ...)
	CshType                  // C-Shell (csh, tcsh)
	BatType                  // BAT (windows)
)

// Environment models the recipe(s) to craft and obtain a given environment
type Environment struct {
	LoadFromSystem bool        // whether to load values from system
	SearchPath     []string    // search paths for XML files (used by 'include' elements)
	Processors     []Processor // list of processors to massage env.vars.
	stack          []Action
	vars           map[string]Var
	loaded         map[string]struct{} // set of XML env files already 'included'
	dirstack       []string            // stack of files being processed
}

func New() *Environment {
	wd, _ := os.Getwd()
	return &Environment{
		LoadFromSystem: true,
		SearchPath:     make([]string, 0),
		Processors:     defaultProcessors(),
		stack:          make([]Action, 0),
		vars: map[string]Var{
			".": Var{
				Name:  ".",
				Value: wd,
				Type:  VarScalar,
				Local: true,
			},
		},
		loaded:   make(map[string]struct{}),
		dirstack: make([]string, 0),
	}
}

// Declare declares a new variable in the Environment
func (env *Environment) Declare(name string, vtype VarType, local bool) error {
	var err error
	v, dup := env.vars[name]
	if dup {
		if v.Local != local {
			return fmt.Errorf("lbenv: redeclaration of %q", name)
		} else {
			if vtype != v.Type {
				return fmt.Errorf("lbenv: redeclaration of %q", name)
			}
		}
	}

	v = Var{
		Name:  name,
		Type:  vtype,
		Local: local,
	}
	if env.LoadFromSystem && !local {
		v.Value = os.Getenv(name)
	}
	v.set(env.process(&v, v.Value))

	env.vars[name] = v
	env.stack = append(env.stack, &DeclareVar{
		Name:  name,
		Type:  vtype,
		Local: local,
	})
	return err
}

// Append appends to an existing environment variable (or create a new one)
func (env *Environment) Append(name, value string) error {
	var err error
	v, ok := env.vars[name]
	if !ok {
		local := false
		err = env.Declare(name, guessType(name), local)
		if err != nil {
			return err
		}
		v = env.vars[name]
	}
	v.append(env.process(&v, value))

	env.vars[name] = v
	env.stack = append(env.stack, &AppendVar{
		Name:  name,
		Value: value,
	})
	return err
}

// Prepend prepends to an existing environment variable (or create a new one)
func (env *Environment) Prepend(name, value string) error {
	var err error
	v, ok := env.vars[name]
	if !ok {
		local := false
		err = env.Declare(name, guessType(name), local)
		if err != nil {
			return err
		}
		v = env.vars[name]
	}
	v.prepend(env.process(&v, value))

	env.vars[name] = v
	env.stack = append(env.stack, &PrependVar{
		Name:  name,
		Value: value,
	})
	return err
}

// Set sets a single variable.
// Set overrides any previous value.
func (env *Environment) Set(name, value string) error {
	var err error
	v, ok := env.vars[name]
	if !ok {
		local := false
		err = env.Declare(name, guessType(name), local)
		if err != nil {
			return err
		}
		v = env.vars[name]
	}
	v.set(env.process(&v, value))

	env.vars[name] = v
	env.stack = append(env.stack, &SetVar{
		Name:  name,
		Value: value,
	})
	return err
}

// Unset unsets a single variable to an empty value
// Unset overrides any previous value.
func (env *Environment) Unset(name string) error {
	var err error
	_, ok := env.vars[name]
	if ok {
		delete(env.vars, name)
	}

	env.stack = append(env.stack, &UnsetVar{
		Name: name,
	})
	return err
}

// Remove removes a value from a variable
func (env *Environment) Remove(name, value string) error {
	var err error
	v, ok := env.vars[name]
	if !ok {
		local := false
		err = env.Declare(name, guessType(name), local)
		if err != nil {
			return err
		}
		v = env.vars[name]
	}
	v.remove(env.process(&v, value))

	env.vars[name] = v
	env.stack = append(env.stack, &RemoveVar{
		Name:  name,
		Value: value,
	})
	return err
}

// RemoveRegexp removes a value from a variable, where value is a regexp
func (env *Environment) RemoveRegexp(name, value string) error {
	var err error
	v, ok := env.vars[name]
	if !ok {
		local := false
		err = env.Declare(name, guessType(name), local)
		if err != nil {
			return err
		}
		v = env.vars[name]
	}
	re, err := regexp.Compile(env.process(&v, value))
	if err != nil {
		return err
	}
	v.remove_regexp(re)

	env.vars[name] = v
	env.stack = append(env.stack, &RemoveRegexp{
		Name:  name,
		Value: value,
	})
	return err
}

// Include includes an XML file's definitions into the environment
func (env *Environment) Include(fname, caller, hints string) error {
	var err error
	/*
		env.stack = append(env.stack, &Include{
			File:   fname,
			Caller: caller,
			Hints:  hints,
		})
	*/
	fname, err = env.locate(fname, caller, hints)
	if err != nil {
		return err
	}

	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	err = env.LoadXML(f)

	return err
}

// Has returns whether the environment has a variable named name.
func (env *Environment) Has(name string) bool {
	_, ok := env.vars[name]
	return ok
}

// Get returns the environment variable name, if any
func (env *Environment) Get(name string) Var {
	v, _ := env.vars[name]
	return v
}

// Keys returns the list of environment variables' names.
func (env *Environment) Keys() []string {
	keys := make([]string, 0, len(env.vars))
	for k := range env.vars {
		if k == "." || k == "_" {
			continue
		}
		keys = append(keys, k)
	}
	return keys
}

// SaveXML writes the current state of the environment to w.
func (env *Environment) SaveXML(w io.Writer) error {
	var err error
	err = Encode(w, env.stack)
	return err
}

// LoadXML reads the input reader into the environment.
func (env *Environment) LoadXML(r io.Reader) error {
	var err error

	fname := ""
	// detect whether this reader has already been processed.
	if rr, ok := r.(interface {
		Name() string
	}); ok {
		fname = rr.Name()
		if _, dup := env.loaded[fname]; dup {
			// ignore recursion.
			return nil
		}
	}

	// guard against recursion
	env.loaded[fname] = struct{}{}

	dot := env.vars["."]
	defer func() {
		env.vars["."] = dot
	}()
	// push the previous value of ${.} onto the stack
	env.dirstack = append(env.dirstack, dot.Value)
	env.vars["."] = Var{
		Name:  ".",
		Type:  VarScalar,
		Value: filepath.Dir(fname),
	}
	actions, err := Decode(r)
	if err != nil {
		if err != io.EOF {
			return err
		}
		err = nil
	}

	for _, action := range actions {
		err = env.load(action)
		if err != nil {
			return err
		}
	}

	return err
}

// LoadXMLByName locates a file by name and runs LoadXML
func (env *Environment) LoadXMLByName(fname string) error {
	fname, err := env.locate(fname, "", "")
	if err != nil {
		return err
	}
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()
	return env.LoadXML(f)
}

// GenScript generates shell script replaying the environment modifications.
// shell can be 'sh' or 'csh'
func (env *Environment) GenScript(shell ShellType, w io.Writer) error {
	var err error

	switch shell {
	case ShType:
		_, err = fmt.Fprintf(w, "#!/bin/sh\n")
		if err != nil {
			return err
		}
		for _, k := range env.Keys() {
			v := env.Get(k)
			_, err = fmt.Fprintf(w, "export %s=%q\n", k, v.Value)
			if err != nil {
				return err
			}
		}
		_, err = fmt.Fprintf(w, "## EOF\n")
		if err != nil {
			return err
		}
	default:
		panic(fmt.Errorf("lbenv: unhandled shell type %v\n", shell))
	}
	return err
}

// locate locates a XML file by name in the internal search path
func (env *Environment) locate(file, caller, hints string) (string, error) {
	var err error
	if filepath.IsAbs(file) {
		return file, nil
	}

	dirs := make([]string, len(g_envxmlpath))
	copy(dirs, g_envxmlpath)
	for _, dir := range env.SearchPath {
		dirs = append(dirs, dir)
	}

	if hints != "" {
		for _, dir := range splitpath(hints) {
			dirs = append(dirs, dir)
		}
	}

	if caller != "" {
		dir := filepath.Dir(caller)
		fname := filepath.Join(dir, file)
		if fi, err := os.Stat(fname); err == nil && !fi.IsDir() {
			return fname, nil
		}
		// allow for relative hints
		if hints != "" {
			for _, hdir := range splitpath(hints) {
				dirs = append(dirs, filepath.Join(dir, hdir))
			}
		}
	}
	for _, dir := range dirs {
		fname := filepath.Join(dir, file)
		fi, err := os.Stat(fname)
		if err != nil || fi.IsDir() {
			continue
		}
		return filepath.Abs(fname)
	}

	// nothing found. re-use os.PathError.
	_, err = os.Stat(file)
	return file, err
}

// load loads an action into the environment
func (env *Environment) load(action Action) error {
	var err error
	switch a := action.(type) {
	case *DeclareVar:
		err = env.Declare(a.Name, a.Type, a.Local)
	case *DefaultVar:
		err = env.Set(a.Name, a.Value)
	case *AppendVar:
		err = env.Append(a.Name, a.Value)
	case *PrependVar:
		err = env.Prepend(a.Name, a.Value)
	case *SetVar:
		err = env.Set(a.Name, a.Value)
	case *UnsetVar:
		err = env.Unset(a.Name)
	case *RemoveVar:
		err = env.Remove(a.Name, a.Value)
	case *RemoveRegexp:
		err = env.RemoveRegexp(a.Name, a.Value)
	case *Include:
		err = env.Include(a.File, a.Caller, a.Hints)
	default:
		panic(fmt.Errorf("lbenv: unknown Action: %[1]v (type=%[1]T)", a))
	}
	return err
}

// process runs all the registered processors on value
func (env *Environment) process(v *Var, value string) string {
	for _, process := range env.Processors {
		value = process(v, value, env)
	}
	return value
}
