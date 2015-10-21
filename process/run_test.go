package process

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"kego.io/kerr/assert"
	"kego.io/process/settings"
	"kego.io/process/tests"
)

func TestRun(t *testing.T) {

	namespace, err := tests.CreateTemporaryNamespace()
	assert.NoError(t, err)
	defer os.RemoveAll(namespace)

	path, dir, _, err := tests.CreateTemporaryPackage(namespace, "d", map[string]string{
		"a.json": `{"type": "system:type", "id": "a"}`,
		"d.go":   `package d`,
	})

	err = Run(C_TYPES, &settings.Settings{Dir: dir, Path: path})
	assert.NoError(t, err)

	bytes, err := ioutil.ReadFile(filepath.Join(dir, "types", "generated-types.go"))
	assert.NoError(t, err)
	source := string(bytes)
	assert.Contains(t, source, "system.Register")
}