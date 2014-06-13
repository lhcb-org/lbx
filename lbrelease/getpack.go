package lbrelease

import (
	"fmt"

	"github.com/lhcb-org/lbx/lbctx"
)

type GetPack struct {
	ReqPkg     string // requested package
	ReqPkgVers string

	pkgs  map[string][]lbctx.RepoInfo
	projs []string
	repos lbctx.RepoInfos

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

	err = gp.initPkgs()
	if err != nil {
		return err
	}

	gp.init = true
	return err
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

	gp.repos = make(lbctx.RepoInfos, 3)

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

	gp.pkgs = make(map[string][]lbctx.RepoInfo)

	for _, repo := range gp.repos {
		for _, p := range repo[0].ListPackages(gp.sel_hat) {
			if _, ok := gp.pkgs[p]; !ok {
				gp.pkgs[p] = make([]lbctx.RepoInfo, 0, 1)
			}
			gp.pkgs[p] = append(gp.pkgs[p], repo[0])
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

	return err
}
