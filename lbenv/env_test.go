package lbenv

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestValues(t *testing.T) {
	env := New()
	if env.Has("MY_PATH") {
		t.Fatalf("expected no 'MY_PATH' env.var")
	}
	err := env.Append("MY_PATH", "newValue")
	if err != nil {
		t.Fatalf("problem appending: %v", err)
	}

	if !env.Has("MY_PATH") {
		t.Fatalf("expected a 'MY_PATH' env.var")
	}

	err = env.Append("MY_PATH", "newValue:secondVal:valval")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	v := splitpath(env.Get("MY_PATH").Value)

	if v[len(v)-1] != "valval" {
		t.Fatalf("expected %q. got=%q", "valval", v[len(v)-1])
	}

	if !in_str_slice("newValue", v) {
		t.Fatalf("expected %q to be in slice %v", "newValue", v)
	}

	err = env.Remove("MY_PATH", "newValue")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	v = splitpath(env.Get("MY_PATH").Value)
	if in_str_slice("newValue", v) {
		t.Fatalf("expected %q NOT to be in slice %v", "newValue", v)
	}

	err = env.Prepend("MY_PATH", "newValue")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	v = splitpath(env.Get("MY_PATH").Value)
	if !in_str_slice("newValue", v) {
		t.Fatalf("expected %q to be in slice %v", "newValue", v)
	}

	err = env.Set("MY_PATH", "hi:hello")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	v = splitpath(env.Get("MY_PATH").Value)
	if !reflect.DeepEqual(v, []string{"hi", "hello"}) {
		t.Fatalf("expected env=%v. got=%v", "newValue", []string{"hi", "hello"}, v)
	}

	err = env.Unset("MY_PATH")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if env.Has("MY_PATH") {
		t.Fatalf("expected no 'MY_PATH' env.var")
	}
}

func TestHiddingDotVar(t *testing.T) {
	env := New()
	err := env.Set(".", "some/dir")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if !env.Has(".") {
		t.Fatalf("expected env with '.'")
	}

	keys := env.Keys()
	if len(keys) != 0 {
		t.Fatalf("expected env with HIDDEN '.'")
	}

	err = env.Set("MY_DIR", "${.}")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	v := env.Get("MY_DIR")
	if v.Value != "some/dir" {
		t.Fatalf("expected MY_DIR=%q. got=%q.\n%v", "some/dir", v.Value, str_actions(env.stack))
	}

	keys = env.Keys()
	exp := []string{"MY_DIR"}
	if !reflect.DeepEqual(keys, exp) {
		t.Fatalf("expected env with keys=%v. got=%v", exp, keys)
	}
}

func TestSaveLoadXML(t *testing.T) {
	env := New()

	err := env.Unset("MY_PATH")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	err = env.Set("MY_PATH", "set:toDelete")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	err = env.Append("MY_PATH", "appended:toDelete")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	err = env.Prepend("MY_PATH", "prepended:toDelete")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	err = env.Remove("MY_PATH", "toDelete")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	f, err := os.Create("testdata/test-output-file.xml")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer f.Close()
	defer os.RemoveAll("testdata/test-output-file.xml")

	err = env.SaveXML(f)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = f.Sync()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	env = New()
	if env.Has("MY_PATH") {
		t.Fatalf("expected NO 'MY_PATH' env.var.")
	}

	err = env.LoadXML(f)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if !env.Has("MY_PATH") {
		t.Fatalf("expected a MY_PATH env.var.")
	}
	v := splitpath(env.Get("MY_PATH").Value)

	exp := []string{"prepended", "set", "appended"}
	if !reflect.DeepEqual(v, exp) {
		t.Fatalf("expected MY_PATH=%v. got=%v", exp, v)
	}
}

