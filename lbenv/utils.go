package lbenv

import (
	"fmt"
	"os"
	"strings"
)

// guessType guesses the type of a variable from its name.
// if the name contains PATH or DIRS, then it is a list, otherwise: a scalar.
func guessType(name string) VarType {
	name = strings.ToUpper(name)
	if strings.Contains(name, "PATH") || strings.Contains(name, "DIRS") {
		return VarList
	}
	return VarScalar
}

// splitpath returns a list of paths from a VarList variable
func splitpath(value string) []string {
	return strings.Split(value, string(string(os.PathListSeparator)))
}

func in_str_slice(value string, slice []string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func str_actions(actions []Action) []string {
	o := make([]string, 0, len(actions))
	for _, action := range actions {
		o = append(o, fmt.Sprintf(
			"%[1]v (type=%[1]T)", action,
		))
	}
	return o
}
