package lbx

import (
	"github.com/lhcb-org/lbx/lbx/vcs"
)

// Repos is the database of known repositories
var Repos = RepoInfos{
	"gaudi": []RepoInfo{
		{
			Cmd:  vcs.Svn,
			Repo: "svn+ssh://svn.cern.ch/reps/gaudi",
			Root: "/gaudi",
		},
		{
			Cmd:  vcs.Svn,
			Repo: "http://svn.cern.ch/guest/gaudi",
			Root: "/gaudi",
		},
	},

	"lbsvn": []RepoInfo{
		{
			Cmd:  vcs.Svn,
			Repo: "svn+ssh://svn.cern.ch/reps/lhcb",
			Root: "/lhcb",
		},
		{
			Cmd:  vcs.Svn,
			Repo: "http://svn.cern.ch/guest/lhcb",
			Root: "/lhcb",
		},
	},

	"dirac": []RepoInfo{
		{
			Cmd:  vcs.Svn,
			Repo: "svn+ssh://svn.cern.ch/reps/dirac",
			Root: "/dirac",
		},
		{
			Cmd:  vcs.Svn,
			Repo: "http://svn.cern.ch/guest/dirac",
			Root: "/dirac",
		},
	},

	"lhcbint": []RepoInfo{
		{
			Cmd:  vcs.Svn,
			Repo: "svn+ssh://svn.cern.ch/reps/lhcbint",
			Root: "/lhcbint",
		},
	},
}

type RepoInfo struct {
	Cmd  *vcs.Cmd
	Repo string
	Root string
}

type RepoInfos map[string][]RepoInfo

// Repositories returns a map of named-repositories
func Repositories(user, protocol string) RepoInfos {
	repos := make(RepoInfos, len(Repos))
	for k := range Repos {
		repos[k] = append([]RepoInfo{}, Repos[k]...)
	}
	if user != "" {

	}
	return repos
}

func (repo *RepoInfo) ListPackages(hat string) []string {
	pkgs := make([]string, 0)
	return pkgs
}
