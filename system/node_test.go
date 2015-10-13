package system_test

import (
	"testing"

	"kego.io/json"
	"kego.io/kerr/assert"
	"kego.io/system"
	_ "kego.io/system/types"
)

func TestFoo(t *testing.T) {
	s := `{
		"description": "Restriction rules for bools",
		"type": "type",
		"embed": ["rule"],
		"fields": {
			"default": {
				"description": "Default value if this is missing or null",
				"type": "@bool",
				"optional": true
			}
		}
	}`
	n := &system.Node{}
	err := json.UnmarshalPlainContext([]byte(s), n, "kego.io/system", map[string]string{})
	assert.NoError(t, err)
}
