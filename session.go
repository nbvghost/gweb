package gweb

import (
	"net/url"
	"sync"
)

var Sessions = &SessionSafeMap{}

type SessionSafeMap struct {
	//sync.RWMutex
	_data sync.Map
}
type Session struct {
	//sync.RWMutex
	Attributes        *Attributes
	CreateTime        int64
	GLSESSIONID       string
	LastOperationTime int64
	LastRequestURL    *url.URL
}

func (s *SessionSafeMap) DeleteSession(GLSESSIONID string) {
	//s.Lock()
	//defer s.Unlock()
	//delete(s.Data, k)
	//db.NotifyAll(&db.Message{db.Socket_Type_2_STC,k})
	s._data.Delete(GLSESSIONID)
}

func (s *SessionSafeMap) AddSession(GLSESSIONID string, session *Session) {

	//s.Lock()
	//defer s.Unlock()
	/*if s.Data==nil{
		s.Data=make(map[string]*Session)
	}
	s.Data[GLSESSIONID] = session*/
	//db.NotifyAll(&db.Message{db.Socket_Type_1_STC,session})
	s._data.Store(GLSESSIONID, session)
}
func (s *SessionSafeMap) Range(f func(key, value interface{}) bool) {
	s._data.Range(f)
}
func (s *SessionSafeMap) GetSession(GLSESSIONID string) *Session {
	//s.RLock()
	//defer s.RUnlock()

	session, ok := s._data.Load(GLSESSIONID)
	if !ok {
		return nil
	}

	//session := s.Data[GLSESSIONID]

	//db.NotifyAll(&db.Message{db.Socket_Type_1_STC,session})
	return session.(*Session)
}
