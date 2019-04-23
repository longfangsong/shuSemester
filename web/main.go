package main

import (
	"log"
	"net/http"
	"os"
	"shuSemester/handler"
)

func main() {
	http.HandleFunc("/ping", handler.PingPongHandler)
	http.HandleFunc("/semester", handler.SemesterHandler)
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
