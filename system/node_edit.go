package system

import "kego.io/json"

func (n *Node) GetEditor() Editor {
	switch n.JsonType {
	case json.J_STRING, json.J_BOOL, json.J_NUMBER:
		return &NodeValueEditor{Node: n}
	case json.J_MAP:
		return &NodeMapEditor{Node: n}
	case json.J_OBJECT:
		return &NodeObjectEditor{Node: n}
	}
	return nil
}
