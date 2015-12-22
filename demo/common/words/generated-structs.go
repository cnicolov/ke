package words

import (
	"reflect"

	"golang.org/x/net/context"
	"kego.io/json"
	"kego.io/system"
)

// Automatically created basic rule for localizer
type LocalizerRule struct {
	*system.Object
	*system.Rule
}

// Automatically created basic rule for simple
type SimpleRule struct {
	*system.Object
	*system.Rule
}

// Automatically created basic rule for translation
type TranslationRule struct {
	*system.Object
	*system.Rule
}
type Simple struct {
	*system.Object
	String *system.String `json:"string"`
}
type SimpleInterface interface {
	GetSimple(ctx context.Context) *Simple
}

func (o *Simple) GetSimple(ctx context.Context) *Simple {
	return o
}

// This represents a translated string
type Translation struct {
	*system.Object
	// The original English string
	English *system.String `kego:"{\"default\":{\"value\":\"http\"}}" json:"english"`
	// The translated strings
	Translations map[string]*system.String `json:"translations"`
}
type TranslationInterface interface {
	GetTranslation(ctx context.Context) *Translation
}

func (o *Translation) GetTranslation(ctx context.Context) *Translation {
	return o
}
func init() {
	json.Register("kego.io/demo/common/words", "@localizer", reflect.TypeOf((*LocalizerRule)(nil)), nil, 517105247315493779)
	json.Register("kego.io/demo/common/words", "@simple", reflect.TypeOf((*SimpleRule)(nil)), nil, 7410426200142012193)
	json.Register("kego.io/demo/common/words", "@translation", reflect.TypeOf((*TranslationRule)(nil)), nil, 2609009266940530265)
	json.Register("kego.io/demo/common/words", "localizer", reflect.TypeOf((*Localizer)(nil)).Elem(), nil, 517105247315493779)
	json.Register("kego.io/demo/common/words", "simple", reflect.TypeOf((*Simple)(nil)), reflect.TypeOf((*SimpleInterface)(nil)).Elem(), 7410426200142012193)
	json.Register("kego.io/demo/common/words", "translation", reflect.TypeOf((*Translation)(nil)), reflect.TypeOf((*TranslationInterface)(nil)).Elem(), 2609009266940530265)
}
