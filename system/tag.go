package system

import (
	"context"
	"fmt"
	"regexp"

	"github.com/davelondon/kerr"
)

type Tag string

func (t *Tag) Value() string {
	return string(*t)
}

func (r *TagRule) Validate(ctx context.Context) (fail bool, messages []string, err error) {
	if r.Pattern != nil {
		if _, err := regexp.Compile(r.Pattern.Value()); err != nil {
			fail = true
			messages = append(messages, fmt.Sprintf("Pattern: regex does not compile: %s", r.Pattern.Value()))
		}
	}
	return
}

var _ Validator = (*TagRule)(nil)

func (r *TagRule) Enforce(ctx context.Context, data interface{}) (fail bool, messages []string, err error) {

	if i, ok := data.(TagInterface); ok && i != nil {
		data = i.GetTag(ctx)
	}

	t, ok := data.(*Tag)
	if !ok && data != nil {
		return true, nil, kerr.New("OHGAJQKYSP", "Tag rule: value %T should be *system.Tag", data)
	}

	// Internals should be a subset of StringRule, so create a dummy StringRule
	// and String.
	sr := StringRule{
		Format:  r.Format,
		Enum:    r.Enum,
		Pattern: r.Pattern,
	}
	s := NewString(t.Value())
	if fail, messages, err = sr.Enforce(ctx, s); err != nil {
		err = kerr.Wrap("WNQJDLCSGF", err)
	}

	return
}

var _ Enforcer = (*TagRule)(nil)
