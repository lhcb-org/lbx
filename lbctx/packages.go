package lbctx

type Package struct {
	Name    string // fullname of the package
	Project string // project holding this package
	Repo    string // repository holding the project
}

type Packages map[string]Package

func (p Packages) Packages() []Package {
	pkgs := make([]Package, 0, len(p))
	for _, pkg := range p {
		pkgs = append(pkgs, pkg)
	}
	return pkgs
}
