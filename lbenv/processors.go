package lbenv

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// regexp for expanding env.vars.
var g_expenv *regexp.Regexp

func init() {
	g_expenv = regexp.MustCompile(`\$([A-Za-z_][A-Za-z0-9_]*)|\$\(([A-Za-z_][A-Za-z0-9_]*)\)|\$\{([A-Za-z_][A-Za-z0-9_]*)\}|\$\{(\.)\}`)

}

// Processor massages an environment variable value
type Processor func(v *Var, value string, env *Environment) string

// ExpandVar expands the env.var value
func ExpandVar(v *Var, value string, env *Environment) string {
	expval := g_expenv.ReplaceAllStringFunc(value,
		func(str string) string {
			// fmt.Printf(">>> %q\n", str)
			k := str[1:] // "$key"
			if k[0] == '{' {
				k = str[2 : len(str)-1] // "${key}"
			}
			// fmt.Printf("<<< %q\n", k)
			if env.Has(k) {
				return env.Get(k).Value
			}
			return str
		},
	)
	// fmt.Printf("::: expandvar(%q) ==> %q\n", value, expval)
	return expval
}

// PathNormalizer calls filepath.Clean on every entry of the environment variable.
func PathNormalizer(v *Var, value string, env *Environment) string {
	if value == "" {
		return value
	}

	switch v.Type {
	case VarList:
		paths := splitpath(value)
		for i, p := range paths {
			if p == "" {
				continue
			}
			if strings.Contains(p, "://") {
				// might be a URL
				continue
			}
			paths[i] = filepath.Clean(p)
		}
		value = strings.Join(paths, string(os.PathListSeparator))

	case VarScalar:
		if !strings.Contains(value, "://") {
			// not a URL
			value = filepath.Clean(value)
		}
	}

	return value
}

// DuplicatesRemover removes duplicate entries from lists
func DuplicatesRemover(v *Var, value string, env *Environment) string {
	if v.Type == VarScalar {
		return value
	}
	paths := splitpath(value)
	dirs := make([]string, 0, len(paths))
	set := make(map[string]struct{}, len(paths))
	for _, p := range paths {
		if _, dup := set[p]; dup {
			continue
		}
		set[p] = struct{}{}
		dirs = append(dirs, p)
	}
	return strings.Join(dirs, string(os.PathListSeparator))
}

func defaultProcessors() []Processor {
	return []Processor{
		ExpandVar,
		PathNormalizer,
		DuplicatesRemover,
	}
}
