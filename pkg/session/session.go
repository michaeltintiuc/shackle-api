package session

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type Session struct {
	store *sessions.CookieStore
	name  string
}

type SessionInfo struct {
	AuthKey []byte
	EncKey  []byte
	Name    string
}

func NewSession(s SessionInfo) *Session {
	asd := &Session{
		store: sessions.NewCookieStore(s.AuthKey, s.EncKey),
		name:  s.Name,
	}

	asd.store.Options.Path = "/"
	asd.store.Options.Secure = true
	asd.store.Options.HttpOnly = true
	return asd
}

func (s *Session) Get(r *http.Request) (*sessions.Session, error) {
	return s.store.Get(r, s.name)
}
