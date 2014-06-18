package lbctx

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/lhcb-org/lbx/lbctx/vcs"
)

// Repos is the database of known repositories
var Repos = RepoDb{
	"gaudi": RepoInfos{
		{
			Cmd:  vcs.Svn,
			Repo: "svn+ssh://svn.cern.ch/reps/gaudi",
		},
		{
			Cmd:  vcs.Svn,
			Repo: "http://svn.cern.ch/guest/gaudi",
		},
	},

	"lbsvn": RepoInfos{
		{
			Cmd:  vcs.Svn,
			Repo: "svn+ssh://svn.cern.ch/reps/lhcb",
		},
		{
			Cmd:  vcs.Svn,
			Repo: "http://svn.cern.ch/guest/lhcb",
		},
	},

	"dirac": RepoInfos{
		{
			Cmd:  vcs.Svn,
			Repo: "svn+ssh://svn.cern.ch/reps/dirac",
		},
		{
			Cmd:  vcs.Svn,
			Repo: "http://svn.cern.ch/guest/dirac",
		},
	},

	"lhcbint": RepoInfos{
		{
			Cmd:  vcs.Svn,
			Repo: "svn+ssh://svn.cern.ch/reps/lhcbint",
		},
	},
}

type RepoInfo struct {
	Cmd  *vcs.Cmd
	Repo string

	pkgs Packages
}

type RepoInfos []RepoInfo
type RepoDb map[string]RepoInfos

// Repositories returns a map of named-repositories
func Repositories(user, protocol string) RepoDb {
	repos := make(RepoDb, len(Repos))
	for k := range Repos {
		repos[k] = append([]RepoInfo{}, Repos[k]...)
	}
	if user != "" {

	}
	return repos
}

func (repos *RepoInfos) ListPackages(hat string) []Package {
	for _, repo := range *repos {
		pkgs := repo.ListPackages(hat)
		if pkgs != nil {
			return pkgs
		}
	}
	return nil
}

func (repo *RepoInfo) ListPackages(hat string) []Package {
	if repo.pkgs == nil {
		err := repo.initPkgs()
		if err != nil {
			return nil
		}
	}
	pkgs := make([]Package, 0)
	for _, pkg := range repo.pkgs {
		if !strings.HasPrefix(pkg.Name, hat) {
			continue
		}
		pkgs = append(pkgs, pkg)
	}
	return pkgs
}

func (repo *RepoInfo) initPkgs() error {
	var err error

	// FIXME: first check 'propget version' >= 2.0

	// assume propget-version >= 2.0
	bout, err := vcs.Run(repo.Cmd, "propget packages {repo}", "repo", repo.Repo)
	if err != nil {
		return nil
	}

	pkgs := make(Packages)
	scan := bufio.NewScanner(bytes.NewReader(bout))
	for scan.Scan() {
		bline := bytes.Trim(scan.Bytes(), " \n")
		if bytes.HasPrefix(bline, []byte("#")) {
			continue
		}
		bline = bytes.Replace(bline, []byte("\t"), []byte(" "), -1)
		tokens := make([]string, 0)
		for _, tok := range strings.Split(string(bline), " ") {
			tok = strings.Trim(tok, " \t\n")
			if tok == "" {
				continue
			}
			tokens = append(tokens, tok)
		}

		if len(tokens) <= 0 {
			continue
		}
		project := ""
		pkgname := tokens[0]
		if len(tokens) > 1 {
			project = tokens[1]
		}
		pkgs[pkgname] = Package{
			Name:    pkgname,
			Project: project,
			Repo:    repo.Repo,
		}
	}
	err = scan.Err()
	if err != nil {
		return err
	}

	repo.pkgs = pkgs
	return err
}
