package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type user struct {
	name     string
	email    string
	password string
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

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
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(w, "ParseForm_error")
			return
		}
		for key, values := range r.Form {
			fmt.Printf("%s: %v\n", key, values)
		}
	}
}

func admin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is admin"))
}

func advert(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Объявление с id", id)
}

func enter(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("../ui/html/enter.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	if r.Method == http.MethodGet {
		log.Println("geeet")
		err = ts.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
		}
	}
	if r.Method == http.MethodPost {
		log.Println("poost")
		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(w, "ParseForm_error")
			return
		}
		//formData := make(map[string][]string)
		for key, values := range r.Form {
			//formData[key] = values
			fmt.Printf("%s: %v\n", key, values)
		}
		targetURL := "/admin"
		http.Redirect(w, r, targetURL, http.StatusFound)
	}
}

func registration(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("../ui/html/regist.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
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

		db, err := sql.Open("sqlite3", "../ui/static/All_db.db")
		if err != nil {
			panic(err)
		}
		defer db.Close()
		fmt.Println("insert into log_pswd (login, password, id_level, fio, tel) values ('" + formData["name"] + "', '" + formData["password"] + "', 1, NULL, NULL)")
		_, err = db.Exec("insert into log_pswd (login, password, id_level, fio, tel) values ('" + formData["name"] + "', '" + formData["password"] + "', 1, NULL, NULL)")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/admin", admin)
	mux.HandleFunc("/admin/advert", advert)
	mux.HandleFunc("/registration", registration)
	mux.HandleFunc("/enter", enter)

	log.Println("Запуск веб-сервера на http://127.0.0.1:4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
