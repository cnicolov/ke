package jsonselect

import (
	"reflect"

	"kego.io/json"

	"kego.io/system"
)

//***********************************************************
//*** @basic ***
//***********************************************************

// Automatically created basic rule for basic
type Basic_rule struct {
	*system.Object
}

//***********************************************************
//*** @c ***
//***********************************************************

// Automatically created basic rule for c
type C_rule struct {
	*system.Object
}

//***********************************************************
//*** @collision ***
//***********************************************************

// Automatically created basic rule for collision
type Collision_rule struct {
	*system.Object
}

//***********************************************************
//*** @diagram ***
//***********************************************************

// Automatically created basic rule for diagram
type Diagram_rule struct {
	*system.Object
}

//***********************************************************
//*** @expr ***
//***********************************************************

// Automatically created basic rule for expr
type Expr_rule struct {
	*system.Object
}

//***********************************************************
//*** @gallery ***
//***********************************************************

// Automatically created basic rule for gallery
type Gallery_rule struct {
	*system.Object
}

//***********************************************************
//*** @image ***
//***********************************************************

// Automatically created basic rule for image
type Image_rule struct {
	*system.Object
}

//***********************************************************
//*** @kid ***
//***********************************************************

// Automatically created basic rule for kid
type Kid_rule struct {
	*system.Object
}

//***********************************************************
//*** @photo ***
//***********************************************************

// Automatically created basic rule for photo
type Photo_rule struct {
	*system.Object
}

//***********************************************************
//*** @polykids ***
//***********************************************************

// Automatically created basic rule for polykids
type Polykids_rule struct {
	*system.Object
}

//***********************************************************
//*** @sibling ***
//***********************************************************

// Automatically created basic rule for sibling
type Sibling_rule struct {
	*system.Object
}

//***********************************************************
//*** @typed ***
//***********************************************************

// Automatically created basic rule for typed
type Typed_rule struct {
	*system.Object
}

//***********************************************************
//*** basic ***
//***********************************************************

type Basic struct {
	*system.Object

	DrinkPreference []system.String

	FavoriteColor system.String

	LanguagesSpoken []map[string]system.String

	Name map[string]system.String

	SeatingPreference []system.String

	Weight system.Number
}

//***********************************************************
//*** c ***
//***********************************************************

type C struct {
	*system.Object

	A system.Number

	B system.Number

	C map[string]system.Number
}

//***********************************************************
//*** collision ***
//***********************************************************

type Collision struct {
	*system.Object

	Number map[string]system.String
}

//***********************************************************
//*** diagram ***
//***********************************************************

type Diagram struct {
	*system.Object

	Url system.String
}

//***********************************************************
//*** expr ***
//***********************************************************

type Expr struct {
	*system.Object

	False system.Bool

	Float system.Number

	Int system.Number

	Null system.String

	String system.String

	String2 system.String

	True system.Bool
}

//***********************************************************
//*** gallery ***
//***********************************************************

// This represents a gallery - it's just a list of images
type Gallery struct {
	*system.Object

	Images map[string]Image
}

//***********************************************************
//*** image ***
//***********************************************************

//***********************************************************
//*** kid ***
//***********************************************************

type Kid struct {
	*system.Object

	Language system.String

	Level system.String

	Preferred system.Bool
}

//***********************************************************
//*** photo ***
//***********************************************************

// This represents an image, and contains path, server and protocol separately
type Photo struct {
	*system.Object

	// The path for the url - e.g. /foo/bar.jpg
	Path system.String

	// The protocol for the url - e.g. http or https
	Protocol system.String `kego:"{\"default\":{\"value\":\"http\"}}"`

	// The server for the url - e.g. www.google.com
	Server system.String
}

//***********************************************************
//*** polykids ***
//***********************************************************

type Polykids struct {
	*system.Object

	A []*Kid
}

//***********************************************************
//*** sibling ***
//***********************************************************

type Sibling struct {
	*system.Object

	A system.Number

	B system.Number

	C *C

	D map[string]system.Number

	E map[string]system.Number
}

//***********************************************************
//*** typed ***
//***********************************************************

type Typed struct {
	*system.Object

	Avatar Image

	DrinkPreference []system.String

	FavoriteColor system.String

	Kids map[string]*Kid

	Name map[string]system.String

	Weight system.Number
}

func init() {

	json.RegisterType("kego.io/jsonselect:@basic", reflect.TypeOf(&Basic_rule{}))

	json.RegisterType("kego.io/jsonselect:@c", reflect.TypeOf(&C_rule{}))

	json.RegisterType("kego.io/jsonselect:@collision", reflect.TypeOf(&Collision_rule{}))

	json.RegisterType("kego.io/jsonselect:@diagram", reflect.TypeOf(&Diagram_rule{}))

	json.RegisterType("kego.io/jsonselect:@expr", reflect.TypeOf(&Expr_rule{}))

	json.RegisterType("kego.io/jsonselect:@gallery", reflect.TypeOf(&Gallery_rule{}))

	json.RegisterType("kego.io/jsonselect:@image", reflect.TypeOf(&Image_rule{}))

	json.RegisterType("kego.io/jsonselect:@kid", reflect.TypeOf(&Kid_rule{}))

	json.RegisterType("kego.io/jsonselect:@photo", reflect.TypeOf(&Photo_rule{}))

	json.RegisterType("kego.io/jsonselect:@polykids", reflect.TypeOf(&Polykids_rule{}))

	json.RegisterType("kego.io/jsonselect:@sibling", reflect.TypeOf(&Sibling_rule{}))

	json.RegisterType("kego.io/jsonselect:@typed", reflect.TypeOf(&Typed_rule{}))

	json.RegisterType("kego.io/jsonselect:basic", reflect.TypeOf(&Basic{}))

	json.RegisterType("kego.io/jsonselect:c", reflect.TypeOf(&C{}))

	json.RegisterType("kego.io/jsonselect:collision", reflect.TypeOf(&Collision{}))

	json.RegisterType("kego.io/jsonselect:diagram", reflect.TypeOf(&Diagram{}))

	json.RegisterType("kego.io/jsonselect:expr", reflect.TypeOf(&Expr{}))

	json.RegisterType("kego.io/jsonselect:gallery", reflect.TypeOf(&Gallery{}))

	json.RegisterType("kego.io/jsonselect:kid", reflect.TypeOf(&Kid{}))

	json.RegisterType("kego.io/jsonselect:photo", reflect.TypeOf(&Photo{}))

	json.RegisterType("kego.io/jsonselect:polykids", reflect.TypeOf(&Polykids{}))

	json.RegisterType("kego.io/jsonselect:sibling", reflect.TypeOf(&Sibling{}))

	json.RegisterType("kego.io/jsonselect:typed", reflect.TypeOf(&Typed{}))

}
