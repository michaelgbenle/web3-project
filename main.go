package main

import "fmt"

func main() {
	r := mux.newRouter()
	r.HandleFunc("/", ).Methods("GET")
	r.HandleFunc("/", ).Methods("POST")
	r.HandleFunc("/", ).Methods("POST")
	fmt.Println("connectring to server")

}