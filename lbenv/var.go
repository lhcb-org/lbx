package lbenv

import (
	"os"
	"strings"
)

// Var is an environment variable.
type Var struct {
	Name  string
	Value string
	Type  VarType
	Local bool
}

func (v *Var) append(value string) {
	switch v.Type {
	case VarList:
		v.Value += string(os.PathListSeparator) + value
	case VarScalar:
		v.Value += value
	}
}

func (v *Var) prepend(value string) {
	switch v.Type {
	case VarList:
		v.Value = value + string(os.PathListSeparator) + v.Value
	case VarScalar:
		v.Value = value + v.Value
	}
}

func (v *Var) remove(value string) {
	switch v.Type {
	case VarList:
		vals := make([]string, 0)
		for _, vv := range splitpath(v.Value) {
			if vv != value {
				vals = append(vals, vv)
			}
		}
		v.Value = strings.Join(vals, string(os.PathListSeparator))
	case VarScalar:
		v.Value = strings.Replace(v.Value, value, "", -1)
	}
}
