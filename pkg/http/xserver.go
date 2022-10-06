package xhttp

import (
	"github.com/gorilla/context"
	xsession "github.com/motai3/xsession/pkg"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	defaultManagerName   = "xsession_manage"
	defaultSessionIdName = "xsessionId"
)

type XServer struct {
	Engine         *gin.Engine
	sessionManager *xsession.Manager
	sessionName    string
}

func Default(name ...string) *XServer {
	n := defaultManagerName
	if len(name) > 0 {
		n = name[0]
	}
	s := &XServer{
		Engine:         gin.Default(),
		sessionManager: xsession.New(time.Second * 60 * 60),
		sessionName:    n,
	}
	s.SetSessionStorage(xsession.NewStorageMemory())
	return s
}

func New(name ...string) *XServer {
	n := defaultManagerName
	if len(name) > 0 {
		n = name[0]
	}
	s := &XServer{
		Engine:         gin.New(),
		sessionManager: xsession.New(time.Second * 60 * 60),
		sessionName:    n,
	}
	s.SetSessionStorage(xsession.NewStorageMemory())
	return s
}

func GetSession(c *gin.Context) (session *xsession.Session, err error) {
	manager, ok := c.Get(defaultManagerName)
	if !ok {
		return nil, err
	}
	m := manager.(*xsession.Manager)

	cookie, err := c.Request.Cookie(c.GetString(defaultSessionIdName))
	if err != nil {
		return m.New(c), err
	}
	sessionId := cookie.Value
	return m.New(c, sessionId), err
}

func (s *XServer) SetSessionStorage(storage xsession.Storage, ttl ...time.Duration) {
	if len(ttl) > 0 {
		s.sessionManager = xsession.New(ttl[0], storage)
	} else {
		s.sessionManager = xsession.New(time.Second*60*60, storage)
	}
	if s.Engine != nil {
		s.Engine.Use(func(c *gin.Context) {
			c.Set(defaultManagerName, s.sessionManager)
			c.Set(defaultSessionIdName, s.sessionName)
			defer context.Clear(c.Request)
			c.Next()
		})
	}
}

func (s *XServer) Run(addr ...string) {
	s.Engine.Run(addr...)
}
