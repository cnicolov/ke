package clientctx

import (
	"sync"

	"context"

	"github.com/davelondon/kerr"
	"kego.io/json"
)

type ctxKeyType int

var ctxKey ctxKeyType = 0

func NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey, &EditorCache{m: map[string]interface{}{}})
}

func FromContext(ctx context.Context) *EditorCache {
	ec, ok := ctx.Value(ctxKey).(*EditorCache)
	if !ok {
		panic(kerr.New("BIUVXISEMA", "No editors in ctx").Error())
	}
	return ec
}

func FromContextOrNil(ctx context.Context) *EditorCache {
	ec, ok := ctx.Value(ctxKey).(*EditorCache)
	if ok {
		return ec
	}
	return nil
}

type EditorCache struct {
	sync.RWMutex
	m map[string]interface{}
}

func (c *EditorCache) Get(path string, name string, t json.Type) (interface{}, bool) {
	k := key{path: path, name: name, typ: t}
	c.RLock()
	defer c.RUnlock()
	e, ok := c.m[k]
	return e, ok
}

func (c *EditorCache) Set(path string, name string, t json.Type, e interface{}) {
	k := key{path: path, name: name, typ: t}
	c.Lock()
	defer c.Unlock()
	c.m[k] = e
}

type key struct {
	path string
	name string
	typ  string
}
