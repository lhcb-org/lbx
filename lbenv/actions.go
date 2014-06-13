package lbenv

type VarType int

const (
	VarList   VarType = 0
	VarScalar VarType = 1
)

func (vt VarType) String() string {
	switch vt {
	case VarList:
		return "list"
	case VarScalar:
		return "scalar"
	}
	panic("unreachable")
}

type Action interface {
	Run(env *Environment) error
}

type DeclareVar struct {
	Name  string
	Local bool
	Type  VarType
}

func (v *DeclareVar) Run(env *Environment) error {
	var err error
	return err
}

type SetVar struct {
	Name  string
	Value string
}

func (v *SetVar) Run(env *Environment) error {
	var err error
	return err
}

type UnsetVar struct {
	Name string
}

func (v *UnsetVar) Run(env *Environment) error {
	var err error
	return err
}

type AppendVar struct {
	Name  string
	Value string
}

func (v *AppendVar) Run(env *Environment) error {
	var err error
	return err
}

type PrependVar struct {
	Name  string
	Value string
}

func (v *PrependVar) Run(env *Environment) error {
	var err error
	return err
}

type RemoveVar struct {
	Name  string
	Value string
}

func (v *RemoveVar) Run(env *Environment) error {
	var err error
	return err
}

type RemoveRegexp struct {
	Name  string
	Value string
}

func (v *RemoveRegexp) Run(env *Environment) error {
	var err error
	return err
}

type Include struct {
	File   string `xml:"include"`
	Caller string
	Hints  string `xml:"hints,attr"`
}

func (v *Include) Run(env *Environment) error {
	var err error
	return err
}
