package lbctx

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

type datapkgType struct {
	Path    string
	Version string
}

func (i datapkgType) less(j datapkgType) bool {
	ii := i.version()
	jj := j.version()
	nn := len(ii)
	if nn > len(jj) {
		nn = len(jj)
	}
	for idx := 0; idx < nn; idx++ {
		if ii[idx] != jj[idx] {
			return ii[idx] < jj[idx]
		}
	}
	return len(ii) < len(jj)
}

func (dp datapkgType) version() []int {
	re := regexp.MustCompile(`\d+`)
	slice := re.Split(dp.Version, -1)
	v := make([]int, 0, len(slice))
	for _, str := range slice {
		vv, err := strconv.Atoi(str)
		if err != nil {
			panic(fmt.Errorf("lbx.datapkgType.version: %v (string=%q)", err, str))
		}
		v = append(v, vv)
	}
	return v
}

type datapkgTypes []datapkgType

func (p datapkgTypes) Len() int           { return len(p) }
func (p datapkgTypes) Less(i, j int) bool { return p[i].Version < p[j].Version }
func (p datapkgTypes) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// FindDataPackage finds a data package among the Context.ProjectsPath,
// using optionally the standard suffixes "DBASE", "PARAM" and "EXTRAPACKAGES".
func (ctx *Context) FindDataPackage(name, version string) (string, error) {
	suffixes := []string{"", "EXTRAPACKAGES", "DBASE", "PARAM"}
	versions := make([]datapkgType, 0, 1)

	for _, path := range ctx.ProjectsPath {
		for _, suffix := range suffixes {
			p := filepath.Join(path, suffix, name)
			if _, err := os.Stat(p); err != nil {
				continue
			}
			list, err := filepath.Glob(p + "/*")
			if err != nil {
				return "", err
			}
			for _, v := range list {
				if v == version {
					// stop searching if we've found an exact match
					return filepath.Join(p, v), nil
				}
				if ok, err := filepath.Match(version, v); ok && err == nil {
					versions = append(versions, datapkgType{Path: p, Version: v})
				}
			}
		}
	}

	if len(versions) <= 0 {
		return "", fmt.Errorf("lbx: could not find data package %[1]q %[2]q in %[3]v",
			name, version, ctx.ProjectsPath,
		)
	}

	sort.Reverse(datapkgTypes(versions))
	sel := versions[0]
	return filepath.Join(sel.Path, sel.Version), nil
}
