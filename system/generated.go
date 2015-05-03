package system

import (
	"reflect"

	"kego.io/json"
)

//***********************************************************
//*** @array ***
//***********************************************************

// Restriction rules for arrays
type Array_rule struct {
	*Object

	// This is a rule object, defining the type and restrictions on the value of the items
	Items Rule

	// This is the maximum number of items alowed in the array
	MaxItems Number

	// This is the minimum number of items alowed in the array
	MinItems Number

	// If this is true, each item must be unique
	UniqueItems Bool `kego:"{\"default\":{\"type\":\"kego.io/system:bool\",\"value\":false,\"path\":\"kego.io/system\"}}"`
}

//***********************************************************
//*** @bool ***
//***********************************************************

// Restriction rules for bools
type Bool_rule struct {
	*Object

	// Default value of this is missing or null
	Default Bool
}

//***********************************************************
//*** @context ***
//***********************************************************

// Automatically created basic rule for context
type Context_rule struct {
	*Object
}

//***********************************************************
//*** @map ***
//***********************************************************

// Restriction rules for maps
type Map_rule struct {
	*Object

	// This is a rule object, defining the type and restrictions on the value of the items.
	Items Rule

	// This is the maximum number of items alowed in the array
	MaxItems Number

	// This is the minimum number of items alowed in the array
	MinItems Number
}

//***********************************************************
//*** @number ***
//***********************************************************

// Restriction rules for numbers
type Number_rule struct {
	*Object

	// Default value if this property is omitted
	Default Number

	// If this is true, the value must be less than maximum. If false or not provided, the value must be less than or equal to the maximum.
	ExclusiveMaximum Bool `kego:"{\"default\":{\"type\":\"kego.io/system:bool\",\"value\":false,\"path\":\"kego.io/system\"}}"`

	// If this is true, the value must be greater than minimum. If false or not provided, the value must be greater than or equal to the minimum.
	ExclusiveMinimum Bool `kego:"{\"default\":{\"type\":\"kego.io/system:bool\",\"value\":false,\"path\":\"kego.io/system\"}}"`

	// This provides an upper bound for the restriction
	Maximum Number

	// This provides a lower bound for the restriction
	Minimum Number

	// This restricts the number to be a multiple of the given number
	MultipleOf Number
}

//***********************************************************
//*** @object ***
//***********************************************************

// Automatically created basic rule for object
type Object_rule struct {
	*Object
}

//***********************************************************
//*** @property ***
//***********************************************************

// Automatically created basic rule for property
type Property_rule struct {
	*Object
}

//***********************************************************
//*** @reference ***
//***********************************************************

// Restriction rules for references
type Reference_rule struct {
	*Object

	// Default value of this is missing or null
	Default Reference
}

//***********************************************************
//*** @rule ***
//***********************************************************

// Automatically created basic rule for rule
type Rule_rule struct {
	*Object
}

//***********************************************************
//*** @string ***
//***********************************************************

// Restriction rules for strings
type String_rule struct {
	*Object

	// Default value of this is missing or null
	Default String

	// The value of this string is restricted to one of the provided values
	Enum []String

	// This restricts the value to one of several built-in formats
	Format String

	// The value must be shorter or equal to the provided maximum length
	MaxLength Number

	// The value must be longer or equal to the provided minimum length
	MinLength Number

	// This is a regex to match the value to
	Pattern String
}

//***********************************************************
//*** @type ***
//***********************************************************

// Automatically created basic rule for type
type Type_rule struct {
	*Object
}

//***********************************************************
//*** array ***
//***********************************************************

// This is the native json array data type
type Array struct {
	*Object
}

//***********************************************************
//*** bool ***
//***********************************************************

//***********************************************************
//*** context ***
//***********************************************************

// Unmarshal context.
type Context struct {

	// A list of imports.
	Imports map[string]string

	// The path of the local package.
	Package string
}

//***********************************************************
//*** map ***
//***********************************************************

// This is the native json object data type.
type Map struct {
	*Object
}

//***********************************************************
//*** number ***
//***********************************************************

//***********************************************************
//*** object ***
//***********************************************************

// This is the most basic type.
type Object struct {

	// Unmarshaling context. This should not be in the json - it's added by the unmarshaler.
	Context *Context

	// Description for the developer
	Description string

	// All global objects should have an id.
	Id string

	// Type of the object.
	Type Reference
}

//***********************************************************
//*** property ***
//***********************************************************

// This is a property of a type
type Property struct {
	*Object

	// This specifies that the field is the default value for a rule
	Defaulter Bool `kego:"{\"default\":{\"type\":\"kego.io/system:bool\",\"value\":false,\"path\":\"kego.io/system\"}}"`

	// This is a rule object, defining the type and restrictions on the value of the this property
	Item Rule

	// This specifies that the field is optional
	Optional Bool `kego:"{\"default\":{\"type\":\"kego.io/system:bool\",\"value\":false,\"path\":\"kego.io/system\"}}"`
}

//***********************************************************
//*** reference ***
//***********************************************************

//***********************************************************
//*** rule ***
//***********************************************************

//***********************************************************
//*** string ***
//***********************************************************

//***********************************************************
//*** type ***
//***********************************************************

// This is the most basic type.
type Type struct {
	*Object

	// Type which this should extend
	Extends Reference `kego:"{\"default\":{\"type\":\"kego.io/system:reference\",\"value\":\"kego.io/system:object\",\"path\":\"kego.io/system\"}}"`

	// Is this type an interface?
	Interface Bool `kego:"{\"default\":{\"type\":\"kego.io/system:bool\",\"value\":false,\"path\":\"kego.io/system\"}}"`

	// Array of interface types that this type should support
	Is []Reference

	// This is the native json type that represents this type. If omitted, default is object.
	Native String `kego:"{\"default\":{\"type\":\"kego.io/system:string\",\"value\":\"object\",\"path\":\"kego.io/system\"}}"`

	// Each field is listed with it's type
	Properties map[string]*Property

	// Embedded type that defines restriction rules for this type. By convention, the ID should be this type prefixed with the @ character.
	Rule *Type
}

