package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/shenghui0779/yiigo"
)

var sessionStore = sessions.NewCookieStore([]byte("t09-secret"))
var client *redis.Client
var templates *template.Template

func main() {
	r := mux.NewRouter()
	templates = template.Must(template.ParseGlob("./templates/*.html"))
	r.HandleFunc("/", indexGetHandler).Methods("GET")
	r.HandleFunc("/", indexPostHandler).Methods("POST")
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
	r.HandleFunc("/register", registerGetHandler).Methods("GET")
	r.HandleFunc("/register", registerPostHandler).Methods("POST")
	http.Handle("/", r)
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	env := yiigo.Env("app.env").String()
	fmt.Println(env)

	http.ListenAndServe(":8080", nil)
}

type Todo struct {
	Title string
	Items []string
}
type  Book struct {
	Book string
	Price int64
}
func indexGetHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("hhhh")

	b1:=Book{}
	err:=yiigo.DB().Get(&b1,"select  * from book where book=?","go-zero")

	if err != nil{
		panic(err)
	}
	fmt.Println(b1)
	values, _ := client.LRange("todos", 0, 10).Result()
	templates.ExecuteTemplate(w, "index.html", Todo{
		Title: "首页",
		Items: values,
	})

}
func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", "欢迎来到登陆页")
}
func registerGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "register.html", "欢迎来到注册页")
}

func indexPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.FormValue("title")

	if len(title) == 0 {
		log.Println(title)
		return
	}
	client.LPush("todos", title)
	http.Redirect(w, r, "/", http.StatusFound)

}
func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	pwd := r.FormValue("pwd")
	log.Println(username, pwd)
	store := sessions.NewSession(sessionStore, username)
	sessionStore.Save(r, w, store)

}
func registerPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	pwd := r.FormValue("pwd")
	client.HSet("user", "username", username)
	client.HSet("user", "pwd", pwd)
	http.Redirect(w, r, "/login", http.StatusFound)

}
