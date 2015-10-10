package system

import "kego.io/json"

var _ json.Unpackable = (*Node)(nil)

func (n *Node) UpType() int {
	if n == nil {
		return json.J_NULL
	}
	return n.JsonType
}
func (n *Node) UpNumber() float64 {
	return n.ValueNumber
}
func (n *Node) UpString() string {
	return n.ValueString
}
func (n *Node) UpBool() bool {
	return n.ValueBool
}
func (n *Node) UpArray() []json.Unpackable {
	out := []json.Unpackable{}
	for _, v := range n.Array {
		out = append(out, v)
	}
	return out
}
func (n *Node) UpMap() map[string]json.Unpackable {
	out := map[string]json.Unpackable{}
	for n, v := range n.Map {
		out[n] = v
	}
	for n, v := range n.Fields {
		if !v.Missing {
			out[n] = v
		}
	}
	return out
}
func (n *Node) UpInterface() interface{} {
	return n.Value
}