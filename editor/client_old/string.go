package editor

import (
	"golang.org/x/net/context"
	"honnef.co/go/js/dom"
	"kego.io/editor/client_old/mdl"
)

type StringEditor struct {
	*Node
	*Editor
	original string
	textbox  *mdl.TextboxStruct
}

var _ EditorInterface = (*StringEditor)(nil)

func NewStringEditor(n *Node) *StringEditor {
	return &StringEditor{Node: n, Editor: &Editor{}}
}

func (e *StringEditor) Layout() Layout {
	return Inline
}

func (e *StringEditor) Initialize(ctx context.Context, holder BranchInterface, layout Layout, fail chan error) error {

	e.Editor.Initialize(ctx, holder, layout, fail)

	e.original = e.ValueString

	e.textbox = mdl.Textbox(e.ValueString, e.Node.Key)
	e.textbox.Style().Set("width", "100%")
	e.AppendChild(e.textbox)
	e.textbox.Input.AddEventListener("input", true, func(ev dom.Event) {
		e.update(e.textbox.Input.Value)
		e.Send(e)
	})

	return nil
}

func (e *StringEditor) update(s string) {
	e.Missing = false
	e.Null = false
	e.ValueString = s
	e.Node.Value = s
}

func (e *StringEditor) Dirty() bool {
	return e.ValueString != e.original
}

func (e *StringEditor) Focus() {
	e.textbox.Input.Focus()
}

func (e *StringEditor) Value() interface{} {
	return e.Node.Value
}