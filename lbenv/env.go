package lbenv

import (
	"fmt"
	"os"
)

// Environment models the recipe(s) to craft and obtain a given environment
type Environment struct {
	LoadFromSystem bool     // whether to load values from system
	SearchPath     []string // search paths for XML files (used by 'include' elements)
	stack          []Action
	vars           map[string]Var
}

func New() *Environment {
	return &Environment{
		SearchPath: make([]string, 0),
		stack:      make([]Action, 0),
		vars: map[string]Var{
			".": Var{
				Name:  ".",
				Value: "",
				Type:  VarScalar,
				Local: true,
			},
		},
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
	v.append(value)

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
	v.prepend(value)

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
	v.Value = value

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
	v.remove(value)

	env.vars[name] = v
	env.stack = append(env.stack, &RemoveVar{
		Name:  name,
		Value: value,
	})
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
		keys = append(keys, k)
	}
	return keys
}
