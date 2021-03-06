// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Represents JSON data structure using native Go types: booleans, floats,
// strings, arrays, and maps.

package json

import (
	"bytes"
	"encoding"
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

	"context"

	"github.com/davelondon/kerr"
	"kego.io/context/envctx"
	"kego.io/context/jsonctx"
)

// UnmarshalPlain parses the JSON-encoded data and stores the result
// in the value pointed to by v.
//
// Unmarshal uses the inverse of the encodings that
// Marshal uses, allocating maps, slices, and pointers as necessary,
// with the following additional rules:
//
// To unmarshal JSON into a pointer, Unmarshal first handles the case of
// the JSON being the JSON literal null.  In that case, Unmarshal sets
// the pointer to nil.  Otherwise, Unmarshal unmarshals the JSON into
// the value pointed at by the pointer.  If the pointer is nil, Unmarshal
// allocates a new value for it to point to.
//
// To unmarshal JSON into a struct, Unmarshal matches incoming object
// keys to the keys used by Marshal (either the struct field name or its tag),
// preferring an exact match but also accepting a case-insensitive match.
//
// To unmarshal JSON into an interface value,
// Unmarshal stores one of these in the interface value:
//
//	bool, for JSON booleans
//	float64, for JSON numbers
//	string, for JSON strings
//	[]interface{}, for JSON arrays
//	map[string]interface{}, for JSON objects
//	nil for JSON null
//
// If a JSON value is not appropriate for a given target type,
// or if a JSON number overflows the target type, Unmarshal
// skips that field and completes the unmarshalling as best it can.
// If no more serious errors are encountered, Unmarshal returns
// an UnmarshalTypeError describing the earliest such error.
//
// The JSON null value unmarshals into an interface, map, pointer, or slice
// by setting that Go value to nil. Because null is often used in JSON to mean
// ``not present,'' unmarshaling a JSON null into any other Go type has no effect
// on the value and produces no error.
//
// When unmarshaling quoted strings, invalid UTF-8 or
// invalid UTF-16 surrogate pairs are not treated as an error.
// Instead, they are replaced by the Unicode replacement
// character U+FFFD.
//
func UnmarshalPlain(data []byte, v interface{}) error {
	// Check for well-formedness.
	// Avoids filling out half a data structure
	// before discovering a JSON syntax error.
	var d decodeState
	err := checkValid(data, &d.scan)
	if err != nil {
		return err
	}

	d.init(nil, data, false)

	return d.unmarshal(v)
}
func UnmarshalUntyped(ctx context.Context, data []byte, v interface{}) error {
	// Check for well-formedness.
	// Avoids filling out half a data structure
	// before discovering a JSON syntax error.
	var d decodeState
	err := checkValid(data, &d.scan)
	if err != nil {
		// ke: {"block": {"notest": true}}
		return kerr.Wrap("BACKEQJSYF", err)
	}
	d.init(ctx, data, false)
	err = d.unmarshal(v)
	if err != nil {
		return kerr.Wrap("YFLBEQIIUY", err)
	}
	return nil
}

type UnknownPackageError struct {
	kerr.Struct
	UnknownPackage string
}

type UnknownTypeError struct {
	kerr.Struct
	UnknownType string
}

func Unmarshal(ctx context.Context, data []byte, v *interface{}) error {
	// Check for well-formedness.
	// Avoids filling out half a data structure
	// before discovering a JSON syntax error.
	var d decodeState
	err := checkValid(data, &d.scan)
	err = d.getError(err)
	if err != nil {
		// ke: {"block": {"notest": true}}
		if err != nil {
			return kerr.Wrap("EBQHWRCBXE", err)
		}
	}

	d.init(ctx, data, true)
	err = d.unmarshalTyped(v)
	err = d.getError(err)
	if err != nil {
		return kerr.Wrap("GMSJBIQRQT", err)
	}
	return nil
}

// Unmarshaler is the interface implemented by objects
// that can unmarshal a JSON description of themselves.
// The input can be assumed to be a valid encoding of
// a JSON value. UnmarshalJSON must copy the JSON data
// if it wishes to retain the data after returning.
type Unmarshaler interface {
	UnmarshalJSON(ctx context.Context, in []byte) error
}

type Unpacker interface {
	Unpack(ctx context.Context, in Packed) error
}

// An UnmarshalTypeError describes a JSON value that was
// not appropriate for a value of a specific Go type.
type UnmarshalTypeError struct {
	Value string       // description of JSON value - "bool", "array", "number -5"
	Type  reflect.Type // type of Go value it could not be assigned to
}

func (e *UnmarshalTypeError) Error() string {
	return "json: cannot unmarshal " + e.Value + " into Go value of type " + e.Type.String()
}

// An UnmarshalFieldError describes a JSON object key that
// led to an unexported (and therefore unwritable) struct field.
// (No longer used; kept for compatibility.)
type UnmarshalFieldError struct {
	Key   string
	Type  reflect.Type
	Field reflect.StructField
}

func (e *UnmarshalFieldError) Error() string {
	return "json: cannot unmarshal object key " + strconv.Quote(e.Key) + " into unexported field " + e.Field.Name + " of type " + e.Type.String()
}

// An InvalidUnmarshalError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "json: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "json: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "json: Unmarshal(nil " + e.Type.String() + ")"
}

func (d *decodeState) unmarshal(v interface{}) error {
	rv := reflect.ValueOf(v)
	if err := d.unmarshalValue(rv); err != nil {
		return kerr.Wrap("BRLGHLJAKV", err)
	}
	return nil
}

func (d *decodeState) unmarshalTyped(v *interface{}) (err error) {
	rv := reflect.ValueOf(v)
	return d.unmarshalValue(rv)
}

func (d *decodeState) unmarshalValue(rv reflect.Value) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				// ke: {"block": {"notest": true}}
				panic(r)
			}
			err = r.(error)
			if err != nil {
				err = kerr.Wrap("RXTXRDOLVG", err)
			}
		}
	}()

	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		// This allows us to get the type of v for the
		// error message, even though we don't have v.
		var nilInterface interface{}
		var inputType reflect.Type
		if rv == (reflect.Value{}) {
			inputType = reflect.TypeOf(nilInterface)
		} else {
			inputType = rv.Type()
		}
		return &InvalidUnmarshalError{inputType}
	}

	d.scan.reset()
	// We decode rv not rv.Elem because the Unmarshaler interface
	// test must be applied at the top level of the value.
	d.value(rv)
	return d.savedError
}

