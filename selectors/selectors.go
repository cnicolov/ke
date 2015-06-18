package selectors // import "kego.io/selectors"

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"kego.io/kerr"
	"kego.io/system"
)

// for tests
type Image interface{}

type Parser struct {
	Data    *Element
	nodes   []*node
	path    string
	imports map[string]string
}

func CreateParser(element *Element, path string, imports map[string]string) (*Parser, error) {
	parser := Parser{Data: element, nodes: nil, path: path, imports: imports}
	if err := parser.mapDocument(); err != nil {
		return nil, kerr.New("SWLNNMSUTM", err, "selectors.CreateParser", "parser.mapDocument")
	}
	return &parser, nil
}

func (p *Parser) evaluateSelector(selector string) ([]*node, error) {
	tokens, err := lex(selector, selectorScanner)
	if err != nil {
		return nil, err
	}

	//for _, t := range tokens {
	//	fmt.Println("Token:", t)
	//}

	//for _, n := range p.nodes {
	//	fmt.Printf("Node: %#v\n", n)
	//}

	nodes, err := p.selectorProduction(tokens, p.nodes, 1)
	if err != nil {
		return nil, err
	}

	logger.Print(len(nodes), " matches found")

	return nodes, nil
}

func (p *Parser) GetElements(selector string) ([]*Element, error) {
	nodes, err := p.evaluateSelector(selector)
	if err != nil {
		return nil, err
	}

	var results = make([]*Element, 0, len(nodes))
	for _, node := range nodes {
		results = append(
			results,
			node.element,
		)
	}
	return results, nil
}

func (p *Parser) GetValues(selector string) ([]interface{}, error) {
	nodes, err := p.evaluateSelector(selector)
	if err != nil {
		return nil, err
	}

	var results = make([]interface{}, 0, len(nodes))
	for _, node := range nodes {
		results = append(
			results,
			node.value,
		)
	}

	return results, nil
}

func (p *Parser) selectorProduction(tokens []*token, documentMap []*node, recursionDepth int) ([]*node, error) {
	var results []*node
	var matched bool
	var value interface{}
	var validator func(*node) (bool, error)
	var validators = make([]func(*node) (bool, error), 0, 10)
	logger.Print("selectorProduction(", recursionDepth, ") starting with ", tokens[0], " - ", len(tokens), " tokens remaining.")

	_, matched, _ = p.peek(tokens, S_NATIVE_TYPE)
	if matched {
		value, tokens, _ = p.match(tokens, S_NATIVE_TYPE)
		validators = append(
			validators,
			p.nativeProduction(value),
		)
	}
	_, matched, _ = p.peek(tokens, S_KEGO_TYPE)
	if matched {
		value, tokens, _ = p.match(tokens, S_KEGO_TYPE)
		validators = append(
			validators,
			p.kegoProduction(value),
		)
	}
	_, matched, _ = p.peek(tokens, S_IDENTIFIER)
	if matched {
		value, tokens, _ = p.match(tokens, S_IDENTIFIER)
		validators = append(
			validators,
			p.keyProduction(value),
		)
	}
	_, matched, _ = p.peek(tokens, S_PCLASS)
	if matched {
		value, tokens, _ = p.match(tokens, S_PCLASS)
		validators = append(
			validators,
			p.pclassProduction(value),
		)
	}
	_, matched, _ = p.peek(tokens, S_NTH_FUNC)
	if matched {
		value, tokens, _ = p.match(tokens, S_NTH_FUNC)
		validator, tokens = p.nthChildProduction(value, tokens)
		validators = append(validators, validator)
	}
	_, matched, _ = p.peek(tokens, S_PCLASS_FUNC)
	if matched {
		value, tokens, _ = p.match(tokens, S_PCLASS_FUNC)
		validator, tokens = p.pclassFuncProduction(value, tokens, documentMap)
		validators = append(validators, validator)
	}
	result, matched, _ := p.peek(tokens, S_OPER)
	if matched && result.(string) == "*" {
		value, tokens, _ = p.match(tokens, S_OPER)
		validator = p.universalProduction(value)
		validators = append(validators, validator)
	}

	if len(validators) < 1 {
		return nil, errors.New("No selector recognized")
	}

	results, err := p.matchNodes(validators, documentMap)
	if err != nil {
		return nil, err
	}
	logger.Print("Applying ", len(validators), " validators to document resulted in ", len(results), " matches")

	_, matched, _ = p.peek(tokens, S_OPER)
	if matched {
		value, tokens, _ = p.match(tokens, S_OPER)
		logger.Print("Recursing selectorProduction(", recursionDepth, ") via operator ", value, " starting with ", tokens[0], " ;", len(tokens), " tokens remaining")
		logger.IncreaseDepth()
		rvals, err := p.selectorProduction(tokens, documentMap, recursionDepth+1)
		logger.DecreaseDepth()
		logger.Print("Recursion completed; returned control to selectorProduction(", recursionDepth, ") via operator ", value, " with ", len(rvals), " matches.")
		if err != nil {
			return nil, err
		}
		switch value {
		case ",":
			logger.Print("(", recursionDepth, ") Operator ',': ", len(results), " => ", len(results)+len(rvals))
			for _, val := range rvals {
				// TODO: This is quite slow
				// it seems like it's probably quite easy to expand
				// the list just once
				results = append(results, val)
			}
		case ">":
			originalLength := len(results)
			results = parents(results, rvals)
			logger.Print("(", recursionDepth, ") Operator '>': ", originalLength, " => ", len(results))
		case "~":
			originalLength := len(results)
			results = siblings(results, rvals)
			logger.Print("(", recursionDepth, ") Operator '~': ", originalLength, " => ", len(results))
		case " ":
			originalLength := len(results)
			results = ancestors(results, rvals)
			logger.Print("(", recursionDepth, ") Operator ' ': ", originalLength, " => ", len(results))
		default:
			return nil, errors.New("Unrecognized operator")
		}
	} else if len(tokens) > 0 {
		logger.Print("Recursing selectorProduction(", recursionDepth, ") for excess tokens starting with ", tokens[0], " ;", len(tokens), " tokens remaining")
		logger.IncreaseDepth()
		rvals, err := p.selectorProduction(tokens, documentMap, recursionDepth+1)
		logger.DecreaseDepth()
		logger.Print("Recursion completed; returned control to selectorProduction(", recursionDepth, ") for excess tokens with ", len(rvals), " matches.")
		if err != nil {
			return nil, err
		}
		results = ancestors(results, rvals)
	}

	logger.Print("selectorProduction(", recursionDepth, ") returning ", len(results), " matches.")
	return results, nil
}

