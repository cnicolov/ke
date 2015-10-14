package system

import (
	"honnef.co/go/js/dom"
	"kego.io/editor/mdl"
	"kego.io/kerr"
)

type NodeObjectEditor struct {
	*Node
	path        string
	aliases     map[string]string
	panel       *dom.HTMLDivElement
	input       *dom.HTMLInputElement
	initialized bool
}

var _ Editor = (*NodeObjectEditor)(nil)

func (e *NodeObjectEditor) Initialized() bool {
	return e.initialized
}

func (e *NodeObjectEditor) Initialize(panel *dom.HTMLDivElement, path string, aliases map[string]string) error {

	e.panel = panel
	e.path = path
	e.aliases = aliases

	table := mdl.Table()

	names := table.Column("name")
	origins := table.Column("origin")
	holds := table.Column("holds")
	values := table.Column("value")

	for name, field := range e.Fields {

		names.Cell(name)

		origin, err := field.Origin.ValueContext(e.path, e.aliases)
		if err != nil {
			return kerr.New("ACQLJXWYQX", err, "ValueContext")
		}
		origins.Cell(origin)

		hold, err := field.Rule.ParentType.Id.ValueContext(e.path, e.aliases)
		if err != nil {
			return kerr.New("XDKOSFJVQV", err, "ValueContext")
		}
		holds.Cell(hold)

		if field.Missing || field.Null {
			values.Cell("")
		} else {
			value, err := field.Type.Id.ValueContext(e.path, e.aliases)
			if err != nil {
				return kerr.New("RWHEKAOPHQ", err, "ValueContext")
			}
			values.Cell(value)
		}

	}
	e.panel.AppendChild(table.Build())

	e.initialized = true
	return nil
}

func (e *NodeObjectEditor) Show() {
	e.panel.Style().Set("display", "block")
}
func (e *NodeObjectEditor) Hide() {
	e.panel.Style().Set("display", "none")
}

func (e *NodeObjectEditor) Update() {
	//if e.Exists {
	//	e.node.ValueString = e.input.Value
	//}
}