func TestGenScript(t *testing.T) {
	env := New()
	err := env.Append("sysVar", "newValue:lala")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	f, err := os.Create("testdata/test-setup.sh")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer f.Close()
	defer os.RemoveAll(f.Name())

	err = env.GenScript(ShType, f)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	err = f.Sync()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	bout, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	out := string(bout)
	exp := `#!/bin/sh
export sysVar="newValue:lala"
## EOF
`
	if out != exp {
		t.Fatalf("error:\nexp=%v\ngot=%v\n", exp, out)
	}
}

func TestVariables(t *testing.T) {
	env := New()
	err := env.Append("MY_PATH", "newValue")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	v := env.Get("MY_PATH")
	if v.Local != false {
		t.Fatalf("expected MY_PATH to be non-local")
	}
	if v.Type != VarList {
		t.Fatalf("expected MY_PATH to be a list")
	}

	err = env.Declare("loc", VarList, true)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	v = env.Get("loc")
	if v.Local != true {
		t.Fatalf("expected loc to be local")
	}
	if v.Type != VarList {
		t.Fatalf("expected loc to be a list")
	}

	err = env.Declare("myVar2", VarScalar, false)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	v = env.Get("myVar2")
	if v.Local != false {
		t.Fatalf("expected myVar2 to be non-local")
	}
	if v.Type != VarScalar {
		t.Fatalf("expected myVar2 to be a scalar")
	}

	err = env.Declare("loc2", VarScalar, true)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	v = env.Get("loc2")
	if v.Local != true {
		t.Fatalf("expected loc2 to be local")
	}
	if v.Type != VarScalar {
		t.Fatalf("expected loc2 to be a scalar")
	}

	{
		name := "MY_PATH"
		err = env.Declare(name, VarList, false)
		if err != nil {
			t.Fatalf("error: %v", err)
		}

		for _, table := range []struct {
			name string
			typ  VarType
			loc  bool
		}{
			{
				name: name,
				typ:  VarList,
				loc:  true,
			},
			{
				name: name,
				typ:  VarScalar,
				loc:  true,
			},
			{
				name: name,
				typ:  VarScalar,
				loc:  false,
			},
		} {
			err = env.Declare(table.name, table.typ, table.loc)
			if err == nil {
				t.Fatalf("expected a redeclaration error (table=%v)", table)
			}
		}
	}

	{
		name := "loc"
		err = env.Declare(name, VarList, true)
		if err != nil {
			t.Fatalf("error: %v", err)
		}

		for _, table := range []struct {
			name string
			typ  VarType
			loc  bool
		}{
			{
				name: name,
				typ:  VarList,
				loc:  false,
			},
			{
				name: name,
				typ:  VarScalar,
				loc:  true,
			},
			{
				name: name,
				typ:  VarScalar,
				loc:  false,
			},
		} {
			err = env.Declare(table.name, table.typ, table.loc)
			if err == nil {
				t.Fatalf("expected a redeclaration error (table=%v)", table)
			}
		}
	}

	{
		name := "myVar2"
		err = env.Declare(name, VarScalar, false)
		if err != nil {
			t.Fatalf("error: %v", err)
		}

		for _, table := range []struct {
			name string
			typ  VarType
			loc  bool
		}{
			{
				name: name,
				typ:  VarList,
				loc:  false,
			},
			{
				name: name,
				typ:  VarList,
				loc:  true,
			},
			{
				name: name,
				typ:  VarScalar,
				loc:  true,
			},
		} {
			err = env.Declare(table.name, table.typ, table.loc)
			if err == nil {
				t.Fatalf("expected a redeclaration error (table=%v)", table)
			}
		}
	}

	{
		name := "loc2"
		err = env.Declare(name, VarScalar, true)
		if err != nil {
			t.Fatalf("error: %v", err)
		}

		for _, table := range []struct {
			name string
			typ  VarType
			loc  bool
		}{
			{
				name: name,
				typ:  VarList,
				loc:  false,
			},
			{
				name: name,
				typ:  VarList,
				loc:  true,
			},
			{
				name: name,
				typ:  VarScalar,
				loc:  false,
			},
		} {
			err = env.Declare(table.name, table.typ, table.loc)
			if err == nil {
				t.Fatalf("expected a redeclaration error (table=%v)", table)
			}
		}
	}

}

