package generate

import (
	"regexp"
	"testing"

	"github.com/davelondon/ktest/assert"
	"kego.io/context/sysctx"
	"kego.io/process/generate/builder"
	"kego.io/process/parser"
	"kego.io/system"
	"kego.io/tests"
)

func TestGenerateSourceErr1(t *testing.T) {
	cb := tests.Context("a.b/c").Sempty()
	_, err := Structs(cb.Ctx(), cb.Env())
	assert.IsError(t, err, "DQVQWTKRSK")
}
func TestGenerateSourceErr2(t *testing.T) {

	cb := tests.Context("b.c/d")

	cb.Stype("a", &system.Type{
		Object: &system.Object{
			Id:   system.NewReference("b.c/d", "a"),
			Type: system.NewReference("kego.io/system", "type")},
		Native: system.NewString("object"),
		Fields: map[string]system.RuleInterface{
			"c": &system.StringRule{
				Object: &system.Object{
					Type: system.NewReference("kego.io/system", "@foo"),
				},
				Rule: &system.Rule{},
			},
		},
	})

	cb.Alias("e", "f.g/h")

	_, err := Structs(cb.Ctx(), cb.Env())
	assert.IsError(t, err, "XKRYMXUIJD")
	assert.HasError(t, err, "TFXFBIRXHN")
	assert.HasError(t, err, "KYCTDXKFYR") // type system:foo not found
}
func TestGenerateSourceNoTypes(t *testing.T) {
	cb := tests.Context("a.b/c").Spkg("a.b/c")
	_, err := Structs(cb.Ctx(), cb.Env())
	assert.NoError(t, err)
}
func TestGenerateSource(t *testing.T) {

	cb := tests.Context("b.c/d")

	// Native collections shouldn't go into the generated file
	cb.Stype("amap", &system.Type{
		Object: &system.Object{
			Id:   system.NewReference("b.c/d", "amap"),
			Type: system.NewReference("kego.io/system", "type")},
		Native: system.NewString("map"),
	})

	// Native collections shouldn't go into the generated file
	cb.Stype("aarr", &system.Type{
		Object: &system.Object{
			Id:   system.NewReference("b.c/d", "aarr"),
			Type: system.NewReference("kego.io/system", "type")},
		Native: system.NewString("array"),
	})

	cb.Stype("a", &system.Type{
		Object: &system.Object{
			Description: "d",
			Id:          system.NewReference("b.c/d", "a"),
			Type:        system.NewReference("kego.io/system", "type")},
		Native: system.NewString("object"),
	})

	cb.Stype("an", &system.Type{
		Object: &system.Object{
			Id:   system.NewReference("b.c/d", "an"),
			Type: system.NewReference("kego.io/system", "type")},
		Native: system.NewString("bool"),
	})

	cb.Stype("ai", &system.Type{
		Object: &system.Object{
			Id:   system.NewReference("b.c/d", "ai"),
			Type: system.NewReference("kego.io/system", "type")},
		Native:    system.NewString("object"),
		Interface: true,
	})

	cb.Alias("e", "f.g/h")

	source, err := Structs(cb.Ctx(), cb.Env())
	assert.NoError(t, err)
	assert.Contains(t, string(source), "package d\n")
	imp := getImports(t, string(source))
	assert.Contains(t, imp, "\t\"context\"\n")
	assert.Contains(t, imp, "\t\"kego.io/context/jsonctx\"\n")
	assert.Contains(t, imp, "\t\"kego.io/system\"\n")
	assert.Contains(t, imp, "\t\"reflect\"\n")
	assert.NotContains(t, imp, "\"f.g/h\"")
	assert.NotContains(t, imp, "Amap") // Native collections shouldn't go into the generated file
	assert.NotContains(t, imp, "Aarr") // Native collections shouldn't go into the generated file
	assert.Contains(t, string(source), "\n// d\ntype A struct {\n\t*system.Object\n}\n")
	assert.Contains(t, string(source), "pkg.InitType(\"a\", reflect.TypeOf((*A)(nil)), reflect.TypeOf((*ARule)(nil)), reflect.TypeOf((*AInterface)(nil)).Elem())\n")

	// Type "Ai" is an interface
	assert.Contains(t, string(source), "pkg.InitType(\"ai\", reflect.TypeOf((*Ai)(nil)).Elem(), reflect.TypeOf((*AiRule)(nil)), nil)\n")

	// Type "An" should not be a struct because it's a bool native. However, it should be registered.
	assert.NotContains(t, string(source), "type An struct")
	assert.Contains(t, string(source), "pkg.InitType(\"an\", reflect.TypeOf((*An)(nil)), reflect.TypeOf((*AnRule)(nil)), reflect.TypeOf((*AnInterface)(nil)).Elem())\n")

}

