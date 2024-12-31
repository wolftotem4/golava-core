package lang

import (
	"testing"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
)

func TestFallbackTranslator_T(t *testing.T) {
	uni := ut.New(en.New(), zh.New())
	enTrans, _ := uni.GetTranslator("en")
	zhTrans, _ := uni.GetTranslator("zh")

	enTrans.Add("hello", "Hello", false)
	zhTrans.Add("hello", "你好", false)
	enTrans.Add("world", "World", false)

	translator := FallbackTranslator{
		Translator: zhTrans,
		Fallback:   enTrans,
	}

	t.Run("TestFallbackTranslator_T", func(t *testing.T) {
		expected := "你好"
		value, err := translator.T("hello")
		if err != nil {
			t.Fatal(err)
		}

		if value != expected {
			t.Errorf("Expected %s, got %s", expected, value)
		}
	})

	t.Run("TestFallbackTranslator_T_Fallback", func(t *testing.T) {
		expected := "World"
		value, err := translator.T("world")
		if err != nil {
			t.Fatal(err)
		}

		if value != expected {
			t.Errorf("Expected %s, got %s", expected, value)
		}
	})
}