// A NumberLiteral represents a JSON number literal.
type NumberLiteral string

// String returns the literal text of the number.
func (n NumberLiteral) String() string { return string(n) }

// Float64 returns the number as a float64.
func (n NumberLiteral) Float64() (float64, error) {
	return strconv.ParseFloat(string(n), 64)
}

// Int64 returns the number as an int64.
func (n NumberLiteral) Int64() (int64, error) {
	return strconv.ParseInt(string(n), 10, 64)
}

// decodeState represents the state while decoding a JSON value.
type decodeState struct {
	ctx            context.Context
	typed          bool
	data           []byte
	off            int // read offset in data
	scan           scanner
	nextscan       scanner // for calls to nextValue
	savedError     error
	useNumber      bool
	unknownType    string // have we encountered an unknown type?
	unknownPackage string // have we encountered an unknown package?
}

func (d *decodeState) getError(err error) error {
	if err != nil {
		// ke: {"block": {"notest": true}}
		return err
	}
	if d.unknownPackage != "" {
		return UnknownPackageError{
			Struct:         kerr.New("TBNIEVUCPL", "Unknown package %s", d.unknownPackage),
			UnknownPackage: d.unknownPackage,
		}
	}
	if d.unknownType != "" {
		return UnknownTypeError{
			Struct:      kerr.New("YEOQSWVFVH", "Unknown type %s", d.unknownType),
			UnknownType: d.unknownType,
		}
	}
	return nil
}

// backup saves the exact state of the decoder so we
// can remember where we were before scanning for the
// type json key, and restore afterwards.
func (d *decodeState) backup() decodeState {
	return decodeState{
		ctx:            d.ctx,
		typed:          d.typed,
		off:            d.off,
		savedError:     d.savedError,
		useNumber:      d.useNumber,
		scan:           d.scan.backup(),
		nextscan:       d.nextscan.backup(),
		unknownType:    d.unknownType,
		unknownPackage: d.unknownPackage,
	}
}

// restore restores the exact state of the decoder so
// we can remember where we were before scanning for
// the type json key, and restore afterwards.
func (d *decodeState) restore(from decodeState) {
	d.ctx = from.ctx
	d.typed = from.typed
	d.off = from.off
	d.savedError = from.savedError
	d.useNumber = from.useNumber
	d.scan = from.scan
	d.nextscan = from.nextscan
	d.unknownType = from.unknownType
	d.unknownPackage = from.unknownPackage
}

// errPhase is used for errors that should not happen unless
// there is a bug in the JSON decoder or something is editing
// the data slice while the decoder executes.
var errPhase = errors.New("JSON decoder out of sync - data changing underfoot?")

func (d *decodeState) init(ctx context.Context, data []byte, typed bool) *decodeState {
	d.ctx = ctx
	d.typed = typed
	d.data = data
	d.off = 0
	d.savedError = nil
	d.unknownType = ""
	d.unknownPackage = ""
	return d
}

// error aborts the decoding by panicking with err.
func (d *decodeState) error(err error) {
	panic(err)
}

// saveError saves the first err it is called with,
// for reporting at the end of the unmarshal.
func (d *decodeState) saveError(err error) {
	if d.savedError == nil {
		d.savedError = err
	}
}

// next cuts off and returns the next full JSON value in d.data[d.off:].
// The next value is known to be an object or array, not a literal.
func (d *decodeState) next() []byte {
	c := d.data[d.off]
	item, rest, err := nextValue(d.data[d.off:], &d.nextscan)
	if err != nil {
		// ke: {"block": {"notest": true}}
		d.error(err)
	}
	d.off = len(d.data) - len(rest)

	// Our scanner has seen the opening brace/bracket
	// and thinks we're still in the middle of the object.
	// invent a closing brace/bracket to get it out.
	if c == '{' {
		d.scan.step(&d.scan, '}')
	} else {
		d.scan.step(&d.scan, ']')
	}

	return item
}

// scanWhile processes bytes in d.data[d.off:] until it
// receives a scan code not equal to op.
// It updates d.off and returns the new scan code.
func (d *decodeState) scanWhile(op int) int {
	var newOp int
	for {
		if d.off >= len(d.data) {
			newOp = d.scan.eof()
			d.off = len(d.data) + 1 // mark processed EOF with len+1
		} else {
			c := int(d.data[d.off])
			d.off++
			newOp = d.scan.step(&d.scan, c)
		}
		if newOp != op {
			break
		}
	}
	return newOp
}

// value decodes a JSON value from d.data[d.off:] into the value.
// it updates d.off to point past the decoded value.
func (d *decodeState) value(v reflect.Value) {
	if !v.IsValid() {
		_, rest, err := nextValue(d.data[d.off:], &d.nextscan)
		if err != nil {
			// ke: {"block": {"notest": true}}
			d.error(err)
		}
		d.off = len(d.data) - len(rest)

		// d.scan thinks we're still at the beginning of the item.
		// Feed in an empty string - the shortest, simplest value -
		// so that it knows we got to the end of the value.
		if d.scan.redo {
			// rewind.
			d.scan.redo = false
			d.scan.step = stateBeginValue
		}
		d.scan.step(&d.scan, '"')
		d.scan.step(&d.scan, '"')

		n := len(d.scan.parseState)
		if n > 0 && d.scan.parseState[n-1] == parseObjectKey {
			// d.scan thinks we just read an object key; finish the object
			d.scan.step(&d.scan, ':')
			d.scan.step(&d.scan, '"')
			d.scan.step(&d.scan, '"')
			d.scan.step(&d.scan, '}')
		}

		return
	}

	switch op := d.scanWhile(scanSkipSpace); op {
	default:
		// ke: {"block": {"notest": true}}
		d.error(errPhase)

	case scanBeginArray:
		d.array(v)

	case scanBeginObject:
		d.object(v)

	case scanBeginLiteral:
		d.literal(v)
	}
}

type unquotedValue struct{}