func TestDelete(t *testing.T) {
	env := New()

	// --- test list ---

	err := env.Append("MY_PATH", "myVal:anotherVal:lastVal")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	err = env.Remove("MY_PATH", "anotherVal")
	if err != nil {
		t.Fatalf("error: %v")
	}

	v := env.Get("MY_PATH").Value
	if strings.Contains(v, "anotherVal") {
		t.Fatalf("expected no 'anotherVal' in MY_PATH. got=%q", v)
	}
	if !strings.Contains(v, "myVal") {
		t.Fatalf("expected 'myVal' in MY_PATH. got=%q", v)
	}
	if !strings.Contains(v, "lastVal") {
		t.Fatalf("expected 'lastVal' in MY_PATH. got=%q", v)
	}

	err = env.Set("MY_PATH", "myVal:anotherVal:lastVal:else")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	err = env.Remove("MY_PATH", "^anotherVal$")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	v = env.Get("MY_PATH").Value
	if !strings.Contains(v, "anotherVal") {
		t.Fatalf("expected 'anotherVal' in MY_PATH. got=%q", v)
	}

	err = env.RemoveRegexp("MY_PATH", "^anotherVal$")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	v = env.Get("MY_PATH").Value
	if exp := "myVal:lastVal:else"; v != exp {
		t.Fatalf("expected MY_PATH=%q. got=%q", exp, v)
	}

	err = env.RemoveRegexp("MY_PATH", "Val")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	v = env.Get("MY_PATH").Value
	if exp := "else"; v != exp {
		t.Fatalf("expected MY_PATH=%q. got=%q", exp, v)
	}

	// --- test scalar ---

	err = env.Declare("myLoc", VarScalar, false)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	err = env.Append("myLoc", "myVal:anotherVal:lastVal")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	err = env.RemoveRegexp("myLoc", "Val:")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	v = env.Get("myLoc").Value
	if exp := "myanotherlastVal"; exp != v {
		t.Fatalf("expected myLoc=%q. got=%q", exp, v)
	}
}

func TestSystemEnvironment(t *testing.T) {
	env := New()

	// --- test list ---
	os.Setenv("MY_PATH", "$myVal")
	err := env.Set("ABC", "anyValue")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	err = env.Declare("MY_PATH", VarList, false)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	err = env.Append("MY_PATH", "$ABC")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	v := env.Get("MY_PATH").Value
	if exp := "$myVal:anyValue"; exp != v {
		t.Fatalf("expected MY_PATH=%q. got=%q", exp, v)
	}

	// --- test scalar ---
	os.Setenv("myScal", "$myVal")

	err = env.Declare("myScal", VarScalar, false)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	err = env.Append("myScal", "$ABC")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	v = env.Get("myScal").Value
	if exp := "$myValanyValue"; v != exp {
		t.Fatalf("expected myScal=%q. got=%q", exp, v)
	}
}

