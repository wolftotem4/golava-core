package lang

import (
	"errors"

	ut "github.com/go-playground/universal-translator"
)

type FallbackTranslator struct {
	ut.Translator
	Fallback ut.Translator
}

func NewFallbackTranslator(translator, fallback ut.Translator) *FallbackTranslator {
	return &FallbackTranslator{
		Translator: translator,
		Fallback:   fallback,
	}
}

func (g *FallbackTranslator) T(key interface{}, params ...string) (string, error) {
	line, err := g.Translator.T(key, params...)
	if errors.Is(err, ut.ErrUnknowTranslation) {
		return g.Fallback.T(key, params...)
	}
	return line, err
}

func (g *FallbackTranslator) C(key interface{}, num float64, digits uint64, param string) (string, error) {
	line, err := g.Translator.C(key, num, digits, param)
	if errors.Is(err, ut.ErrUnknowTranslation) {
		return g.Fallback.C(key, num, digits, param)
	}
	return line, err
}

func (g *FallbackTranslator) O(key interface{}, num float64, digits uint64, param string) (string, error) {
	line, err := g.Translator.O(key, num, digits, param)
	if errors.Is(err, ut.ErrUnknowTranslation) {
		return g.Fallback.O(key, num, digits, param)
	}
	return line, err
}

func (g *FallbackTranslator) R(key interface{}, num1 float64, digits1 uint64, num2 float64, digits2 uint64, param1, param2 string) (string, error) {
	line, err := g.Translator.R(key, num1, digits1, num2, digits2, param1, param2)
	if errors.Is(err, ut.ErrUnknowTranslation) {
		return g.Fallback.R(key, num1, digits1, num2, digits2, param1, param2)
	}
	return line, err
}
