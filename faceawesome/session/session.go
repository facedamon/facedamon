package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

/* 把行为和结构按照OOP逻辑拆分 */

// Session 操作接口
type Session interface {
	//set session value
	Set(key, value interface{}) error
	//get session value
	Get(key interface{}) interface{}
	//delete session value
	Delete(key interface{}) error
	//back current sessionID
	SessionID() string
}

// session底层存储结构
type Provider interface {
	SessionInit(sid string) (Session, error)

	SessionRead(sid string) (Session, error)

	SessionDestroy(sid string) error

	// GC 过期session
	SessionGC(maxLifeTime int64)
}

// 全局 session 管理器
type Manager struct {
	// private cookie name
	cookieName string
	//protects session
	lock sync.Mutex
	// session 具体实现者
	provider Provider

	maxlifetime int64
}

/* 以上设计思路来源于database/sql/driver,先定义好接口，
   然后具体的存储session的结构实现相应的接口并注册后，
   相应功能这样就可以使用了
*/

// 全局session容器
var provides = make(map[string]Provider)

// session注册器
func Register(name string, provider Provider) {
	if nil == provider {
		panic("session: Register provide is nil")
	}
	if _, dup := provides[name]; dup {
		panic("session: Register called twice for provide " + name)
	}
	provides[name] = provider
}

// 生成全局唯一sessionid
func (manager *Manager) sessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func NewManager(provideName, cookieName string, maxlifetime int64) (*Manager, error) {
	provider, ok := provides[provideName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
	}
	return &Manager{provider: provider, cookieName: cookieName, maxlifetime: maxlifetime}, nil
}

// 创建session
func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session Session) {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		sid := manager.sessionId()
		session, _ = manager.provider.SessionInit(sid)
		cookie := http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(manager.maxlifetime)}
		http.SetCookie(w, &cookie)
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session, _ = manager.provider.SessionRead(sid)
	}
	return
}