// valueQuoted is like value but decodes a
// quoted string literal or literal null into an interface value.
// If it finds anything other than a quoted string literal or null,
// valueQuoted returns unquotedValue{}.
func (d *decodeState) valueQuoted() interface{} {
	switch op := d.scanWhile(scanSkipSpace); op {
	default:
		// ke: {"block": {"notest": true}}
		d.error(errPhase)

	case scanBeginArray:
		// ke: {"block": {"notest": true}}
		d.array(reflect.Value{})

	case scanBeginObject:
		// ke: {"block": {"notest": true}}
		d.object(reflect.Value{})

	case scanBeginLiteral:
		switch v := d.literalInterface().(type) {
		case nil, string:
			return v
		}
	}

	// ke: {"block": {"notest": true}}
	return unquotedValue{}
}

// indirect walks down v allocating pointers as needed,
// until it gets to a non-pointer.
// if it encounters an Unmarshaler, indirect stops and returns that.
// if decodingNull is true, indirect stops at the last pointer so it can be set to nil.
func indirect(v reflect.Value, decodingNull bool, unmarshalers bool, unpackers bool) (Unmarshaler, encoding.TextUnmarshaler, Unpacker, reflect.Value) {

	// If v is a named type and is addressable,
	// start with its address, so that if the type has pointer methods,
	// we find them.
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		v = v.Addr()
	}
	for {
		// Load value from interface, but only if the result will be
		// usefully addressable.
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() && (!decodingNull || e.Elem().Kind() == reflect.Ptr) {
				v = e
				continue
			}
		}

		if v.Kind() != reflect.Ptr {
			break
		}

		if v.Elem().Kind() != reflect.Ptr && decodingNull && v.CanSet() {
			break
		}
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		if v.Type().NumMethod() > 0 {
			if unmarshalers {
				if u, ok := v.Interface().(Unmarshaler); ok {
					return u, nil, nil, reflect.Value{}
				}
				if ut, ok := v.Interface().(encoding.TextUnmarshaler); ok {
					return nil, ut, nil, reflect.Value{}
				}
			}
			if unpackers {
				if up, ok := v.Interface().(Unpacker); ok {
					return nil, nil, up, reflect.Value{}
				}
			}
		}
		v = v.Elem()
	}
	return nil, nil, nil, v
}

// array consumes an array from d.data[d.off-1:], decoding into the value v.
// the first byte of the array ('[') has been read already.
func (d *decodeState) array(v reflect.Value) {
	// Check for unmarshaler.
	u, ut, up, pv := indirect(v, false, true, true)
	if u != nil {
		d.off--
		err := u.UnmarshalJSON(d.ctx, d.next())
		if err != nil {
			// ke: {"block": {"notest": true}}
			d.error(err)
		}
		return
	}
	if ut != nil {
		// ke: {"block": {"notest": true}}
		d.saveError(&UnmarshalTypeError{"array", v.Type()})
		d.off--
		d.next()
		return
	}
	if up != nil {
		d.off--
		var i interface{}
		if err := UnmarshalUntyped(d.ctx, d.next(), &i); err != nil {
			// ke: {"block": {"notest": true}}
			d.error(err)
			return
		}
		if err := up.Unpack(d.ctx, Pack(i)); err != nil {
			// ke: {"block": {"notest": true}}
			d.error(err)
			return
		}
		return
	}

	v = pv

	// Check type of target.
	switch v.Kind() {
	case reflect.Interface:
		if v.NumMethod() == 0 {
			// Decoding into nil interface?  Switch to non-reflect code.
			v.Set(reflect.ValueOf(d.arrayInterface()))
			return
		}
		// Otherwise it's invalid.
		fallthrough
	default:
		d.saveError(&UnmarshalTypeError{"array", v.Type()})
		d.off--
		d.next()
		return
	case reflect.Array:
	case reflect.Slice:
		break
	}

	i := 0
	for {
		// Look ahead for ] - can only happen on first iteration.
		op := d.scanWhile(scanSkipSpace)
		if op == scanEndArray {
			break
		}

		// Back up so d.value can have the byte we just read.
		d.off--
		d.scan.undo(op)

		// Get element of array, growing if necessary.
		if v.Kind() == reflect.Slice {
			// Grow slice if necessary
			if i >= v.Cap() {
				newcap := v.Cap() + v.Cap()/2
				if newcap < 4 {
					newcap = 4
				}
				newv := reflect.MakeSlice(v.Type(), v.Len(), newcap)
				reflect.Copy(newv, v)
				v.Set(newv)
			}
			if i >= v.Len() {
				v.SetLen(i + 1)
			}
		}

		if i < v.Len() {
			// Decode into element.
			d.value(v.Index(i))
		} else {
			// Ran out of fixed array: skip.
			d.value(reflect.Value{})
		}
		i++

		// Next token must be , or ].
		op = d.scanWhile(scanSkipSpace)
		if op == scanEndArray {
			break
		}
		if op != scanArrayValue {
			// ke: {"block": {"notest": true}}
			d.error(errPhase)
		}
	}

	if i < v.Len() {
		if v.Kind() == reflect.Array {
			// Array.  Zero the rest.
			z := reflect.Zero(v.Type().Elem())
			for ; i < v.Len(); i++ {
				v.Index(i).Set(z)
			}
		} else {
			// ke: {"block": {"notest": true}}
			v.SetLen(i)
		}
	}
	if i == 0 && v.Kind() == reflect.Slice {
		v.Set(reflect.MakeSlice(v.Type(), 0, 0))
	}
}

var nullLiteral = []byte("null")

func (d *decodeState) scanForType(v reflect.Value) {
	if typeName, _, found := d.scanForAttribute("type", true); found {
		d.setType(typeName, v)
	} else {
		d.unknownType = "(no type attribute)"
	}
}

