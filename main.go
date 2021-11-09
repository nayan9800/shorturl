package main

import (
	"log"
	"net/http"
	"time"

	"github.com/nayan9800/shorturl/router"
)

func main() {

	mux := http.NewServeMux()

	files := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", files))
	mux.HandleFunc("/", router.Method(router.Index, "GET"))
	mux.HandleFunc("/dashboard", router.Method(router.Dashboard, "GET"))
	mux.HandleFunc("/login", router.Loginpage)
	mux.HandleFunc("/logout", router.Logout)
	mux.HandleFunc("/signup", router.SignupPage)
	mux.HandleFunc("/addservice", router.AddService)
	mux.HandleFunc("/deletesvc", router.DeleteService)
	server := http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    time.Duration(10 * time.Second),
		WriteTimeout:   time.Duration(600 * time.Second),
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("Server is starting", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}

}
