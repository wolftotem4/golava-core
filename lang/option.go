package lang

import ut "github.com/go-playground/universal-translator"

type TranslatorOption func(args *TranslatorArgs)

type TranslatorArgs struct {
	Soft     bool
	Fallback ut.Translator
}

func (args *TranslatorArgs) Apply(ut ut.Translator) ut.Translator {
	if args.Fallback != nil {
		if args.Fallback.Locale() != ut.Locale() {
			ut = NewFallbackTranslator(ut, args.Fallback)
		}
	}

	if args.Soft {
		ut = NewSoftTranslator(ut)
	}

	return ut
}

func Soft(args *TranslatorArgs) {
	args.Soft = true
}

func Fallback(fallback ut.Translator) TranslatorOption {
	return func(args *TranslatorArgs) {
		args.Fallback = fallback
	}
}