func (d *decodeState) scanForAttribute(attribute string, decodeString bool) (string, []byte, bool) {
	// We should scan forwards to find the "type" field, then reset the scanner.
	stateBackup := d.backup()
	var valueString string
	var valueBytes []byte
	var found bool
	for {
		// Read opening " of string key or closing }.
		op := d.scanWhile(scanSkipSpace)
		if op == scanEndObject {
			// ke: {"block": {"notest": true}}
			// closing } - can only happen on first iteration.
			break
		}
		if op != scanBeginLiteral {
			// ke: {"block": {"notest": true}}
			d.error(errPhase)
		}
		// Read key.
		start := d.off - 1
		op = d.scanWhile(scanContinue)
		item := d.data[start : d.off-1]
		key, ok := unquoteBytes(item)
		if !ok {
			// ke: {"block": {"notest": true}}
			d.error(errPhase)
		}

		// Read : before value.
		if op == scanSkipSpace {
			// ke: {"block": {"notest": true}}
			op = d.scanWhile(scanSkipSpace)
		}
		if op != scanObjectKey {
			// ke: {"block": {"notest": true}}
			d.error(errPhase)
		}

		if string(key) == attribute {
			// Read value.
			if decodeString {
				var str string
				strv := reflect.ValueOf(&str).Elem()
				d.value(strv)
				valueString = str
			} else {
				valueBytes = d.next()
			}
			found = true
			break
		} else {
			// skip value
			d.value(reflect.Value{})
		}

		// Next token must be , or }.
		op = d.scanWhile(scanSkipSpace)
		if op == scanEndObject {
			// ke: {"block": {"notest": true}}
			break
		}
		if op != scanObjectValue {
			// ke: {"block": {"notest": true}}
			d.error(errPhase)
		}

	}
	// Rewind and restore the state of the decoder
	d.restore(stateBackup)
	return valueString, valueBytes, found
}

func (d *decodeState) setType(typeName string, v reflect.Value) {

	path, name, err := GetReferencePartsFromTypeString(d.ctx, typeName)
	if unk, ok := kerr.Source(err).(UnknownPackageError); ok {
		// We don't want to throw an error here, because when we're scanning for
		// aliases we need to tolerate unknown packages
		d.unknownPackage = unk.UnknownPackage
	}

	jcache := jsonctx.FromContext(d.ctx)

	var typ reflect.Type

	// We should look the type up in the type resolver
	typ, ok := jcache.GetType(path, name)
	if !ok && v.Kind() == reflect.Interface {
		// If we can't find the type in the resolver, and
		// we're unmarshaling into an interface, then look
		// the interface in the dummy interface resolver.
		typ, _ = jcache.Dummies.Get(v.Type())
	}

	if typ == nil {
		d.unknownType = path + ":" + name
		return
	}

	if err := setType(v, typ); err != nil {
		d.saveError(kerr.Wrap("IOBHNASQUE", err))
	}

}

func setType(v reflect.Value, typ reflect.Type) error {

	if !v.CanSet() &&
		v.Kind() == reflect.Ptr &&
		v.Elem().Kind() == reflect.Interface &&
		v.Elem().Type().NumMethod() == 0 {
		// If we're unmarshaling into a *interface{} we must set
		// the interface{} bit (the Elem).
		v = v.Elem()
	}

	if v.Type() == typ {
		return nil
	}

	if !v.CanSet() {
		return kerr.New("HDJOTUDTIR", "Can't set type %s", typ.String())
	}

	val := getEmptyValue(typ)

	v.Set(val)

	return nil
}

func getEmptyValue(typ reflect.Type) reflect.Value {

	val := reflect.New(typ).Elem()

	// this ensures value is set
	indirect(val, false, false, false)

	return val
}

func GetReferencePartsFromTypeString(ctx context.Context, typeString string) (path string, name string, err error) {

	env := envctx.FromContext(ctx)

	if strings.Contains(typeString, "/") {
		// If the type name contains a slash, I'm assuming it's a fully
		// qualified type name of the form "kego.io/system:type".
		// TODO: Improve this with a regex?
		parts := strings.Split(typeString, ":")

		// We hard-code system and json to prevent them having to always be
		// specified in the aliases
		if parts[0] == "kego.io/system" {
			return "kego.io/system", parts[1], nil
		} else if parts[0] == "kego.io/json" {
			return "kego.io/json", parts[1], nil
		}

		_, found := findKey(env.Aliases, parts[0])
		if !found && parts[0] != env.Path {
			return "", "", UnknownPackageError{
				Struct:         kerr.New("KJSOXDESFD", "Unknown package %s", parts[0]),
				UnknownPackage: parts[0],
			}
		}
		return parts[0], parts[1], nil
	} else if strings.Contains(typeString, ":") {
		// If the type name contains a colon, I'm assuming it's an abreviated
		// qualified type name of the form "system:type". We should look the
		// package name up in the aliases map.
		// TODO: Improve this with a regex?
		parts := strings.Split(typeString, ":")

		// We hard-code system and json to prevent them having to always be
		// specified in the aliases
		if parts[0] == "system" {
			return "kego.io/system", parts[1], nil
		} else if parts[0] == "json" {
			return "kego.io/json", parts[1], nil
		}

		packagePath, ok := env.Aliases[parts[0]]
		if !ok {
			return "", "", UnknownPackageError{
				Struct:         kerr.New("DKKFLKDKYI", "Unknown package %s", parts[0]),
				UnknownPackage: parts[0],
			}
		}
		return packagePath, parts[1], nil
	} else {
		return env.Path, typeString, nil
	}
}
func findKey(m map[string]string, value string) (string, bool) {
	for k, v := range m {
		if value == v {
			return k, true
		}
	}
	return "", false
}

