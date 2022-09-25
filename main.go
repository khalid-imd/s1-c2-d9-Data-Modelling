package main

import (
	"context"
	"fmt"
	"log"

	// "math"
	"net/http"
	"personal-project/connection"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	route := mux.NewRouter()

	connection.DatabaseConnect()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/hi", helloworld).Methods("GET")

	route.HandleFunc("/home", home).Methods("GET")

	route.HandleFunc("/project", project).Methods("GET")

	route.HandleFunc("/submit", submit).Methods("POST")

	route.HandleFunc("/contact", contact).Methods("GET")

	route.HandleFunc("/delete/{index}", delete).Methods("GET")

	route.HandleFunc("/detail", detail).Methods("GET")

	fmt.Println("server is Running")
	http.ListenAndServe("localhost:8000", route)
}

func helloworld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset-utf8")

	var tmplt, err = template.ParseFiles("pages/home.html")

	if err != nil {
		w.Write([]byte("file doesn't exist: " + err.Error()))
		return
	}

	// response := map[string]interface{}{
	// 	"Project": projectData,
	// }

	// w.Write([]byte("home"))
	//w.WriteHeader(http.StatusAccepted)

	data, _ := connection.Conn.Query(context.Background(), "SELECT title, description FROM tb_project")
	fmt.Println(data)

	var result []Project
	for data.Next() {
		var each = Project{}

		var err = data.Scan(&each.Title, &each.Description)
		if err != nil {
			fmt.Println(err.Error)
			return
		}
		result = append(result, each)
	}

	resData := map[string]interface{}{
		"Project": result,
	}

	fmt.Println(result)

	tmplt.Execute(w, resData)
}

func project(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset-utf8")

	var tmplt, err = template.ParseFiles("pages/project.html")

	if err != nil {
		w.Write([]byte("file doesn't exist: " + err.Error()))
		return
	}

	// w.Write([]byte("home"))
	//w.WriteHeader(http.StatusAccepted)
	tmplt.Execute(w, "")
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset-utf8")

	var tmplt, err = template.ParseFiles("pages/contact.html")

	if err != nil {
		w.Write([]byte("file doesn't exist: " + err.Error()))
		return
	}

	// w.Write([]byte("home"))
	//w.WriteHeader(http.StatusAccepted)
	tmplt.Execute(w, "")
}

func detail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset-utf8")

	var tmplt, err = template.ParseFiles("pages/detail.html")

	if err != nil {
		w.Write([]byte("file doesn't exist: " + err.Error()))
		return
	}

	// w.Write([]byte("home"))
	//w.WriteHeader(http.StatusAccepted)
	tmplt.Execute(w, "")
}

type Project struct {
	Title       string
	StartDate   string
	EndDate     string
	Description string
	Duration    string
}

var projectData = []Project{}

func submit(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("title: " + r.PostForm.Get("addTitle"))
	// fmt.Println("start date: " + r.PostForm.Get("addStartDate"))
	// fmt.Println("end date: " + r.PostForm.Get("addEndDate"))
	// fmt.Println("description: " + r.PostForm.Get("addDescription"))
	// fmt.Println(r.PostForm.Get("addNode"))
	// fmt.Println(r.PostForm.Get("addReact"))
	// fmt.Println(r.PostForm.Get("addNext"))
	// fmt.Println(r.PostForm.Get("addTypeScript"))

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)

	title := r.PostForm.Get("addTitle")
	startDate := r.PostForm.Get("addStartDate")
	endDate := r.PostForm.Get("addEndDate")
	description := r.PostForm.Get("addDescription")

	layout := "2006-01-02"
	parsingstartdate, _ := time.Parse(layout, startDate)
	parsingenddate, _ := time.Parse(layout, endDate)

	hours := parsingenddate.Sub(parsingstartdate).Hours()
	days := hours / 24
	// weeks := math.Round(days / 7)
	// month := math.Round(days / 30)
	// year := math.Round(days / 365)

	var duration string
	// if year > 0 {
	// 	duration = strconv.FormatFloat(year, 'f', 0, 64) + " years"
	// } else if month > 0 {
	// 	duration = strconv.FormatFloat(month, 'f', 0, 64) + " month"
	// } else if weeks > 0 {
	// 	duration = strconv.FormatFloat(weeks, 'f', 0, 64) + " weeks"
	// } else
	if days > 0 {
		duration = strconv.FormatFloat(days, 'f', 0, 64) + " days"
	}

	newProject := Project{
		Title:       title,
		StartDate:   startDate,
		EndDate:     endDate,
		Duration:    duration,
		Description: description,
	}
	projectData = append(projectData, newProject)

	// fmt.Println(projectData)
}

func delete(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(mux.Vars(r)["index"])
	projectData = append(projectData[:index], projectData[index+1:]...)
	http.Redirect(w, r, "/home", http.StatusFound)
}
