package system

import (
	"reflect"

	"kego.io/json"
	"kego.io/kerr"
)

type Rule interface {
	GetRuleBase() *RuleBase
}

// Enforcer is a rule with properties that need to be enforced against data.
type Enforcer interface {
	Enforce(data interface{}, path string, imports map[string]string) (bool, string, error)
}

func (b *RuleBase) GetRuleBase() *RuleBase {
	if b == nil {
		return &RuleBase{}
	}
	return b
}

func init() {
	type dummyRule struct {
		*Base
		*RuleBase
		Default interface{}
	}
	json.RegisterInterface(reflect.TypeOf((*Rule)(nil)).Elem(), reflect.TypeOf(&dummyRule{}))
}

type RuleHolder struct {
	Rule       Rule
	RuleType   *Type
	ParentType *Type
	Path       string
	Imports    map[string]string
}

func NewMinimalRuleHolder(t *Type, path string, imports map[string]string) *RuleHolder {
	return &RuleHolder{Rule: nil, RuleType: nil, ParentType: t, Path: path, Imports: imports}
}
func NewRuleHolder(r Rule, path string, imports map[string]string) (*RuleHolder, error) {
	rt, pt, err := ruleTypes(r, path, imports)
	if err != nil {
		return nil, kerr.New("VRCWUGOTMA", err, "NewRuleHolder", "ruleTypes")
	}
	return &RuleHolder{Rule: r, RuleType: rt, ParentType: pt, Path: path, Imports: imports}, nil
}

func ruleTypes(r Rule, path string, imports map[string]string) (ruleType *Type, parentType *Type, err error) {
	ruleReference, err := ruleTypeReference(r, path, imports)
	if err != nil {
		return nil, nil, kerr.New("BNEKIFYDDL", err, "ruleTypes", "ruleTypeReference")
	}
	rt, ok := ruleReference.GetType()
	if !ok {
		return nil, nil, kerr.New("PFGWISOHRR", nil, "ruleTypes", "ruleReference.GetType: type %v not found", ruleReference.Value())
	}
	typeReference, err := ruleReference.RuleToParentType()
	if err != nil {
		return nil, nil, kerr.New("NXRCPQMUIE", err, "ruleTypes", "ruleReference.RuleToParentType")
	}
	pt, ok := typeReference.GetType()
	if !ok {
		return nil, nil, kerr.New("KYCTDXKFYR", nil, "ruleTypes", "typeReference.GetType: type %v not found", typeReference.Value())
	}
	return rt, pt, nil
}

func ruleTypeReference(r Rule, path string, imports map[string]string) (*Reference, error) {
	ob, ok := r.(Object)
	if !ok {
		return nil, kerr.New("VKFNPJDNVB", nil, "system.ruleTypeReference", "r does not implement Object")
	}
	return &ob.GetBase().Type, nil

}

// ItemsRule returns Items rule for a collection Rule.
func (r *RuleHolder) ItemsRule() (*RuleHolder, error) {
	if !r.ParentType.IsNativeCollection() {
		return nil, kerr.New("VPAGXSTQHM", nil, "RuleHolder.ItemsRule", "parentType %s is not a collection", r.ParentType.Id)
	}
	items, _, ok, err := ruleFieldByReflection(r.Rule, "Items")
	if err != nil {
		return nil, kerr.New("LIDXIQYGJD", err, "RuleHolder.ItemsRule", "ruleFieldByReflection")
	}
	if !ok {
		return nil, kerr.New("VYTHGJTSNJ", nil, "RuleHolder.ItemsRule", "ruleFieldByReflection could not find Items field")
	}
	rule, ok := items.(Rule)
	if !ok {
		return nil, kerr.New("DIFVRMVWMC", nil, "RuleHolder.ItemsRule", "items is not a rule")
	}
	rh, err := NewRuleHolder(rule, r.Path, r.Imports)
	if err != nil {
		return nil, kerr.New("FGYMQPNBQJ", err, "RuleHolder.ItemsRule", "NewRuleHolder")
	}
	return rh, nil
}

func ruleFieldByReflection(object interface{}, name string) (value interface{}, pointer interface{}, ok bool, err error) {
	v := reflect.ValueOf(object)
	if v.Kind() != reflect.Ptr {
		return nil, nil, false, kerr.New("QOYMWPXWUO", nil, "ruleFieldByReflection", "val.Kind (%s) is not a Ptr: %v", name, v.Kind())
	}
	if v.Elem().Kind() != reflect.Struct {
		return nil, nil, false, kerr.New("IGOUOBGXAN", nil, "ruleFieldByReflection", "val.Elem().Kind (%s) is not a Struct: %v", name, v.Elem().Kind())
	}
	value, pointer, _, found, zero, err := GetObjectField(v.Elem(), name)

	// zero => !ok
	return value, pointer, found && !zero, err
}

func GetMapMember(v reflect.Value, name string) (object interface{}, pointer interface{}, value reflect.Value, found bool, zero bool, err error) {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	member := v.MapIndex(reflect.ValueOf(name))
	return returnValue(member)
}
func GetArrayMember(v reflect.Value, index int) (object interface{}, pointer interface{}, value reflect.Value, found bool, zero bool, err error) {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	member := v.Index(index)
	return returnValue(member)
}
func GetObjectField(v reflect.Value, name string) (object interface{}, pointer interface{}, value reflect.Value, found bool, zero bool, err error) {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	field := v.FieldByName(name)
	return returnValue(field)
}
func returnValue(v reflect.Value) (object interface{}, pointer interface{}, value reflect.Value, found bool, zero bool, err error) {
	value = v
	empty := reflect.Value{}
	if v == empty {
		// Field does not exist
		return
	}
	if v.Kind() == reflect.Ptr {
		// If it's a pointer we should only return not found if
		// it's nil:
		if v.IsNil() {
			zero = true
			return
		}
	} else if v.Kind() == reflect.Map || v.Kind() == reflect.Slice {
		if v.Len() == 0 {
			zero = true
		}
	} else {
		// If it's not a pointer, we return not found if it's an
		// zero value
		nilValue := reflect.Zero(v.Type())
		if v.Interface() == nilValue.Interface() {
			zero = true
		}
	}
	found = true
	object = v.Interface()
	// This prevents **Foo being returned for pointer when value is already *Foo
	if v.Kind() == reflect.Ptr {
		pointer = v.Interface()
	} else if v.CanAddr() {
		pointer = v.Addr().Interface()
	}
	return
}
