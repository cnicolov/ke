package validate

import (
	"testing"

	"kego.io/process/parser"

	"kego.io/kerr/assert"
	"kego.io/process/tests"
)

func TestValidate_NeedsTypes(t *testing.T) {

	cb := tests.NewContextBuilder().TempGopath(false)
	defer cb.Cleanup()

	files := map[string]string{
		"a.json": `{
			"description": "a",
			"type": "system:type",
			"id": "a",
			"fields": {
				"a": {
					"type": "system:@string"
				}
			}
		}`,
	}
	path, dir, _ := cb.TempPackage("a", files)

	cb.Path(path).Dir(dir).Jsystem().Sauto(parser.Parse)

	err := ValidatePackage(cb.Ctx())
	assert.NoError(t, err)

}

func TestValidate_error1(t *testing.T) {

	cb := tests.NewContextBuilder().TempGopath(false)
	defer cb.Cleanup()

	files := map[string]string{
		"b.json": `{
			"description": "b",
			"type": "system:type",
			"id": "b",
			"fields": {
				"b": {
					"type": "system:@string",
					"minLength": 10,
					"maxLength": 5
				}
			}
		}`,
	}
	path, dir, _ := cb.TempPackage("b", files)

	cb.Path(path).Dir(dir).Jsystem().Sauto(parser.Parse)

	err := ValidatePackage(cb.Ctx())
	// @string is invalid because minLength > maxLength
	assert.HasError(t, err, "YLONAMFUAG")

}
