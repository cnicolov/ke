package selectors

import (
	"log"
	"strconv"

	"fmt"

	"kego.io/json"
	"kego.io/system"
)

func nodeIsMemberOfHaystack(needle *system.Node, haystack map[*system.Node]*system.Node) bool {
	_, ok := haystack[needle]
	return ok
}

func nodeIsMemberOfList(needle *system.Node, haystack []*system.Node) bool {
	for _, element := range haystack {
		if element == needle {
			return true
		}
	}
	return false
}

func (p *Parser) appendAncestorsToHaystack(n *system.Node, haystack map[*system.Node]*system.Node) {
	if !p.isRoot(n) {
		haystack[n.Parent] = n.Parent
		p.appendAncestorsToHaystack(n.Parent, haystack)
	}
}

func (p *Parser) nodeIsChildOfHaystackMember(needle *system.Node, haystack map[*system.Node]*system.Node) bool {
	if nodeIsMemberOfHaystack(needle, haystack) {
		return true
	}
	if p.isRoot(needle) {
		return false
	}
	return p.nodeIsChildOfHaystackMember(needle.Parent, haystack)
}

func (p *Parser) parents(lhs []*system.Node, rhs []*system.Node) []*system.Node {
	var results []*system.Node

	lhsHaystack := getHaystackFromNodeList(lhs)

	for _, element := range rhs {
		if !p.isRoot(element) && nodeIsMemberOfHaystack(element.Parent, lhsHaystack) {
			results = append(results, element)
		}
	}

	return results
}
func (p *Parser) isRoot(n *system.Node) bool {
	if n.Parent == nil {
		return true
	}
	if n == p.root {
		return true
	}
	return false
}
func (p *Parser) getKey(n *system.Node) string {
	if p.isRoot(n) {
		return ""
	}
	return n.Key
}
func (p *Parser) getIndex(n *system.Node) int {
	if p.isRoot(n) {
		return -1
	}
	return n.Index + 1
}
func (p *Parser) getSiblings(n *system.Node) int {
	if p.isRoot(n) {
		return 0
	}
	return n.ArraySiblings
}

func (p *Parser) ancestors(lhs []*system.Node, rhs []*system.Node) []*system.Node {
	var results []*system.Node
	haystack := getHaystackFromNodeList(lhs)

	for _, element := range rhs {
		if p.nodeIsChildOfHaystackMember(element, haystack) {
			results = append(results, element)
		}
	}

	return results
}

func (p *Parser) siblings(lhs []*system.Node, rhs []*system.Node) []*system.Node {
	var results []*system.Node
	parents := make(map[*system.Node]*system.Node, len(lhs))

	for _, element := range lhs {
		if !p.isRoot(element) {
			parents[element.Parent] = element.Parent
		}
	}

	for _, element := range rhs {
		if !p.isRoot(element) && nodeIsMemberOfHaystack(element.Parent, parents) {
			results = append(results, element)
		}
	}

	return results
}

func getFloat64(in interface{}) float64 {
	if as_node, ok := in.(*system.Node); ok {
		return as_node.ValueNumber
	}
	as_nativeNum, ok := in.(system.NativeNumber)
	if ok {
		num, exists := as_nativeNum.NativeNumber()
		if !exists {
			return 0.0
		}
		return num
	}
	as_float, ok := in.(float64)
	if ok {
		return as_float
	}
	as_int, ok := in.(int64)
	if ok {
		value := float64(as_int)
		return value
	}
	as_string, ok := in.(string)
	if ok {
		parsed_float_string, err := strconv.ParseFloat(as_string, 64)
		if err == nil {
			value := parsed_float_string
			return value
		}
		parsed_int_string, err := strconv.ParseInt(as_string, 10, 32)
		if err == nil {
			value := float64(parsed_int_string)
			return value
		}
	}
	panic(fmt.Sprintf("ERR: %T", in))
	result := float64(-1)
	log.Print("Error transforming ", in, " into Float64")
	return result
}

func getInt32(in interface{}) int32 {
	value := int32(getFloat64(in))
	if value == -1 {
		log.Print("Error transforming ", in, " into Int32")
	}
	return value
}

func isNull(in interface{}) bool {
	if in == nil {
		return true
	}
	if e, ok := in.(json.EmptyAware); ok {
		return e.Empty()
	}
	if n, ok := in.(*system.Node); ok {
		return n.Null || n.Missing
	}
	return false
}

func getBool(in interface{}) bool {
	as_node, ok := in.(*system.Node)
	if ok {
		return as_node.ValueBool
	}
	as_nativeBool, ok := in.(system.NativeBool)
	if ok {
		b, exists := as_nativeBool.NativeBool()
		if !exists {
			return false
		}
		return b
	}
	as_bool, ok := in.(bool)
	if ok {
		return as_bool
	}
	return false
}

func getJsonString(in interface{}) string {
	as_node, ok := in.(*system.Node)
	if ok {
		return as_node.ValueString
	}
	as_nativeStr, ok := in.(system.NativeString)
	if ok {
		s, exists := as_nativeStr.NativeString()
		if !exists {
			return ""
		}
		return s
	}
	as_string, ok := in.(string)
	if ok {
		return as_string
	}
	marshaled_result, err := json.Marshal(in)
	if err != nil {
		log.Print("Error transforming ", in, " into JSON string")
	}
	result := string(marshaled_result)
	return result
}

func exprElementIsTruthy(e exprElement) (bool, error) {
	switch e.typ {
	case json.J_STRING:
		return len(e.value.(string)) > 0, nil
	case json.J_NUMBER:
		return e.value.(float64) > 0, nil
	case json.J_OBJECT:
		return true, nil
	case json.J_MAP:
		return true, nil
	case json.J_ARRAY:
		return true, nil
	case json.J_BOOL:
		return e.value.(bool), nil
	case json.J_NULL:
		return false, nil
	default:
		return false, nil
	}
}

func exprElementsMatch(lhs exprElement, rhs exprElement) bool {
	return lhs.typ == rhs.typ
}

func getHaystackFromNodeList(nodes []*system.Node) map[*system.Node]*system.Node {
	hashmap := make(map[*system.Node]*system.Node, len(nodes))
	for _, node := range nodes {
		hashmap[node] = node
	}
	return hashmap
}