func TestDependencies(t *testing.T) {
	env := New()

	err := env.Declare("myVar", VarList, false)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	err = env.Declare("loc", VarList, true)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = env.Append("loc", "locVal")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = env.Append("loc", "locVal2")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	err = env.Declare("scal", VarScalar, false)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = env.Append("scal", "scalVal")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = env.Append("scal", "scalVal2")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	err = env.Declare("scal2", VarScalar, true)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = env.Append("scal2", "locScal")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = env.Append("scal2", "locScal2")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	for _, table := range []struct {
		name  string
		setv  string
		value string
	}{
		{
			name:  "myVar",
			setv:  "newValue:$loc:endValue",
			value: "newValue:locVal:locVal2:endValue",
		},
		{
			name:  "myVar",
			setv:  "newValue:$scal:endValue",
			value: "newValue:scalValscalVal2:endValue",
		},
		{
			name:  "myVar",
			setv:  "new${scal}Value:endValue",
			value: "newscalValscalVal2Value:endValue",
		},
		{
			name:  "myVar",
			setv:  "bla:$myVar:Value",
			value: "bla:newscalValscalVal2Value:endValue:Value",
		},
		{
			name:  "scal",
			setv:  "new${scal2}Value",
			value: "newlocScallocScal2Value",
		},
		{
			name:  "scal",
			setv:  "new${loc}Value",
			value: "newlocVal:locVal2Value",
		},
		{
			name:  "scal2",
			setv:  "new${scal2}Value",
			value: "newlocScallocScal2Value",
		},
	} {
		err = env.Set(table.name, table.setv)
		if err != nil {
			t.Fatalf("error: %v", err)
		}

		v := env.Get(table.name).Value
		if v != table.value {
			t.Fatalf("expected %s=%q. got=%q (table=%v)", table.name, table.value, v, table)
		}

	}
}

