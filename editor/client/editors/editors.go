package editors

import (
	"bytes"
	"math/rand"

	"context"

	"github.com/davelondon/vecty"
	"github.com/davelondon/vecty/elem"
	"github.com/davelondon/vecty/prop"
	"kego.io/editor/client/clientctx"
	"kego.io/editor/client/models"
	"kego.io/json"
	"kego.io/system"
	"kego.io/system/node"
)

func Register(ctx context.Context) {
	// Don't do this. Implement the Editable interface instead. We can't do this
	// for system types so we use this method instead.
	editors := clientctx.FromContext(ctx)

	editors.Set("string", "", json.J_OBJECT, new(StringEditor))
	editors.Set("kego.io/json", "string", json.J_OBJECT, new(StringEditor))
	editors.Set("kego.io/system", "string", json.J_OBJECT, new(StringEditor))

	editors.Set("number", "", json.J_OBJECT, new(NumberEditor))
	editors.Set("kego.io/json", "number", json.J_OBJECT, new(NumberEditor))
	editors.Set("kego.io/system", "number", json.J_OBJECT, new(NumberEditor))
	editors.Set("kego.io/system", "int", json.J_OBJECT, new(NumberEditor))

	editors.Set("bool", "", json.J_OBJECT, new(BoolEditor))
	editors.Set("kego.io/json", "bool", json.J_OBJECT, new(BoolEditor))
	editors.Set("kego.io/system", "bool", json.J_OBJECT, new(BoolEditor))

	editors.Set("kego.io/system", "object", json.J_OBJECT, new(ObjectEditor))
	editors.Set("kego.io/system", "tag", json.J_ARRAY, new(TagEditor))
}

func helpBlock(ctx context.Context, n *node.Node) vecty.Markup {
	if n.Rule == nil {
		return vecty.List{}
	}
	description := n.Rule.Interface.(system.ObjectInterface).GetObject(ctx).Description
	if description == "" {
		return vecty.List{}
	}
	return elem.Paragraph(
		prop.Class("help-block"),
		vecty.Text(description),
	)
}

func errorBlock(ctx context.Context, m *models.NodeModel) vecty.Markup {
	if !m.Invalid {
		return vecty.List{}
	}

	errors := vecty.List{}
	for _, e := range m.Errors {
		errors = append(errors, elem.ListItem(vecty.Text(e.Description)))
	}
	return elem.Paragraph(
		prop.Class("help-block text-danger"),
		elem.UnorderedList(errors),
	)
}

func randomId() string {
	randInt := func(min int, max int) int {
		return min + rand.Intn(max-min)
	}
	var result bytes.Buffer
	var temp string
	for i := 0; i < 20; {
		if string(randInt(65, 90)) != temp {
			temp = string(randInt(65, 90))
			result.WriteString(temp)
			i++
		}
	}
	return result.String()
}