func (p *Parser) peek(tokens []*token, typ tokenType) (interface{}, bool, error) {
	if len(tokens) < 1 {
		return nil, false, errors.New("No more tokens")
	}
	if tokens[0].typ == typ {
		return tokens[0].val, true, nil
	}
	return nil, false, nil
}

func (p *Parser) match(tokens []*token, typ tokenType) (interface{}, []*token, error) {
	value, matched, _ := p.peek(tokens, typ)
	if !matched {
		return nil, tokens, errors.New("Match not successful")
	}
	_, tokens = tokens[0], tokens[1:]
	return value, tokens, nil
}

func (p *Parser) matchNodes(validators []func(*node) (bool, error), documentMap []*node) ([]*node, error) {
	var matches []*node
	nodeCount := 0
	if logger.Enabled {
		nodeCount = len(documentMap)
	}
	for idx, node := range documentMap {
		var passed = true
		if logger.Enabled {
			logger.SetPrefix("[Node ", idx, "/", nodeCount, "] ")
		}
		for _, validator := range validators {
			result, err := validator(node)
			if err != nil {
				return nil, kerr.New("GUOFXMISKX", err, "selectors.matchNodes", "validator")
			}
			if !result {
				passed = false
				break
			}
		}
		if passed {
			logger.Print("MATCHED: ", node)
			matches = append(matches, node)
		}
	}
	return matches, nil
}

func (p *Parser) nativeProduction(value interface{}) func(*node) (bool, error) {
	logger.Print("Creating nativeProduction validator ", value)
	return func(n *node) (bool, error) {
		logger.Print("nativeProduction ? ", n.native, " == ", value)
		return string(n.native) == value, nil
	}
}

func (p *Parser) kegoProduction(value interface{}) func(*node) (bool, error) {
	logger.Print("Creating kegoProduction validator ", value)
	return func(n *node) (bool, error) {
		tokenString := value.(string)
		kegoType := tokenString[1 : len(tokenString)-1]
		r, err := system.NewReferenceFromString(strconv.Quote(kegoType), p.path, p.imports)
		if err != nil {
			return false, kerr.New("RWDOYBBDVK", err, "selectors.kegoProduction", "NewReferenceFromString")
		}
		logger.Print("kegoProduction ? ", n.ktyperef.Value, " == ", r.Value)

		if n.ktyperef.Value == r.Value {
			return true, nil
		}

		for _, ref := range n.ktype.Is {
			if ref.Value == r.Value {
				return true, nil
			}
		}

		return false, nil
	}
}

