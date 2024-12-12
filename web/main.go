package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/ncruces/go-sqlite3"
	_ "github.com/ncruces/go-sqlite3/embed"
)

var store *sessions.CookieStore

func init() {
	os.Setenv("SECRET_SESSION", "secret_key")
	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	session, erro := store.Get(r, "session-name")
	if erro != nil {
		http.Error(w, erro.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["login"] = ""
	session.Values["id"] = ""
	// Используем функцию template.ParseFiles() для чтения файла шаблона.
	// Если возникла ошибка, мы запишем детальное сообщение ошибки и
	// используя функцию http.Error() мы отправим пользователю
	// ответ: 500 Internal Server Error (Внутренняя ошибка на сервере)
	ts, err := template.ParseFiles("../ui/html/home.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// Затем мы используем метод Execute() для записи содержимого
	// шаблона в тело HTTP ответа. Последний параметр в Execute() предоставляет
	// возможность отправки динамических данных в шаблон.
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func user(w http.ResponseWriter, r *http.Request) {
	session, erro := store.Get(r, "session-name")
	if erro != nil {
		http.Error(w, erro.Error(), http.StatusInternalServerError)
		return
	}
	data := map[string]string{
		"Login": session.Values["login"].(string),
	}
	ts, err := template.ParseFiles("../ui/html/user.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	ts.ExecuteTemplate(w, "user.tmpl", data)
}

func redir(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func enter(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("../ui/html/enter.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	if r.Method == http.MethodGet {
		err = ts.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
		}
	}
	if r.Method == http.MethodPost {
		session, erro := store.Get(r, "session-name")
		if erro != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(w, "ParseForm_error")
			return
		}
		formData := make(map[string][]string)
		for key, values := range r.Form {
			formData[key] = values
			fmt.Printf("%s: %v\n", key, values)
		}

		db, err := sqlite3.Open("../ui/static/All_db.db")
		if err != nil {
			panic(err)
		}
		defer db.Close()
		us_exist, _, err := db.Prepare("SELECT id_user, login, password from log_pswd WHERE login = '" + formData["name"][0] + "'")
		if err != nil {
			log.Fatal(err)
		}
		defer us_exist.Close()
		var ent_pass string
		var id_user string
		for us_exist.Step() {
			ent_pass = us_exist.ColumnText(2)
			id_user = us_exist.ColumnText(0)
		}
		if ent_pass == formData["password"][0] {
			session.Values["login"] = formData["name"][0]
			session.Values["id"] = id_user
			err := session.Save(r, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			redir(w, r, "/user")
		} else {
			redir(w, r, "/enter")
		}
	}
}

func registration(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("../ui/html/regist.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	if r.Method == http.MethodGet {
		err = ts.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
		}
	}
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(w, "ParseForm_error")
			return
		}
		formData := make(map[string]string)
		for key, values := range r.Form {
			formData[key] = values[0]
			fmt.Printf("%s: %v\n", key, values)
		}

		db, err := sqlite3.Open("../ui/static/All_db.db")
		if err != nil {
			panic(err)
		}
		defer db.Close()
		err = db.Exec("insert into log_pswd (login, password, id_level, fio, tel) values ('" + formData["name"] + "', '" + formData["password"] + "', 2, NULL, NULL)")
		if err != nil {
			log.Fatal(err)
		}
		redir(w, r, "/")
	}
}

func add_anim(w http.ResponseWriter, r *http.Request) {
	session, erro := store.Get(r, "session-name")
	if erro != nil {
		http.Error(w, erro.Error(), http.StatusInternalServerError)
		return
	}
	ts, err := template.ParseFiles("../ui/html/add_anim.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	if r.Method == http.MethodGet {
		err = ts.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
		}
	}
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(w, "ParseForm_error")
			return
		}
		formData := make(map[string]string)
		for key, values := range r.Form {
			formData[key] = values[0]
			fmt.Printf("%s: %v\n", key, values)
		}

		db, err := sqlite3.Open("../ui/static/All_db.db")
		if err != nil {
			panic(err)
		}
		defer db.Close()
		err = db.Exec(fmt.Sprintf("insert into anim (id_user, anim, poroda, klicka, vozrast) values ('%s', '%s', '%s', '%s', '%s')", session.Values["id"].(string), formData["animal"], formData["poroda"], formData["name"], formData["age"]))
		if err != nil {
			log.Fatal(err)
		}
		redir(w, r, "/user")
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/user", user)
	mux.HandleFunc("/registration", registration)
	mux.HandleFunc("/enter", enter)
	mux.HandleFunc("/add_anim", add_anim)

	log.Println("Запуск веб-сервера на http://127.0.0.1:4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
