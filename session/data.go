package session

type ClientData struct {
	UserID    any
	IPAddress string
	UserAgent string
}

type SessionData struct {
	ClientData
	Payload []byte
}
