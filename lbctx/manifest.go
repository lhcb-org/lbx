package lbctx

import (
	"encoding/xml"
	"io"
	"os"
	"path/filepath"
)

type Manifest struct {
	XMLName xml.Name `xml:"manifest"`
	Project struct {
		Name    string `xml:"name,attr"`
		Version string `xml:"version,attr"`
	} `xml:"project"`

	HepTools struct {
		Version   string `xml:"version"`
		BinTag    string `xml:"binary_tag"`
		LcgSystem string `xml:"lcg_system"`
	} `xml:"heptools"`

	UsedProjects []struct {
		Name    string `xml:"name,attr"`
		Version string `xml:"version,attr"`
	} `xml:"used_projects>project"`

	UsedDataPkgs []struct {
		Name    string `xml:"name,attr"`
		Version string `xml:"version,attr"`
	} `xml:"used_data_pkgs>package"`
}

// ParseManifest parses a CMake-generated manifest.xml file
func ParseManifest(r io.Reader) (*Manifest, error) {
	var m Manifest
	err := xml.NewDecoder(r).Decode(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// EnvXMLPath returns the list of directories to be added to the XML search
// path for a given project
func (ctx *Context) EnvXMLPath(project, version, platform string) ([]string, error) {
	projdir, err := ctx.FindProject(project, version, platform)
	if err != nil {
		return nil, err
	}

	// xml paths to return
	paths := []string{projdir}
	// unique list of paths to return
	pset := map[string]struct{}{
		projdir: struct{}{},
	}

	manifests := []string{filepath.Join(projdir, "manifest.xml")}
	for len(manifests) > 0 {
		fname := manifests[len(manifests)-1]
		manifests = manifests[:len(manifests)-1]
		if _, err := os.Stat(fname); err != nil {
			continue
		}
		f, err := os.Open(fname)
		if err != nil {
			return nil, err
		}
		m, err := ParseManifest(f)
		if err != nil {
			return nil, err
		}
		// add the data package directories
		for _, dpkg := range m.UsedDataPkgs {
			dir, err := ctx.FindDataPackage(dpkg.Name, dpkg.Version)
			if err != nil {
				return nil, err
			}
			if _, dup := pset[dir]; !dup {
				paths = append(paths, dir)
				pset[dir] = struct{}{}
			}
		}

		// add the project directories
		for _, proj := range m.UsedProjects {
			dir, err := ctx.FindProject(proj.Name, proj.Version, platform)
			if err != nil {
				return nil, err
			}
			if _, dup := pset[dir]; !dup {
				paths = append(paths, dir)
				pset[dir] = struct{}{}
			}
			// add project's manifest to the list of manifests to parse
			manifests = append(manifests, filepath.Join(dir, "manifest.xml"))
		}
	}

	return paths, nil
}
