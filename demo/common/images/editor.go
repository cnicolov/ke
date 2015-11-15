package images

import (
	"kego.io/editor"
	"kego.io/editor/mdl"
)

func (i *Icon) GetEditor(n *editor.Node) editor.Editor {
	return &IconEditor{Icon: n.Value.(*Icon), Node: n, Common: &editor.Common{}}
}

var _ editor.Editable = (*Icon)(nil)

type IconEditor struct {
	*Icon
	*editor.Common
	*editor.Node
	image *mdl.ImageStruct
	url   *editor.StringEditor
}

var _ editor.Editor = (*IconEditor)(nil)

func (e *IconEditor) Layout() editor.Layout {
	return editor.Block
}

func (e *IconEditor) Initialize(holder editor.Holder, layout editor.Layout, path string, aliases map[string]string) error {

	e.Common.Initialize(holder, layout, path, aliases)

	e.image = mdl.Image(e.Url.Value())
	e.AppendChild(e.image)

	e.url = editor.NewStringEditor(e.Node.Map["url"])
	e.url.Initialize(holder, editor.Block, path, aliases)
	e.Editors = append(e.Editors, e.url)
	e.AppendChild(e.url)

	go func() {
		for se := range e.url.Listen().Ch {
			e.update(se.(*editor.StringEditor).ValueString)
			e.Notify(e)
		}
	}()

	e.update(e.Url.Value())

	return nil
}

func (e *IconEditor) update(url string) {
	e.Url.Set(url)
	e.image.Src = url
	e.image.Visibility(url != "")
}

func (e *IconEditor) AddChildTreeEntry(child editor.Editor) bool {
	return false
}

func (e *IconEditor) Focus() {
	e.url.Focus()
}

func (e *IconEditor) Value() interface{} {
	return e.Icon
}