func (p *Parser) keyProduction(value interface{}) func(*node) (bool, error) {
	logger.Print("Creating keyProduction validator ", value)
	return func(n *node) (bool, error) {
		logger.Print("keyProduction ? ", n.parent_key, " == ", value)
		if n.parent_key == "" {
			return false, nil
		}

		// must make sure the node exists
		if !nodeExists(n) {
			return false, nil
		}

		return string(n.parent_key) == value, nil
	}
}

func nodeExists(n *node) bool {
	switch n.native {
	case J_NULL:
		return false
	}
	return n.value != nil
	/*
		switch node.json.Rule.ParentType.Native.Value {
		case "string":
			if _, exists := node.json.Data.(system.NativeString).NativeString(); !exists {
				return false
			}
		case "number":
			if _, exists := node.json.Data.(system.NativeNumber).NativeNumber(); !exists {
				return false
			}
		case "bool":
			if _, exists := node.json.Data.(system.NativeBool).NativeBool(); !exists {
				return false
			}
		case "map", "array", "object":
			if node.json.Data == nil {
				return false
			}
		}
		return true*/
}

func (p *Parser) universalProduction(value interface{}) func(*node) (bool, error) {
	operator := value.(string)
	if operator == "*" {
		return func(n *node) (bool, error) {
			logger.Print("universalProduction ? true")
			return true, nil
		}
	} else {
		logger.Print("Error: Unexpected operator: ", operator)
		return func(n *node) (bool, error) {
			logger.Print("Asserting false due to failed universalProduction")
			return false, nil
		}
	}
}

func (p *Parser) pclassProduction(value interface{}) func(*node) (bool, error) {
	pclass := value.(string)
	logger.Print("Creating pclassProduction validator ", pclass)
	if pclass == "first-child" {
		return func(n *node) (bool, error) {
			logger.Print("pclassProduction first-child ? ", n.idx, " == 1")
			return n.idx == 1, nil
		}
	} else if pclass == "last-child" {
		return func(n *node) (bool, error) {
			logger.Print("pclassProduction last-child ? ", n.siblings, " > 0 AND ", n.idx, " == ", n.siblings)
			return n.siblings > 0 && n.idx == n.siblings, nil
		}
	} else if pclass == "only-child" {
		return func(n *node) (bool, error) {
			logger.Print("pclassProduction ony-child ? ", n.siblings, " == 1")
			return n.siblings == 1, nil
		}
	} else if pclass == "root" {
		return func(n *node) (bool, error) {
			logger.Print("pclassProduction root ? ", n.parent, " == nil")
			return n.parent == nil, nil
		}
	} else if pclass == "empty" {
		return func(n *node) (bool, error) {
			logger.Print("pclassProduction empty ? ", n.native, " == ", J_ARRAY, " AND ", len(n.value.(string)), " < 1")
			return n.native == J_ARRAY && len(n.value.(string)) < 1, nil
		}
	}
	logger.Print("Error: Unknown pclass: ", pclass)
	return func(n *node) (bool, error) {
		logger.Print("Asserting false due to failed pclassProduction")
		return false, nil
	}
}

