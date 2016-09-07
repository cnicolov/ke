package models

import (
	"context"

	"reflect"

	"github.com/davelondon/kerr"
	"kego.io/context/jsonctx"
	"kego.io/editor/client/clientctx"
	"kego.io/editor/client/editable"
	"kego.io/json"
	"kego.io/system"
	"kego.io/system/node"
)

type EditorModel struct {
	Node    *node.Node
	Deleted bool
}

func NewEditor(n *node.Node) *EditorModel {
	return &EditorModel{Node: n}
}

func GetEditable(ctx context.Context, node *node.Node, embed *system.Reference) (editable.Editable, error) {
	ed, err := getEditable(ctx, node, embed, json.J_OBJECT)
	if err != nil {
		return kerr.Wrap("LCFUAFAOQU", err)
	}
	return ed.(editable.Editable)
}

func GetEditableMap(ctx context.Context, node *node.Node, embed *system.Reference) (editable.EditableMap, error) {
	ed, err := getEditable(ctx, node, embed, json.J_MAP)
	if err != nil {
		return kerr.Wrap("XOYXNHCLBY", err)
	}
	return ed.(editable.EditableMap)
}

func GetEditableArray(ctx context.Context, node *node.Node, embed *system.Reference) (editable.EditableArray, error) {
	ed, err := getEditable(ctx, node, embed, json.J_ARRAY)
	if err != nil {
		return kerr.Wrap("WAGWUFNDVQ", err)
	}
	return ed.(editable.EditableArray)
}

func getEditable(ctx context.Context, node *node.Node, embed *system.Reference, editorType json.Type) (interface{}, error) {

	if node == nil || node.Null || node.Missing {
		return nil, nil
	}

	// If we're after the editor for the type of the whole struct (not an
	// embedded struct), or the provided embedded type is nil...
	if embed == nil || *node.Type.Id == *embed {

		if ed := getEditableMulti(node.Value, editorType); ed != nil {
			return ed, nil
		}

		editors := clientctx.FromContext(ctx)

		// Don't do this. Implement the Editable interface instead. We can't do
		// this for system types so we use this method instead.
		if e, ok := editors.Get(node.Type.Id.Package, node.Type.Id.Name, editorType); ok {
			return e, nil
		}

		if node.JsonType != "" {
			if e, ok := editors.Get(string(node.JsonType), "", editorType); ok {
				return e, nil
			}
		}

		return nil, nil
	}

	jcache := jsonctx.FromContext(ctx)
	t, ok := jcache.GetType(embed.Package, embed.Name)
	if !ok {
		return nil, kerr.New("DGWDERFPVV", "Can't find %s in jsonctx", embed.String())
	}

	v := node.Val
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	var field reflect.Value
	for i := 0; i < v.Type().NumField(); i++ {
		f := v.Type().Field(i)
		if f.Anonymous && f.Type == t {
			field = v.Field(i)
			break
		}
	}
	if field == (reflect.Value{}) {
		return nil, kerr.New("UDBOWYUBER", "Can't find %s field in struct", t)
	}

	if ed := getEditableMulti(field.Interface(), editorType); ed != nil {
		return ed, nil
	}

	editors := clientctx.FromContext(ctx)

	// Don't do this. Implement the Editable interface instead. We can't do this
	// for system types so we use this method instead.
	if e, ok := editors.Get(embed.Package, embed.Name, editorType); ok {
		return e, nil
	}

	return nil, nil

}

func getEditableMulti(v interface{}, t json.Type) interface{} {
	switch t {
	case json.J_OBJECT:
		if ed, ok := v.(editable.Editable); ok {
			return ed, nil
		}
	case json.J_ARRAY:
		if ed, ok := v.(editable.EditableArray); ok {
			return ed, nil
		}
	case json.J_MAP:
		if ed, ok := v.(editable.EditableMap); ok {
			return ed, nil
		}
	}
	return nil
}