// object consumes an object from d.data[d.off-1:], decoding into the value v.
// the first byte ('{') of the object has been read already.
func (d *decodeState) object(v reflect.Value) {

	hasConcreteType := false
	concreteTypePath := ""
	concreteTypeName := ""
	if d.typed {
		_, _, _, val := indirect(v, false, false, false)

		// If the type we're unmarshaling into is an interface, we should scan
		// for a "type" attribute and initialise the correct type.
		switch val.Kind() {
		case reflect.Interface:
			// If needed, this sets the value of v to the correct type based on
			// the "type" attribute.
			d.scanForType(v)
			val1 := val
			for val1.Kind() == reflect.Ptr || val1.Kind() == reflect.Interface {
				val1 = val1.Elem()
			}
			if isKeNative(val1.Kind()) {
				var found bool
				var data []byte
				if _, data, found = d.scanForAttribute("value", false); !found {
					d.error(kerr.New("HNKBOLXUWU", "Can't find value field"))
					return
				}
				d.off--
				d.next()
				d.literalStore(data, v, false)
				return
			}
		case reflect.Struct:
			// If we're unmarshaling into a concrete type, we want to be able to omit the "type"
			// attribute, so we should add it back in if it's missing so the system:base object is
			// correct.
			path, name, ok := jsonctx.FromContext(d.ctx).GetTypeByReflectType(val.Type())
			if ok {
				hasConcreteType = true
				concreteTypePath = path
				concreteTypeName = name
			}
		}
	}

	u, ut, up, pv := indirect(v, false, true, true)
	if ut != nil {
		// ke: {"block": {"notest": true}}
		d.saveError(&UnmarshalTypeError{"object", v.Type()})
		d.off--
		d.next() // skip over { } in input
		return
	}
	if u != nil {
		d.off--
		if err := u.UnmarshalJSON(d.ctx, d.next()); err != nil {
			// ke: {"block": {"notest": true}}
			d.error(err)
		}
		return
	}
	if up != nil {
		d.off--
		var i interface{}
		if err := UnmarshalUntyped(d.ctx, d.next(), &i); err != nil {
			// ke: {"block": {"notest": true}}
			d.error(err)
			return
		}
		if err := up.Unpack(d.ctx, Pack(i)); err != nil {
			// ke: {"block": {"notest": true}}
			d.error(err)
			return
		}
		return
	}
	v = pv

	// Decoding into nil interface?  Switch to non-reflect code.
	if v.Kind() == reflect.Interface && v.NumMethod() == 0 {
		v.Set(reflect.ValueOf(d.objectInterface()))
		return
	}

	// Check type of target: struct or map[string]T
	switch v.Kind() {
	case reflect.Map:
		// map must have string kind
		t := v.Type()
		if t.Key().Kind() != reflect.String {
			d.saveError(&UnmarshalTypeError{"object", v.Type()})
			d.off--
			d.next() // skip over { } in input
			return
		}
		if v.IsNil() {
			v.Set(reflect.MakeMap(t))
		}
	case reflect.Struct:

	default:
		if !d.typed {
			// we should only avoid this error if we're in typed mode,
			// or we'll break the untyped behaviour.
			d.saveError(&UnmarshalTypeError{"object", v.Type()})
		}
		d.off--
		d.next() // skip over { } in input
		return
	}

	var mapElem reflect.Value

	foundFields := make([]field, 10)
	for {
		// Read opening " of string key or closing }.
		op := d.scanWhile(scanSkipSpace)
		if op == scanEndObject {
			// closing } - can only happen on first iteration.
			break
		}
		if op != scanBeginLiteral {
			// ke: {"block": {"notest": true}}
			d.error(errPhase)
		}

		// Read key.
		start := d.off - 1
		op = d.scanWhile(scanContinue)
		item := d.data[start : d.off-1]
		key, ok := unquoteBytes(item)
		if !ok {
			// ke: {"block": {"notest": true}}
			d.error(errPhase)
		}

		// Figure out field corresponding to key.
		var subv reflect.Value
		destring := false // whether the value is wrapped in a string to be decoded first

		if v.Kind() == reflect.Map {
			elemType := v.Type().Elem()
			if !mapElem.IsValid() {
				mapElem = reflect.New(elemType).Elem()
			} else {
				mapElem.Set(reflect.Zero(elemType))
			}
			subv = mapElem
		} else {
			var f *field
			fields := cachedTypeFields(v.Type())
			for i := range fields {
				ff := &fields[i]
				if bytes.Equal(ff.nameBytes, key) {
					f = ff
					break
				}
				if f == nil && ff.equalFold(ff.nameBytes, key) {
					f = ff
				}
			}
			if f != nil {
				subv = v
				destring = f.quoted
				for _, i := range f.index {
					if subv.Kind() == reflect.Ptr {
						if subv.IsNil() {
							subv.Set(reflect.New(subv.Type().Elem()))
						}
						subv = subv.Elem()
					}
					subv = subv.Field(i)
				}
				foundFields = append(foundFields, *f)
			}
		}

		// Read : before value.
		if op == scanSkipSpace {
			op = d.scanWhile(scanSkipSpace)
		}
		if op != scanObjectKey {
			// ke: {"block": {"notest": true}}
			d.error(errPhase)
		}

		// Read value.
		if destring {
			switch qv := d.valueQuoted().(type) {
			case nil:
				d.literalStore(nullLiteral, subv, false)
			case string:
				d.literalStore([]byte(qv), subv, true)
			default:
				// ke: {"block": {"notest": true}}
				d.saveError(fmt.Errorf("json: invalid use of ,string struct tag, trying to unmarshal unquoted value into %v", item, v.Type()))
			}
		} else {
			d.value(subv)
		}

		// Write value back to map;
		// if using struct, subv points into struct already.
		if v.Kind() == reflect.Map {
			kv := reflect.ValueOf(key).Convert(v.Type().Key())
			v.SetMapIndex(kv, subv)
		}

		// Next token must be , or }.
		op = d.scanWhile(scanSkipSpace)
		if op == scanEndObject {
			break
		}
		if op != scanObjectValue {
			// ke: {"block": {"notest": true}}
			d.error(errPhase)
		}
	}
	if err := initialiseUnmarshaledObject(d.ctx, v, foundFields, d.typed, hasConcreteType, concreteTypePath, concreteTypeName); err != nil {
		// ke: {"block": {"notest": true}}
		d.saveError(err)
	}
}

