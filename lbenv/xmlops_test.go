package lbenv

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestXMLOps(t *testing.T) {
	const data = `<?xml version="1.0" ?>
<env:config xmlns:env="EnvSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="EnvSchema ./EnvSchema.xsd ">
<env:declare local="false" type="list" variable="varToUnset"/>
<env:unset variable="varToUnset"/>
<env:declare local="true" type="list" variable="myVar"/>
<env:set variable="myVar">setVal:$local</env:set>
<env:append variable="myVar">appVal:appVal2</env:append>
<env:prepend variable="myVar">prepVal:prepVal2</env:prepend>
<env:declare local="false" type="scalar" variable="myScalar"/>
<env:set variable="myScalar">setValscal</env:set>
<env:append variable="myScalar">appValscal</env:append>
<env:prepend variable="myScalar">prepValscal</env:prepend>
<env:declare local="true" type="scalar" variable="myScalar2"/>
<env:include>some_file.xml</env:include>
<env:include hints="some:place">another_file.xml</env:include>
</env:config>
`

	r := bytes.NewBufferString(data)

	actions, err := Decode(r)

	if err != nil && err != io.EOF {
		panic(err)
	}

	expected := []Action{
		&DeclareVar{
			Name:  "varToUnset",
			Local: false,
			Type:  VarList,
		},
		&UnsetVar{
			Name: "varToUnset",
		},
		&DeclareVar{
			Name:  "myVar",
			Local: true,
			Type:  VarList,
		},
		&SetVar{
			Name:  "myVar",
			Value: "setVal:$local",
		},
		&AppendVar{
			Name:  "myVar",
			Value: "appVal:appVal2",
		},
		&PrependVar{
			Name:  "myVar",
			Value: "prepVal:prepVal2",
		},
		&DeclareVar{
			Name:  "myScalar",
			Local: false,
			Type:  VarScalar,
		},
		&SetVar{
			Name:  "myScalar",
			Value: "setValscal",
		},
		&AppendVar{
			Name:  "myScalar",
			Value: "appValscal",
		},
		&PrependVar{
			Name:  "myScalar",
			Value: "prepValscal",
		},
		&DeclareVar{
			Name:  "myScalar2",
			Local: true,
			Type:  VarScalar,
		},
		&Include{
			File: "some_file.xml",
		},
		&Include{
			File:  "another_file.xml",
			Hints: "some:place",
		},
	}
	if len(actions) != len(expected) {
		t.Fatalf("expected %d actions. got=%d", len(expected), len(actions))
	}

	for i := 0; i < len(expected); i++ {
		if !reflect.DeepEqual(actions[i], expected[i]) {
			t.Fatalf("actions[%d]=%v\nexpected=%v", actions[i], expected[i])
		}
	}
}
