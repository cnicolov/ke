package images

import (
	"reflect"

	"kego.io/demo/common/units"
	"kego.io/json"
	"kego.io/system"
)

// Automatically created basic rule for icon
type Icon_rule struct {
	*system.Object_base
	*system.Rule_base
}

// Restriction rules for images
type Image_rule struct {
	*system.Object_base
	*system.Rule_base
	Secure system.Bool `json:"secure"`
}

// Automatically created basic rule for photo
type Photo_rule struct {
	*system.Object_base
	*system.Rule_base
}

// This is a type of image, which just contains the url of the image
type Icon struct {
	*system.Object_base
	Url system.String `json:"url"`
}

// This represents an image, and contains path, server and protocol separately
type Photo struct {
	*system.Object_base
	// The path for the url - e.g. /foo/bar.jpg
	Path system.String `json:"path"`
	// The protocol for the url - e.g. http or https
	Protocol system.String `kego:"{\"default\":{\"value\":\"http\"}}" json:"protocol"`
	// The server for the url - e.g. www.google.com
	Server system.String    `json:"server"`
	Size   *units.Rectangle `json:"size"`
}

func init() {
	json.Register("kego.io/demo/common/images", "@icon", reflect.TypeOf(&Icon_rule{}), 0x7c5035ca01145c14)
	json.Register("kego.io/demo/common/images", "@image", reflect.TypeOf(&Image_rule{}), 0x1117cdef6aab34e4)
	json.Register("kego.io/demo/common/images", "@photo", reflect.TypeOf(&Photo_rule{}), 0xf2b64533e434a543)
	json.Register("kego.io/demo/common/images", "icon", reflect.TypeOf(&Icon{}), 0x7c5035ca01145c14)
	json.Register("kego.io/demo/common/images", "image", reflect.TypeOf((*Image)(nil)).Elem(), 0x1117cdef6aab34e4)
	json.Register("kego.io/demo/common/images", "photo", reflect.TypeOf(&Photo{}), 0xf2b64533e434a543)
}
