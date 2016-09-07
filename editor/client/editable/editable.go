package editable

import (
	"context"

	"github.com/davelondon/vecty"
	"kego.io/system"
	"kego.io/system/node"
)

type Format string

const (
	Branch Format = "branch"
	Block  Format = "block"
	Inline Format = "inline"
)

type Editable interface {
	EditorView(ctx context.Context, node *node.Node, format Format) vecty.Component
	EditorFormat(rule *system.RuleWrapper) Format
}

type EditableArray interface {
	EditorViewArray(ctx context.Context, node *node.Node, format Format) vecty.Component
	EditorFormatArray(rule *system.RuleWrapper) Format
}

type EditableMap interface {
	EditorViewMap(ctx context.Context, node *node.Node, format Format) vecty.Component
	EditorFormatMap(rule *system.RuleWrapper) Format
}

type EditsExtraEmbeddedTypes interface {
	ExtraEmbeddedTypes() []*system.Reference
}
