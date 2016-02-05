package json_test

import (
	"fmt"
	"reflect"
	"testing"

	. "kego.io/json"
	"kego.io/kerr"
	"kego.io/kerr/assert"
	"kego.io/process/tests"
	"kego.io/process/tests/unpacker"
)

func TestDecodeSimple(t *testing.T) {
	testDecodeSimple(t, unpacker.Unmarshal)
	testDecodeSimple(t, unpacker.Unpack)
}
func testDecodeSimple(t *testing.T, up unpacker.Interface) {

	test := func(data string) {

		type Foo struct {
			Foo string
			Bar float64
			Baz bool
		}

		ctx := tests.Context("kego.io/json").Jtype("foo", reflect.TypeOf(&Foo{})).Ctx()

		var i interface{}
		err := up.Process(ctx, []byte(data), &i)
		assert.NoError(t, err)
		f, ok := i.(*Foo)
		assert.True(t, ok, "Type %T not correct", i)
		assert.NotNil(t, f)
		assert.Equal(t, "a", f.Foo)
		assert.Equal(t, 2.0, f.Bar)
		assert.Equal(t, true, f.Baz)

	}

	// Standard
	test(`{
		"type": "foo",
		"foo": "a",
		"bar": 2,
		"baz": true
	}`)

	// Type not first
	test(`{
		"foo": "a",
		"bar": 2,
		"type": "foo",
		"baz": true
	}`)

	// Extra attributes
	test(`{
		"foo": "a",
		"qux": "extra value",
		"bar": 2,
		"type": "foo",
		"baz": true
	}`)

}

func TestDecodeDefaults(t *testing.T) {
	testDecodeDefaults(t, unpacker.Unmarshal)
	testDecodeDefaults(t, unpacker.Unpack)
}
func testDecodeDefaults(t *testing.T, up unpacker.Interface) {

	test := func(data string, strExpected string, numExpected float64, boolExpected bool) {

		type Foo struct {
			Foo string  `kego:"{\"default\":{\"value\":\"a\"}}"`
			Bar float64 `kego:"{\"default\":{\"value\":2}}"`
			Baz bool    `kego:"{\"default\":{\"value\":true}}"`
		}

		ctx := tests.Context("kego.io/json").Jtype("foo", reflect.TypeOf(&Foo{})).Ctx()

		var i interface{}
		err := up.Process(ctx, []byte(data), &i)
		assert.NoError(t, err)
		f, ok := i.(*Foo)
		assert.True(t, ok, "Type %T not correct", i)
		assert.NotNil(t, f)
		assert.Equal(t, strExpected, f.Foo)
		assert.Equal(t, numExpected, f.Bar)
		assert.Equal(t, boolExpected, f.Baz)

	}

	test(`{
		"type": "foo",
		"foo": "b",
		"bar": 3,
		"baz": false
	}`, "b", 3.0, false)

	test(`{
		"type": "foo",
		"bar": 3,
		"baz": false
	}`, "a", 3.0, false)

	test(`{
		"type": "foo",
		"baz": false
	}`, "a", 2.0, false)

	test(`{
		"type": "foo"
	}`, "a", 2.0, true)

}

func TestDecodeCollections(t *testing.T) {
	testDecodeCollections(t, unpacker.Unmarshal)
	testDecodeCollections(t, unpacker.Unpack)
}
func testDecodeCollections(t *testing.T, up unpacker.Interface) {

	type Foo struct {
		StringMap   map[string]string
		NumberMap   map[string]float64
		BoolMap     map[string]bool
		StringArray []string
		NumberArray []float64
		BoolArray   []bool
	}

	data := `{
		"type": "foo",
		"stringMap": {"a": "aa", "b": "bb", "c": "cc"},
		"numberMap": {"d": 1, "e": 1.5, "f": 2},
		"boolMap": {"g": true, "h": false, "i": true},
		"stringArray": ["a", "b", "c"],
		"numberArray": [1, 1.5, 2],
		"boolArray": [true, false, true]
	}`

	ctx := tests.Context("kego.io/json").Jtype("foo", reflect.TypeOf(&Foo{})).Ctx()

	var i interface{}
	err := up.Process(ctx, []byte(data), &i)
	assert.NoError(t, err)
	f, ok := i.(*Foo)
	assert.True(t, ok, "Type %T not correct", i)
	assert.NotNil(t, f)
	assert.Equal(t, "aa", f.StringMap["a"])
	assert.Equal(t, "bb", f.StringMap["b"])
	assert.Equal(t, "cc", f.StringMap["c"])
	assert.Equal(t, 1.0, f.NumberMap["d"])
	assert.Equal(t, 1.5, f.NumberMap["e"])
	assert.Equal(t, 2.0, f.NumberMap["f"])
	assert.Equal(t, true, f.BoolMap["g"])
	assert.Equal(t, false, f.BoolMap["h"])
	assert.Equal(t, true, f.BoolMap["i"])
	assert.Equal(t, []string{"a", "b", "c"}, f.StringArray)
	assert.Equal(t, []float64{1.0, 1.5, 2.0}, f.NumberArray)
	assert.Equal(t, []bool{true, false, true}, f.BoolArray)

}

