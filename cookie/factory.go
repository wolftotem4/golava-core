package cookie

import "net/http"

type CookieFactory struct {
	Manager func() IEncryptableCookieManager
}

func (f *CookieFactory) Make(request *http.Request, responseWriter http.ResponseWriter) IEncryptableCookieManager {
	manager := f.Manager()
	manager.SetRequest(request)
	manager.SetResponseWriter(responseWriter)
	return manager
}
