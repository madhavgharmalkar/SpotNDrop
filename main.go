package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {

	db, err := ConnectDB("bcc49feaf6f6d4", "74b6ab4e", "us-cdbr-iron-east-03.cleardb.net", "3306", "ad_afa2c365c758c76")
	if err != nil {
		log.Panic(err)
	}

	log.Println("Database sucsuessfully connected!")

	r := mux.NewRouter()

	r.HandleFunc("/ws", websocketHandler)
	r.HandleFunc("/get/", db.DBGet()).Methods("GET")
	r.HandleFunc("/", indexHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	http.Handle("/", r)

	go updateListen(db)

	port := os.Getenv("VCAP_APP_PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("http server listening at localhost:" + port)
	http.ListenAndServe(":"+port, nil)

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, nil)
}

func handleWebErr(w http.ResponseWriter, err error) {
	http.Error(w, "Internal server error: "+err.Error(), 500)
}