func TestDecodeEmbed(t *testing.T) {
	testDecodeEmbed(t, unpacker.Unmarshal)
	testDecodeEmbed(t, unpacker.Unpack)
}
func testDecodeEmbed(t *testing.T, up unpacker.Interface) {

	test := func(data string) {

		type Bar struct {
			String string
			Number float64
			Bool   bool
		}

		type Foo struct {
			Embed Bar
		}

		ctx := tests.Context("kego.io/json").
			Jtype("foo", reflect.TypeOf(&Foo{})).
			Jtype("bar", reflect.TypeOf(&Bar{})).Ctx()

		var i interface{}
		err := up.Process(ctx, []byte(data), &i)
		assert.NoError(t, err)
		f, ok := i.(*Foo)
		assert.True(t, ok, "Type %T not correct", i)
		assert.NotNil(t, f)
		assert.Equal(t, "a", f.Embed.String)
		assert.Equal(t, 2.0, f.Embed.Number)
		assert.Equal(t, true, f.Embed.Bool)

	}

	// Standard
	test(`{
		"type": "foo",
		"embed": {
			"type": "bar",
			"string": "a",
			"number": 2,
			"bool": true
		}
	}`)

	// Type not first
	test(`{
		"embed": {
			"string": "a",
			"number": 2,
			"type": "bar",
			"bool": true
		},
		"type": "foo"
	}`)

}

func TestDecodeEmbedCollections(t *testing.T) {
	testDecodeEmbedCollections(t, unpacker.Unmarshal)
	testDecodeEmbedCollections(t, unpacker.Unpack)
}
func testDecodeEmbedCollections(t *testing.T, up unpacker.Interface) {

	type Bar struct {
		String string
	}

	type Foo struct {
		MapEmbed   map[string]Bar
		ArrayEmbed []Bar
	}

	data := `{
		"type": "foo",
		"mapEmbed": {
			"a": {
				"type": "bar",
				"string": "a"
			},
			"b": {
				"type": "bar",
				"string": "b"
			}
		},
		"arrayEmbed": [
			{
				"type": "bar",
				"string": "c"
			},
			{
				"type": "bar",
				"string": "d"
			}
		]
	}`

	ctx := tests.Context("kego.io/json").
		Jtype("foo", reflect.TypeOf(&Foo{})).
		Jtype("bar", reflect.TypeOf(&Bar{})).Ctx()

	var i interface{}
	err := up.Process(ctx, []byte(data), &i)
	assert.NoError(t, err)
	f, ok := i.(*Foo)
	assert.True(t, ok, "Type %T not correct", i)
	assert.NotNil(t, f)
	assert.Equal(t, "a", f.MapEmbed["a"].String)
	assert.Equal(t, "b", f.MapEmbed["b"].String)
	assert.Equal(t, "c", f.ArrayEmbed[0].String)
	assert.Equal(t, "d", f.ArrayEmbed[1].String)

}

func TestDecodeComposition(t *testing.T) {
	testDecodeComposition(t, unpacker.Unmarshal)
	testDecodeComposition(t, unpacker.Unpack)
}
func testDecodeComposition(t *testing.T, up unpacker.Interface) {

	type Base struct {
		BaseString string
	}

	type Foo struct {
		*Base
		FooString string
	}

	data := `{
		"type": "foo",
		"baseString": "a",
		"fooString": "b"
	}`

	ctx := tests.Context("kego.io/json").
		Jtype("foo", reflect.TypeOf(&Foo{})).
		Jtype("base", reflect.TypeOf(&Base{})).Ctx()

	var i interface{}
	err := up.Process(ctx, []byte(data), &i)
	assert.NoError(t, err)
	f, ok := i.(*Foo)
	assert.True(t, ok, "Type %T not correct", i)
	assert.NotNil(t, f)
	assert.Equal(t, "a", f.BaseString)
	assert.Equal(t, "b", f.FooString)

}

type Photo struct {
	Id string
}
type Diagram struct {
	Key string
}

func (p *Photo) Url() string {
	return fmt.Sprintf("http://www.photos.com/%s.jpg", p.Id)
}
func (d *Diagram) Url() string {
	return fmt.Sprintf("http://www.diagrams.com/%s.jpg", d.Key)
}

