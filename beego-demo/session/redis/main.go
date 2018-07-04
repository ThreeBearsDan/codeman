package main

import (
	"fmt"
	"net/http"

	"sync"

	"github.com/astaxie/beego/session"
	_ "github.com/astaxie/beego/session/redis"
)

var mgr *session.Manager
var once sync.Once
var count int

func init() {
	// Here save session id in cookie, certainly you can also save it in URL query or request header anyway.
	mgrConf := &session.ManagerConfig{CookieName: "beego-session", ProviderConfig: "127.0.0.1:6379", EnableSetCookie: true, Maxlifetime: 3600}
	var err error
	mgr, err = session.NewManager("redis", mgrConf)
	if err != nil {
		panic(err)
	}
}

// Implement http.Handler interface.
var handle = func(w http.ResponseWriter, r *http.Request) {
	store, err := mgr.SessionStart(w, r)
	if err != nil {
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}
	// Don't forget this statement, it save session to redis.
	defer store.SessionRelease(w)

	// This only can do once
	once.Do(func() {
		if err := store.Set("name", "I am codeman"); err != nil {
			http.Error(w, "failed", http.StatusInternalServerError)
			return
		}

		if err := store.Set("count", count); err != nil {
			http.Error(w, "failed", http.StatusInternalServerError)
			return
		}
	})

	count = store.Get("count").(int)
	count++
	// Record in session every times.
	if err := store.Set("count", count); err != nil {
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}

	// Get the value from redis session every times after one excute Set function.
	name := store.Get("name").(string)
	w.Write([]byte(fmt.Sprintf("%s, you guys called me %d times", name, count)))
}

func main() {
	http.ListenAndServe(":9000", http.HandlerFunc(handle))
}
