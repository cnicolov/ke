package node_test

import (
	"testing"

	"context"

	"github.com/davelondon/ktest/assert"
	"kego.io/context/envctx"
	_ "kego.io/demo/common/images"
	_ "kego.io/demo/common/units"
	_ "kego.io/demo/common/words"
	_ "kego.io/demo/site"
	"kego.io/process"
	"kego.io/process/scanner"
	"kego.io/system/node"
)

func TestUnpack(t *testing.T) {
	testUnpack(t, "kego.io/system")
	testUnpack(t, "kego.io/demo/site")
}
func testUnpack(t *testing.T, path string) {
	ctx, _, err := process.Initialise(context.Background(), &process.Options{
		Path: path,
	})
	assert.NoError(t, err)

	env := envctx.FromContext(ctx)

	files := scanner.ScanDirToFiles(ctx, env.Dir, env.Recursive)
	bytes := scanner.ScanFilesToBytes(ctx, files)
	for b := range bytes {
		_, err := node.Unmarshal(ctx, b.Bytes)
		assert.NoError(t, err, b.File)
	}
}

func TestNodeUnpack(t *testing.T) {

	ctx, _, err := process.Initialise(context.Background(), &process.Options{
		Path: "kego.io/demo/site",
	})
	assert.NoError(t, err)

	j := `{"type":"system:package","aliases":{"images":"kego.io/demo/common/images","units":"kego.io/demo/common/units","words":"kego.io/demo/common/words"}}`

	_, err = node.Unmarshal(ctx, []byte(j))
	assert.NoError(t, err)

}
