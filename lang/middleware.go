package lang

import (
	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/instance"
	"golang.org/x/text/language"
)

func SetLocale(key string, supportedLocales map[language.Tag]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		instance := instance.MustGetInstance(c)

		if !setLocaleWithQuery(c, instance, key) {
			setLocaleWithAcceptLanguage(c, instance, supportedLocales)
		}

		c.Next()
	}
}

func SetLocaleWithQuery(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		instance := instance.MustGetInstance(c)
		setLocaleWithQuery(c, instance, key)
		c.Next()
	}
}

func SetLocaleWithAcceptLanguage(supportedLocales map[language.Tag]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		instance := instance.MustGetInstance(c)
		setLocaleWithAcceptLanguage(c, instance, supportedLocales)
		c.Next()
	}
}

func setLocaleWithQuery(c *gin.Context, instance *instance.Instance, key string) bool {
	value := c.Query(key)
	if value == "" {
		return false
	}

	_, ok := instance.App.Base().Translation.GetTranslator(value)
	if ok {
		instance.Locale = value
		return true
	}

	return false
}

func setLocaleWithAcceptLanguage(c *gin.Context, instance *instance.Instance, supportedLocales map[language.Tag]string) bool {
	userPrefs, _, err := language.ParseAcceptLanguage(c.GetHeader("Accept-Language"))
	if err != nil {
		userPrefs = []language.Tag{language.English}
	}

	serverLangs := serverLangs(supportedLocales)
	_, index, confidence := language.NewMatcher(serverLangs).Match(userPrefs...)
	if confidence == language.Exact || confidence == language.High {
		tag := serverLangs[index]
		instance.Locale = supportedLocales[tag]
		return true
	}

	return false
}

func serverLangs(supportedLocales map[language.Tag]string) []language.Tag {
	langs := make([]language.Tag, 0, len(supportedLocales))
	for lang := range supportedLocales {
		langs = append(langs, lang)
	}
	return langs
}
