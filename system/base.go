package system

type Object interface {
	GetBase() *Base
}

func (b *Base) GetBase() *Base {
	if b == nil {
		return &Base{}
	}
	return b
}

// SetContext satisfies the json.Contexter interface, which allows the json unmarshal
// function to store the unmarshal context in every object.
func (o *Base) SetContext(path string, imports map[string]string) {
	o.Context = &Context{Package: path, Imports: imports}
}

type Ruler interface {
	GetRules() []Rule
	RulesApply() rulesApplication
}

func (o *Base) GetRules() []Rule {
	return []Rule{}
}

func (o *Base) RulesApply() rulesApplication {
	if o.Type.Value == "kego.io/system:type" {
		return RULES_APPLY_TO_TYPES
	} else if o.Type.Type[0:0] == "@" {
		return RULES_APPLY_TO_TYPES
	}
	return RULES_APPLY_TO_OBJECTS
}

type rulesApplication string

const (
	RULES_APPLY_TO_TYPES   rulesApplication = "types"
	RULES_APPLY_TO_OBJECTS                  = "objects"
)
