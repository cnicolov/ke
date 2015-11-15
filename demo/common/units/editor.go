package units

import "kego.io/editor"

func (i *Rectangle) GetEditor(n *editor.Node) editor.Editor {
	return &RectangleEditor{Rectangle: n.Value.(*Rectangle), Node: n, Common: &editor.Common{}}
}

var _ editor.Editable = (*Rectangle)(nil)

type RectangleEditor struct {
	*Rectangle
	*editor.Common
	*editor.Node
	height *editor.NumberEditor
	width  *editor.NumberEditor
}

var _ editor.Editor = (*RectangleEditor)(nil)

func (e *RectangleEditor) Layout() editor.Layout {
	return editor.Inline
}

func (e *RectangleEditor) Initialize(holder editor.Holder, layout editor.Layout, path string, aliases map[string]string) error {

	e.Common.Initialize(holder, layout, path, aliases)

	e.height = editor.NewNumberEditor(e.Node.Map["height"])
	e.height.Initialize(holder, editor.Inline, path, aliases)
	e.height.Style().Set("width", "50%")

	e.width = editor.NewNumberEditor(e.Node.Map["width"])
	e.width.Initialize(holder, editor.Inline, path, aliases)
	e.width.Style().Set("width", "50%")

	e.Editors = append(e.Editors, e.height)
	e.Editors = append(e.Editors, e.width)

	e.AppendChild(e.height)
	e.AppendChild(e.width)

	go func() {
		width := e.width.Listen().Ch
		height := e.height.Listen().Ch
		for {
			select {
			case ne := <-width:
				e.Width.Set(int(ne.(*editor.NumberEditor).ValueNumber))
			case ne := <-height:
				e.Height.Set(int(ne.(*editor.NumberEditor).ValueNumber))
			}
			e.Notify(e)
		}
	}()

	return nil
}

func (e *RectangleEditor) Focus() {
	e.height.Focus()
}

func (e *RectangleEditor) Value() interface{} {
	return e.Rectangle
}
