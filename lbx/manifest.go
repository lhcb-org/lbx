package lbx

import (
	"encoding/xml"
	"io"
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
