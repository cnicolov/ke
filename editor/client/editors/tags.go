package editors

import (
	"time"

	"context"

	"github.com/davelondon/vecty"
	"github.com/davelondon/vecty/elem"
	"github.com/davelondon/vecty/event"
	"github.com/davelondon/vecty/prop"
	"kego.io/editor/client/actions"
	"kego.io/editor/client/common"
	"kego.io/editor/client/editable"
	"kego.io/editor/client/models"
	"kego.io/editor/client/stores"
	"kego.io/editor/client/views"
	"kego.io/flux"
	"kego.io/system"
	"kego.io/system/node"
)

var _ editable.EditableArray = (*TagEditor)(nil)

type TagEditor struct{}

func (e *TagEditor) EditorFormatArray(rule *system.RuleWrapper) editable.Format {
	return editable.Block
}

func (e *TagEditor) EditorViewArray(ctx context.Context, node *node.Node, format editable.Format) vecty.Component {
	return NewTagEditorView(ctx, node, format)
}

type TagEditorView struct {
	*views.View

	model  *models.EditorModel
	node   *models.NodeModel
	input  *vecty.Element
	format editable.Format
}

func NewTagEditorView(ctx context.Context, node *node.Node, format editable.Format) *TagEditorView {
	v := &TagEditorView{}
	v.View = views.New(ctx, v)
	v.model = v.App.Editors.Get(node)
	v.node = v.App.Nodes.Get(node)
	v.format = format
	v.Watch(v.model.Node,
		stores.NodeFocus,
		stores.NodeValueChanged,
		stores.NodeErrorsChanged,
	)
	return v
}

func (v *TagEditorView) Reconcile(old vecty.Component) {
	if old, ok := old.(*TagEditorView); ok {
		v.Body = old.Body
	}
	v.ReconcileBody()
}

func (v *TagEditorView) Receive(notif flux.NotifPayload) {
	defer close(notif.Done)
	v.ReconcileBody()
	if notif.Type == stores.NodeFocus {
		v.Focus()
	}
}

func (v *TagEditorView) Focus() {
	v.input.Node().Call("focus")
}

func (v *TagEditorView) Render() vecty.Component {

	contents := vecty.List{
		prop.Value(v.model.Node.ValueString),
		prop.Class("form-control"),
		event.KeyUp(func(e *vecty.Event) {
			getVal := func() interface{} {
				return e.Target.Get("value").String()
			}
			val := getVal()
			changed := func() bool {
				return val != getVal()
			}
			go func() {
				<-time.After(common.EditorKeyboardDebounceShort)
				if changed() {
					return
				}
				v.App.Dispatch(&actions.Modify{
					Undoer:  &actions.Undoer{},
					Editor:  v.model,
					Before:  v.model.Node.NativeValue(),
					After:   val,
					Changed: changed,
				})
			}()
		}),
	}

	if sr, ok := v.model.Node.Rule.Interface.(*system.StringRule); ok && sr.Long {
		v.input = elem.TextArea(
			contents,
		)
	} else {
		v.input = elem.Input(
			prop.Type(prop.TypeText),
			contents,
		)
	}

	return views.NewEditorView(v.Ctx, v.model.Node).Controls(
		v.input,
	)
}