func (p *Parser) nthChildProduction(value interface{}, tokens []*token) (func(*node) (bool, error), []*token) {
	nthChildRegexp := regexp.MustCompile(`^\s*\(\s*(?:([+\-]?)([0-9]*)n\s*(?:([+\-])\s*([0-9]))?|(odd|even)|([+\-]?[0-9]+))\s*\)`)
	args, tokens, _ := p.match(tokens, S_EXPR)
	var a int
	var b int
	var reverse bool = false

	pattern := nthChildRegexp.FindStringSubmatch(args.(string))

	logger.Print("Creating nthChildProduction validator ", pattern)
	if logger.Enabled {
		for idx, pat := range pattern {
			logger.Print("[", idx, "] ", pat)
		}
	}

	if pattern[5] != "" {
		a = 2
		if pattern[5] == "odd" {
			b = 1
		} else {
			b = 0
		}
	} else if pattern[6] != "" {
		a = 0
		b_parsed, _ := strconv.ParseInt(pattern[6], 10, 64)
		b = int(b_parsed)
	} else {
		// Expression like +n-3
		sign := "+"
		if pattern[1] != "" {
			sign = pattern[1]
		}
		coeff := "1"
		if pattern[2] != "" {
			coeff = pattern[2]
		}
		a_parsed, _ := strconv.ParseInt(coeff, 10, 64)
		a = int(a_parsed)
		if sign == "-" {
			a = -1 * a
		}

		if pattern[3] != "" {
			b_sign := pattern[3]
			b_parsed, _ := strconv.ParseInt(pattern[4], 10, 64)
			b = int(b_parsed)
			if b_sign == "-" {
				b = -1 * b
			}
		} else {
			b = 0
		}
	}

	if value.(string) == "nth-last-child" {
		reverse = true
	}

	return func(n *node) (bool, error) {
		var b_str string
		if b > 0 {
			b_str = "+"
		} else {
			b_str = "-"
		}
		logger.Print("nthChildProduction ", a, "n", b_str, b)
		logger.Print("nthChildProduction ? ", n.siblings, " == 0")
		if n.siblings == 0 {
			return false, nil
		}

		idx := n.idx - 1
		if reverse {
			idx = n.siblings - idx
		} else {
			idx++
		}

		logger.Print("nthChildProduction (continued-1) ? ", a, " == 0")
		if a == 0 {
			return b == idx, nil
		} else {
			logger.Print("nthChildProduction (continued-2) ? ", idx-b%a, " == 0 AND ", idx*a+b, " >= 0")
			return ((idx-b)%a) == 0 && (idx*a+b) >= 0, nil
		}
	}, tokens
}

func (p *Parser) pclassFuncProduction(value interface{}, tokens []*token, documentMap []*node) (func(*node) (bool, error), []*token) {
	sargs, tokens, _ := p.match(tokens, S_EXPR)
	pclass := value.(string)

	logger.Print("Creating pclassFuncProduction validator ", pclass)

	if pclass == "expr" {
		tokens, err := lex(sargs.(string), expressionScanner)
		if err != nil {
			panic(err)
		}
		var tokens_to_return []*token
		return func(n *node) (bool, error) {
			result := p.parseExpression(tokens, n)
			logger.Print("pclassFuncProduction expr ? ", result)
			return exprElementIsTruthy(result)
		}, tokens_to_return
	}

	lexString := sargs.(string)[1 : len(sargs.(string))-1]
	args, _ := lex(lexString, selectorScanner)

	logger.Print("pclassFuncProduction lex results for [", lexString, "]: (follow)")
	if logger.Enabled {
		for i, arg := range args {
			logger.Print("[", i, "]: ", arg)
		}
	}

	if pclass == "has" {
		return func(n *node) (bool, error) {
			newMap, err := p.getFlooredDocumentMap(n)
			if err != nil {
				return false, kerr.New("HFVKFIUIUT", err, "selectors.pclassFuncProduction", "getFlooredDocumentMap")
			}
			logger.Print("pclassFuncProduction recursing into selectorProduction(-100) starting with ", args[0], "; ", len(args), " tokens remaining.")
			logger.IncreaseDepth()
			rvals, _ := p.selectorProduction(args, newMap, -100)
			logger.DecreaseDepth()
			logger.Print("pclassFuncProduction resursion completed with ", len(rvals), " results.")
			ancestors := make(map[*Element]*node, len(rvals))
			for _, n := range rvals {
				if n.parent != nil {
					ancestors[n.parent.element] = n.parent
				}
			}
			logger.Print("pclassFuncProduction has ? ", n, " ∈ ", getFormattedNodeMap(ancestors))
			return nodeIsMemberOfHaystack(n, ancestors), nil
		}, tokens
	} else if pclass == "contains" {
		return func(n *node) (bool, error) {
			logger.Print("pclassFuncProduction contains ? ", n.native, " == ", J_STRING, " AND ", strings.Count(n.value.(string), args[0].val.(string)), " > 0")
			return n.native == J_STRING && strings.Count(n.value.(string), args[0].val.(string)) > 0, nil
		}, tokens
	} else if pclass == "val" {
		return func(n *node) (bool, error) {
			lhsString := getJsonString(n.value)
			rhsString := getJsonString(args[0].val)
			logger.Print("pclassFuncProduction val ? ", lhsString, " == ", rhsString)
			return lhsString == rhsString, nil
		}, tokens
	}

	// If we didn't find a known pclass, do not match anything.
	logger.Print("Error: Unknown pclass: ", pclass)
	return func(n *node) (bool, error) {
		logger.Print("Asserting false due to failed pclassFuncProduction")
		return false, nil
	}, tokens
}
