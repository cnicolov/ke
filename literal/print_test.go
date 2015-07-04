package literal

import (
	"reflect"
	"testing"

	"kego.io/kerr/assert"
	"kego.io/process/generator"
)

func TestGetName(t *testing.T) {

	imp := generator.NewImports_test(map[string]string{"e.f/g": "d"})

	s := ""
	b := true
	i := 1
	f := 1.2

	assert.Equal(t, "string", GetName(reflect.TypeOf(s), "a.b/c", imp.Add_test))
	assert.Equal(t, "bool", GetName(reflect.TypeOf(b), "a.b/c", imp.Add_test))
	assert.Equal(t, "int", GetName(reflect.TypeOf(i), "a.b/c", imp.Add_test))
	assert.Equal(t, "float64", GetName(reflect.TypeOf(f), "a.b/c", imp.Add_test))

	assert.Equal(t, "*string", GetName(reflect.TypeOf(&s), "a.b/c", imp.Add_test))
	assert.Equal(t, "*bool", GetName(reflect.TypeOf(&b), "a.b/c", imp.Add_test))
	assert.Equal(t, "*int", GetName(reflect.TypeOf(&i), "a.b/c", imp.Add_test))
	assert.Equal(t, "*float64", GetName(reflect.TypeOf(&f), "a.b/c", imp.Add_test))

	impk := generator.NewImports_test(map[string]string{"e.f/g": "d", "kego.io/literal": "h"})
	type MyType struct{}
	m := MyType{}

	assert.Equal(t, "MyType", GetName(reflect.TypeOf(m), "kego.io/literal", imp.Add_test))
	assert.Equal(t, "h.MyType", GetName(reflect.TypeOf(m), "a.b/c", impk.Add_test))
	assert.Equal(t, "*MyType", GetName(reflect.TypeOf(&m), "kego.io/literal", imp.Add_test))
	assert.Equal(t, "*h.MyType", GetName(reflect.TypeOf(&m), "a.b/c", impk.Add_test))

	asp := "kego.io/kerr/assert"
	impa := generator.NewImports_test(map[string]string{"e.f/g": "d", asp: "as"})
	ass := assert.Assertions{} // Just using a random struct from a non system package

	assert.Equal(t, "Assertions", GetName(reflect.TypeOf(ass), asp, imp.Add_test))
	assert.Equal(t, "as.Assertions", GetName(reflect.TypeOf(ass), "a.b/c", impa.Add_test))
	assert.Equal(t, "*Assertions", GetName(reflect.TypeOf(&ass), asp, imp.Add_test))
	assert.Equal(t, "*as.Assertions", GetName(reflect.TypeOf(&ass), "a.b/c", impa.Add_test))

	as := []string{}
	ab := []bool{}
	abp := []*bool{}
	am := []MyType{}
	aass := []assert.Assertions{}
	aassp := []*assert.Assertions{}
	assert.Equal(t, "[]string", GetName(reflect.TypeOf(as), "a.b/c", imp.Add_test))
	assert.Equal(t, "*[]bool", GetName(reflect.TypeOf(&ab), "a.b/c", imp.Add_test))
	assert.Equal(t, "[]*bool", GetName(reflect.TypeOf(abp), "a.b/c", imp.Add_test))
	assert.Equal(t, "*[]*bool", GetName(reflect.TypeOf(&abp), "a.b/c", imp.Add_test))
	assert.Equal(t, "[]h.MyType", GetName(reflect.TypeOf(am), "a.b/c", impk.Add_test))
	assert.Equal(t, "[]as.Assertions", GetName(reflect.TypeOf(aass), "a.b/c", impa.Add_test))
	assert.Equal(t, "*[]as.Assertions", GetName(reflect.TypeOf(&aass), "a.b/c", impa.Add_test))
	assert.Equal(t, "*[]*as.Assertions", GetName(reflect.TypeOf(&aassp), "a.b/c", impa.Add_test))

	ms := map[string]string{}
	mb := map[string]bool{}
	mbp := map[string]*bool{}
	mm := map[string]MyType{}
	mass := map[string]assert.Assertions{}
	massp := map[string]*assert.Assertions{}
	assert.Equal(t, "map[string]string", GetName(reflect.TypeOf(ms), "a.b/c", imp.Add_test))
	assert.Equal(t, "*map[string]bool", GetName(reflect.TypeOf(&mb), "a.b/c", imp.Add_test))
	assert.Equal(t, "map[string]*bool", GetName(reflect.TypeOf(mbp), "a.b/c", imp.Add_test))
	assert.Equal(t, "*map[string]*bool", GetName(reflect.TypeOf(&mbp), "a.b/c", imp.Add_test))
	assert.Equal(t, "map[string]h.MyType", GetName(reflect.TypeOf(mm), "a.b/c", impk.Add_test))
	assert.Equal(t, "map[string]as.Assertions", GetName(reflect.TypeOf(mass), "a.b/c", impa.Add_test))
	assert.Equal(t, "*map[string]as.Assertions", GetName(reflect.TypeOf(&mass), "a.b/c", impa.Add_test))
	assert.Equal(t, "*map[string]*as.Assertions", GetName(reflect.TypeOf(&massp), "a.b/c", impa.Add_test))

	assert.Equal(t, "map[bool]string", GetName(reflect.TypeOf(map[bool]string{}), "a.b/c", imp.Add_test))
	assert.Equal(t, "map[*int]string", GetName(reflect.TypeOf(map[*int]string{}), "a.b/c", imp.Add_test))

	assert.Equal(t, "map[MyType]string", GetName(reflect.TypeOf(map[MyType]string{}), "kego.io/literal", imp.Add_test))
	assert.Equal(t, "map[h.MyType]string", GetName(reflect.TypeOf(map[MyType]string{}), "a.b/c", impk.Add_test))
	assert.Equal(t, "map[*MyType]string", GetName(reflect.TypeOf(map[*MyType]string{}), "kego.io/literal", imp.Add_test))
	assert.Equal(t, "map[*h.MyType]string", GetName(reflect.TypeOf(map[*MyType]string{}), "a.b/c", impk.Add_test))

	assert.Equal(t, "map[Assertions]string", GetName(reflect.TypeOf(map[assert.Assertions]string{}), asp, imp.Add_test))
	assert.Equal(t, "map[as.Assertions]string", GetName(reflect.TypeOf(map[assert.Assertions]string{}), "a.b/c", impa.Add_test))
	assert.Equal(t, "map[*Assertions]string", GetName(reflect.TypeOf(map[*assert.Assertions]string{}), asp, imp.Add_test))
	assert.Equal(t, "map[*as.Assertions]string", GetName(reflect.TypeOf(map[*assert.Assertions]string{}), "a.b/c", impa.Add_test))
}
