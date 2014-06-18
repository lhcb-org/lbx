package lbrelease

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/gonuts/toml"
	"github.com/lhcb-org/lbx/lbctx"
	"github.com/lhcb-org/lbx/lbctx/vcs"
)

type GetPack struct {
	ReqPkg     string // requested package
	ReqPkgVers string

	pkgs  lbctx.Packages
	projs []string
	repos lbctx.RepoDb

	sel_repo string // selected repository
	sel_hat  string // selected repository hat

	proj_name string
	proj_vers string

	init bool
}

func (gp *GetPack) setup() error {
	var err error
	if gp.init {
		return err
	}

	err = gp.initRepos(nil, "", "")
	if err != nil {
		return err
	}

	const fname = ".lbx/packages-db.toml"
	if _, err := os.Stat(fname); err == nil {
		return gp.loadPkgs(fname)
	}

	err = gp.initPkgs()
	if err != nil {
		return err
	}

	gp.init = true

	return gp.savePkgs(fname)
}

func (gp *GetPack) initRepos(excludes []string, user, protocol string) error {
	var err error
	if gp.repos != nil {
		return err
	}

	excl := map[string]struct{}{}
	for _, v := range excludes {
		excl[v] = struct{}{}
	}

	gp.repos = make(lbctx.RepoDb, 3)

	// prepare repositories urls
	// filter the requested protocols for the known repositories
	for k, v := range lbctx.Repositories(user, protocol) {
		if _, dup := excl[k]; dup {
			continue
		}
		gp.repos[k] = v
	}

	if len(gp.repos) <= 0 {
		return fmt.Errorf("getpack: unable to find a repository for the specified protocol")
	}

	return err
}

func (gp *GetPack) initPkgs() error {
	var err error
	if gp.pkgs != nil {
		return err
	}

	gp.pkgs = make(lbctx.Packages)

	pkgs := make(chan []lbctx.Package, len(gp.repos))
	for repo := range gp.repos {
		go func(n string) {
			repo := gp.repos[n]
			pkgs <- repo.ListPackages(gp.sel_hat)
		}(repo)
	}

	for _ = range gp.repos {
		ps := <-pkgs
		for _, pkg := range ps {
			gp.pkgs[pkg.Name] = pkg
		}
	}
	return err
}

func (gp *GetPack) Run() error {
	var err error
	err = gp.setup()
	if err != nil {
		return err
	}

	pkg, ok := gp.pkgs[gp.ReqPkg]
	if !ok {
		return fmt.Errorf("lbrelease: no such package [%s]", gp.ReqPkg)
	}

	var url []string
	switch gp.ReqPkgVers {
	case "", "head", "trunk":
		url = []string{pkg.Repo, pkg.Project, "trunk", pkg.Name}
	default:
		url = []string{pkg.Repo, pkg.Project, "tags", pkg.Name, gp.ReqPkgVers}
	}

	var repo *lbctx.RepoInfo
	for _, r := range gp.repos {
		if r[0].Repo == pkg.Repo {
			repo = &r[0]
			break
		}
	}

	bout, err := vcs.Run(repo.Cmd, "checkout {url} ./{dir}", "url", strings.Join(url, "/"), "dir", pkg.Name)
	if err != nil {
		scan := bufio.NewScanner(bytes.NewReader(bout))
		for scan.Scan() {
			fmt.Fprintf(os.Stderr, "%s\n", scan.Text())
		}
		return err
	}
	return err
}

func (gp *GetPack) loadPkgs(fname string) error {
	ctx := struct {
		Packages lbctx.Packages
	}{}
	_, err := toml.DecodeFile(fname, &ctx)
	if err != nil {
		return err
	}
	gp.pkgs = ctx.Packages
	gp.init = true
	return err
}

func (gp *GetPack) savePkgs(fname string) error {
	ctx := struct {
		Packages lbctx.Packages
	}{
		Packages: gp.pkgs,
	}
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	return toml.NewEncoder(f).Encode(&ctx)
}
