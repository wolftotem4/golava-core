package lang

import ut "github.com/go-playground/universal-translator"

type SilentTranslator struct {
	ut.Translator
}

func (g SilentTranslator) T(key interface{}, params ...string) (string, error) {
	line, err := g.Translator.T(key, params...)
	if err != nil {
		return key.(string), nil
	}
	return line, nil
}

func (g SilentTranslator) C(key interface{}, num float64, digits uint64, param string) (string, error) {
	line, err := g.Translator.C(key, num, digits, param)
	if err != nil {
		return key.(string), nil
	}
	return line, nil
}

func (g SilentTranslator) O(key interface{}, num float64, digits uint64, param string) (string, error) {
	line, err := g.Translator.O(key, num, digits, param)
	if err != nil {
		return key.(string), nil
	}
	return line, nil
}

func (g SilentTranslator) R(key interface{}, num1 float64, digits1 uint64, num2 float64, digits2 uint64, param1, param2 string) (string, error) {
	line, err := g.Translator.R(key, num1, digits1, num2, digits2, param1, param2)
	if err != nil {
		return key.(string), nil
	}
	return line, nil
}
