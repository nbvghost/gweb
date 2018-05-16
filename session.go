package gweb

import (
	"sync"
	"net/url"
)

var Sessions = &sessionMap{Data: make(map[string]*Session)}

type sessionMap struct {
	sync.RWMutex
	Data map[string]*Session
}
type Session struct {
	sync.RWMutex
	Attributes  *Attributes
	CreateTime  int64
	ActionTime  int64
	Operation   int64
	GLSESSIONID string
	LastRequestURL *url.URL
}

func (s *sessionMap) DelectSession(k string) {
	s.Lock()
	delete(s.Data, k)
	defer s.Unlock()
	//db.NotifyAll(&db.Message{db.Socket_Type_2_STC,k})
}

func (s *sessionMap) addSession(GLSESSIONID string, session *Session) {
	s.Lock()
	s.Data[GLSESSIONID] = session
	defer s.Unlock()
	//db.NotifyAll(&db.Message{db.Socket_Type_1_STC,session})
}
func (s *sessionMap) GetSession(GLSESSIONID string) *Session {
	s.RLock()
	session := s.Data[GLSESSIONID]
	defer s.RUnlock()
	//db.NotifyAll(&db.Message{db.Socket_Type_1_STC,session})
	return session
}
