package images

// ke: {"package": {"notest": true}}

import (
	"github.com/davelondon/vecty"
	"github.com/davelondon/vecty/elem"
	"github.com/davelondon/vecty/prop"
	"github.com/davelondon/vecty/style"
	"golang.org/x/net/context"
	"kego.io/editor/client/editable"
	"kego.io/editor/client/models"
	"kego.io/editor/client/stores"
	"kego.io/flux"
	"kego.io/system"
	"kego.io/system/editors"
	"kego.io/system/node"
)

func (s *Icon) Format(rule *system.RuleWrapper) editable.Format {
	return editable.Block
}

func (s *Icon) EditorView(ctx context.Context, node *node.Node) vecty.Component {
	return NewIconEditorView(ctx, node)
}

type IconEditorView struct {
	vecty.Composite
	ctx    context.Context
	app    *stores.App
	notifs chan flux.NotifPayload

	model *models.EditorModel
	icon  *Icon
	input *vecty.Element
}

func NewIconEditorView(ctx context.Context, node *node.Node) *IconEditorView {
	v := &IconEditorView{
		ctx: ctx,
		app: stores.FromContext(ctx),
	}
	v.model = v.app.Editors.Get(node)
	v.icon = v.model.Node.Value.(*Icon)
	v.Mount()
	return v
}

func (v *IconEditorView) Reconcile(old vecty.Component) {
	if old, ok := old.(*IconEditorView); ok {
		v.Body = old.Body
	}
	v.RenderFunc = v.render
	v.ReconcileBody()
}

// Apply implements the vecty.Markup interface.
func (v *IconEditorView) Apply(element *vecty.Element) {
	element.AddChild(v)
}

func (v *IconEditorView) Mount() {
	v.notifs = v.app.Watch(v.model.Node,
		stores.NodeValueChanged,
		stores.NodeDescendantValueChanged,
	)
	go func() {
		for notif := range v.notifs {
			v.reaction(notif)
		}
	}()
}

func (v *IconEditorView) reaction(notif flux.NotifPayload) {
	defer close(notif.Done)
	v.icon = v.model.Node.Value.(*Icon)
	v.ReconcileBody()
}

func (v *IconEditorView) Unmount() {
	if v.notifs != nil {
		v.app.Delete(v.notifs)
		v.notifs = nil
	}
	v.Body.Unmount()
}

func (v *IconEditorView) render() vecty.Component {
	return elem.Div(
		elem.Div(
			elem.Image(
				style.MaxWidth("200px"),
				prop.Src(v.icon.Url.Value()),
			),
		),
		editors.NewStringEditorView(v.ctx, v.model.Node.Map["url"]),
	)
}