func TestInclude(t *testing.T) {
	err := os.MkdirAll("testdata/test-includes", 0755)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer os.RemoveAll("testdata/test-includes")

	for _, table := range []struct {
		name string
		cont string
	}{
		{
			name: "testdata/test-includes/first.xml",
			cont: `<?xml version="1.0" ?>
<env:config xmlns:env="EnvSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="EnvSchema ./EnvSchema.xsd ">
<env:set variable="main">first</env:set>
<env:append variable="test_path">data1</env:append>
<env:include>first_inc.xml</env:include>
</env:config>`,
		},
		{
			name: "testdata/test-includes/second.xml",
			cont: `<?xml version="1.0" ?>
<env:config xmlns:env="EnvSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="EnvSchema ./EnvSchema.xsd ">
<env:set variable="main">second</env:set>
<env:include>second_inc.xml</env:include>
<env:append variable="test_path">data1</env:append>
</env:config>`,
		},
		{
			name: "testdata/test-includes/third.xml",
			cont: `<?xml version="1.0" ?>
<env:config xmlns:env="EnvSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="EnvSchema ./EnvSchema.xsd ">
<env:set variable="main">third</env:set>
<env:append variable="test_path">data1</env:append>
<env:include>subdir/first_inc.xml</env:include>
</env:config>`,
		},
		{
			name: "testdata/test-includes/fourth.xml",
			cont: `<?xml version="1.0" ?>
<env:config xmlns:env="EnvSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="EnvSchema ./EnvSchema.xsd ">
<env:set variable="main">fourth</env:set>
<env:include hints="subdir2">fourth_inc.xml</env:include>
</env:config>`,
		},
		{
			name: "testdata/test-includes/recursion.xml",
			cont: `<?xml version="1.0" ?>
<env:config xmlns:env="EnvSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="EnvSchema ./EnvSchema.xsd ">
<env:set variable="main">recursion</env:set>
<env:include>recursion.xml</env:include>
</env:config>`,
		},
		{
			name: "testdata/test-includes/first_inc.xml",
			cont: `<?xml version="1.0" ?>
<env:config xmlns:env="EnvSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="EnvSchema ./EnvSchema.xsd ">
<env:append variable="test_path">data2</env:append>
<env:append variable="derived">another_${main}</env:append>
</env:config>`,
		},
		{
			name: "testdata/test-includes/subdir/second_inc.xml",
			cont: `<?xml version="1.0" ?>
<env:config xmlns:env="EnvSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="EnvSchema ./EnvSchema.xsd ">
<env:append variable="test_path">data0</env:append>
<env:set variable="map">this_is_second_inc</env:set>
</env:config>`,
		},
		{
			name: "testdata/test-includes/subdir/first_inc.xml",
			cont: `<?xml version="1.0" ?>
<env:config xmlns:env="EnvSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="EnvSchema ./EnvSchema.xsd ">
<env:append variable="derived">second_${main}</env:append>
</env:config>`,
		},
		{
			name: "testdata/test-includes/subdir/fourth_inc.xml",
			cont: `<?xml version="1.0" ?>
<env:config xmlns:env="EnvSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="EnvSchema ./EnvSchema.xsd ">
<env:append variable="included">from subdir</env:append>
</env:config>`,
		},
		{
			name: "testdata/test-includes/subdir2/fourth_inc.xml",
			cont: `<?xml version="1.0" ?>
<env:config xmlns:env="EnvSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="EnvSchema ./EnvSchema.xsd ">
<env:append variable="included">from subdir2</env:append>
</env:config>`,
		},
	} {
		dir := filepath.Dir(table.name)
		if dir != "" {
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				t.Fatalf("error: %v", err)
			}
		}
		err = ioutil.WriteFile(table.name, []byte(table.cont), 0644)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}

	// set the basic search path to the minimal default
	os.Setenv("ENVXMLPATH", "")
	var saved_path []string
	mini_path := []string{"."}
	saved_path, g_envxmlpath = g_envxmlpath, mini_path
	defer func() {
		g_envxmlpath = saved_path
	}()

	type pair struct {
		name string
		val  string
	}
	for _, table := range []struct {
		search []string
		fname  string
		err    bool
		tests  []pair
		pre    func()
		post   func()
	}{
		{
			search: []string{},
			fname:  "testdata/test-includes/first.xml",
			tests: []pair{
				{
					name: "main",
					val:  "first",
				},
				{
					name: "test_path",
					val:  "data1:data2",
				},
				{
					name: "derived",
					val:  "another_first",
				},
			},
		},
		{
			search: []string{"testdata/test-includes"},
			fname:  "testdata/test-includes/first.xml",
			tests: []pair{
				{
					name: "main",
					val:  "first",
				},
				{
					name: "test_path",
					val:  "data1:data2",
				},
				{
					name: "derived",
					val:  "another_first",
				},
			},
		},
		{
			search: []string{"testdata/test-includes"},
			fname:  "first.xml",
			tests: []pair{
				{
					name: "main",
					val:  "first",
				},
				{
					name: "test_path",
					val:  "data1:data2",
				},
				{
					name: "derived",
					val:  "another_first",
				},
			},
		},
		{
			search: []string{},
			fname:  "testdata/test-includes/second.xml",
			err:    true,
		},
		{
			search: []string{"testdata/test-includes", "testdata/test-includes/subdir"},
			fname:  "testdata/test-includes/second.xml",
			tests: []pair{
				{
					name: "main",
					val:  "second",
				},
				{
					name: "test_path",
					val:  "data0:data1",
				},
				{
					name: "map",
					val:  "this_is_second_inc",
				},
			},
		},
		{
			search: []string{"testdata/test-includes", "testdata/test-includes/subdir"},
			fname:  "testdata/test-includes/first.xml",
			tests: []pair{
				{
					name: "main",
					val:  "first",
				},
				{
					name: "test_path",
					val:  "data1:data2",
				},
				{
					name: "derived",
					val:  "another_first",
				},
			},
		},
		{
			search: []string{"testdata/test-includes/subdir", "testdata/test-includes"},
			fname:  "testdata/test-includes/first.xml",
			tests: []pair{
				{
					name: "main",
					val:  "first",
				},
				{
					name: "test_path",
					val:  "data1:data2",
				},
				{
					name: "derived",
					val:  "another_first",
				},
			},
		},
		{
			search: []string{"testdata/test-includes/subdir", "testdata/test-includes"},
			fname:  "first.xml",
			tests: []pair{
				{
					name: "main",
					val:  "first",
				},
				{
					name: "test_path",
					val:  "data1:data2",
				},
				{
					name: "derived",
					val:  "another_first",
				},
			},
		},
		{
			search: []string{},
			fname:  "testdata/test-includes/second.xml",
			tests: []pair{
				{
					name: "main",
					val:  "second",
				},
				{
					name: "test_path",
					val:  "data0:data1",
				},
				{
					name: "map",
					val:  "this_is_second_inc",
				},
			},
			pre: func() {
				g_envxmlpath = []string{".", "testdata/test-includes", "testdata/test-includes/subdir"}
			},
			post: func() {
				g_envxmlpath = []string{"."}
			},
		},
		{
			search: []string{},
			fname:  "testdata/test-includes/third.xml",
			tests: []pair{
				{
					name: "main",
					val:  "third",
				},
				{
					name: "test_path",
					val:  "data1",
				},
				{
					name: "derived",
					val:  "second_third",
				},
			},
		},
		{
			search: []string{"testdata/test-includes/subdir"},
			fname:  "testdata/test-includes/fourth.xml",
			tests: []pair{
				{
					name: "main",
					val:  "fourth",
				},
				{
					name: "included",
					val:  "from subdir",
				},
			},
		},
		{
			search: []string{},
			fname:  "testdata/test-includes/fourth.xml",
			tests: []pair{
				{
					name: "main",
					val:  "fourth",
				},
				{
					name: "included",
					val:  "from subdir2",
				},
			},
		},
		{
			search: []string{},
			fname:  "testdata/test-includes/recursion.xml",
		},
	} {
		if table.pre != nil {
			table.pre()
		}
		env := New()
		env.SearchPath = table.search

		err = env.LoadXMLByName(table.fname)
		if (err != nil) != table.err {
			t.Fatalf("error loading file [%s]: %v", table.fname, err)
		}
		if table.err {
			continue
		}
		for _, tt := range table.tests {
			v := env.Get(tt.name).Value
			if v != tt.val {
				t.Fatalf("expected %s=%q. got=%q (table=%v)", tt.name, tt.val, v, table)
			}
		}

		if table.post != nil {
			table.post()
		}
	}
}

