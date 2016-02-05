package parser_test

import (
	"testing"

	"kego.io/process/parser"

	"kego.io/context/sysctx"
	"kego.io/kerr/assert"
	"kego.io/process/generate"
	"kego.io/process/tests"
	"kego.io/system"
)

func TestScan_errors(t *testing.T) {

	cb := tests.Context("").Cmd().Sempty().Jsystem().TempGopath(false)
	defer cb.Cleanup()

	_, err := parser.Parse(cb.Ctx(), "does-not-exist")
	assert.IsError(t, err, "GJRHNGGWFD")  // Error from parser.ScanForEnv
	assert.HasError(t, err, "SUTCWEVRXS") // Dir not found from packages.getDirFromEmptyPackage

	path, _, _ := cb.TempPackage("a", map[string]string{
		"a.json": `{
			"type": "system:package",
			"recursive": false
		}`,
		"b.json": "foo",
	})

	cb.Path(path)

	_, err = parser.Parse(cb.Ctx(), path)
	assert.IsError(t, err, "VFUNPHUFHD")  // Error from parser.scanForTypes
	assert.HasError(t, err, "HCYGNBDFFA") // Error trying to unmarshal a type

	path, _, _ = cb.TempPackage("a", map[string]string{
		"a.json": "foo",
	})

	cb.Path(path)

	_, err = parser.Parse(cb.Ctx(), path)
	assert.IsError(t, err, "GJRHNGGWFD")  // Error from parser.ScanForEnv
	assert.HasError(t, err, "MTDCXBYBEJ") // Error trying to scan for packages

}

func TestParseRule(t *testing.T) {

	cb := tests.NewContextBuilder().TempGopath(false)
	defer cb.Cleanup()

	path, dir, _ := cb.TempPackage("a", map[string]string{
		"a.json": `{
			"description": "a",
			"type": "system:type",
			"id": "b",
			"fields": {
				"c": {
					"type": "system:@string"
				}
			},
			"rule": {
				"description": "d",
				"type": "system:type",
				"embed": ["system:rule"],
				"fields": {
					"e": {
						"description": "f",
						"type": "system:@string"
					}
				}
			}
		}`,
	})

	cb.Path(path).Dir(dir).Cmd().Sempty().Jsystem()

	_, err := parser.Parse(cb.Ctx(), path)
	assert.NoError(t, err)

	scache := sysctx.FromContext(cb.Ctx())
	i, ok := scache.GetType(path, "b")
	assert.True(t, ok)
	ty, ok := i.(*system.Type)
	assert.True(t, ok)
	assert.Equal(t, "@b", ty.Rule.Id.Name)

}

func TestParse1(t *testing.T) {

	cb := tests.NewContextBuilder().TempGopath(false)
	defer cb.Cleanup()

	path, dir, _ := cb.TempPackage("a", map[string]string{
		"a.json": `{
			"description": "a",
			"type": "system:type",
			"id": "b",
			"fields": {
				"c": {
					"type": "system:@string"
				}
			}
		}`,
	})

	cb.Path(path).Dir(dir).Cmd().Sempty().Jsystem()

	_, err := parser.Parse(cb.Ctx(), path)
	assert.NoError(t, err)

	scache := sysctx.FromContext(cb.Ctx())
	i, ok := scache.GetType(path, "b")
	assert.True(t, ok)
	ty, ok := i.(*system.Type)
	assert.True(t, ok)
	assert.Equal(t, "a", ty.Description)

}

func TestParse(t *testing.T) {

	cb := tests.NewContextBuilder().TempGopath(false)
	defer cb.Cleanup()

	files := map[string]string{
		"a.json": `{
			"type": "system:type",
			"id": "a",
			"rule": {
				"type": "system:type",
				"embed": ["system:rule"],
				"fields": {
					"a": {
						"type": "system:@bool"
					}
				}
			}
		}`,
		"b.json": `{
			"type": "system:type",
			"id": "b",
			"fields": {
				"c": {
					"type": "@a",
					"a": true,
					"optional": true
				},
				"d": {
					"type": "@b"
				}
			}
		}`,
	}
	path, dir, _ := cb.TempPackage("a", files)

	path = "kego.io/system"

	cb.Path(path).Dir(dir).Cmd().Sempty().Jsystem()

	pi, err := parser.Parse(cb.Ctx(), path)
	assert.NoError(t, err)

	_, err = generate.Structs(cb.Ctx(), pi.Environment)
	assert.NoError(t, err)

}
