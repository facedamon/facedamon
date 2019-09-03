package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"log"
	"net/http"
)

var db *sql.DB

func errHandler(err error) {
	if err != nil {
		log.Print(err)
	}
}

func SimplePage(w http.ResponseWriter, r *http.Request, template string) {
	req := render.New(render.Options{Directory: "templates"})
	req.HTML(w, http.StatusOK, template, nil)
}

func SimpleAuthenticatedPage(w http.ResponseWriter, r *http.Request, template string) {
	session := sessions.GetSession(r)
	sess := session.Get("hello")

	if sess == nil {
		http.Redirect(w, r, "/notauthenticated", 301)
	}
	req := render.New(render.Options{Directory: "templates"})
	req.HTML(w, http.StatusOK, template, nil)
}

func LoginPost(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	session.Set("hello", "world")
	username := r.FormValue("username")
	password := r.FormValue("password")

	var email string

	err := db.QueryRow("select email from users where username = ? and password = ?", username, password).Scan(&email)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(username + "用户不存在或密码错误")
		} else {
			log.Fatal(err)
		}
		http.Redirect(w, r, "/failedquery", 301)
	}
	fmt.Println(email)
	http.Redirect(w, r, "/home", 302)
}

func SignupPost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	stmt, err := db.Prepare("insert into users (username, password, email) value (?, ?, ?)")
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(username, password, email)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/login", 302)
}

func APIHandler(w http.ResponseWriter, r *http.Request) {
	data, _ := json.Marshal("{'API Test':'works!'}")
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write(data)
}

func StaticServer(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))).ServeHTTP(w, r)
}

func main() {
	db, _ = sql.Open("mysql", "damon:damon@tcp(127.0.0.1:3306)/negroni")
	db.Ping()
	defer db.Close()

	mux := http.NewServeMux()
	n := negroni.Classic()

	store := cookiestore.New([]byte("secret"))
	n.Use(sessions.Sessions("global_session_store", store))

	mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		StaticServer(w, r)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		SimplePage(w, r, "mainpage")
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			SimplePage(w, r, "login")
		} else if r.Method == http.MethodPost {
			LoginPost(w, r)
		}
	})

	mux.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			SimplePage(w, r, "signup")
		} else if r.Method == http.MethodPost {
			SignupPost(w, r)
		}
	})

	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		SimplePage(w, r, "logout")
	})

	mux.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		SimpleAuthenticatedPage(w, r, "home")
	})

	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		APIHandler(w, r)
	})

	n.UseHandler(mux)
	n.Run(":3000")
}