func TestInterface(t *testing.T) {
	testInterface(t, unpacker.Unmarshal)
	testInterface(t, unpacker.Unpack)
}
func testInterface(t *testing.T, up unpacker.Interface) {

	type Image interface {
		Url() string
	}

	type Foo struct {
		Img Image
	}

	data := `{
		"type": "foo",
		"img": {
			"type": "photo",
			"id": "a"
		}
	}`

	ctx := tests.Context("kego.io/json").
		Jtype("foo", reflect.TypeOf(&Foo{})).
		Jtype("photo", reflect.TypeOf(&Photo{})).
		Jtype("diagram", reflect.TypeOf(&Diagram{})).Ctx()

	var i interface{}
	err := up.Process(ctx, []byte(data), &i)
	assert.NoError(t, err)
	f, ok := i.(*Foo)
	assert.True(t, ok, "Type %T not correct", i)
	assert.NotNil(t, f)
	assert.Equal(t, "http://www.photos.com/a.jpg", f.Img.Url())

}

func TestNilInterface(t *testing.T) {
	testNilInterface(t, unpacker.Unmarshal)
	testNilInterface(t, unpacker.Unpack)
}
func testNilInterface(t *testing.T, up unpacker.Interface) {

	type Image interface {
		Url() string
	}

	type Foo struct {
		Iface interface{}
	}

	data := `{
		"type": "foo",
		"iface": {
			"type": "photo",
			"id": "a"
		}
	}`

	ctx := tests.Context("kego.io/json").
		Jtype("foo", reflect.TypeOf(&Foo{})).
		Jtype("photo", reflect.TypeOf(&Photo{})).
		Jtype("diagram", reflect.TypeOf(&Diagram{})).Ctx()

	var i interface{}
	err := up.Process(ctx, []byte(data), &i)
	assert.NoError(t, err)
	f, ok := i.(*Foo)
	assert.True(t, ok, "Type %T not correct", i)
	assert.NotNil(t, f)
	p, ok := f.Iface.(*Photo)
	assert.True(t, ok, "Type %T not correct", f.Iface)
	assert.Equal(t, "http://www.photos.com/a.jpg", p.Url())

}

func TestInterfaceCollections(t *testing.T) {
	testInterfaceCollections(t, unpacker.Unmarshal)
	testInterfaceCollections(t, unpacker.Unpack)
}
func testInterfaceCollections(t *testing.T, up unpacker.Interface) {

	type Image interface {
		Url() string
	}

	type Foo struct {
		ImageArray []Image
		ImageMap   map[string]Image
	}

	data := `{
		"type": "foo",
		"imageArray": [
			{
				"type": "photo",
				"id": "a"
			},
			{
				"type": "diagram",
				"key": "b"
			}
		],
		"imageMap": {
			"c": {
				"type": "photo",
				"id": "c"
			},
			"d": {
				"type": "diagram",
				"key": "d"
			}
		}
	}`

	ctx := tests.Context("kego.io/json").
		Jtype("foo", reflect.TypeOf(&Foo{})).
		Jtype("photo", reflect.TypeOf(&Photo{})).
		Jtype("diagram", reflect.TypeOf(&Diagram{})).Ctx()

	var i interface{}
	err := up.Process(ctx, []byte(data), &i)
	assert.NoError(t, err)
	f, ok := i.(*Foo)
	assert.True(t, ok, "Type %T not correct", i)
	assert.NotNil(t, f)
	assert.Equal(t, "http://www.photos.com/a.jpg", f.ImageArray[0].Url())
	assert.Equal(t, "http://www.diagrams.com/b.jpg", f.ImageArray[1].Url())
	assert.Equal(t, "http://www.photos.com/c.jpg", f.ImageMap["c"].Url())
	assert.Equal(t, "http://www.diagrams.com/d.jpg", f.ImageMap["d"].Url())

}

func TestInterfaceCollectionsComplex(t *testing.T) {
	testInterfaceCollectionsComplex(t, unpacker.Unmarshal)
	testInterfaceCollectionsComplex(t, unpacker.Unpack)
}
func testInterfaceCollectionsComplex(t *testing.T, up unpacker.Interface) {

	type Image interface {
		Url() string
	}

	type Foo struct {
		ImageMap map[string][]Image
	}

	data := `{
		"type": "foo",
		"imageMap": {
			"a": [
				{
					"type": "photo",
					"id": "c"
				},
				{
					"type": "diagram",
					"key": "d"
				}
			],
			"b": [
				{
					"type": "photo",
					"id": "e"
				},
				{
					"type": "diagram",
					"key": "f"
				}
			]
		}
	}`

	ctx := tests.Context("kego.io/json").
		Jtype("foo", reflect.TypeOf(&Foo{})).
		Jtype("photo", reflect.TypeOf(&Photo{})).
		Jtype("diagram", reflect.TypeOf(&Diagram{})).Ctx()

	var i interface{}
	err := up.Process(ctx, []byte(data), &i)
	assert.NoError(t, err)
	f, ok := i.(*Foo)
	assert.True(t, ok, "Type %T not correct", i)
	assert.NotNil(t, f)
	assert.Equal(t, "http://www.photos.com/c.jpg", f.ImageMap["a"][0].Url())
	assert.Equal(t, "http://www.diagrams.com/d.jpg", f.ImageMap["a"][1].Url())
	assert.Equal(t, "http://www.photos.com/e.jpg", f.ImageMap["b"][0].Url())
	assert.Equal(t, "http://www.diagrams.com/f.jpg", f.ImageMap["b"][1].Url())

}

