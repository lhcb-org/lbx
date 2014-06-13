package lbctx

import (
	"strings"
)

// list of original project names
// please keep this list in the hierarchical order of the dependencies
var ProjectNames = []string{
	"Gaudi", "LHCb", "Lbcom", "Rec", "Boole", "Brunel",
	"Gauss", "Phys", "Analysis", "Hlt", "Alignment", "Moore",
	"Online", "Euler", "Geant4", "DaVinci", "Bender", "Orwell",
	"Panoramix", "LbScripts", "Dirac", "LHCbGrid", "Panoptes",
	"Curie", "Vetra", "VetraTB", "Compat", "VanDerMeer", "Ganga",
	"LHCbDirac", "Integration", "Erasmus", "Feicim",
	"Stripping", "LHCbExternals", "Urania", "VMDirac", "LHCbVMDirac", "Noether", "Tesla",
	"MooreOnline", "BeautyDirac", "Kepler",
}

// FixProjectCase converts the case of the project name to the correct one,
// based on a list of known project names.
// If the project is not known, the name is returned unchanged
func FixProjectCase(project string) string {
	proj := strings.ToLower(project)
	for _, p := range ProjectNames {
		if strings.ToLower(p) == proj {
			return p
		}
	}
	return project
}
