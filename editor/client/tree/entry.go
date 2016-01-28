package tree

import (
	"fmt"

	"kego.io/editor"
	"kego.io/kerr"
)

// Entry items are nodes. Each branch inside a source branch are entry.
type entry struct {
	*Branch
	*editor.Node
	name  string
	index int
}

var _ BranchInterface = (*entry)(nil)
var _ Editable = (*entry)(nil)

func addEntryFromBytes(tree *Tree, name string, bytes []byte, parentBranch BranchInterface, parentEditor editor.EditorInterface) (*entry, error) {

	n, err := editor.UnmarshalNode(tree.ctx, bytes)
	if err != nil {
		return nil, kerr.Wrap("EAMNRBNRCE", err)
	}

	e, err := addEntry(name, -1, n, parentBranch, parentEditor)
	if err != nil {
		return nil, kerr.Wrap("FMIAASWRPY", err)
	}
	return e, nil

}

func addEntry(name string, index int, node *editor.Node, parentBranch BranchInterface, parentEditor editor.EditorInterface) (*entry, error) {

	ed := node.Editor()
	if parentEditor != nil && !parentEditor.AddChildTreeEntry(ed) {
		return nil, nil
	}

	e := &entry{name: name, index: index, Node: node}
	e.Branch = NewBranch(e, parentBranch)

	if index > -1 {
		e.setLabel(fmt.Sprint("[", index, "]"))
	} else {
		e.setLabel(name)
	}

	parentBranch.Append(e)

	if err := ed.Initialize(e.tree.ctx, e, editor.Page, e.tree.Fail); err != nil {
		return nil, kerr.Wrap("PMLOGADEVK", err)
	}

	e.ListenForEditorChanges(ed.Listen().Ch)

	if err := addEntryChildren(node, e, ed); err != nil {
		return nil, kerr.Wrap("UPRWSRECVR", err)
	}

	e.close()

	return e, nil
}

func addEntryChildren(parentNode *editor.Node, parentBranch BranchInterface, parentEditor editor.EditorInterface) error {

	if parentNode == nil {
		return nil
	}
	switch parentNode.Type.Native.Value() {
	case "array":
		for i, v := range parentNode.Array {
			child := v.(*editor.Node)
			if _, err := addEntry("", i, child, parentBranch, parentEditor); err != nil {
				return kerr.Wrap("IOXSWBQDXH", err)
			}
		}
	case "map":
		for name, v := range parentNode.Map {
			child := v.(*editor.Node)
			if _, err := addEntry(name, -1, child, parentBranch, parentEditor); err != nil {
				return kerr.Wrap("YVTQCADGJF", err)
			}
		}
	case "object":
		for name, v := range parentNode.Map {
			child := v.(*editor.Node)
			if child.Missing {
				continue
			}
			if _, err := addEntry(name, -1, child, parentBranch, parentEditor); err != nil {
				return kerr.Wrap("SIBWLRIXRG", err)
			}
		}
	}
	return nil
}

func (parent *holder) addTypes(types map[string][]byte) error {
	for name, bytes := range types {
		if _, err := addEntryFromBytes(parent.tree, name, bytes, parent, parent.editor); err != nil {
			return kerr.Wrap("CBLBCUUJFH", err)
		}
	}
	return nil
}
