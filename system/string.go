package system

import (
	"fmt"
	"regexp"
	"strconv"

	"kego.io/json"
	"kego.io/kerr"
)

type String struct {
	Value  string
	Exists bool
}

func NewString(s string) String {
	return String{Value: s, Exists: true}
}

func (r *StringRule) Validate(path string, aliases map[string]string) (ok bool, message string, err error) {
	if r.MaxLength.Exists && r.MinLength.Exists {
		if r.MaxLength.Value < r.MinLength.Value {
			return false, fmt.Sprintf("MaxLength %d must not be less than MinLength %d", r.MaxLength.Value, r.MinLength.Value), nil
		}
	}
	if r.Pattern.Exists {
		if _, err := regexp.Compile(r.Pattern.Value); err != nil {
			return false, fmt.Sprintf("Pattern: regex does not compile: %s", r.Pattern.Value), nil
		}
	}
	return true, "", nil
}

func (r *StringRule) Enforce(data interface{}, path string, aliases map[string]string) (bool, string, error) {
	s, ok := data.(String)
	if !ok {
		return false, "", kerr.New("SXFBXGQSEA", nil, "data %T should be system.String.", data)
	}

	// TODO: This restricts the value to one of several built-in formats
	// TODO: Format String
	if r.Format.Exists {
		panic("TODO: Format rule implementation")
	}

	// This is a regex to match the value to
	// Pattern String
	if r.Pattern.Exists {
		if !s.Exists && !r.Optional {
			return false, "Pattern: value must exist", nil
		}
		if s.Exists {
			reg, err := regexp.Compile(r.Pattern.Value)
			if err != nil {
				return false, fmt.Sprintf("Pattern: regex does not compile: %s", r.Pattern.Value), nil
			}
			if !reg.Match([]byte(s.Value)) {
				return false, fmt.Sprintf("Pattern: value must match %s", r.Pattern.Value), nil
			}
		}
	}

	// This is a string that the value must match
	// Equal String
	if r.Equal.Exists {
		if !s.Exists && !r.Optional {
			return false, "Equal: value must exist", nil
		}
		if s.Exists && s.Value != r.Equal.Value {
			return false, fmt.Sprintf("Equal: value must equal '%s'", r.Equal.Value), nil
		}
	}

	// The value must be longer or equal to the provided minimum length
	// MinLength Int
	if r.MinLength.Exists {
		if !s.Exists && !r.Optional {
			return false, "MinLength: value must exist", nil
		}
		if s.Exists && len(s.Value) < r.MinLength.Value {
			return false, fmt.Sprintf("MinLength: length must not be less than %d", r.MinLength.Value), nil
		}
	}

	// The value must be shorter or equal to the provided maximum length
	// MaxLength Int
	if r.MaxLength.Exists {
		if !s.Exists && !r.Optional {
			return false, "MaxLength: value must exist", nil
		}
		if s.Exists && len(s.Value) > r.MaxLength.Value {
			return false, fmt.Sprintf("MaxLength: length must not be greater than %d", r.MaxLength.Value), nil
		}
	}

	// The value of this string is restricted to one of the provided values
	// Enum []string
	if len(r.Enum) > 0 {
		if !s.Exists && !r.Optional {
			return false, "Enum: value must exist", nil
		}
		if s.Exists {
			found := false
			for _, e := range r.Enum {
				if e == s.Value {
					found = true
				}
			}
			if !found {
				return false, fmt.Sprintf("Enum: value must be one of: %v", r.Enum), nil
			}
		}
	}

	return true, "", nil
}

func (out *String) Unpack(in json.Unpackable) error {
	if in == nil || in.UpType() == json.J_NULL {
		out.Exists = false
		out.Value = ""
		return nil
	}
	if in.UpType() != json.J_STRING {
		return kerr.New("IXASCXOPMG", nil, "Can't unpack %s into system.String", in.UpType())
	}
	out.Exists = true
	out.Value = in.UpString()
	return nil
}

var _ json.Unpacker = (*String)(nil)

func (s String) MarshalJSON() ([]byte, error) {
	if !s.Exists {
		return []byte("null"), nil
	}
	return []byte(strconv.Quote(s.Value)), nil
}

var _ json.Marshaler = (*String)(nil)

func (s *String) String() string {
	if !s.Exists {
		return ""
	}
	return s.Value
}

type NativeString interface {
	NativeString() (string, bool)
}

func (s String) NativeString() (value string, exists bool) {
	return s.Value, s.Exists
}

// We satisfy the json.EmptyAware interface to allow intelligent omission of
// empty values when marshalling
func (s String) Empty() bool {
	return !s.Exists
}

var _ json.EmptyAware = (*String)(nil)