func init() {

	json.RegisterType("kego.io/system:@array", reflect.TypeOf(&Array_rule{}))

	RegisterType("kego.io/system:@array", &Type{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "Restriction rules for arrays", Id: "@array", Type: Reference{Value: "kego.io/system:type", Package: "kego.io/system", Type: "type", Exists: true}}, Extends: Reference{Value: "kego.io/system:object", Package: "kego.io/system", Type: "object", Exists: true}, Interface: Bool{Value: false, Exists: true}, Is: []Reference{Reference{Value: "kego.io/system:rule", Package: "kego.io/system", Type: "rule", Exists: true}}, Native: String{Value: "object", Exists: true}, Properties: map[string]*Property{"items": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This is a rule object, defining the type and restrictions on the value of the items", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Rule_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@rule", Package: "kego.io/system", Type: "@rule", Exists: true}}}, Optional: Bool{Value: false, Exists: true}}, "minItems": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This is the minimum number of items alowed in the array", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Number_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@number", Package: "kego.io/system", Type: "@number", Exists: true}}, Default: Number{Value: 0, Exists: false}, ExclusiveMaximum: Bool{Value: false, Exists: true}, ExclusiveMinimum: Bool{Value: false, Exists: true}, Maximum: Number{Value: 0, Exists: false}, Minimum: Number{Value: 0, Exists: true}, MultipleOf: Number{Value: 1, Exists: true}}, Optional: Bool{Value: true, Exists: true}}, "maxItems": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This is the maximum number of items alowed in the array", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Number_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@number", Package: "kego.io/system", Type: "@number", Exists: true}}, Default: Number{Value: 0, Exists: false}, ExclusiveMaximum: Bool{Value: false, Exists: true}, ExclusiveMinimum: Bool{Value: false, Exists: true}, Maximum: Number{Value: 0, Exists: false}, Minimum: Number{Value: 0, Exists: true}, MultipleOf: Number{Value: 1, Exists: true}}, Optional: Bool{Value: true, Exists: true}}, "uniqueItems": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "If this is true, each item must be unique", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Bool_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@bool", Package: "kego.io/system", Type: "@bool", Exists: true}}, Default: Bool{Value: false, Exists: true}}, Optional: Bool{Value: true, Exists: true}}}, Rule: (*Type)(nil)})

	json.RegisterType("kego.io/system:@bool", reflect.TypeOf(&Bool_rule{}))

	RegisterType("kego.io/system:@bool", &Type{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "Restriction rules for bools", Id: "@bool", Type: Reference{Value: "kego.io/system:type", Package: "kego.io/system", Type: "type", Exists: true}}, Extends: Reference{Value: "kego.io/system:object", Package: "kego.io/system", Type: "object", Exists: true}, Interface: Bool{Value: false, Exists: true}, Is: []Reference{Reference{Value: "kego.io/system:rule", Package: "kego.io/system", Type: "rule", Exists: true}}, Native: String{Value: "object", Exists: true}, Properties: map[string]*Property{"default": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "Default value of this is missing or null", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: true, Exists: true}, Item: &Bool_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@bool", Package: "kego.io/system", Type: "@bool", Exists: true}}, Default: Bool{Value: false, Exists: false}}, Optional: Bool{Value: true, Exists: true}}}, Rule: (*Type)(nil)})

	json.RegisterType("kego.io/system:@context", reflect.TypeOf(&Context_rule{}))

	RegisterType("kego.io/system:@context", &Type{Object: &Object{Context: (*Context)(nil), Description: "Automatically created basic rule for context", Id: "@context", Type: Reference{Value: "kego.io/system:type", Package: "kego.io/system", Type: "type", Exists: true}}, Extends: Reference{Value: "kego.io/system:object", Package: "kego.io/system", Type: "object", Exists: true}, Interface: Bool{Value: false, Exists: true}, Is: []Reference{Reference{Value: "kego.io/system:rule", Package: "kego.io/system", Type: "rule", Exists: true}}, Native: String{Value: "object", Exists: true}, Properties: map[string]*Property(nil), Rule: (*Type)(nil)})

	json.RegisterType("kego.io/system:@map", reflect.TypeOf(&Map_rule{}))

	RegisterType("kego.io/system:@map", &Type{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "Restriction rules for maps", Id: "@map", Type: Reference{Value: "kego.io/system:type", Package: "kego.io/system", Type: "type", Exists: true}}, Extends: Reference{Value: "kego.io/system:object", Package: "kego.io/system", Type: "object", Exists: true}, Interface: Bool{Value: false, Exists: true}, Is: []Reference{Reference{Value: "kego.io/system:rule", Package: "kego.io/system", Type: "rule", Exists: true}}, Native: String{Value: "object", Exists: true}, Properties: map[string]*Property{"items": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This is a rule object, defining the type and restrictions on the value of the items.", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Rule_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@rule", Package: "kego.io/system", Type: "@rule", Exists: true}}}, Optional: Bool{Value: false, Exists: true}}, "minItems": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This is the minimum number of items alowed in the array", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Number_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@number", Package: "kego.io/system", Type: "@number", Exists: true}}, Default: Number{Value: 0, Exists: false}, ExclusiveMaximum: Bool{Value: false, Exists: true}, ExclusiveMinimum: Bool{Value: false, Exists: true}, Maximum: Number{Value: 0, Exists: false}, Minimum: Number{Value: 0, Exists: true}, MultipleOf: Number{Value: 1, Exists: true}}, Optional: Bool{Value: true, Exists: true}}, "maxItems": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This is the maximum number of items alowed in the array", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Number_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@number", Package: "kego.io/system", Type: "@number", Exists: true}}, Default: Number{Value: 0, Exists: false}, ExclusiveMaximum: Bool{Value: false, Exists: true}, ExclusiveMinimum: Bool{Value: false, Exists: true}, Maximum: Number{Value: 0, Exists: false}, Minimum: Number{Value: 0, Exists: true}, MultipleOf: Number{Value: 1, Exists: true}}, Optional: Bool{Value: true, Exists: true}}}, Rule: (*Type)(nil)})

	json.RegisterType("kego.io/system:@number", reflect.TypeOf(&Number_rule{}))

	RegisterType("kego.io/system:@number", &Type{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "Restriction rules for numbers", Id: "@number", Type: Reference{Value: "kego.io/system:type", Package: "kego.io/system", Type: "type", Exists: true}}, Extends: Reference{Value: "kego.io/system:object", Package: "kego.io/system", Type: "object", Exists: true}, Interface: Bool{Value: false, Exists: true}, Is: []Reference{Reference{Value: "kego.io/system:rule", Package: "kego.io/system", Type: "rule", Exists: true}}, Native: String{Value: "object", Exists: true}, Properties: map[string]*Property{"minimum": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This provides a lower bound for the restriction", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Number_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@number", Package: "kego.io/system", Type: "@number", Exists: true}}, Default: Number{Value: 0, Exists: false}, ExclusiveMaximum: Bool{Value: false, Exists: true}, ExclusiveMinimum: Bool{Value: false, Exists: true}, Maximum: Number{Value: 0, Exists: false}, Minimum: Number{Value: 0, Exists: false}, MultipleOf: Number{Value: 0, Exists: false}}, Optional: Bool{Value: true, Exists: true}}, "exclusiveMinimum": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "If this is true, the value must be greater than minimum. If false or not provided, the value must be greater than or equal to the minimum.", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Bool_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@bool", Package: "kego.io/system", Type: "@bool", Exists: true}}, Default: Bool{Value: false, Exists: true}}, Optional: Bool{Value: true, Exists: true}}, "maximum": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This provides an upper bound for the restriction", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Number_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@number", Package: "kego.io/system", Type: "@number", Exists: true}}, Default: Number{Value: 0, Exists: false}, ExclusiveMaximum: Bool{Value: false, Exists: true}, ExclusiveMinimum: Bool{Value: false, Exists: true}, Maximum: Number{Value: 0, Exists: false}, Minimum: Number{Value: 0, Exists: false}, MultipleOf: Number{Value: 0, Exists: false}}, Optional: Bool{Value: true, Exists: true}}, "exclusiveMaximum": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "If this is true, the value must be less than maximum. If false or not provided, the value must be less than or equal to the maximum.", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Bool_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@bool", Package: "kego.io/system", Type: "@bool", Exists: true}}, Default: Bool{Value: false, Exists: true}}, Optional: Bool{Value: true, Exists: true}}, "default": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "Default value if this property is omitted", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: true, Exists: true}, Item: &Number_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@number", Package: "kego.io/system", Type: "@number", Exists: true}}, Default: Number{Value: 0, Exists: false}, ExclusiveMaximum: Bool{Value: false, Exists: true}, ExclusiveMinimum: Bool{Value: false, Exists: true}, Maximum: Number{Value: 0, Exists: false}, Minimum: Number{Value: 0, Exists: false}, MultipleOf: Number{Value: 0, Exists: false}}, Optional: Bool{Value: true, Exists: true}}, "multipleOf": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This restricts the number to be a multiple of the given number", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Number_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@number", Package: "kego.io/system", Type: "@number", Exists: true}}, Default: Number{Value: 0, Exists: false}, ExclusiveMaximum: Bool{Value: false, Exists: true}, ExclusiveMinimum: Bool{Value: false, Exists: true}, Maximum: Number{Value: 0, Exists: false}, Minimum: Number{Value: 1, Exists: true}, MultipleOf: Number{Value: 1, Exists: true}}, Optional: Bool{Value: true, Exists: true}}}, Rule: (*Type)(nil)})

	json.RegisterType("kego.io/system:@object", reflect.TypeOf(&Object_rule{}))

	RegisterType("kego.io/system:@object", &Type{Object: &Object{Context: (*Context)(nil), Description: "Automatically created basic rule for object", Id: "@object", Type: Reference{Value: "kego.io/system:type", Package: "kego.io/system", Type: "type", Exists: true}}, Extends: Reference{Value: "kego.io/system:object", Package: "kego.io/system", Type: "object", Exists: true}, Interface: Bool{Value: false, Exists: true}, Is: []Reference{Reference{Value: "kego.io/system:rule", Package: "kego.io/system", Type: "rule", Exists: true}}, Native: String{Value: "object", Exists: true}, Properties: map[string]*Property(nil), Rule: (*Type)(nil)})

	json.RegisterType("kego.io/system:@property", reflect.TypeOf(&Property_rule{}))

	RegisterType("kego.io/system:@property", &Type{Object: &Object{Context: (*Context)(nil), Description: "Automatically created basic rule for property", Id: "@property", Type: Reference{Value: "kego.io/system:type", Package: "kego.io/system", Type: "type", Exists: true}}, Extends: Reference{Value: "kego.io/system:object", Package: "kego.io/system", Type: "object", Exists: true}, Interface: Bool{Value: false, Exists: true}, Is: []Reference{Reference{Value: "kego.io/system:rule", Package: "kego.io/system", Type: "rule", Exists: true}}, Native: String{Value: "object", Exists: true}, Properties: map[string]*Property(nil), Rule: (*Type)(nil)})

	json.RegisterType("kego.io/system:@reference", reflect.TypeOf(&Reference_rule{}))

	RegisterType("kego.io/system:@reference", &Type{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "Restriction rules for references", Id: "@reference", Type: Reference{Value: "kego.io/system:type", Package: "kego.io/system", Type: "type", Exists: true}}, Extends: Reference{Value: "kego.io/system:object", Package: "kego.io/system", Type: "object", Exists: true}, Interface: Bool{Value: false, Exists: true}, Is: []Reference{Reference{Value: "kego.io/system:rule", Package: "kego.io/system", Type: "rule", Exists: true}}, Native: String{Value: "object", Exists: true}, Properties: map[string]*Property{"default": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "Default value of this is missing or null", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: true, Exists: true}, Item: &Reference_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@reference", Package: "kego.io/system", Type: "@reference", Exists: true}}, Default: Reference{Value: "", Package: "", Type: "", Exists: false}}, Optional: Bool{Value: true, Exists: true}}}, Rule: (*Type)(nil)})

	json.RegisterType("kego.io/system:@rule", reflect.TypeOf(&Rule_rule{}))

	RegisterType("kego.io/system:@rule", &Type{Object: &Object{Context: (*Context)(nil), Description: "Automatically created basic rule for rule", Id: "@rule", Type: Reference{Value: "kego.io/system:type", Package: "kego.io/system", Type: "type", Exists: true}}, Extends: Reference{Value: "kego.io/system:object", Package: "kego.io/system", Type: "object", Exists: true}, Interface: Bool{Value: false, Exists: true}, Is: []Reference{Reference{Value: "kego.io/system:rule", Package: "kego.io/system", Type: "rule", Exists: true}}, Native: String{Value: "object", Exists: true}, Properties: map[string]*Property(nil), Rule: (*Type)(nil)})

	json.RegisterType("kego.io/system:@string", reflect.TypeOf(&String_rule{}))

	RegisterType("kego.io/system:@string", &Type{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "Restriction rules for strings", Id: "@string", Type: Reference{Value: "kego.io/system:type", Package: "kego.io/system", Type: "type", Exists: true}}, Extends: Reference{Value: "kego.io/system:object", Package: "kego.io/system", Type: "object", Exists: true}, Interface: Bool{Value: false, Exists: true}, Is: []Reference{Reference{Value: "kego.io/system:rule", Package: "kego.io/system", Type: "rule", Exists: true}}, Native: String{Value: "object", Exists: true}, Properties: map[string]*Property{"default": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "Default value of this is missing or null", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: true, Exists: true}, Item: &String_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@string", Package: "kego.io/system", Type: "@string", Exists: true}}, Default: String{Value: "", Exists: false}, Enum: []String(nil), Format: String{Value: "", Exists: false}, MaxLength: Number{Value: 0, Exists: false}, MinLength: Number{Value: 0, Exists: false}, Pattern: String{Value: "", Exists: false}}, Optional: Bool{Value: true, Exists: true}}, "enum": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "The value of this string is restricted to one of the provided values", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Array_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@array", Package: "kego.io/system", Type: "@array", Exists: true}}, Items: &String_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@string", Package: "kego.io/system", Type: "@string", Exists: true}}, Default: String{Value: "", Exists: false}, Enum: []String(nil), Format: String{Value: "", Exists: false}, MaxLength: Number{Value: 0, Exists: false}, MinLength: Number{Value: 0, Exists: false}, Pattern: String{Value: "", Exists: false}}, MaxItems: Number{Value: 0, Exists: false}, MinItems: Number{Value: 0, Exists: false}, UniqueItems: Bool{Value: false, Exists: true}}, Optional: Bool{Value: true, Exists: true}}, "minLength": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "The value must be longer or equal to the provided minimum length", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Number_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@number", Package: "kego.io/system", Type: "@number", Exists: true}}, Default: Number{Value: 0, Exists: false}, ExclusiveMaximum: Bool{Value: false, Exists: true}, ExclusiveMinimum: Bool{Value: false, Exists: true}, Maximum: Number{Value: 0, Exists: false}, Minimum: Number{Value: 0, Exists: false}, MultipleOf: Number{Value: 1, Exists: true}}, Optional: Bool{Value: true, Exists: true}}, "maxLength": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "The value must be shorter or equal to the provided maximum length", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Number_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@number", Package: "kego.io/system", Type: "@number", Exists: true}}, Default: Number{Value: 0, Exists: false}, ExclusiveMaximum: Bool{Value: false, Exists: true}, ExclusiveMinimum: Bool{Value: false, Exists: true}, Maximum: Number{Value: 0, Exists: false}, Minimum: Number{Value: 0, Exists: false}, MultipleOf: Number{Value: 1, Exists: true}}, Optional: Bool{Value: true, Exists: true}}, "pattern": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This is a regex to match the value to", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &String_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@string", Package: "kego.io/system", Type: "@string", Exists: true}}, Default: String{Value: "", Exists: false}, Enum: []String(nil), Format: String{Value: "", Exists: false}, MaxLength: Number{Value: 0, Exists: false}, MinLength: Number{Value: 0, Exists: false}, Pattern: String{Value: "", Exists: false}}, Optional: Bool{Value: true, Exists: true}}, "format": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This restricts the value to one of several built-in formats", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &String_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@string", Package: "kego.io/system", Type: "@string", Exists: true}}, Default: String{Value: "", Exists: false}, Enum: []String{String{Value: "date-time", Exists: true}, String{Value: "email", Exists: true}, String{Value: "hostname", Exists: true}, String{Value: "ipv4", Exists: true}, String{Value: "ipv6", Exists: true}, String{Value: "uri", Exists: true}}, Format: String{Value: "", Exists: false}, MaxLength: Number{Value: 0, Exists: false}, MinLength: Number{Value: 0, Exists: false}, Pattern: String{Value: "", Exists: false}}, Optional: Bool{Value: true, Exists: true}}}, Rule: (*Type)(nil)})

	json.RegisterType("kego.io/system:@type", reflect.TypeOf(&Type_rule{}))

	RegisterType("kego.io/system:@type", &Type{Object: &Object{Context: (*Context)(nil), Description: "Automatically created basic rule for type", Id: "@type", Type: Reference{Value: "kego.io/system:type", Package: "kego.io/system", Type: "type", Exists: true}}, Extends: Reference{Value: "kego.io/system:object", Package: "kego.io/system", Type: "object", Exists: true}, Interface: Bool{Value: false, Exists: true}, Is: []Reference{Reference{Value: "kego.io/system:rule", Package: "kego.io/system", Type: "rule", Exists: true}}, Native: String{Value: "object", Exists: true}, Properties: map[string]*Property(nil), Rule: (*Type)(nil)})

	json.RegisterType("kego.io/system:array", reflect.TypeOf(&Array{}))

	RegisterType("kego.io/system:array", &Type{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This is the native json array data type", Id: "array", Type: Reference{Value: "kego.io/system:type", Package: "kego.io/system", Type: "type", Exists: true}}, Extends: Reference{Value: "kego.io/system:object", Package: "kego.io/system", Type: "object", Exists: true}, Interface: Bool{Value: false, Exists: true}, Is: []Reference(nil), Native: String{Value: "array", Exists: true}, Properties: map[string]*Property(nil), Rule: &Type{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "Restriction rules for arrays", Id: "@array", Type: Reference{Value: "kego.io/system:type", Package: "kego.io/system", Type: "type", Exists: true}}, Extends: Reference{Value: "kego.io/system:object", Package: "kego.io/system", Type: "object", Exists: true}, Interface: Bool{Value: false, Exists: true}, Is: []Reference{Reference{Value: "kego.io/system:rule", Package: "kego.io/system", Type: "rule", Exists: true}}, Native: String{Value: "object", Exists: true}, Properties: map[string]*Property{"items": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This is a rule object, defining the type and restrictions on the value of the items", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Rule_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@rule", Package: "kego.io/system", Type: "@rule", Exists: true}}}, Optional: Bool{Value: false, Exists: true}}, "minItems": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This is the minimum number of items alowed in the array", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Number_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@number", Package: "kego.io/system", Type: "@number", Exists: true}}, Default: Number{Value: 0, Exists: false}, ExclusiveMaximum: Bool{Value: false, Exists: true}, ExclusiveMinimum: Bool{Value: false, Exists: true}, Maximum: Number{Value: 0, Exists: false}, Minimum: Number{Value: 0, Exists: true}, MultipleOf: Number{Value: 1, Exists: true}}, Optional: Bool{Value: true, Exists: true}}, "maxItems": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This is the maximum number of items alowed in the array", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Number_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@number", Package: "kego.io/system", Type: "@number", Exists: true}}, Default: Number{Value: 0, Exists: false}, ExclusiveMaximum: Bool{Value: false, Exists: true}, ExclusiveMinimum: Bool{Value: false, Exists: true}, Maximum: Number{Value: 0, Exists: false}, Minimum: Number{Value: 0, Exists: true}, MultipleOf: Number{Value: 1, Exists: true}}, Optional: Bool{Value: true, Exists: true}}, "uniqueItems": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "If this is true, each item must be unique", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Bool_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@bool", Package: "kego.io/system", Type: "@bool", Exists: true}}, Default: Bool{Value: false, Exists: true}}, Optional: Bool{Value: true, Exists: true}}}, Rule: (*Type)(nil)}})

	json.RegisterType("kego.io/system:bool", reflect.TypeOf(&Bool{}))

	RegisterType("kego.io/system:bool", &Type{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This is the native json bool data type", Id: "bool", Type: Reference{Value: "kego.io/system:type", Package: "kego.io/system", Type: "type", Exists: true}}, Extends: Reference{Value: "kego.io/system:object", Package: "kego.io/system", Type: "object", Exists: true}, Interface: Bool{Value: false, Exists: true}, Is: []Reference(nil), Native: String{Value: "bool", Exists: true}, Properties: map[string]*Property(nil), Rule: &Type{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "Restriction rules for bools", Id: "@bool", Type: Reference{Value: "kego.io/system:type", Package: "kego.io/system", Type: "type", Exists: true}}, Extends: Reference{Value: "kego.io/system:object", Package: "kego.io/system", Type: "object", Exists: true}, Interface: Bool{Value: false, Exists: true}, Is: []Reference{Reference{Value: "kego.io/system:rule", Package: "kego.io/system", Type: "rule", Exists: true}}, Native: String{Value: "object", Exists: true}, Properties: map[string]*Property{"default": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "Default value of this is missing or null", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: true, Exists: true}, Item: &Bool_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@bool", Package: "kego.io/system", Type: "@bool", Exists: true}}, Default: Bool{Value: false, Exists: false}}, Optional: Bool{Value: true, Exists: true}}}, Rule: (*Type)(nil)}})

	json.RegisterType("kego.io/system:context", reflect.TypeOf(&Context{}))

	RegisterType("kego.io/system:context", &Type{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "Unmarshal context.", Id: "context", Type: Reference{Value: "kego.io/system:type", Package: "kego.io/system", Type: "type", Exists: true}}, Extends: Reference{Value: "", Package: "", Type: "", Exists: false}, Interface: Bool{Value: false, Exists: true}, Is: []Reference(nil), Native: String{Value: "object", Exists: true}, Properties: map[string]*Property{"imports": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "A list of imports.", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Map_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@map", Package: "kego.io/system", Type: "@map", Exists: true}}, Items: &JsonString_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/json:@string", Package: "kego.io/json", Type: "@string", Exists: true}}}, MaxItems: Number{Value: 0, Exists: false}, MinItems: Number{Value: 0, Exists: false}}, Optional: Bool{Value: false, Exists: true}}, "package": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "The path of the local package.", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &JsonString_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/json:@string", Package: "kego.io/json", Type: "@string", Exists: true}}}, Optional: Bool{Value: false, Exists: true}}}, Rule: (*Type)(nil)})

	json.RegisterType("kego.io/system:map", reflect.TypeOf(&Map{}))

	RegisterType("kego.io/system:map", &Type{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This is the native json object data type.", Id: "map", Type: Reference{Value: "kego.io/system:type", Package: "kego.io/system", Type: "type", Exists: true}}, Extends: Reference{Value: "kego.io/system:object", Package: "kego.io/system", Type: "object", Exists: true}, Interface: Bool{Value: false, Exists: true}, Is: []Reference(nil), Native: String{Value: "map", Exists: true}, Properties: map[string]*Property(nil), Rule: &Type{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "Restriction rules for maps", Id: "@map", Type: Reference{Value: "kego.io/system:type", Package: "kego.io/system", Type: "type", Exists: true}}, Extends: Reference{Value: "kego.io/system:object", Package: "kego.io/system", Type: "object", Exists: true}, Interface: Bool{Value: false, Exists: true}, Is: []Reference{Reference{Value: "kego.io/system:rule", Package: "kego.io/system", Type: "rule", Exists: true}}, Native: String{Value: "object", Exists: true}, Properties: map[string]*Property{"maxItems": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This is the maximum number of items alowed in the array", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Number_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@number", Package: "kego.io/system", Type: "@number", Exists: true}}, Default: Number{Value: 0, Exists: false}, ExclusiveMaximum: Bool{Value: false, Exists: true}, ExclusiveMinimum: Bool{Value: false, Exists: true}, Maximum: Number{Value: 0, Exists: false}, Minimum: Number{Value: 0, Exists: true}, MultipleOf: Number{Value: 1, Exists: true}}, Optional: Bool{Value: true, Exists: true}}, "items": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This is a rule object, defining the type and restrictions on the value of the items.", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Rule_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@rule", Package: "kego.io/system", Type: "@rule", Exists: true}}}, Optional: Bool{Value: false, Exists: true}}, "minItems": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This is the minimum number of items alowed in the array", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Number_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@number", Package: "kego.io/system", Type: "@number", Exists: true}}, Default: Number{Value: 0, Exists: false}, ExclusiveMaximum: Bool{Value: false, Exists: true}, ExclusiveMinimum: Bool{Value: false, Exists: true}, Maximum: Number{Value: 0, Exists: false}, Minimum: Number{Value: 0, Exists: true}, MultipleOf: Number{Value: 1, Exists: true}}, Optional: Bool{Value: true, Exists: true}}}, Rule: (*Type)(nil)}})

	json.RegisterType("kego.io/system:number", reflect.TypeOf(&Number{}))

	RegisterType("kego.io/system:number", &Type{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This is the native json number data type", Id: "number", Type: Reference{Value: "kego.io/system:type", Package: "kego.io/system", Type: "type", Exists: true}}, Extends: Reference{Value: "kego.io/system:object", Package: "kego.io/system", Type: "object", Exists: true}, Interface: Bool{Value: false, Exists: true}, Is: []Reference(nil), Native: String{Value: "number", Exists: true}, Properties: map[string]*Property(nil), Rule: &Type{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "Restriction rules for numbers", Id: "@number", Type: Reference{Value: "kego.io/system:type", Package: "kego.io/system", Type: "type", Exists: true}}, Extends: Reference{Value: "kego.io/system:object", Package: "kego.io/system", Type: "object", Exists: true}, Interface: Bool{Value: false, Exists: true}, Is: []Reference{Reference{Value: "kego.io/system:rule", Package: "kego.io/system", Type: "rule", Exists: true}}, Native: String{Value: "object", Exists: true}, Properties: map[string]*Property{"maximum": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This provides an upper bound for the restriction", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Number_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@number", Package: "kego.io/system", Type: "@number", Exists: true}}, Default: Number{Value: 0, Exists: false}, ExclusiveMaximum: Bool{Value: false, Exists: true}, ExclusiveMinimum: Bool{Value: false, Exists: true}, Maximum: Number{Value: 0, Exists: false}, Minimum: Number{Value: 0, Exists: false}, MultipleOf: Number{Value: 0, Exists: false}}, Optional: Bool{Value: true, Exists: true}}, "exclusiveMaximum": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "If this is true, the value must be less than maximum. If false or not provided, the value must be less than or equal to the maximum.", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Bool_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@bool", Package: "kego.io/system", Type: "@bool", Exists: true}}, Default: Bool{Value: false, Exists: true}}, Optional: Bool{Value: true, Exists: true}}, "default": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "Default value if this property is omitted", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: true, Exists: true}, Item: &Number_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@number", Package: "kego.io/system", Type: "@number", Exists: true}}, Default: Number{Value: 0, Exists: false}, ExclusiveMaximum: Bool{Value: false, Exists: true}, ExclusiveMinimum: Bool{Value: false, Exists: true}, Maximum: Number{Value: 0, Exists: false}, Minimum: Number{Value: 0, Exists: false}, MultipleOf: Number{Value: 0, Exists: false}}, Optional: Bool{Value: true, Exists: true}}, "multipleOf": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This restricts the number to be a multiple of the given number", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Number_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@number", Package: "kego.io/system", Type: "@number", Exists: true}}, Default: Number{Value: 0, Exists: false}, ExclusiveMaximum: Bool{Value: false, Exists: true}, ExclusiveMinimum: Bool{Value: false, Exists: true}, Maximum: Number{Value: 0, Exists: false}, Minimum: Number{Value: 1, Exists: true}, MultipleOf: Number{Value: 1, Exists: true}}, Optional: Bool{Value: true, Exists: true}}, "minimum": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This provides a lower bound for the restriction", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Number_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@number", Package: "kego.io/system", Type: "@number", Exists: true}}, Default: Number{Value: 0, Exists: false}, ExclusiveMaximum: Bool{Value: false, Exists: true}, ExclusiveMinimum: Bool{Value: false, Exists: true}, Maximum: Number{Value: 0, Exists: false}, Minimum: Number{Value: 0, Exists: false}, MultipleOf: Number{Value: 0, Exists: false}}, Optional: Bool{Value: true, Exists: true}}, "exclusiveMinimum": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "If this is true, the value must be greater than minimum. If false or not provided, the value must be greater than or equal to the minimum.", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Bool_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@bool", Package: "kego.io/system", Type: "@bool", Exists: true}}, Default: Bool{Value: false, Exists: true}}, Optional: Bool{Value: true, Exists: true}}}, Rule: (*Type)(nil)}})

	json.RegisterType("kego.io/system:object", reflect.TypeOf(&Object{}))

	RegisterType("kego.io/system:object", &Type{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This is the most basic type.", Id: "object", Type: Reference{Value: "kego.io/system:type", Package: "kego.io/system", Type: "type", Exists: true}}, Extends: Reference{Value: "", Package: "", Type: "", Exists: false}, Interface: Bool{Value: false, Exists: true}, Is: []Reference(nil), Native: String{Value: "object", Exists: true}, Properties: map[string]*Property{"type": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "Type of the object.", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Reference_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@reference", Package: "kego.io/system", Type: "@reference", Exists: true}}, Default: Reference{Value: "", Package: "", Type: "", Exists: false}}, Optional: Bool{Value: false, Exists: true}}, "id": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "All global objects should have an id.", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &JsonString_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/json:@string", Package: "kego.io/json", Type: "@string", Exists: true}}}, Optional: Bool{Value: true, Exists: true}}, "description": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "Description for the developer", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &JsonString_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/json:@string", Package: "kego.io/json", Type: "@string", Exists: true}}}, Optional: Bool{Value: true, Exists: true}}, "context": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "Unmarshaling context. This should not be in the json - it's added by the unmarshaler.", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Context_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@context", Package: "kego.io/system", Type: "@context", Exists: true}}}, Optional: Bool{Value: true, Exists: true}}}, Rule: (*Type)(nil)})

	json.RegisterType("kego.io/system:property", reflect.TypeOf(&Property{}))

	RegisterType("kego.io/system:property", &Type{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This is a property of a type", Id: "property", Type: Reference{Value: "kego.io/system:type", Package: "kego.io/system", Type: "type", Exists: true}}, Extends: Reference{Value: "kego.io/system:object", Package: "kego.io/system", Type: "object", Exists: true}, Interface: Bool{Value: false, Exists: true}, Is: []Reference(nil), Native: String{Value: "object", Exists: true}, Properties: map[string]*Property{"optional": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This specifies that the field is optional", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Bool_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@bool", Package: "kego.io/system", Type: "@bool", Exists: true}}, Default: Bool{Value: false, Exists: true}}, Optional: Bool{Value: true, Exists: true}}, "defaulter": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This specifies that the field is the default value for a rule", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Bool_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@bool", Package: "kego.io/system", Type: "@bool", Exists: true}}, Default: Bool{Value: false, Exists: true}}, Optional: Bool{Value: true, Exists: true}}, "item": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This is a rule object, defining the type and restrictions on the value of the this property", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Rule_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@rule", Package: "kego.io/system", Type: "@rule", Exists: true}}}, Optional: Bool{Value: false, Exists: true}}}, Rule: (*Type)(nil)})

	json.RegisterType("kego.io/system:reference", reflect.TypeOf(&Reference{}))

	RegisterType("kego.io/system:reference", &Type{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This is a reference to another object, of the form: [local id] or [import name]:[id] or [full package path]:[id]", Id: "reference", Type: Reference{Value: "kego.io/system:type", Package: "kego.io/system", Type: "type", Exists: true}}, Extends: Reference{Value: "kego.io/system:object", Package: "kego.io/system", Type: "object", Exists: true}, Interface: Bool{Value: false, Exists: true}, Is: []Reference(nil), Native: String{Value: "string", Exists: true}, Properties: map[string]*Property(nil), Rule: &Type{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "Restriction rules for references", Id: "@reference", Type: Reference{Value: "kego.io/system:type", Package: "kego.io/system", Type: "type", Exists: true}}, Extends: Reference{Value: "kego.io/system:object", Package: "kego.io/system", Type: "object", Exists: true}, Interface: Bool{Value: false, Exists: true}, Is: []Reference{Reference{Value: "kego.io/system:rule", Package: "kego.io/system", Type: "rule", Exists: true}}, Native: String{Value: "object", Exists: true}, Properties: map[string]*Property{"default": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "Default value of this is missing or null", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: true, Exists: true}, Item: &Reference_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@reference", Package: "kego.io/system", Type: "@reference", Exists: true}}, Default: Reference{Value: "", Package: "", Type: "", Exists: false}}, Optional: Bool{Value: true, Exists: true}}}, Rule: (*Type)(nil)}})

	RegisterType("kego.io/system:rule", &Type{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This is one of several rule types, derived from the rules property of other types", Id: "rule", Type: Reference{Value: "kego.io/system:type", Package: "kego.io/system", Type: "type", Exists: true}}, Extends: Reference{Value: "kego.io/system:object", Package: "kego.io/system", Type: "object", Exists: true}, Interface: Bool{Value: true, Exists: true}, Is: []Reference(nil), Native: String{Value: "object", Exists: true}, Properties: map[string]*Property(nil), Rule: (*Type)(nil)})

	json.RegisterType("kego.io/system:string", reflect.TypeOf(&String{}))

	RegisterType("kego.io/system:string", &Type{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This is the native json string data type", Id: "string", Type: Reference{Value: "kego.io/system:type", Package: "kego.io/system", Type: "type", Exists: true}}, Extends: Reference{Value: "kego.io/system:object", Package: "kego.io/system", Type: "object", Exists: true}, Interface: Bool{Value: false, Exists: true}, Is: []Reference(nil), Native: String{Value: "string", Exists: true}, Properties: map[string]*Property(nil), Rule: &Type{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "Restriction rules for strings", Id: "@string", Type: Reference{Value: "kego.io/system:type", Package: "kego.io/system", Type: "type", Exists: true}}, Extends: Reference{Value: "kego.io/system:object", Package: "kego.io/system", Type: "object", Exists: true}, Interface: Bool{Value: false, Exists: true}, Is: []Reference{Reference{Value: "kego.io/system:rule", Package: "kego.io/system", Type: "rule", Exists: true}}, Native: String{Value: "object", Exists: true}, Properties: map[string]*Property{"default": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "Default value of this is missing or null", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: true, Exists: true}, Item: &String_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@string", Package: "kego.io/system", Type: "@string", Exists: true}}, Default: String{Value: "", Exists: false}, Enum: []String(nil), Format: String{Value: "", Exists: false}, MaxLength: Number{Value: 0, Exists: false}, MinLength: Number{Value: 0, Exists: false}, Pattern: String{Value: "", Exists: false}}, Optional: Bool{Value: true, Exists: true}}, "enum": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "The value of this string is restricted to one of the provided values", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Array_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@array", Package: "kego.io/system", Type: "@array", Exists: true}}, Items: &String_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@string", Package: "kego.io/system", Type: "@string", Exists: true}}, Default: String{Value: "", Exists: false}, Enum: []String(nil), Format: String{Value: "", Exists: false}, MaxLength: Number{Value: 0, Exists: false}, MinLength: Number{Value: 0, Exists: false}, Pattern: String{Value: "", Exists: false}}, MaxItems: Number{Value: 0, Exists: false}, MinItems: Number{Value: 0, Exists: false}, UniqueItems: Bool{Value: false, Exists: true}}, Optional: Bool{Value: true, Exists: true}}, "minLength": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "The value must be longer or equal to the provided minimum length", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Number_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@number", Package: "kego.io/system", Type: "@number", Exists: true}}, Default: Number{Value: 0, Exists: false}, ExclusiveMaximum: Bool{Value: false, Exists: true}, ExclusiveMinimum: Bool{Value: false, Exists: true}, Maximum: Number{Value: 0, Exists: false}, Minimum: Number{Value: 0, Exists: false}, MultipleOf: Number{Value: 1, Exists: true}}, Optional: Bool{Value: true, Exists: true}}, "maxLength": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "The value must be shorter or equal to the provided maximum length", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Number_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@number", Package: "kego.io/system", Type: "@number", Exists: true}}, Default: Number{Value: 0, Exists: false}, ExclusiveMaximum: Bool{Value: false, Exists: true}, ExclusiveMinimum: Bool{Value: false, Exists: true}, Maximum: Number{Value: 0, Exists: false}, Minimum: Number{Value: 0, Exists: false}, MultipleOf: Number{Value: 1, Exists: true}}, Optional: Bool{Value: true, Exists: true}}, "pattern": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This is a regex to match the value to", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &String_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@string", Package: "kego.io/system", Type: "@string", Exists: true}}, Default: String{Value: "", Exists: false}, Enum: []String(nil), Format: String{Value: "", Exists: false}, MaxLength: Number{Value: 0, Exists: false}, MinLength: Number{Value: 0, Exists: false}, Pattern: String{Value: "", Exists: false}}, Optional: Bool{Value: true, Exists: true}}, "format": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This restricts the value to one of several built-in formats", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &String_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@string", Package: "kego.io/system", Type: "@string", Exists: true}}, Default: String{Value: "", Exists: false}, Enum: []String{String{Value: "date-time", Exists: true}, String{Value: "email", Exists: true}, String{Value: "hostname", Exists: true}, String{Value: "ipv4", Exists: true}, String{Value: "ipv6", Exists: true}, String{Value: "uri", Exists: true}}, Format: String{Value: "", Exists: false}, MaxLength: Number{Value: 0, Exists: false}, MinLength: Number{Value: 0, Exists: false}, Pattern: String{Value: "", Exists: false}}, Optional: Bool{Value: true, Exists: true}}}, Rule: (*Type)(nil)}})

	json.RegisterType("kego.io/system:type", reflect.TypeOf(&Type{}))

	RegisterType("kego.io/system:type", &Type{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This is the most basic type.", Id: "type", Type: Reference{Value: "kego.io/system:type", Package: "kego.io/system", Type: "type", Exists: true}}, Extends: Reference{Value: "kego.io/system:object", Package: "kego.io/system", Type: "object", Exists: true}, Interface: Bool{Value: false, Exists: true}, Is: []Reference(nil), Native: String{Value: "object", Exists: true}, Properties: map[string]*Property{"extends": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "Type which this should extend", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Reference_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@reference", Package: "kego.io/system", Type: "@reference", Exists: true}}, Default: Reference{Value: "kego.io/system:object", Package: "kego.io/system", Type: "object", Exists: true}}, Optional: Bool{Value: true, Exists: true}}, "is": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "Array of interface types that this type should support", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Array_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@array", Package: "kego.io/system", Type: "@array", Exists: true}}, Items: &Reference_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@reference", Package: "kego.io/system", Type: "@reference", Exists: true}}, Default: Reference{Value: "", Package: "", Type: "", Exists: false}}, MaxItems: Number{Value: 0, Exists: false}, MinItems: Number{Value: 0, Exists: false}, UniqueItems: Bool{Value: false, Exists: true}}, Optional: Bool{Value: true, Exists: true}}, "native": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "This is the native json type that represents this type. If omitted, default is object.", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &String_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@string", Package: "kego.io/system", Type: "@string", Exists: true}}, Default: String{Value: "object", Exists: true}, Enum: []String{String{Value: "string", Exists: true}, String{Value: "number", Exists: true}, String{Value: "bool", Exists: true}, String{Value: "array", Exists: true}, String{Value: "object", Exists: true}, String{Value: "map", Exists: true}}, Format: String{Value: "", Exists: false}, MaxLength: Number{Value: 0, Exists: false}, MinLength: Number{Value: 0, Exists: false}, Pattern: String{Value: "", Exists: false}}, Optional: Bool{Value: true, Exists: true}}, "interface": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "Is this type an interface?", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Bool_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@bool", Package: "kego.io/system", Type: "@bool", Exists: true}}, Default: Bool{Value: false, Exists: true}}, Optional: Bool{Value: true, Exists: true}}, "properties": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "Each field is listed with it's type", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Map_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@map", Package: "kego.io/system", Type: "@map", Exists: true}}, Items: &Property_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@property", Package: "kego.io/system", Type: "@property", Exists: true}}}, MaxItems: Number{Value: 0, Exists: false}, MinItems: Number{Value: 0, Exists: false}}, Optional: Bool{Value: true, Exists: true}}, "rule": &Property{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "Embedded type that defines restriction rules for this type. By convention, the ID should be this type prefixed with the @ character.", Id: "", Type: Reference{Value: "kego.io/system:property", Package: "kego.io/system", Type: "property", Exists: true}}, Defaulter: Bool{Value: false, Exists: true}, Item: &Type_rule{Object: &Object{Context: &Context{Imports: map[string]string{}, Package: "kego.io/system"}, Description: "", Id: "", Type: Reference{Value: "kego.io/system:@type", Package: "kego.io/system", Type: "@type", Exists: true}}}, Optional: Bool{Value: true, Exists: true}}}, Rule: (*Type)(nil)})

}