func TestFileDir(t *testing.T) {

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	tmpdir := filepath.Join(cwd, "testdata/test-file-dir")

	err = os.MkdirAll(tmpdir, 0755)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer os.RemoveAll(tmpdir)

	for _, table := range []struct {
		name string
		cont string
	}{
		{
			name: "testdata/test-file-dir/env.xml",
			cont: `<?xml version="1.0" ?>
<env:config xmlns:env="EnvSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="EnvSchema ./EnvSchema.xsd ">
<env:set variable="mydirs">${.}</env:set>
<env:set variable="myparent">${.}/..</env:set>
</env:config>`,
		},
	} {
		dir := filepath.Dir(table.name)
		if dir != "" {
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				t.Fatalf("error: %v", err)
			}
		}
		err = ioutil.WriteFile(table.name, []byte(table.cont), 0644)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}

	type pair struct {
		name string
		val  string
	}

	for _, table := range []struct {
		dir   string
		name  string
		tests []pair
	}{
		// {
		// 	dir:  ".",
		// 	name: filepath.Join(tmpdir, "env.xml"),
		// 	tests: []pair{
		// 		{
		// 			name: "mydirs",
		// 			val:  tmpdir,
		// 		},
		// 		{
		// 			name: "myparent",
		// 			val:  filepath.Dir(tmpdir),
		// 		},
		// 	},
		// },
		{
			dir:  tmpdir,
			name: "env.xml",
			tests: []pair{
				{
					name: "mydirs",
					val:  tmpdir,
				},
				{
					name: "myparent",
					val:  filepath.Dir(tmpdir),
				},
			},
		},
	} {
		err := os.Chdir(table.dir)
		if err != nil {
			t.Fatalf("error chdir to [%s]: %v", table.dir, err)
		}
		env := New()

		err = env.LoadXMLByName(table.name)
		if err != nil {
			t.Fatalf("error opening file [%s]: %v", table.name, err)

		}

		for _, test := range table.tests {
			v := env.Get(test.name).Value
			if v != test.val {
				t.Fatalf("expected %s=%q. got=%q (table=%v)", test.name, test.val, v, table)
			}
		}
	}

}
