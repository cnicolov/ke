package system

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"kego.io/json"
	"kego.io/uerr"
)

func TestRuleTypes(t *testing.T) {

	type ruleStruct struct {
		*Object
	}
	parentType := &Type{
		Object: &Object{Id: "a", Type: NewReference("kego.io/system", "type")},
	}
	ruleType := &Type{
		Object: &Object{Id: "@a", Type: NewReference("kego.io/system", "type")},
	}
	RegisterType("a.b/c:a", parentType)
	RegisterType("a.b/c:@a", ruleType)
	defer UnregisterType("a.b/c:a")
	defer UnregisterType("a.b/c:@a")

	r := &ruleStruct{
		&Object{Type: NewReference("a.b/c", "@a")},
	}
	rt, pt, err := ruleTypes(r, "", map[string]string{})
	assert.NoError(t, err)
	assert.Equal(t, "a", pt.Id)
	assert.Equal(t, "@a", rt.Id)

	r1 := ruleStruct{}
	rt, pt, err = ruleTypes(r1, "", map[string]string{})
	// A non pointer rule will cause ruleTypeReference to return an error
	uerr.Assert(t, err, "BNEKIFYDDL")

	r = &ruleStruct{
		&Object{Type: NewReference("a.b/c", "unregistered")},
	}
	rt, pt, err = ruleTypes(r, "", map[string]string{})
	// An unregistered type will cause ruleReference.GetType to return an error
	uerr.Assert(t, err, "PFGWISOHRR")

	r = &ruleStruct{
		&Object{Type: NewReference("a.b/c", "a")},
	}
	rt, pt, err = ruleTypes(r, "", map[string]string{})
	// A rule with a non rule type will cause ruleReference.RuleToParentType to error
	uerr.Assert(t, err, "NXRCPQMUIE")

	RegisterType("a.b/c:@b", ruleType)
	r = &ruleStruct{
		&Object{Type: NewReference("a.b/c", "@b")},
	}
	rt, pt, err = ruleTypes(r, "", map[string]string{})
	// An rule type with an unregistered parent type typeReference.GetType to return an error
	uerr.Assert(t, err, "KYCTDXKFYR")

}

func TestRuleTypeReference(t *testing.T) {

	type ruleStruct struct {
		*Object
	}
	rs := &ruleStruct{
		&Object{Type: NewReference("a.b/c", "@a")},
	}
	r, err := ruleTypeReference(rs, "", map[string]string{})
	assert.NoError(t, err)
	assert.Equal(t, "a.b/c:@a", r.Value)

	ri := map[string]interface{}{}
	r, err = ruleTypeReference(ri, "", map[string]string{})
	uerr.Assert(t, err, "OLHOVKXEXN")

	ri = map[string]interface{}{
		"type": 1, //not a string
	}
	r, err = ruleTypeReference(ri, "", map[string]string{})
	uerr.Assert(t, err, "IILEXGQDXL")

	ri = map[string]interface{}{
		"type": "a:b", // package will not be registered so UnmarshalJSON will error
	}
	r, err = ruleTypeReference(ri, "", map[string]string{})
	uerr.Assert(t, err, "QBTHPRVBWN")

	ri = map[string]interface{}{
		"type": "a.b/c:@a",
	}
	r, err = ruleTypeReference(rs, "", map[string]string{})
	assert.NoError(t, err)
	assert.Equal(t, "a.b/c:@a", r.Value)

	rsp := ruleStruct{}
	r, err = ruleTypeReference(rsp, "", map[string]string{})
	// rsp is not a pointer so ruleFieldByReflection will error
	uerr.Assert(t, err, "QJQAIGPYXC")

	type structWithoutType struct{}
	rwt := &structWithoutType{}
	r, err = ruleTypeReference(rwt, "", map[string]string{})
	uerr.Assert(t, err, "NXYRAJITEV")

	type structWithIntType struct {
		Type int
	}
	rwi := &structWithIntType{
		Type: 1,
	}
	r, err = ruleTypeReference(rwi, "", map[string]string{})
	uerr.Assert(t, err, "FHUPSRTRFE")
}

func TestRuleHolderItemsRule(t *testing.T) {
	type parentStruct struct {
		*Object
	}
	type ruleStruct struct {
		*Object
	}
	parentType := &Type{
		Object: &Object{Context: &Context{Package: "a.b/c"}, Id: "a", Type: NewReference("kego.io/system", "type")},
	}
	ruleType := &Type{
		Object: &Object{Context: &Context{Package: "a.b/c"}, Id: "@a", Type: NewReference("kego.io/system", "type")},
	}
	json.RegisterType("a.b/c:a", reflect.TypeOf(&parentStruct{}))
	json.RegisterType("a.b/c:@a", reflect.TypeOf(&ruleStruct{}))
	RegisterType("a.b/c:a", parentType)
	RegisterType("a.b/c:@a", ruleType)
	defer json.UnregisterType("a.b/c:a")
	defer json.UnregisterType("a.b/c:@a")
	defer UnregisterType("a.b/c:a")
	defer UnregisterType("a.b/c:@a")

	rh := &RuleHolder{
		rule:       &ruleStruct{},
		ruleType:   ruleType,
		parentType: parentType,
		path:       "d.e/f",
		imports:    map[string]string{},
	}
	_, err := rh.ItemsRule()
	uerr.Assert(t, err, "VPAGXSTQHM")

	parentType.Native = NewString("array")
	rh.rule = "a"
	_, err = rh.ItemsRule()
	// rh.rule must be a pointer or ruleFieldByReflection will error
	uerr.Assert(t, err, "LIDXIQYGJD")

	rh.rule = &struct{}{}
	_, err = rh.ItemsRule()
	// rh.rule needs an Items field
	uerr.Assert(t, err, "VYTHGJTSNJ")

	rh.rule = &struct{ Items int }{Items: 1}
	_, err = rh.ItemsRule()
	// Items must be a rule or NewRuleHolder will error
	uerr.Assert(t, err, "FGYMQPNBQJ")

}