func initialiseUnmarshaledObject(ctx context.Context, v reflect.Value, foundFields []field, typed bool, hasConcreteType bool, concreteTypePath string, concreteTypeName string) error {
	// When we finish decoding the json for an object, we should loop round the fields
	// of the struct we're decoding into and assign any default values for missing items.
	if v.Kind() == reflect.Struct {
		fields := cachedTypeFields(v.Type())
	outer:
		for _, f := range fields {
			for _, found := range foundFields {
				if f.name == found.name {
					continue outer
				}
			}

			// field not found
			if f.kego != nil && f.kego.Default != nil {
				// Figure out field corresponding to key.
				subv := v
				for _, i := range f.index {
					if subv.Kind() == reflect.Ptr {
						// I can't seem to get in here. I'm sure it's possible.
						// TODO: Get coverage in here!
						// ke: {"block": {"notest": true}}
						if subv.IsNil() {
							// ke: {"block": {"notest": true}}
							subv.Set(reflect.New(subv.Type().Elem()))
						}
						// ke: {"block": {"notest": true}}
						subv = subv.Elem()
					}
					subv = subv.Field(i)
				}
				def := f.kego.Default
				path := def.Path
				if path == "" {
					path = "kego.io/json"
				}
				unpackCtx := envctx.Dummy(ctx, path, def.Aliases)
				var d decodeState
				err := checkValid(*def.Value, &d.scan)
				if err != nil {
					return kerr.Wrap("XCLKSFIKEE", err)
				}
				d.init(unpackCtx, *def.Value, true)
				if def.Type != "" {
					d.setType(def.Type, subv)
				}
				err = d.unmarshalValue(subv.Addr())
				if err != nil {
					return kerr.Wrap("PHUDKAKHMN", err)
				}
			}
		}
	}

	// We loop round the fields, and initialise any anonymous
	// fields that are nil.
	// This is tested in system.TestInitialiseAnonymousFields
	if typed && v.Kind() == reflect.Struct {
		for i := 0; i < v.Type().NumField(); i++ {
			sf := v.Type().Field(i)
			fv := v.FieldByName(sf.Name)
			if sf.Anonymous && fv.IsNil() {
				fv.Set(reflect.New(sf.Type.Elem()))
			}
		}
	}

	// If we are unpacking this object into a concrete type defined by the schema, we should set
	// the type with the InitializableType interface. This is implemented by system:object. This
	// enables us to allow the type to be omitted from the json.
	if typed && hasConcreteType && v.Kind() == reflect.Struct {
		it, ok := v.Interface().(InitializableType)
		if !ok && v.CanAddr() {
			it, ok = v.Addr().Interface().(InitializableType)
		}
		if ok {
			if err := it.InitializeType(concreteTypePath, concreteTypeName); err != nil {
				if ite, ok := err.(InitializableTypeError); ok {
					return &UnmarshalTypeError{fmt.Sprint(ite.UnmarshalledPath, ":", ite.UnmarshalledName), v.Type()}
				} else {
					return &UnmarshalTypeError{"unknown object", v.Type()}
				}
			}
		}
	}
	return nil
}

type InitializableType interface {
	InitializeType(path string, name string) error
}

// if we tried to unmarshal an incompatible type, we return this from the InitializeType
// funciton. We include the package path and name of the unmarshaled object.
type InitializableTypeError struct {
	UnmarshalledPath string
	UnmarshalledName string
	IntoPath         string
	IntoName         string
}

func (i InitializableTypeError) Error() string {
	return fmt.Sprintf("Tried to unmarshal %s:%s into %s:%s", i.UnmarshalledPath, i.UnmarshalledName, i.IntoPath, i.IntoName)
}

// literal consumes a literal from d.data[d.off-1:], decoding into the value v.
// The first byte of the literal has been read already
// (that's how the caller knows it's a literal).
func (d *decodeState) literal(v reflect.Value) {
	// All bytes inside literal return scanContinue op code.
	start := d.off - 1
	op := d.scanWhile(scanContinue)

	// Scan read one byte too far; back up.
	d.off--
	d.scan.undo(op)

	d.literalStore(d.data[start:d.off], v, false)
}

// convertNumber converts the number literal s to a float64 or a NumberLiteral
// depending on the setting of d.useNumber.
func (d *decodeState) convertNumber(s string) (interface{}, error) {
	if d.useNumber {
		return NumberLiteral(s), nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		// ke: {"block": {"notest": true}}
		return nil, &UnmarshalTypeError{"number " + s, reflect.TypeOf(0.0)}
	}
	return f, nil
}

var numberType = reflect.TypeOf(NumberLiteral(""))