func TestPrintStructDefinition(t *testing.T) {
	cb := tests.Context("b.c/d").Ssystem(parser.Parse)

	r := &system.StringRule{
		Object: &system.Object{
			Description: "e",
			Type:        system.NewReference("kego.io/system", "@string"),
		},
		Rule: &system.Rule{},
	}

	ty := &system.Type{
		Object: &system.Object{
			Description: "d",
			Id:          system.NewReference("b.c/d", "a"),
			Type:        system.NewReference("kego.io/system", "type")},
		Embed:  []*system.Reference{system.NewReference("b.c/d", "b")},
		Native: system.NewString("object"),
		Fields: map[string]system.RuleInterface{"c": r},
	}

	b := builder.New("b.c/d")

	err := printStructDefinition(cb.Ctx(), cb.Env(), b, ty)
	assert.NoError(t, err)

	source, err := b.Build()
	assert.NoError(t, err)
	assert.Contains(t, string(source), `package d

import (
	"kego.io/system"
)

// d
type A struct {
	*system.Object
	*B
	// e
	C *system.String `+"`"+`json:"c"`+"`"+`
}`)

}

func TestPrintInitFunction(t *testing.T) {
	cb := tests.Context("b.c/d").Ssystem(parser.Parse)

	cb.Stype("a", &system.Type{
		Object: &system.Object{
			Description: "d",
			Id:          system.NewReference("b.c/d", "a"),
			Type:        system.NewReference("kego.io/system", "type")},
		Embed:  []*system.Reference{system.NewReference("b.c/d", "b")},
		Native: system.NewString("object"),
		Fields: map[string]system.RuleInterface{
			"c": &system.StringRule{
				Object: &system.Object{
					Description: "e",
					Type:        system.NewReference("kego.io/system", "@string"),
				},
				Rule: &system.Rule{},
			},
		},
	})

	cb.Stype("@a", &system.Type{
		Object: &system.Object{
			Id:   system.NewReference("b.c/d", "@a"),
			Type: system.NewReference("kego.io/system", "type")},
		Native: system.NewString("object"),
	})

	cb.Stype("amap", &system.Type{
		Object: &system.Object{
			Id:   system.NewReference("b.c/d", "amap"),
			Type: system.NewReference("kego.io/system", "type")},
		Native: system.NewString("map"),
	})

	cb.Stype("@amap", &system.Type{
		Object: &system.Object{
			Id:   system.NewReference("b.c/d", "@amap"),
			Type: system.NewReference("kego.io/system", "type")},
		Native: system.NewString("object"),
	})

	cb.Stype("aiface", &system.Type{
		Object: &system.Object{
			Id:   system.NewReference("b.c/d", "aiface"),
			Type: system.NewReference("kego.io/system", "type")},
		Interface: true,
		Native:    system.NewString("object"),
	})

	scache := sysctx.FromContext(cb.Ctx())
	pi, ok := scache.Get("b.c/d")
	assert.True(t, ok)

	b := builder.New("b.c/d")

	printInitFunction(cb.Env(), b, pi.Types)

	source, err := b.Build()
	assert.NoError(t, err)
	assert.Contains(t, string(source), `package d

import (
	"kego.io/context/jsonctx"
	"reflect"
)

func init() {
	pkg := jsonctx.InitPackage("b.c/d", 0)
	pkg.InitType("a", reflect.TypeOf((*A)(nil)), reflect.TypeOf((*ARule)(nil)), reflect.TypeOf((*AInterface)(nil)).Elem())
	pkg.InitType("aiface", reflect.TypeOf((*Aiface)(nil)).Elem(), reflect.TypeOf((*AifaceRule)(nil)), nil)
	pkg.InitType("amap", nil, reflect.TypeOf((*AmapRule)(nil)), nil)
}
`)

}

func getImports(t *testing.T, source string) string {
	r := regexp.MustCompile(`(?m)import \(([^)]*)\)`)
	matches := r.FindStringSubmatch(source)
	assert.Equal(t, 2, len(matches))
	return matches[1]
}

func TestGenerateCommand_errors(t *testing.T) {

	cb := tests.Context("a.b/\"").Spkg("a.b/\"")

	_, err := Structs(cb.Ctx(), cb.Env())
	// Quote in the path will generate malformed source
	assert.HasError(t, err, "CRBYOUOHPG")

	ty := &system.Type{
		Object: &system.Object{
			Id:   system.NewReference("b.c/d", "a corrupt"),
			Type: system.NewReference("kego.io/system", "type")},
		Native: system.NewString("object"),
	}
	cb = tests.Context("b.c/d").Stype("a", ty)

	_, err = Structs(cb.Ctx(), cb.Env())
	// Corrupt type ID causes error from source formatter
	assert.HasError(t, err, "CRBYOUOHPG")

}
