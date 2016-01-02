package cachectx // import "kego.io/context/cachectx"

import (
	"sync"

	"sort"

	"golang.org/x/net/context"
	"kego.io/context/envctx"
)

// Env is the type of value stored in the Contexts.
type PackageCache struct {
	sync.RWMutex
	m map[string]*PackageInfo
}

type TypeCache struct {
	sync.RWMutex
	m map[string]interface{}
}

type TypeSourceCache struct {
	sync.RWMutex
	m map[string][]byte
}

type GlobalCache struct {
	sync.RWMutex
	m map[string]GlobalInfo
}

type PackageInfo struct {
	Environment  *envctx.Env
	PackageBytes []byte
	Types        *TypeCache
	TypeSource   *TypeSourceCache
	Globals      *GlobalCache
}

type GlobalInfo struct {
	File string
	Name string
}

func (c *PackageCache) Set(env *envctx.Env) *PackageInfo {
	c.Lock()
	defer c.Unlock()
	p := &PackageInfo{
		Environment: env,
		Types:       &TypeCache{m: map[string]interface{}{}},
		TypeSource:  &TypeSourceCache{m: map[string][]byte{}},
		Globals:     &GlobalCache{m: map[string]GlobalInfo{}},
	}
	c.m[env.Path] = p
	return p
}

func (c *PackageCache) Get(path string) (*PackageInfo, bool) {
	c.RLock()
	defer c.RUnlock()
	info, ok := c.m[path]
	return info, ok
}

func (c *PackageCache) All() chan *PackageInfo {
	out := make(chan *PackageInfo)

	go func() {
		c.RLock()
		defer c.RUnlock()
		defer close(out)
		for _, v := range c.m {
			out <- v
		}
	}()

	return out
}

func (c *TypeCache) Set(id string, t interface{}) {
	c.Lock()
	defer c.Unlock()
	c.m[id] = t
}

func (c *TypeCache) Get(id string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()
	t, ok := c.m[id]
	return t, ok
}

func (c *TypeCache) Len() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.m)
}

func (c *TypeCache) All() chan interface{} {
	out := make(chan interface{})

	go func() {
		c.RLock()
		defer c.RUnlock()
		defer close(out)
		keys := []string{}
		for key, _ := range c.m {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			out <- c.m[key]
		}
	}()

	return out
}

func (c *GlobalCache) Set(id string, t GlobalInfo) {
	c.Lock()
	defer c.Unlock()
	c.m[id] = t
}

func (c *GlobalCache) Get(id string) (GlobalInfo, bool) {
	c.RLock()
	defer c.RUnlock()
	t, ok := c.m[id]
	return t, ok
}

func (c *GlobalCache) Len() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.m)
}

func (c *GlobalCache) All() chan GlobalInfo {
	out := make(chan GlobalInfo)

	go func() {
		c.RLock()
		defer c.RUnlock()
		defer close(out)
		keys := []string{}
		for key, _ := range c.m {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			out <- c.m[key]
		}
	}()

	return out
}

func (c *TypeSourceCache) Set(id string, b []byte) {
	c.Lock()
	defer c.Unlock()
	c.m[id] = b
}

func (c *TypeSourceCache) Get(id string) ([]byte, bool) {
	c.RLock()
	defer c.RUnlock()
	t, ok := c.m[id]
	return t, ok
}

func (c *TypeSourceCache) Len() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.m)
}

type Source struct {
	Name  string
	Bytes []byte
}

func (c *TypeSourceCache) All() chan Source {
	out := make(chan Source)

	go func() {
		c.RLock()
		defer c.RUnlock()
		defer close(out)
		keys := []string{}
		for key, _ := range c.m {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			out <- Source{key, c.m[key]}
		}
	}()

	return out
}

// key is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type key int

// cacheKey is the key for cachectx.Cache values in Contexts.  It is
// unexported; clients use cachectx.NewContext and cachectx.FromContext
// instead of using this key directly.
var cacheKey key = 0

// NewContext returns a new Context that carries value u.
func NewContext(ctx context.Context) context.Context {
	c := &PackageCache{m: map[string]*PackageInfo{}}
	return context.WithValue(ctx, cacheKey, c)
}

// FromContext returns the Cache value stored in ctx, and panics if it's not found.
func FromContext(ctx context.Context) *PackageCache {
	e, ok := ctx.Value(cacheKey).(*PackageCache)
	if !ok {
		panic("No cache in ctx")
	}
	return e
}