// literalStore decodes a literal stored in item into v.
//
// fromQuoted indicates whether this literal came from unwrapping a
// string from the ",string" struct tag option. this is used only to
// produce more helpful error messages.
func (d *decodeState) literalStore(item []byte, v reflect.Value, fromQuoted bool) {
	// Check for unmarshaler.
	if len(item) == 0 {
		//Empty string given
		d.saveError(fmt.Errorf("json: invalid use of ,string struct tag, trying to unmarshal %q into %v", item, v.Type()))
		return
	}
	wantptr := item[0] == 'n' // null
	u, ut, up, pv := indirect(v, wantptr, true, true)
	if u != nil {
		err := u.UnmarshalJSON(d.ctx, item)
		if err != nil {
			d.error(err)
		}
		return
	}
	if ut != nil {
		if item[0] != '"' {
			// ke: {"block": {"notest": true}}
			if fromQuoted {
				// ke: {"block": {"notest": true}}
				d.saveError(fmt.Errorf("json: invalid use of ,string struct tag, trying to unmarshal %q into %v", item, v.Type()))
			} else {
				// ke: {"block": {"notest": true}}
				d.saveError(&UnmarshalTypeError{"string", v.Type()})
			}
		}
		s, ok := unquoteBytes(item)
		if !ok {
			// ke: {"block": {"notest": true}}
			if fromQuoted {
				// ke: {"block": {"notest": true}}
				d.error(fmt.Errorf("json: invalid use of ,string struct tag, trying to unmarshal %q into %v", item, v.Type()))
			} else {
				// ke: {"block": {"notest": true}}
				d.error(errPhase)
			}
		}
		err := ut.UnmarshalText(s)
		if err != nil {
			// ke: {"block": {"notest": true}}
			d.error(err)
		}
		return
	}
	if up != nil {
		var i interface{}
		if err := UnmarshalUntyped(d.ctx, item, &i); err != nil {
			// ke: {"block": {"notest": true}}
			d.error(err)
			return
		}
		if err := up.Unpack(d.ctx, Pack(i)); err != nil {
			// ke: {"block": {"notest": true}}
			d.error(err)
			return
		}
		return
	}

	v = pv

	switch c := item[0]; c {
	case 'n': // null
		switch v.Kind() {
		case reflect.Interface, reflect.Ptr, reflect.Map, reflect.Slice:
			v.Set(reflect.Zero(v.Type()))
			// otherwise, ignore null for primitives/string
		}
	case 't', 'f': // true, false
		value := c == 't'
		switch v.Kind() {
		default:
			if fromQuoted {
				d.saveError(fmt.Errorf("json: invalid use of ,string struct tag, trying to unmarshal %q into %v", item, v.Type()))
			} else {
				// ke: {"block": {"notest": true}}
				d.saveError(&UnmarshalTypeError{"bool", v.Type()})
			}
		case reflect.Bool:
			v.SetBool(value)
		case reflect.Interface:
			if v.NumMethod() == 0 {
				v.Set(reflect.ValueOf(value))
			} else if d.typed {
				if err := setDefaultNativeValue(d.ctx, v, item); err != nil {
					d.saveError(kerr.Wrap("WTDRKGKGCC", err))
				}
			} else {
				d.saveError(&UnmarshalTypeError{"bool", v.Type()})
			}
		}

	case '"': // string
		s, ok := unquoteBytes(item)
		if !ok {
			// ke: {"block": {"notest": true}}
			if fromQuoted {
				// ke: {"block": {"notest": true}}
				d.error(fmt.Errorf("json: invalid use of ,string struct tag, trying to unmarshal %q into %v", item, v.Type()))
			} else {
				// ke: {"block": {"notest": true}}
				d.error(errPhase)
			}
		}
		switch v.Kind() {
		default:
			// ke: {"block": {"notest": true}}
			d.saveError(&UnmarshalTypeError{"string", v.Type()})
		case reflect.Slice:
			if v.Type().Elem().Kind() != reflect.Uint8 {
				// ke: {"block": {"notest": true}}
				d.saveError(&UnmarshalTypeError{"string", v.Type()})
				break
			}
			b := make([]byte, base64.StdEncoding.DecodedLen(len(s)))
			n, err := base64.StdEncoding.Decode(b, s)
			if err != nil {
				// ke: {"block": {"notest": true}}
				d.saveError(err)
				break
			}
			v.Set(reflect.ValueOf(b[0:n]))
		case reflect.String:
			v.SetString(string(s))
		case reflect.Interface:
			if v.NumMethod() == 0 {
				v.Set(reflect.ValueOf(string(s)))
			} else if d.typed {
				if err := setDefaultNativeValue(d.ctx, v, item); err != nil {
					d.saveError(kerr.Wrap("MNUCNWUWVF", err))
				}
			} else {
				d.saveError(&UnmarshalTypeError{"string", v.Type()})
			}
		}

	default: // number
		if c != '-' && (c < '0' || c > '9') {
			if fromQuoted {
				d.error(fmt.Errorf("json: invalid use of ,string struct tag, trying to unmarshal %q into %v", item, v.Type()))
			} else {
				// ke: {"block": {"notest": true}}
				d.error(errPhase)
			}
		}
		s := string(item)
		switch v.Kind() {
		default:
			if v.Kind() == reflect.String && v.Type() == numberType {
				v.SetString(s)
				break
			}
			if fromQuoted {
				d.error(fmt.Errorf("json: invalid use of ,string struct tag, trying to unmarshal %q into %v", item, v.Type()))
			} else {
				// ke: {"block": {"notest": true}}
				d.error(&UnmarshalTypeError{"number", v.Type()})
			}
		case reflect.Interface:
			n, err := d.convertNumber(s)
			if err != nil {
				// ke: {"block": {"notest": true}}
				d.saveError(err)
				break
			}
			if v.NumMethod() != 0 {
				if d.typed {
					if err := setDefaultNativeValue(d.ctx, v, item); err != nil {
						d.saveError(kerr.Wrap("XDFCYMETOK", err))
					}
				} else {
					d.saveError(&UnmarshalTypeError{"number", v.Type()})
				}
				break
			}
			v.Set(reflect.ValueOf(n))

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(s, 10, 64)
			if err != nil || v.OverflowInt(n) {
				// ke: {"block": {"notest": true}}
				d.saveError(&UnmarshalTypeError{"number " + s, v.Type()})
				break
			}
			v.SetInt(n)

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			n, err := strconv.ParseUint(s, 10, 64)
			if err != nil || v.OverflowUint(n) {
				// ke: {"block": {"notest": true}}
				d.saveError(&UnmarshalTypeError{"number " + s, v.Type()})
				break
			}
			v.SetUint(n)

		case reflect.Float32, reflect.Float64:
			n, err := strconv.ParseFloat(s, v.Type().Bits())
			if err != nil || v.OverflowFloat(n) {
				// ke: {"block": {"notest": true}}
				d.saveError(&UnmarshalTypeError{"number " + s, v.Type()})
				break
			}
			v.SetFloat(n)
		}
	}
}

func setDefaultNativeValue(ctx context.Context, v reflect.Value, value []byte) error {
	t, ok := jsonctx.FromContext(ctx).GetTypeByInterface(v.Type())
	if !ok {
		return kerr.New("LDURBUTUXQ", "No type found for %s", v.Type().Name())
	}
	p := reflect.New(t.Elem())

	vali := p.Interface()

	if err := UnmarshalUntyped(ctx, value, vali); err != nil {
		return kerr.Wrap("IILXDLDADR", err)
	}

	v.Set(reflect.ValueOf(vali))
	return nil
}

// The xxxInterface routines build up a value to be stored
// in an empty interface.  They are not strictly necessary,
// but they avoid the weight of reflection in this common case.

// valueInterface is like value but returns interface{}
func (d *decodeState) valueInterface() interface{} {
	switch d.scanWhile(scanSkipSpace) {
	default:
		// ke: {"block": {"notest": true}}
		d.error(errPhase)
		panic("unreachable")
	case scanBeginArray:
		return d.arrayInterface()
	case scanBeginObject:
		return d.objectInterface()
	case scanBeginLiteral:
		return d.literalInterface()
	}
}

// arrayInterface is like array but returns []interface{}.
func (d *decodeState) arrayInterface() []interface{} {
	var v = make([]interface{}, 0)
	for {
		// Look ahead for ] - can only happen on first iteration.
		op := d.scanWhile(scanSkipSpace)
		if op == scanEndArray {
			break
		}

		// Back up so d.value can have the byte we just read.
		d.off--
		d.scan.undo(op)

		v = append(v, d.valueInterface())

		// Next token must be , or ].
		op = d.scanWhile(scanSkipSpace)
		if op == scanEndArray {
			break
		}
		if op != scanArrayValue {
			// ke: {"block": {"notest": true}}
			d.error(errPhase)
		}
	}
	return v
}

