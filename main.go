package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type ViewData struct {
	Student string
	Text    string
}

type LogData struct {
	Studenttext []Studenttext
}

type ResponeFishText struct {
	Status    string `json:"status"`
	Text      string `json:"text"`
	ErrorCode string `json:"errorCode"`
}

type Studenttext struct {
	Id   int
	Name string
	Text string
	Ip   string
}

type ResponseGen struct {
	Msg string `json:"msg"`
}

func main() {
	r := mux.NewRouter()
	//r.HandleFunc("/variant", VariantPage)
	r.HandleFunc("/variant/{login}", VariantPage)
	r.HandleFunc("/view", ViewVariantPage)
	http.Handle("/", r)

	http.ListenAndServe(":80", nil)
}

func ViewVariantPage(w http.ResponseWriter, r *http.Request) {
	var data LogData

	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})

	if err != nil {
		log.Println(err.Error())
	}

	var st []Studenttext

	db.Find(&st)

	data.Studenttext = st

	tmpl, err := template.ParseFiles("logs.html")

	if err != nil {
		log.Println(err.Error())
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err.Error())
	}
}

func VariantPage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.Form)
	var data ViewData
	data.Student = strings.Join(r.Form["name"], "")
	//resp, err := http.Get("https://fish-text.ru/get?type=sentence&number=1&format=json")
	resp, err := http.Post("https://randomall.ru/api/gens/2127", "", nil)
	if err != nil {
		data.Text = "Ошибка получения текста. Попробуйте позже"
	}
	var response ResponeFishText
	//var response ResponseGen
	json.NewDecoder(resp.Body).Decode(&response)
	data.Text = response.Text

	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	var txt Studenttext
	txt.Name = data.Student
	txt.Text = data.Text
	txt.Ip = r.RemoteAddr
	db.Create(txt)
	tmpl, _ := template.ParseFiles("index.html")
	tmpl.Execute(w, data)
}
