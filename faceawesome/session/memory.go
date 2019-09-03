package session

import (
	"container/list"
	"sync"
	"time"
)

/* session 存储 */
/* 实现session接口以及provider接口 */

// 这里采用list数据结构存储
type provider struct {
	lock sync.Mutex
	//存储在内存
	sessions map[string]*list.Element
	// 准备的GC
	list *list.List `sessionstore`
}

type sessionstore struct {
	sid string
	// 最后访问时间
	timeAccessed time.Time
	//session 里面存储的值
	value map[interface{}]interface{}
}

var pder = &provider{list: list.New()}

func (st *sessionstore) Set(key, value interface{}) error {
	st.value[key] = value
	pder.SessionUpdate(st.sid)
	return nil
}

func (st *sessionstore) Get(key interface{}) interface{} {
	pder.SessionUpdate(st.sid)
	if v, ok := st.value[key]; ok {
		return v
	}
	return nil
}

func (st *sessionstore) Delete(key interface{}) error {
	delete(st.value, key)
	pder.SessionUpdate(st.sid)
	return nil
}

func (st *sessionstore) SessionID() string {
	return st.sid
}

// 更新session访问时间
func (pder *provider) SessionUpdate(sid string) error {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	if element, ok := pder.sessions[sid]; ok {
		element.Value.(*sessionstore).timeAccessed = time.Now()
		pder.list.MoveToFront(element)
		return nil
	}
	return nil
}

func (pder *provider) SessionInit(sid string) (Session, error) {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	v := make(map[interface{}]interface{}, 0)
	newsess := &sessionstore{sid: sid, timeAccessed: time.Now(), value: v}
	element := pder.list.PushBack(newsess)
	pder.sessions[sid] = element
	return newsess, nil
}

func (pder *provider) SessionRead(sid string) (Session, error) {
	if element, ok := pder.sessions[sid]; ok {
		return element.Value.(*sessionstore), nil
	} else {
		sess, err := pder.SessionInit(sid)
		return sess, err
	}
	return nil, nil
}

func (pder *provider) SessionDestroy(sid string) error {
	if element, ok := pder.sessions[sid]; ok {
		delete(pder.sessions, sid)
		pder.list.Remove(element)
		return nil
	}
	return nil
}

func (pder *provider) SessionGC(maxlifetime int64) {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	for {
		element := pder.list.Back()
		if element == nil {
			break
		}
		if element.Value.(*sessionstore).timeAccessed.Unix()+maxlifetime < time.Now().Unix() {
			pder.list.Remove(element)
			delete(pder.sessions, element.Value.(*sessionstore).sid)
		} else {
			break
		}
	}
}
