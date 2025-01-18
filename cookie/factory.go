package cookie

type CookieFactory struct {
	Manager func() IEncryptableCookieManager
}

func (f *CookieFactory) Make() IEncryptableCookieManager {
	return f.Manager()
}
