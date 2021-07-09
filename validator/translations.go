package validator

import ut "github.com/jinycoo/jinygo/text/i18n/universal-translator"

// TranslationFunc is the function type used to register or override
// custom translations
type TranslationFunc func(ut ut.Translator, fe FieldError) string

// RegisterTranslationsFunc allows for registering of translations
// for a 'ut.Translator' for use within the 'TranslationFunc'
type RegisterTranslationsFunc func(ut ut.Translator) error