// objectInterface is like object but returns map[string]interface{}.
func (d *decodeState) objectInterface() map[string]interface{} {
	m := make(map[string]interface{})
	for {
		// Read opening " of string key or closing }.
		op := d.scanWhile(scanSkipSpace)
		if op == scanEndObject {
			// ke: {"block": {"notest": true}}
			// closing } - can only happen on first iteration.
			break
		}
		if op != scanBeginLiteral {
			// ke: {"block": {"notest": true}}
			d.error(errPhase)
		}

		// Read string key.
		start := d.off - 1
		op = d.scanWhile(scanContinue)
		item := d.data[start : d.off-1]
		key, ok := unquote(item)
		if !ok {
			// ke: {"block": {"notest": true}}
			d.error(errPhase)
		}

		// Read : before value.
		if op == scanSkipSpace {
			// ke: {"block": {"notest": true}}
			op = d.scanWhile(scanSkipSpace)
		}
		if op != scanObjectKey {
			// ke: {"block": {"notest": true}}
			d.error(errPhase)
		}

		// Read value.
		m[key] = d.valueInterface()

		// Next token must be , or }.
		op = d.scanWhile(scanSkipSpace)
		if op == scanEndObject {
			break
		}
		if op != scanObjectValue {
			// ke: {"block": {"notest": true}}
			d.error(errPhase)
		}
	}
	return m
}

// literalInterface is like literal but returns an interface value.
func (d *decodeState) literalInterface() interface{} {
	// All bytes inside literal return scanContinue op code.
	start := d.off - 1
	op := d.scanWhile(scanContinue)

	// Scan read one byte too far; back up.
	d.off--
	d.scan.undo(op)
	item := d.data[start:d.off]

	switch c := item[0]; c {
	case 'n': // null
		return nil

	case 't', 'f': // true, false
		return c == 't'

	case '"': // string
		s, ok := unquote(item)
		if !ok {
			// ke: {"block": {"notest": true}}
			d.error(errPhase)
		}
		return s

	default: // number
		if c != '-' && (c < '0' || c > '9') {
			// ke: {"block": {"notest": true}}
			d.error(errPhase)
		}
		n, err := d.convertNumber(string(item))
		if err != nil {
			// ke: {"block": {"notest": true}}
			d.saveError(err)
		}
		return n
	}
}

// getu4 decodes \uXXXX from the beginning of s, returning the hex value,
// or it returns -1.
func getu4(s []byte) rune {
	if len(s) < 6 || s[0] != '\\' || s[1] != 'u' {
		return -1
	}
	r, err := strconv.ParseUint(string(s[2:6]), 16, 64)
	if err != nil {
		// ke: {"block": {"notest": true}}
		return -1
	}
	return rune(r)
}

// unquote converts a quoted JSON string literal s into an actual string t.
// The rules are different than for Go, so cannot use strconv.Unquote.
func unquote(s []byte) (t string, ok bool) {
	s, ok = unquoteBytes(s)
	t = string(s)
	return
}

func unquoteBytes(s []byte) (t []byte, ok bool) {
	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		// ke: {"block": {"notest": true}}
		return
	}
	s = s[1 : len(s)-1]

	// Check for unusual characters. If there are none,
	// then no unquoting is needed, so return a slice of the
	// original bytes.
	r := 0
	for r < len(s) {
		c := s[r]
		if c == '\\' || c == '"' || c < ' ' {
			break
		}
		if c < utf8.RuneSelf {
			r++
			continue
		}
		rr, size := utf8.DecodeRune(s[r:])
		if rr == utf8.RuneError && size == 1 {
			break
		}
		r += size
	}
	if r == len(s) {
		return s, true
	}

	b := make([]byte, len(s)+2*utf8.UTFMax)
	w := copy(b, s[0:r])
	for r < len(s) {
		// Out of room?  Can only happen if s is full of
		// malformed UTF-8 and we're replacing each
		// byte with RuneError.
		if w >= len(b)-2*utf8.UTFMax {
			nb := make([]byte, (len(b)+utf8.UTFMax)*2)
			copy(nb, b[0:w])
			b = nb
		}
		switch c := s[r]; {
		case c == '\\':
			r++
			if r >= len(s) {
				// ke: {"block": {"notest": true}}
				return
			}
			switch s[r] {
			default:
				// ke: {"block": {"notest": true}}
				return
			case '"', '\\', '/', '\'':
				b[w] = s[r]
				r++
				w++
			case 'b':
				// ke: {"block": {"notest": true}}
				b[w] = '\b'
				r++
				w++
			case 'f':
				// ke: {"block": {"notest": true}}
				b[w] = '\f'
				r++
				w++
			case 'n':
				b[w] = '\n'
				r++
				w++
			case 'r':
				b[w] = '\r'
				r++
				w++
			case 't':
				b[w] = '\t'
				r++
				w++
			case 'u':
				r--
				rr := getu4(s[r:])
				if rr < 0 {
					// ke: {"block": {"notest": true}}
					return
				}
				r += 6
				if utf16.IsSurrogate(rr) {
					rr1 := getu4(s[r:])
					if dec := utf16.DecodeRune(rr, rr1); dec != unicode.ReplacementChar {
						// A valid pair; consume.
						r += 6
						w += utf8.EncodeRune(b[w:], dec)
						break
					}
					// Invalid surrogate; fall back to replacement rune.
					rr = unicode.ReplacementChar
				}
				w += utf8.EncodeRune(b[w:], rr)
			}

		// Quote, control characters are invalid.
		case c == '"', c < ' ':
			// ke: {"block": {"notest": true}}
			return

		// ASCII
		case c < utf8.RuneSelf:
			b[w] = c
			r++
			w++

		// Coerce to well-formed UTF-8.
		default:
			rr, size := utf8.DecodeRune(s[r:])
			r += size
			w += utf8.EncodeRune(b[w:], rr)
		}
	}
	return b[0:w], true
}
