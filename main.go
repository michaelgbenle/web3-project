package main

import "fmt"

func main() {
	r := mux.newRouter()
	r.HandleFunc("/", GetBlockchain ).Methods("GET")
	r.HandleFunc("/", WriteBlock ).Methods("POST")
	r.HandleFunc("/", NewBook ).Methods("POST")

	fmt.Println("connecting to server")

}