package main

import (
	"GoSimpleProxy/lib"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func init() {
	if err := lib.ENV.Init(); err != nil {
		lib.LogFatalf("Error init env: %v \n", err)
	}

	lib.SetLogFileName("donor.log")
	lib.OpenLogFile()
}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hi, dear client"))
	})

	r.HandleFunc("/author", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Name: Igor\nSurname: Stcherbina\nProfession: Software engineer"))
	})

	http.Handle("/", r)
	lib.LogAppRun(lib.ENV.DONAR_PORT)
	log.Fatal(http.ListenAndServe(":"+lib.ENV.DONAR_PORT, lib.LogRequest(http.DefaultServeMux)))
}
