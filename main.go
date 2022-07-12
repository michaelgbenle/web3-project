package main

import "log"

func main() {
	r := mux.newRouter()
	r.HandleFunc("/", GetBlockchain ).Methods("GET")
	r.HandleFunc("/", WriteBlock ).Methods("POST")
	r.HandleFunc("/", NewBook ).Methods("POST")

	log.Println("listening on port 2000")

}