func TestDummyInterfaceNotFound(t *testing.T) {
	testDummyInterfaceNotFound(t, unpacker.Unmarshal)
	testDummyInterfaceNotFound(t, unpacker.Unpack)
}
func testDummyInterfaceNotFound(t *testing.T, up unpacker.Interface) {

	type Image interface {
		Url() string
	}

	type Foo struct {
		Img Image
	}

	cb := tests.Context("kego.io/json").
		Jtype("foo", reflect.TypeOf(&Foo{})).
		Jtype("photo", reflect.TypeOf(&Photo{})).
		Jtype("diagram", reflect.TypeOf(&Diagram{}))

	data := `{
		"type": "foo",
		"img": {
			"type": "bar",
			"id": "a"
		}
	}`
	var i interface{}
	err := up.Process(cb.Ctx(), []byte(data), &i)
	if up == unpacker.Unpack {
		assert.IsError(t, err, "LJBTNGPVSY")
		assert.HasError(t, err, "VUEFNKSTLG")
	} else {
		assert.IsError(t, err, "YEOQSWVFVH")
	}
	ut, ok := kerr.Source(err).(UnknownTypeError)
	assert.True(t, ok)
	assert.Equal(t, "kego.io/json:bar", ut.UnknownType)
	assert.True(t, i.(*Foo).Img == nil)

	data = `{
		"type": "foo",
		"img": {
			"type": "foo:bar",
			"id": "a"
		}
	}`
	err = up.Process(cb.Ctx(), []byte(data), &i)
	if up == unpacker.Unpack {
		assert.IsError(t, err, "LJBTNGPVSY")
		assert.HasError(t, err, "WLKNMHPWJN")
	} else {
		assert.IsError(t, err, "TBNIEVUCPL")
	}
	upe, ok := kerr.Source(err).(UnknownPackageError)
	assert.True(t, ok)
	assert.Equal(t, "foo", upe.UnknownPackage)

	data = `{
		"type": "foo",
		"img": {
			"type": "a.b/c:bar",
			"id": "a"
		}
	}`
	err = up.Process(cb.Ctx(), []byte(data), &i)
	if up == unpacker.Unpack {
		assert.IsError(t, err, "LJBTNGPVSY")
		assert.HasError(t, err, "WLKNMHPWJN")
	} else {
		assert.IsError(t, err, "TBNIEVUCPL")
	}
	upe, ok = kerr.Source(err).(UnknownPackageError)
	assert.True(t, ok)
	assert.Equal(t, "a.b/c", upe.UnknownPackage)

	err = up.Process(cb.Alias("a.b/c", "d").Ctx(), []byte(data), &i)
	if up == unpacker.Unpack {
		assert.IsError(t, err, "LJBTNGPVSY")
		assert.HasError(t, err, "VUEFNKSTLG")
	} else {
		assert.IsError(t, err, "YEOQSWVFVH")
	}
	ut, ok = kerr.Source(err).(UnknownTypeError)
	assert.True(t, ok)
	assert.Equal(t, "a.b/c:bar", ut.UnknownType)

}

type dummyImage struct{}

func (d *dummyImage) Url() string {
	return ""
}

func TestDummyInterface(t *testing.T) {
	testDummyInterface(t, unpacker.Unmarshal)
	testDummyInterface(t, unpacker.Unpack)
}
func testDummyInterface(t *testing.T, up unpacker.Interface) {

	type Image interface {
		Url() string
	}

	type Foo struct {
		Img Image
	}

	ctx := tests.Context("kego.io/json").
		Jtype("foo", reflect.TypeOf(&Foo{})).
		Jtype("photo", reflect.TypeOf(&Photo{})).
		Jtype("diagram", reflect.TypeOf(&Diagram{})).
		Dummy(reflect.TypeOf((*Image)(nil)).Elem(), reflect.TypeOf(&dummyImage{})).
		Ctx()

	data := `{
		"type": "foo",
		"img": {
			"type": "bar",
			"id": "a"
		}
	}`
	var i interface{}
	err := up.Process(ctx, []byte(data), &i)
	assert.NoError(t, err)
}